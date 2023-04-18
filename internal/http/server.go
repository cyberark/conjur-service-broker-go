// Package http implements the http communication layer of the conjur service broker
package http

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/cyberark/conjur-api-go/conjurapi/logging"
	"github.com/cyberark/conjur-service-broker/internal/ctxutil"
	"github.com/cyberark/conjur-service-broker/internal/servicebroker"
	"github.com/cyberark/conjur-service-broker/pkg/conjur"
	"github.com/gin-contrib/requestid"
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"go.uber.org/zap"
	"go.uber.org/zap/zapio"
)

const (
	httpTimeout     = time.Minute      // http timeout
	httpIdleTimeout = 15 * time.Minute // keep-alive timeout
	serviceName     = "conjure-service-broker"
)

// StartHTTPServer starts a new http server to handle requests supported by the service broker
func StartHTTPServer() {
	if logger, err := startServer(); err != nil {
		logger.Sugar().Fatal("failed to start http server: ", err)
	}
}

func startServer() (logger *zap.Logger, err error) {
	ctx := ctxutil.NewContext()
	logCfg := zap.NewProductionConfig()
	cfg, err := newConfig()
	if err != nil {
		err = fmt.Errorf("failed to parse configuration: %w", err)
		logger, _ = logCfg.Build()
		return
	}
	if !cfg.Debug {
		gin.SetMode(gin.ReleaseMode)
	} else {
		logging.ApiLog.Level = logrus.DebugLevel
		logCfg.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	}

	logger, err = logCfg.Build(zap.Fields(zap.String("service", serviceName)))
	if err != nil {
		err = fmt.Errorf("failed to init logger: %w", err)
		logger, _ = logCfg.Build()
		return
	}
	undo := zap.RedirectStdLog(logger)
	defer undo()
	writer := &zapio.Writer{Log: logger, Level: zap.DebugLevel}
	defer writer.Close()
	gin.DefaultWriter = writer
	logging.ApiLog.Out = writer

	ctx = ctx.WithLogger(logger.Sugar())
	ctx = ctx.WithEnableSpaceIdentity(cfg.EnableSpaceIdentity)
	client, err := conjur.NewClient(&cfg.Config)
	if err != nil {
		err = fmt.Errorf("failed to initialize conjur client: %w", err)
		return
	}
	if err = client.ValidateConnectivity(); err != nil {
		err = fmt.Errorf("failed to validate conjur client: %w", err)
		return
	}
	srv := servicebroker.NewServerImpl(client)

	r := gin.New()
	r.Use(ginzap.Ginzap(logger, time.RFC3339, true), ginzap.RecoveryWithZap(logger, true), requestid.New(), errorsMiddleware)
	// TODO: trusted proxies, caching headers

	if len(cfg.SecurityUserName) > 0 { // gin basic auth middleware will fail on empty username
		r.Use(gin.BasicAuth(gin.Accounts{cfg.SecurityUserName: cfg.SecurityUserPassword}))
	}
	validator, err := validatorMiddleware(ctx)
	if err != nil {
		err = fmt.Errorf("failed to initialize validateion middleware: %w", err)
		return
	}
	r.Use(ctx.Inject(), validator)

	r = servicebroker.RegisterHandlers(r, srv)
	httpSrv := &http.Server{
		Addr:              ":" + cfg.Port,
		Handler:           r,
		ReadTimeout:       httpTimeout,
		WriteTimeout:      httpTimeout,
		ReadHeaderTimeout: httpTimeout,
		IdleTimeout:       httpIdleTimeout,
		MaxHeaderBytes:    1 << 20, // TODO: is 1MB enough?
	}
	logger.Info("server starting...")
	go func() {
		if err = httpSrv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Sugar().Fatal("failed to start server: ", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	// kill (no param) default send syscanll.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall. SIGKILL but can"t be catch, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Info("shutdown server...")

	c, cancel := context.WithTimeout(ctx, 5*time.Second)
	go func() {
		defer cancel()
		if err = httpSrv.Shutdown(c); err != nil {
			logger.Sugar().Fatal("failed on server shutdown: ", err)
		}
	}()

	// catching ctx.Done(). timeout of 5 seconds.
	<-c.Done()
	logger.Info("server exit")
	return
}

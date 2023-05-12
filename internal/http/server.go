// Package http implements the http communication layer of the conjur service broker
package http

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/cyberark/conjur-api-go/conjurapi/logging"
	"github.com/gin-contrib/requestid"
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.cyberng.com/Conjur-Enterprise/conjur-service-broker-go/internal/ctxutil"
	"github.cyberng.com/Conjur-Enterprise/conjur-service-broker-go/internal/servicebroker"
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
	cfg, err := newConfig()
	checkFatalErr(err)
	logger, cleanup, err := initLogger(cfg)
	checkFatalErr(err)
	defer cleanup()

	ctx := initCtx(logger, cfg)

	srv, err := initServer(cfg)
	checkFatalErr(err)

	err = startServer(ctx, cfg, srv, logger)
	checkFatalErr(err)
}

func initServer(cfg *config) (servicebroker.ServerInterface, error) {
	client, err := cfg.NewClient()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize conjur client: %w", err)
	}
	if err = client.ValidateConnectivity(); err != nil {
		return nil, fmt.Errorf("failed to validate conjur client: %w", err)
	}
	return servicebroker.NewServerImpl(client), nil
}

func initCtx(logger *zap.Logger, cfg *config) ctxutil.Context {
	ctx := ctxutil.NewContext()
	if logger != nil {
		ctx = ctx.WithLogger(logger.Sugar())
	}
	if cfg != nil {
		ctx = ctx.WithEnableSpaceIdentity(cfg.EnableSpaceIdentity)
	}
	return ctx
}

func initLogger(cfg *config) (logger *zap.Logger, cleanup func(), err error) {
	logCfg := zap.NewProductionConfig()
	logCfg.Encoding = "console"
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
	writer := &zapio.Writer{Log: logger, Level: zap.DebugLevel}
	gin.DefaultWriter = writer
	logging.ApiLog.Out = writer

	cleanup = func() {
		_ = logger.Sync()
		undo()
		_ = writer.Close()
	}

	return logger, cleanup, err
}

func startServer(ctx ctxutil.Context, cfg *config, srv servicebroker.ServerInterface, logger *zap.Logger) error {
	r := gin.New()
	r.Use(ginzap.Ginzap(logger, time.RFC3339, true), ginzap.RecoveryWithZap(logger, true), requestid.New(), errorsMiddleware)

	if len(cfg.SecurityUserName) > 0 { // gin basic auth middleware will fail on empty username
		r.Use(gin.BasicAuth(gin.Accounts{cfg.SecurityUserName: cfg.SecurityUserPassword}))
	}
	validator, err := validatorMiddleware(ctx)
	if err != nil {
		return fmt.Errorf("failed to initialize validateion middleware: %w", err)
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
		MaxHeaderBytes:    1 << 20,
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
	// kill -9 is syscall. SIGKILL can't be caught, no need to add it
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
	return nil
}

func checkFatalErr(err error) {
	if err != nil {
		log.Fatal(fmt.Errorf("failed to start server: %w", err))
	}
}

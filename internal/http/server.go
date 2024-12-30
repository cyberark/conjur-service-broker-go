// Package http implements the http communication layer of the conjur service broker
package http

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/cyberark/conjur-api-go/conjurapi/logging"
	"github.com/cyberark/conjur-service-broker-go/internal/ctxutil"
	"github.com/cyberark/conjur-service-broker-go/internal/servicebroker"
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
)

// StartHTTPServer starts a new http server to handle requests supported by the service broker
func StartHTTPServer(httpClient *http.Client) {
	cfg, err := newConfig()
	checkFatalErr(err)
	logger, cleanup, err := initLogger(cfg)
	checkFatalErr(err)
	defer cleanup()

	ctx := initCtx(logger, cfg)

	srv, err := initServer(cfg, httpClient)
	checkFatalErr(err)

	err = startServer(ctx, cfg, srv, logger)
	checkFatalErr(err)
}

func initServer(cfg *config, httpClient *http.Client) (servicebroker.ServerInterface, error) {
	client, err := cfg.NewClient(httpClient)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize conjur client: %s", err)
	}
	if err = client.ValidateConnectivity(); err != nil {
		return nil, fmt.Errorf("failed to validate conjur client: %s", err)
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
	logger, err = logCfg.Build()
	if err != nil {
		err = fmt.Errorf("failed to init logger: %s", err)
		logger, _ = logCfg.Build()
		return
	}
	undo := zap.RedirectStdLog(logger)
	// we set the debug level for all the logs that would be written on the writer since writer has no mean to understand
	// the actual level, in non DEBUG mode (production) nothing from the writer would get logged
	writer := &zapio.Writer{Log: logger, Level: zap.DebugLevel}
	gin.DefaultWriter = writer

	// Configure conjur-api-go logging (logrus) to go through zap
	logging.ApiLog.ReportCaller = true   // So Zap reports the right caller
	logging.ApiLog.SetOutput(io.Discard) // Prevent logrus from writing its logs
	hook, _ := NewZapHook(logger)
	logging.ApiLog.AddHook(hook)

	cleanup = func() {
		_ = logger.Sync()
		undo()
		_ = writer.Close()
	}

	return logger, cleanup, err
}

func startServer(ctx ctxutil.Context, cfg *config, srv servicebroker.ServerInterface, logger *zap.Logger) error {
	r := gin.New()
	r.Use(
		ginzap.Ginzap(logger, time.RFC3339, true),
		ginzap.RecoveryWithZap(logger, true),
		requestid.New(),
		errorsMiddleware,
	)

	if len(cfg.SecurityUserName) > 0 { // gin basic auth middleware will fail on empty username
		r.Use(gin.BasicAuth(gin.Accounts{cfg.SecurityUserName: cfg.SecurityUserPassword}))
	}
	validator, err := validatorMiddleware(ctx)
	if err != nil {
		return fmt.Errorf("failed to initialize validator middleware: %s", err)
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
		log.Fatal(fmt.Errorf("failed to start server: %s", err))
	}
}

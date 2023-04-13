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

	"github.com/cyberark/conjur-service-broker/internal/ctxutil"
	"github.com/cyberark/conjur-service-broker/internal/servicebroker"
	"github.com/cyberark/conjur-service-broker/pkg/conjur"
	"github.com/gin-contrib/requestid"
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

const (
	httpTimeout     = time.Minute      // http timeout
	httpIdleTimeout = 15 * time.Minute // keep-alive timeout
)

// StartHTTPServer starts a new http server to handle requests supported by the service broker
func StartHTTPServer(logger *zap.Logger) error {
	ctx := ctxutil.NewContext()
	cfg, err := newConfig()
	if err != nil {
		return fmt.Errorf("failed to parse configuration: %w", err)
	}
	ctx = ctx.WithLogger(logger.Sugar())
	ctx = ctx.WithEnableSpaceIdentity(cfg.EnableSpaceIdentity)
	client, err := conjur.NewClient(&cfg.Config)
	if err != nil {
		return fmt.Errorf("failed to initialize conjur client: %w", err)
	}
	if err := client.ValidateConnectivity(); err != nil {
		return fmt.Errorf("failed to validate conjur client: %w", err)
	}
	srv := servicebroker.NewServerImpl(client)

	if !cfg.Debug {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.New()
	r.Use(ginzap.Ginzap(logger, time.RFC3339, true), ginzap.RecoveryWithZap(logger, true), requestid.New(), errorsMiddleware)
	// TODO: trusted proxies, caching headers

	if len(cfg.SecurityUserName) > 0 { // gin basic auth middleware will fail on empty username
		r.Use(gin.BasicAuth(gin.Accounts{cfg.SecurityUserName: cfg.SecurityUserPassword}))
	}
	validator, err := validatorMiddleware(ctx)
	if err != nil {
		return err
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
	go func() {
		if err = httpSrv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(fmt.Errorf("failed to start server: %w", err))
		}
	}()

	quit := make(chan os.Signal, 1)
	// kill (no param) default send syscanll.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall. SIGKILL but can"t be catch, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("shutdown server ...")

	c, cancel := context.WithTimeout(ctx, 5*time.Second)
	if err := httpSrv.Shutdown(c); err != nil {
		cancel()
		log.Fatal("failed on server shutdown:", err)
	}
	// catching ctx.Done(). timeout of 5 seconds.
	<-c.Done()
	log.Println("server exit")
	cancel()
	return nil
}

package http

import (
	"fmt"

	"github.com/cyberark/conjur-service-broker/internal/ctxutil"
	"github.com/cyberark/conjur-service-broker/internal/servicebroker"
	"github.com/cyberark/conjur-service-broker/pkg/conjur"
	"github.com/gin-gonic/gin"
)

// StartHTTPServer starts a new http server to handle requests supported by the service broker
func StartHTTPServer() error {
	ctx := ctxutil.NewContext()
	//TODO: add logger???
	cfg, err := newConfig()
	if err != nil {
		return fmt.Errorf("failed to parse configuration: %w", err)
	}
	ctx.WithEnableSpaceIdentity(cfg.EnableSpaceIdentity)
	client, err := conjur.NewClient(&cfg.Config)
	if err != nil {
		return fmt.Errorf("failed to initialize conjur client: %w", err)
	}
	if err := client.ValidateConnectivity(); err != nil {
		return fmt.Errorf("failed to validate conjur client: %w", err)
	}
	srv := servicebroker.NewServerImpl(client)

	// TODO: make this production grade gin
	if !cfg.Debug {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.Default()

	r.Use(errorsMiddleware)

	if len(cfg.SecurityUserName) > 0 { // gin basic auth middleware will fail on empty username
		r.Use(gin.BasicAuth(gin.Accounts{cfg.SecurityUserName: cfg.SecurityUserPassword}))
	}
	validator, err := validatorMiddleware(ctx)
	if err != nil {
		return err
	}
	r.Use(ctx.Inject(), validator)

	r = servicebroker.RegisterHandlers(r, srv)
	// TODO: graceful shutdown
	err = r.Run()
	if err != nil {
		return fmt.Errorf("failed to start server: %w", err)
	}
	return nil
}

package http

import (
	"context"
	"fmt"

	"github.com/cyberark/conjur-service-broker/internal/servicebroker"
	"github.com/cyberark/conjur-service-broker/pkg/conjur"
	"github.com/gin-gonic/gin"
)

// StartHTTPServer starts a new http server to handle requests supported by the service broker
func StartHTTPServer() error {
	//TODO: add logger???
	cfg, err := newConfig()
	if err != nil {
		return fmt.Errorf("failed to parse configuration: %w", err)
	}
	ctx := context.Background()
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
	openAPI, err := servicebroker.GetSwagger()
	if err != nil {
		return fmt.Errorf("failed to initialize open api validator: %w", err)
	}
	openAPI.Servers = nil // temporary workaround: https://github.com/deepmap/oapi-codegen/issues/710

	err = openAPI.Validate(ctx)
	if err != nil {
		return fmt.Errorf("failed to validate open api definition: %w", err)
	}

	validator, err := OpenAPIValidator(openAPI)
	if err != nil {
		return fmt.Errorf("failed to create open api validator: %w", err)
	}

	r.Use(validator, errorsMiddleware)
	if len(cfg.SecurityUserName) > 0 { // gin basic auth middleware will fail on empty username
		r.Use(gin.BasicAuth(gin.Accounts{cfg.SecurityUserName: cfg.SecurityUserPassword}))
	}
	r = servicebroker.RegisterHandlers(r, srv)
	// TODO: graceful shutdown
	err = r.Run()
	if err != nil {
		return fmt.Errorf("failed to start server: %w", err)
	}
	return nil
}

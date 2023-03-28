package main

import (
	"fmt"
	"log"

	"github.com/cyberark/conjur-service-broker/conjur"
	// "net/http"

	middleware "github.com/deepmap/oapi-codegen/pkg/gin-middleware"
	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/gin-gonic/gin"
)

// TODO: move service code to ./cmd

//go:generate oapi-codegen --config ./oapi-codegen.yaml ./api/openapi.yaml
func main() {
	//TODO: add logger???
	cfg, err := newConfig()
	if err != nil {
		log.Fatal(err)
	}
	//ctx := context.Background()
	client, err := conjur.NewClient(&cfg.Config)
	if err != nil {
		log.Fatal(err)
	}
	if err := client.Validate(); err != nil {
		log.Fatal(err)
	}
	srv := ServerImpl{
		client: client,
	}

	// TODO: make this production grade gin
	if !cfg.Debug {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.Default()
	openAPI, err := GetSwagger()
	if err != nil {
		log.Fatal(err)
	}
	openAPI.Servers = nil // temporary workaround: https://github.com/deepmap/oapi-codegen/issues/710

	r.Use(middleware.OapiRequestValidatorWithOptions(openAPI, validatorOpts()))
	if len(cfg.SecurityUserName) > 0 { // gin basic auth middleware will fail on empty username
		r.Use(gin.BasicAuth(gin.Accounts{cfg.SecurityUserName: cfg.SecurityUserPassword}))
	}
	r = RegisterHandlers(r, &srv)
	// TODO: graceful shutdown
	err = r.Run()
	if err != nil {
		log.Fatal(fmt.Errorf("failed to start server: %w", err))
	}
}

func validatorOpts() *middleware.Options {
	// this is needed to satisfy schema validator since it requires authentication func,
	// the actual authorization is gone on gin, due to the issues on handling http error codes
	// https://github.com/getkin/kin-openapi/issues/479
	return &middleware.Options{
		Options: openapi3filter.Options{
			AuthenticationFunc: openapi3filter.NoopAuthenticationFunc,
		},
	}
}

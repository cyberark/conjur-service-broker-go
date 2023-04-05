package http

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/cyberark/conjur-service-broker/internal/servicebroker"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/getkin/kin-openapi/routers"
	"github.com/getkin/kin-openapi/routers/gorillamux"
	"github.com/gin-gonic/gin"
)

func validatorMiddleware(ctx context.Context) (gin.HandlerFunc, error) {
	openAPI, err := servicebroker.GetSwagger()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize open api validator: %w", err)
	}
	openAPI.Servers = nil // temporary workaround: https://github.com/deepmap/oapi-codegen/issues/710

	err = openAPI.Validate(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to validate open api definition: %w", err)
	}

	validator, err := openAPIValidator(openAPI)
	if err != nil {
		return nil, fmt.Errorf("failed to create open api validator: %w", err)
	}
	return validator, nil
}

func openAPIValidator(spec *openapi3.T) (gin.HandlerFunc, error) {
	ctx := context.Background()

	router, err := gorillamux.NewRouter(spec)
	if err != nil {
		return nil, err
	}
	return func(c *gin.Context) {

		route, pathParams, err := router.FindRoute(c.Request)

		if err != nil {
			if errors.Is(err, routers.ErrMethodNotAllowed) {
				c.AbortWithStatusJSON(http.StatusMethodNotAllowed, gin.H{"error": "methodNotAllowed"})
				return
			}
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "notFound"})
			return
		}
		err = openapi3filter.ValidateRequest(ctx, &openapi3filter.RequestValidationInput{
			Request:    c.Request,
			PathParams: pathParams,
			Route:      route,
			Options:    validatorOpts(),
		})
		if err != nil {
			errMsg := errMsg(err)
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "validationError", "description": errMsg})
			return
		}
		c.Next()
	}, nil
}

func errMsg(err error) string {
	switch e := err.(type) {
	case openapi3.MultiError:
		errMsgs := make([]string, len(e))
		for _, errs := range e {
			errMsgs = append(errMsgs, errMsg(errs))
		}
		return strings.Join(errMsgs, " ")
	case *openapi3filter.RequestError:
		if e.Err == nil {
			return e.Error()
		}
		return fmt.Sprintf("%s %s", e.Reason, errMsg(e.Err))
	case *openapi3.SchemaError:
		return e.Reason
	default:
		return err.Error()
	}
}

func validatorOpts() *openapi3filter.Options {
	// this is needed to satisfy schema validator since it requires authentication func,
	// the actual authorization is done in gin, due to the issues on handling http error codes
	// https://github.com/getkin/kin-openapi/issues/479
	return &openapi3filter.Options{
		IncludeResponseStatus: true,
		MultiError:            true,
		AuthenticationFunc:    openapi3filter.NoopAuthenticationFunc,
	}
}

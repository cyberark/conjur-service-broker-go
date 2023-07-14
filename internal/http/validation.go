// Package http implements the http communication layer of the conjur service broker
package http

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/getkin/kin-openapi/routers"
	"github.com/getkin/kin-openapi/routers/gorillamux"
	"github.com/gin-gonic/gin"
	"github.cyberng.com/Conjur-Enterprise/conjur-service-broker-go/internal/servicebroker"
)

const (
	expectedServiceID = "c024e536-6dc4-45c6-8a53-127e7f8275ab"
	expectedPlanID    = "3a116ac2-fc8b-496f-a715-e9a1b205d05c.community"
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

	validator, err := openAPIValidator(addCustomValidations(openAPI))
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
				c.AbortWithStatusJSON(http.StatusMethodNotAllowed, gin.H{"error": "MethodNotAllowed"})
				return
			}
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "NotFound"})
			return
		}
		err = openapi3filter.ValidateRequest(ctx, &openapi3filter.RequestValidationInput{
			Request:    c.Request,
			PathParams: pathParams,
			Route:      route,
			Options:    validatorOpts(),
		})
		if err != nil {
			c.AbortWithStatusJSON(errorCode(err), gin.H{"error": "ValidationError", "description": err.Error()})
			return
		}
		c.Next()
	}, nil
}

func errorCode(err error) int {
	var e *openapi3filter.RequestError
	if errors.Is(err, openapi3filter.ErrInvalidRequired) && errors.As(err, &e) {
		if e.Parameter != nil && e.Parameter.Name == "X-Broker-API-Version" {
			return http.StatusPreconditionFailed
		}
	}
	return http.StatusBadRequest
}

func validatorOpts() *openapi3filter.Options {
	// this is needed to satisfy schema validator since it requires authentication func,
	// the actual authorization is done in gin, due to the issues on handling http error codes
	// https://github.com/getkin/kin-openapi/issues/479
	opts := &openapi3filter.Options{
		IncludeResponseStatus: true,
		MultiError:            true,
		AuthenticationFunc:    openapi3filter.NoopAuthenticationFunc,
	}
	opts.WithCustomSchemaErrorFunc(func(err *openapi3.SchemaError) string {
		return err.Reason
	})
	return opts
}

func addCustomValidations(api *openapi3.T) *openapi3.T {
	api = addCustomQueryParamValidation(api, "/v2/service_instances/{instance_id}", http.MethodDelete)
	api = addCustomBodyParamValidation(api, "/v2/service_instances/{instance_id}", http.MethodPut)
	api = addCustomBodyParamValidation(api, "/v2/service_instances/{instance_id}", http.MethodPatch)

	api = addCustomQueryParamValidation(api, "/v2/service_instances/{instance_id}/service_bindings/{binding_id}", http.MethodDelete)
	api = addCustomBodyParamValidation(api, "/v2/service_instances/{instance_id}/service_bindings/{binding_id}", http.MethodPut)

	api = addCustomBodyParamEmptyValidation(api, "/v2/service_instances/{instance_id}", http.MethodPut, "parameters")
	api = addCustomBodyParamEmptyValidation(api, "/v2/service_instances/{instance_id}", http.MethodPatch, "parameters")
	return api
}

func addCustomQueryParamValidation(api *openapi3.T, path, method string) *openapi3.T {
	parameters := api.Paths[path].GetOperation(method).Parameters
	parameters.GetByInAndName("query", "service_id").Schema.Value.Enum = []interface{}{expectedServiceID}
	parameters.GetByInAndName("query", "plan_id").Schema.Value.Enum = []interface{}{expectedPlanID}
	return api
}

func addCustomBodyParamValidation(api *openapi3.T, path, method string) *openapi3.T {
	content := api.Paths[path].GetOperation(method).RequestBody.Value.Content
	content.Get("application/json").Schema.Value.Properties["service_id"].Value.Enum = []interface{}{expectedServiceID}
	content.Get("application/json").Schema.Value.Properties["plan_id"].Value.Enum = []interface{}{expectedPlanID}
	return api
}

func addCustomBodyParamEmptyValidation(api *openapi3.T, path, method, param string) *openapi3.T {
	var zero uint64
	content := api.Paths[path].GetOperation(method).RequestBody.Value.Content
	content.Get("application/json").Schema.Value.Properties[param] = &openapi3.SchemaRef{
		Value: &openapi3.Schema{
			Type:     "object",
			MaxProps: &zero,
		},
	}
	return api
}

package main

import (
	"context"
	"io"
	"net/http"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/getkin/kin-openapi/routers"
	"github.com/getkin/kin-openapi/routers/gorillamux"
)

type validator struct {
	ignoreRouteError bool
	router           routers.Router
}

func newValidator() *validator {
	spec, err := openapi3.NewLoader().LoadFromFile("openapi.yaml")
	if err != nil {
		panic(err)
	}
	spec.Servers = nil
	router, err := gorillamux.NewRouter(spec)
	if err != nil {
		panic(err)
	}
	return &validator{
		router: router,
	}
}

func (v *validator) validateResponse(req *http.Request, resp *http.Response, body string) error {
	route, pathParams, err := v.router.FindRoute(req)
	if err != nil {
		if v.ignoreRouteError {
			return nil
		}
		return err
	}
	err = openapi3filter.ValidateResponse(context.Background(), &openapi3filter.ResponseValidationInput{
		RequestValidationInput: &openapi3filter.RequestValidationInput{
			Request:    req,
			PathParams: pathParams,
			// QueryParams:  nil,
			Route: route,
			Options: &openapi3filter.Options{
				MultiError: true,
			},
			// ParamDecoder: nil,
		},
		Status: resp.StatusCode,
		Header: resp.Header,
		Body:   io.NopCloser(strings.NewReader(body)),
		Options: &openapi3filter.Options{
			MultiError: true,
		},
	})
	if err != nil {
		return err
	}
	return nil
}

func (v *validator) iIgnoreRouteNotExists() error {
	v.ignoreRouteError = true
	return nil
}

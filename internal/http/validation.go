package http

import (
	"context"
	"errors"
	"net/http"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/getkin/kin-openapi/routers"
	"github.com/getkin/kin-openapi/routers/gorillamux"
	"github.com/gin-gonic/gin"
)

func OpenAPIValidator(spec *openapi3.T) (gin.HandlerFunc, error) {
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
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "badRequest", "description": err.Error()})
			return
		}
		c.Next()
	}, nil
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

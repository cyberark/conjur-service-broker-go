package http

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func Test_validatorMiddleware_header(t *testing.T) {
	type want struct {
		code int
		body string
	}
	tests := []struct {
		name   string
		header http.Header
		want   want
	}{{
		"positive",
		http.Header{"X-Broker-Api-Version": []string{"2.17"}},
		want{
			code: 200,
			body: "{}",
		},
	}, {
		"negative",
		http.Header{},
		want{
			code: 412,
			body: `{"description":"parameter \"X-Broker-API-Version\" in header has an error: value is required but missing", "error":"ValidationError"}`,
		},
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gin.SetMode(gin.TestMode)
			router := gin.New()
			middleware, err := validatorMiddleware(context.Background())
			require.NoError(t, err)
			router.Use(middleware)
			router.NoRoute(func(c *gin.Context) { // handle all traffic
				c.JSON(200, gin.H{})
			})

			w := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodGet, "/v2/catalog", nil)
			req.Header = tt.header
			req.Header.Set("Origin", "https://example.com")
			router.ServeHTTP(w, req)
			require.JSONEq(t, tt.want.body, w.Body.String())
			require.Equal(t, tt.want.code, w.Result().StatusCode)
		})
	}
}

func Test_validatorMiddleware_body(t *testing.T) {
	type want struct {
		code int
		body string
	}
	type request struct {
		method string
		path   string
		body   string
	}
	tests := []struct {
		name    string
		request request
		want    want
	}{{
		"provision - positive",
		request{
			method: "PUT",
			path:   "/v2/service_instances/9b292a9c-af66-4797-8d98-b30801f32ax7",
			body: `{
      "context": {
        "organization_guid": "e027f3f6-80fe-4d22-9374-da23a035ba0b",
        "space_guid": "8c56f85c-c16e-4158-be79-5dac74f970de",
        "organization_name": "my-organization",
        "space_name": "my-space"
      },
      "service_id": "c024e536-6dc4-45c6-8a53-127e7f8275ab",
      "plan_id": "3a116ac2-fc8b-496f-a715-e9a1b205d05c.community",
      "organization_guid": "e027f3f6-80fe-4d22-9374-da23a035ba0a",
      "space_guid": "8c56f85c-c16e-4158-be79-5dac74f970db",
      "parameters": {
      }
    }`,
		},
		want{
			code: 200,
			body: "{}",
		},
	}, {
		"provision - invalid plan and service",
		request{
			method: "PUT",
			path:   "/v2/service_instances/9b292a9c-af66-4797-8d98-b30801f32ax7",
			body: `{
      "context": {
        "organization_guid": "e027f3f6-80fe-4d22-9374-da23a035ba0b",
        "space_guid": "8c56f85c-c16e-4158-be79-5dac74f970de"
      },
      "service_id": "service",
      "plan_id": "plan",
      "organization_guid": "e027f3f6-80fe-4d22-9374-da23a035ba0a",
      "space_guid": "8c56f85c-c16e-4158-be79-5dac74f970db",
      "parameters": {
      }
    }`,
		},
		want{
			code: 400,
			body: `{"description":"request body has an error: doesn't match schema #/components/schemas/ServiceInstanceProvisionRequestBody: value is not one of the allowed values [\"3a116ac2-fc8b-496f-a715-e9a1b205d05c.community\"] | value is not one of the allowed values [\"c024e536-6dc4-45c6-8a53-127e7f8275ab\"]", "error": "ValidationError"}`,
		},
	}, {
		"provision - missing org and space id",
		request{
			method: "PUT",
			path:   "/v2/service_instances/9b292a9c-af66-4797-8d98-b30801f32ax7",
			body: `{
      "context": {},
      "service_id": "c024e536-6dc4-45c6-8a53-127e7f8275ab",
      "plan_id": "3a116ac2-fc8b-496f-a715-e9a1b205d05c.community",
      "parameters": {
      }
    }`,
		},
		want{
			code: 400,
			body: `{"description":"request body has an error: doesn't match schema #/components/schemas/ServiceInstanceProvisionRequestBody: property \"organization_guid\" is missing | property \"space_guid\" is missing", "error":"ValidationError"}`,
		},
	}, {
		"provision - empty body",
		request{
			method: "PUT",
			path:   "/v2/service_instances/9b292a9c-af66-4797-8d98-b30801f32ax7",
			body:   "{}",
		},
		want{
			code: 400,
			body: `{"description":"request body has an error: doesn't match schema #/components/schemas/ServiceInstanceProvisionRequestBody: property \"service_id\" is missing | property \"plan_id\" is missing | property \"organization_guid\" is missing | property \"space_guid\" is missing", "error":"ValidationError"}`,
		},
	}, {
		"bind - positive",
		request{
			method: "PUT",
			path:   "/v2/service_instances/9b292a9c-af66-4797-8d98-b30801f32ax8/service_bindings/bb841d2b-8287-47a9-ac8f-eef4c16106f2",
			body: `{
    "service_id": "c024e536-6dc4-45c6-8a53-127e7f8275ab",
    "plan_id": "3a116ac2-fc8b-496f-a715-e9a1b205d05c.community",
      "bind_resource": {
        "app_guid": "bb841d2b-8287-47a9-ac8f-eef4c16106f2"
      },
      "context": {
          "organization_guid": "e027f3f6-80fe-4d22-9374-da23a035ba0b",
          "space_guid": "8c56f85c-c16e-4158-be79-5dac74f970de"
      },
      "parameters": {
        "parameter1": 1,
        "parameter2": "foo"
      }
}`,
		},
		want{
			code: 200,
			body: "{}",
		},
	}, {
		"bind - empty body",
		request{
			method: "PUT",
			path:   "/v2/service_instances/9b292a9c-af66-4797-8d98-b30801f32ax8/service_bindings/bb841d2b-8287-47a9-ac8f-eef4c16106f2",
			body:   "{}",
		},
		want{
			code: 400,
			body: `{"description":"request body has an error: doesn't match schema #/components/schemas/ServiceBindingRequest: property \"service_id\" is missing | property \"plan_id\" is missing", "error":"ValidationError"}`,
		},
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gin.SetMode(gin.TestMode)
			router := gin.New()
			middleware, err := validatorMiddleware(context.Background())
			require.NoError(t, err)
			router.Use(middleware)
			router.NoRoute(func(c *gin.Context) {
				c.JSON(200, gin.H{})
			})

			w := httptest.NewRecorder()
			var body io.Reader
			if len(tt.request.body) > 0 {
				body = strings.NewReader(tt.request.body)
			}
			req, _ := http.NewRequest(tt.request.method, tt.request.path, body)
			req.Header = validHeader()
			router.ServeHTTP(w, req)
			require.JSONEq(t, tt.want.body, w.Body.String())
			require.Equal(t, tt.want.code, w.Result().StatusCode)
		})
	}
}

func Test_validatorMiddleware_route(t *testing.T) {
	type want struct {
		code int
		body string
	}
	type request struct {
		method string
		path   string
	}
	tests := []struct {
		name    string
		request request
		want    want
	}{{
		"invalid path",
		request{
			method: "PUT",
			path:   "/invalid/path",
		},
		want{
			code: 404,
			body: `{"error":"NotFound"}`,
		},
	}, {
		"invalid method",
		request{
			method: "POST",
			path:   "/v2/service_instances/9b292a9c-af66-4797-8d98-b30801f32ax7",
		},
		want{
			code: 405,
			body: `{"error":"MethodNotAllowed"}`,
		},
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gin.SetMode(gin.TestMode)
			router := gin.New()
			middleware, err := validatorMiddleware(context.Background())
			require.NoError(t, err)
			router.Use(middleware)
			router.NoRoute(func(c *gin.Context) {
				c.JSON(200, gin.H{})
			})

			w := httptest.NewRecorder()
			req, _ := http.NewRequest(tt.request.method, tt.request.path, nil)
			req.Header = validHeader()
			router.ServeHTTP(w, req)
			require.JSONEq(t, tt.want.body, w.Body.String())
			require.Equal(t, tt.want.code, w.Result().StatusCode)
		})
	}
}

func validHeader() http.Header {
	return http.Header{"Origin": []string{"https://example.com"}, "X-Broker-Api-Version": []string{"2.17"}, "Content-Type": []string{"application/json"}}
}

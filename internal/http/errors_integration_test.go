//go:build integration

package http

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func Test_errorsMiddleware(t *testing.T) {
	// TODO: validate responses
	type want struct {
		code int
		body string
	}
	tests := []struct {
		name    string
		handler gin.HandlerFunc
		want    want
	}{{
		"no error",
		func(c *gin.Context) { c.JSON(http.StatusOK, gin.H{}) },
		want{http.StatusOK, `{}`},
	}, {
		"with error - status not changed",
		func(c *gin.Context) { _ = c.Error(fmt.Errorf("error")) },
		want{http.StatusOK, `{"error": "Ok", "description": "error"}`},
	}, {
		"with error",
		func(c *gin.Context) { _ = c.AbortWithError(http.StatusTeapot, fmt.Errorf("error")) },
		want{http.StatusTeapot, `{"error": "IMATeapot", "description": "error"}`},
	}, {
		"aborted without error",
		func(c *gin.Context) { c.AbortWithStatus(http.StatusTeapot) },
		want{http.StatusTeapot, `{"error": "IMATeapot"}`},
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gin.SetMode(gin.TestMode)
			router := gin.New()
			router.Use(errorsMiddleware)
			router.GET("/", tt.handler)

			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/", nil)
			req.Header.Set("Origin", "https://example.com")
			router.ServeHTTP(w, req)
			require.JSONEq(t, tt.want.body, w.Body.String())
			require.Equal(t, tt.want.code, w.Result().StatusCode)
		})
	}
}

//go:build !integration

package ctxutil

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func TestCtx(t *testing.T) {
	tests := []struct {
		name string
		gCtx *gin.Context
		want Context
	}{{
		"nil",
		nil,
		NewContext(),
	}, {
		"empty",
		emptyGinContext(),
		NewContext(),
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Ctx(tt.gCtx)
			require.Equal(t, tt.want, got)
		})
	}
}

func emptyGinContext() *gin.Context {
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	return c
}

func Test_ctx_Inject(t *testing.T) {
	tests := []struct {
		name string
		ctx  Context
	}{{
		"empty",
		NewContext(),
	}, {
		"with config",
		NewContext().WithEnableSpaceIdentity(true),
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got Context
			gin.SetMode(gin.TestMode)
			router := gin.New()
			router.Use(tt.ctx.Inject())
			router.GET("/", func(c *gin.Context) {
				got = Ctx(c)
			})

			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/", nil)
			req.Header.Set("Origin", "https://example.com")
			router.ServeHTTP(w, req)
			require.Equal(t, tt.ctx, got)
		})
	}
}

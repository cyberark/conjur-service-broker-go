// Package ctxutil implements custom context tailored for service provider needs
package ctxutil

import (
	"context"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

const ctxKey = "service-broker-context"

// Context is a wrapper around context providing ability to check feature toggle and middleware injecting itself to gin context
type Context interface {
	context.Context
	WithEnableSpaceIdentity(enabled bool) Context
	IsEnableSpaceIdentity() bool
	WithLogger(log *zap.SugaredLogger) Context
	Logger() *zap.SugaredLogger
	Inject() gin.HandlerFunc
}

type ctx struct {
	context.Context
}

// NewContext creates new custom conjur service broker context
func NewContext() Context {
	return &ctx{context.Background()}
}

func (c *ctx) Inject() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Set(ctxKey, c)
		ctx.Next()
	}
}

// Ctx unwraps previously injected custom context from gin context
func Ctx(c *gin.Context) Context {
	context, ok := c.Value(ctxKey).(*ctx)
	if !ok {
		context = &ctx{} // avoid nil pointer
	}
	return context
}

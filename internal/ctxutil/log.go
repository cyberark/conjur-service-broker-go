// Package ctxutil implements custom context tailored for service provider needs
package ctxutil

import (
	"context"

	"go.uber.org/zap"
)

type logKey struct{}

func (c *ctx) WithLogger(log *zap.SugaredLogger) Context {
	return &ctx{
		Context: context.WithValue(c.Context, logKey{}, log),
	}
}

func (c *ctx) Logger() *zap.SugaredLogger {
	if log, ok := context.Context(*c).Value(logKey{}).(*zap.SugaredLogger); ok {
		return log
	}
	return nil
}

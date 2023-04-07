// Package ctxutil implements custom context tailored for service provider needs
package ctxutil

import "context"

type configKey struct{}

type config struct {
	enableSpaceIdentity bool
}

func (c *ctx) WithEnableSpaceIdentity(enabled bool) Context {
	cfg := c.config()
	cfg.enableSpaceIdentity = enabled
	return &ctx{
		Context: context.WithValue(c.Context, configKey{}, cfg),
	}
}

func (c *ctx) IsEnableSpaceIdentity() bool {
	return c.config().enableSpaceIdentity
}

func (c *ctx) config() *config {
	if cfg, ok := context.Context(*c).Value(configKey{}).(*config); ok {
		return cfg
	}
	return &config{}
}

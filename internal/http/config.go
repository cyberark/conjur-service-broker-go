// Package http implements the http communication layer of the conjur service broker
package http

import (
	"errors"
	"fmt"

	"github.com/caarlos0/env/v7"
	"github.com/cyberark/conjur-service-broker/pkg/conjur"
)

type config struct {
	conjur.Config

	SecurityUserName     string `env:"SECURITY_USER_NAME,unset"`
	SecurityUserPassword string `env:"SECURITY_USER_PASSWORD,unset"`

	// ENABLE_SPACE_IDENTITY: When set to true, the service broker provides applications with a Space-level host identity, rather than create a new host identity for each application in Conjur at bind time. This allows the broker to use a Conjur follower for application binding, rather than the Conjur master.
	EnableSpaceIdentity bool `env:"ENABLE_SPACE_IDENTITY" envDefault:"false"`

	// PORT: sets the http server listening port
	Port string `env:"PORT" envDefault:"8080"`

	// DEBUG: Enables debug mode
	Debug bool `env:"DEBUG" envDefault:"false"`

	//TODO: logger config
}

func newConfig() (*config, error) {
	cfg := config{}
	if err := env.Parse(&cfg); err != nil {
		return nil, fmt.Errorf("failed to init config: %w", err)
	}
	if err := validate(cfg); err != nil {
		return nil, fmt.Errorf("failed to validate config: %w", err)
	}
	return &cfg, nil
}

func validate(cfg config) error {
	if cfg.ConjurVersion != 5 {
		return errors.New("conjur enterprise v4 is no longer supported, please use conjur service broker v1.1.4 or earlier")
	}
	//TODO: validate certificate if present
	return nil
}

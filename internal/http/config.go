// Package http implements the http communication layer of the conjur service broker
package http

import (
	"errors"
	"fmt"
	"net/url"
	"strings"

	"github.com/caarlos0/env/v7"
	"github.com/cyberark/conjur-service-broker-go/pkg/conjur"
)

// ErrInvalidConjurVersion error indicating invalid conjur version parameter value
var ErrInvalidConjurVersion = errors.New("conjur enterprise v4 is no longer supported, please use conjur service broker v1.1.4 or earlier")

type config struct {
	// TODO: trim backslash in url
	conjur.Config

	SecurityUserName     string `env:"SECURITY_USER_NAME,unset"`
	SecurityUserPassword string `env:"SECURITY_USER_PASSWORD,unset"`

	// ENABLE_SPACE_IDENTITY: When set to true, the service broker provides applications with a Space-level host identity, rather than create a new host identity for each application in Conjur at bind time. This allows the broker to use a Conjur follower for application binding, rather than the Conjur master.
	EnableSpaceIdentity bool `env:"ENABLE_SPACE_IDENTITY" envDefault:"false"`

	// PORT: sets the http server listening port
	Port string `env:"PORT" envDefault:"8080"`

	// DEBUG: Enables debug mode
	Debug bool `env:"DEBUG" envDefault:"false"`
}

func newConfig() (*config, error) {
	cfg := config{}
	if err := env.Parse(&cfg); err != nil {
		return nil, fmt.Errorf("failed to init config: %w", err)
	}
	cfg = normalizeURLs(cfg)
	if err := validate(cfg); err != nil {
		return nil, fmt.Errorf("failed to validate config: %w", err)
	}
	return &cfg, nil
}

func validate(cfg config) error {
	if cfg.ConjurVersion != 5 {
		return ErrInvalidConjurVersion
	}
	if err := validateURL(cfg.ConjurApplianceURL); err != nil {
		return fmt.Errorf("conjur appliance url validation failed: %w", err)
	}
	if len(cfg.ConjurFollowerURL) == 0 {
		return nil
	}
	if err := validateURL(cfg.ConjurFollowerURL); err != nil {
		return fmt.Errorf("conjur follower url validation failed: %w", err)
	}
	return nil
}

func validateURL(urlStr string) error {
	parsed, err := url.Parse(urlStr)
	if err != nil {
		return fmt.Errorf("invalid URL '%s': %w", urlStr, err)
	}
	if len(parsed.Scheme) == 0 {
		return fmt.Errorf("invalid URL '%s': missing scheme", urlStr)
	}
	if parsed.Scheme != "http" && parsed.Scheme != "https" {
		return fmt.Errorf("invalid URL '%s': unsupported scheme, expecting http or https", urlStr)
	}
	if len(parsed.Host) == 0 {
		return fmt.Errorf("invalid URL '%s': missing host", urlStr)
	}
	return nil
}

func normalizeURLs(cfg config) config {
	cfg.ConjurApplianceURL = strings.TrimRight(cfg.ConjurApplianceURL, "/")
	cfg.ConjurFollowerURL = strings.TrimRight(cfg.ConjurFollowerURL, "/")
	return cfg
}

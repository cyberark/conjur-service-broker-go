// Package main provides integration tests for conjur service broker
package main

type cfg struct {
	// Cloud Foundry
	CFUser     string `env:"CF_USERNAME,notEmpty"`
	CFPassword string `env:"CF_PASSWORD,notEmpty"`
	CFURL      string `env:"CF_API_URL,notEmpty"`

	// Conjur
	ConjurAccount      string `env:"PCF_CONJUR_ACCOUNT,notEmpty"`
	ConjurApplianceURL string `env:"PCF_CONJUR_APPLIANCE_URL,notEmpty"`
	ConjurUser         string `env:"PCF_CONJUR_USERNAME,notEmpty"`
	ConjurAPIKey       string `env:"PCF_CONJUR_API_KEY,notEmpty"`
}

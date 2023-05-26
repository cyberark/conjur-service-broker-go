// Package main provides integration tests for conjur service broker
package main

type cfg struct {
	// Cloud Foundry
	CFUser     string `env:"CF_USERNAME,notEmpty"`
	CFPassword string `env:"CF_PASSWORD,notEmpty"`
	CFURL      string `env:"CF_API_URL,notEmpty"`

	// Conjur
	ConjurAccount           string `env:"PCF_CONJUR_ACCOUNT,notEmpty"`
	ConjurApplianceURL      string `env:"PCF_CONJUR_APPLIANCE_URL,notEmpty"`
	ConjurUser              string `env:"PCF_CONJUR_USERNAME,notEmpty"`
	ConjurServiceBrokerUser string `env:"PCF_CONJUR_SERVICE_BROKER_USERNAME" envDefault:"host/pcf/service-broker"`
	ConjurAPIKey            string `env:"PCF_CONJUR_API_KEY,notEmpty"`
	ConjurPolicy            string `env:"PCF_CONJUR_POLICY" envDefault:"pcf/ci"`

	// Service Broker
	ServiceBrokerUser     string `env:"SB_BASIC_AUTH_USER" envDefault:"sb_user"`
	ServiceBrokerPassword string `env:"SB_BASIC_AUTH_PWD" envDefault:"sb_password"`
}

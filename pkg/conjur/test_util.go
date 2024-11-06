//go:build integration

// Package conjur provides a wrapper around conjur go SDK
package conjur

import (
	"github.com/cyberark/conjur-service-broker-go/pkg/conjur/api/mocks"
	"github.com/stretchr/testify/mock"
)

// NewMockClient creates a new client with mocked methods and a handle to the mockery mock to allow testing
func NewMockClient() (Client, *mock.Mock) {
	c := &mocks.Client{}
	return &client{
		client:   c,
		roClient: c,
		config: &Config{
			ConjurAccount:      "dev",
			ConjurApplianceURL: "https://conjur.local",
			ConjurFollowerURL:  "https://follower.local",
			ConjurPolicy:       "cf",
			ConjurAuthNLogin:   "test",
			ConjurAuthNAPIKey:  "test-api-key",
		},
	}, &c.Mock
}

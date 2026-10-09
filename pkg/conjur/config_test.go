//go:build !integration

package conjur

import (
	"testing"

	"github.com/cyberark/conjur-api-go/conjurapi"
	"github.com/stretchr/testify/assert"
)

func TestConfig_mergeConfig(t *testing.T) {
	tests := []struct {
		name     string
		config   Config
		original conjurapi.Config
		want     conjurapi.Config
	}{{
		"positive",
		Config{
			ConjurApplianceURL:   "appliance_url",
			ConjurAccount:        "account",
			ConjurSSLCertificate: "cert",
		},
		conjurapi.Config{
			ApplianceURL: "another_url",
			Account:      "other_account",
			SSLCert:      "different_cert",
		},
		conjurapi.Config{
			ApplianceURL: "appliance_url",
			Account:      "account",
			SSLCert:      "cert",
		},
	}, {
		"empty",
		Config{},
		conjurapi.Config{},
		conjurapi.Config{},
	}, {
		"cert auth - inline PEM fields are merged",
		Config{
			ConjurClientCert:    "test-client-cert-pem",
			ConjurClientCertKey: "test-client-cert-key-pem",
		},
		conjurapi.Config{},
		conjurapi.Config{
			ClientCert:    "test-client-cert-pem",
			ClientCertKey: "test-client-cert-key-pem",
		},
	}, {
		"cert auth - inline PEM overrides existing",
		Config{
			ConjurClientCert:    "new-cert",
			ConjurClientCertKey: "new-key",
		},
		conjurapi.Config{
			ClientCert:    "old-cert",
			ClientCertKey: "old-key",
		},
		conjurapi.Config{
			ClientCert:    "new-cert",
			ClientCertKey: "new-key",
		},
	}, {
		"cert auth - empty inline PEM does not override existing",
		Config{},
		conjurapi.Config{
			ClientCert:    "existing-cert",
			ClientCertKey: "existing-key",
		},
		conjurapi.Config{
			ClientCert:    "existing-cert",
			ClientCertKey: "existing-key",
		},
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.config.mergeConfig(tt.original)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestConfig_AppURL(t *testing.T) {
	tests := []struct {
		name   string
		config Config
		want   string
	}{{
		"returns follower URL when set",
		Config{ConjurApplianceURL: "https://leader.example.com", ConjurFollowerURL: "https://follower.example.com"},
		"https://follower.example.com",
	}, {
		"falls back to appliance URL when follower not set",
		Config{ConjurApplianceURL: "https://leader.example.com"},
		"https://leader.example.com",
	}, {
		"falls back to appliance URL when follower is empty string",
		Config{ConjurApplianceURL: "https://leader.example.com", ConjurFollowerURL: ""},
		"https://leader.example.com",
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.config.AppURL())
		})
	}
}

func TestConfig_AppSSLCertificate(t *testing.T) {
	tests := []struct {
		name   string
		config Config
		want   string
	}{{
		"returns follower cert when set",
		Config{ConjurSSLCertificate: "leader-cert", ConjurFollowerSSLCertificate: "follower-cert"},
		"follower-cert",
	}, {
		"falls back to main cert when follower cert not set",
		Config{ConjurSSLCertificate: "leader-cert"},
		"leader-cert",
	}, {
		"falls back to main cert when follower cert is empty string",
		Config{ConjurSSLCertificate: "leader-cert", ConjurFollowerSSLCertificate: ""},
		"leader-cert",
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.config.AppSSLCertificate())
		})
	}
}

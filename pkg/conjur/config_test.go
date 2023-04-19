package conjur

import (
	"testing"

	"github.com/cyberark/conjur-api-go/conjurapi"
	"github.com/stretchr/testify/require"
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
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.config.mergeConfig(tt.original)
			require.Equal(t, tt.want, got)
		})
	}
}

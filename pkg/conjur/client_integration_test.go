//go:build integration

package conjur_test

import (
	"testing"

	"github.com/cyberark/conjur-service-broker/pkg/conjur"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func Test_client_ValidateConnectivity(t *testing.T) {
	tests := []struct {
		name string
		// fields  fields
		wantErr bool
	}{
		{"positive", false},
		{"negative", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, mockAPI := conjur.NewMockClient()
			mockAPI.On("CheckPermission", mock.Anything, mock.Anything).Return(!tt.wantErr, nil)
			err := client.ValidateConnectivity()
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

//go:build integration

package conjur_test

import (
	"testing"

	"github.com/cyberark/conjur-service-broker/pkg/conjur"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func Test_client_ValidateConnectivity(t *testing.T) {
	tests := []struct {
		name          string
		hasPermission bool
		wantErr       assert.ErrorAssertionFunc
	}{
		{"positive", true, assert.NoError},
		{"negative", false, assert.Error},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, mockAPI := conjur.NewMockClient()
			mockAPI.On("CheckPermission", mock.Anything, mock.Anything).Return(tt.hasPermission, nil)
			err := client.ValidateConnectivity()
			tt.wantErr(t, err)
		})
	}
}

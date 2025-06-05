//go:build integration

package conjur

import (
	"errors"
	"testing"

	"github.com/cyberark/conjur-api-go/conjurapi"
	"github.com/cyberark/conjur-service-broker-go/pkg/conjur/api/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func Test_client_FromBindingID(t *testing.T) {
	type p []interface{}
	type mockParams struct {
		args    p
		returns p
	}
	type want struct {
		orgID     string
		spaceID   string
		bindingID string
		hostID    string
	}
	tests := []struct {
		name    string
		client  mockParams
		want    *want
		wantErr assert.ErrorAssertionFunc
	}{{
		"positive",
		mockParams{
			args:    p{&conjurapi.ResourceFilter{Kind: "host", Search: "bindingID^", Limit: 0, Offset: 0, Role: ""}},
			returns: p{[]map[string]interface{}{{"id": "dev:host:cf/orgID/spaceID/bindingID"}}, nil},
		},
		&want{
			orgID:     "orgID",
			spaceID:   "spaceID",
			bindingID: "bindingID",
			hostID:    "host:policy/orgID/spaceID/bindingID",
		},
		assert.NoError,
	}, {
		"error from resources",
		mockParams{
			args:    p{&conjurapi.ResourceFilter{Kind: "host", Search: "bindingID^", Limit: 0, Offset: 0, Role: ""}},
			returns: p{nil, errors.New("error")},
		},
		nil,
		assert.Error,
	}, {
		"empty response from resources",
		mockParams{
			args:    p{&conjurapi.ResourceFilter{Kind: "host", Search: "bindingID^", Limit: 0, Offset: 0, Role: ""}},
			returns: p{nil, nil},
		},
		&want{orgID: "", spaceID: "", bindingID: "bindingID", hostID: "host:policy/bindingID"},
		assert.NoError,
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := mocks.NewMockClient(t)
			c.On("Resources", tt.client.args...).Return(tt.client.returns...).Once()

			client := client{roClient: c, config: &Config{ConjurPolicy: "policy"}}
			b, err := client.FromBindingID("bindingID")
			c.AssertExpectations(t)

			if tt.want == nil {
				require.Nil(t, b)
				return
			}
			got, ok := b.(*bind)
			require.True(t, ok)

			tt.wantErr(t, err)
			assert.Equal(t, tt.want.orgID, got.orgID)
			assert.Equal(t, tt.want.spaceID, got.spaceID)
			assert.Equal(t, tt.want.bindingID, got.bindingID)
			assert.Equal(t, tt.want.hostID, got.hostID)
		})
	}
}

func Test_client_ValidateConnectivity(t *testing.T) {
	tests := []struct {
		name                string
		hasIdentityResource bool
		hasPermissionRW     bool
		hasPermissionRO     bool
		withErrIdentity     error
		withErrRW           error
		withErrRO           error
		wantErr             assert.ErrorAssertionFunc
	}{
		{"positive", true, true, true, nil, nil, nil, assert.NoError},
		{"missing permission rw", true, false, true, nil, nil, nil, assert.Error},
		{"missing permission ro", true, true, false, nil, nil, nil, assert.Error},
		{"with error on rw client", true, true, true, nil, errors.New("error"), nil, assert.Error},
		{"with error on ro client", true, true, true, nil, nil, errors.New("error"), assert.Error},
		{"missing identity resource", false, true, true, nil, nil, nil, assert.Error},
		{"with identity resource error", true, true, true, errors.New("error"), nil, nil, assert.Error},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, mockAPI := NewMockConjurClient()
			mockAPI.On("ResourceExists", "dev:user:test").Return(tt.hasIdentityResource, tt.withErrIdentity).Once()
			mockAPI.On("CheckPermission", mock.Anything, mock.Anything).Return(tt.hasPermissionRW, tt.withErrRW).Once()
			mockAPI.On("CheckPermission", mock.Anything, mock.Anything).Return(tt.hasPermissionRO, tt.withErrRO).Once()
			err := client.ValidateConnectivity()
			tt.wantErr(t, err)
		})
	}
}

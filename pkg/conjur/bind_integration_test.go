//go:build integration

package conjur

import (
	"errors"
	"testing"

	"github.com/cyberark/conjur-api-go/conjurapi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func Test_bind_BindHostPolicy(t *testing.T) {
	type p []interface{}
	type m struct {
		args    p
		returns p
	}
	type mockParams map[string]m
	policyResp := p{&conjurapi.PolicyResponse{
		CreatedRoles: map[string]conjurapi.CreatedRole{"role": {
			ID:     "role",
			APIKey: "my-api-key",
		}},
	}, nil}
	tests := []struct {
		name    string
		client  mockParams
		want    *Policy
		wantErr assert.ErrorAssertionFunc
	}{{
		"positive",
		mockParams{
			"Resource": m{args: p{"dev:host:test"}, returns: p{map[string]interface{}{}, nil}},
			"LoadPolicy": m{
				args:    p{conjurapi.PolicyModePost, "cf/orgID/spaceID", mock.Anything},
				returns: policyResp},
		},
		&Policy{Account: "dev", ApplianceURL: "https://conjur.local", AuthnAPIKey: "my-api-key"},
		assert.NoError,
	}, {
		"error from load policy",
		mockParams{
			"Resource": m{args: p{"dev:host:test"}, returns: p{map[string]interface{}{}, nil}},
			"LoadPolicy": m{
				args:    p{conjurapi.PolicyModePost, "cf/orgID/spaceID", mock.Anything},
				returns: p{nil, errors.New("error")}},
		},
		nil,
		assert.Error,
	}, {
		"empty response from load policy",
		mockParams{
			"Resource": m{args: p{"dev:host:test"}, returns: p{map[string]interface{}{}, nil}},
			"LoadPolicy": m{
				args:    p{conjurapi.PolicyModePost, "cf/orgID/spaceID", mock.Anything},
				returns: p{nil, nil}},
		},
		nil,
		assert.Error,
	}, {
		"error from resource",
		mockParams{
			"Resource": m{args: p{"dev:host:test"}, returns: p{map[string]interface{}{}, errors.New("error")}},
			"LoadPolicy": m{
				args:    p{conjurapi.PolicyModePost, "cf/orgID/spaceID", mock.Anything},
				returns: policyResp},
		},
		&Policy{Account: "dev", ApplianceURL: "https://conjur.local", AuthnAPIKey: "my-api-key"},
		assert.NoError,
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, mockAPI := NewMockClient()
			for method, v := range tt.client {
				mockAPI.On(method, v.args...).Return(v.returns...).Once()
			}
			b := c.NewBind("orgID", "spaceID", "bindingID", true)
			got, err := b.BindHostPolicy()
			tt.wantErr(t, err)
			assert.Equal(t, tt.want, got)
			mockAPI.AssertExpectations(t)
		})
	}
}

func Test_bind_BindSpacePolicy(t *testing.T) {
	type p []interface{}
	type m struct {
		args    p
		returns p
	}
	type mockParams map[string]m
	tests := []struct {
		name    string
		client  mockParams
		want    *Policy
		wantErr assert.ErrorAssertionFunc
	}{{
		"positive",
		mockParams{
			"RetrieveSecret": m{args: p{"cf/orgID/spaceID/space-host-api-key"}, returns: p{[]byte("secret"), nil}},
		},
		&Policy{Account: "dev", ApplianceURL: "https://conjur.local", AuthnAPIKey: "secret", AuthnLogin: "host/cf/orgID/spaceID"},
		assert.NoError,
	}, {
		"error for retrieve secret",
		mockParams{
			"RetrieveSecret": m{args: p{"cf/orgID/spaceID/space-host-api-key"}, returns: p{nil, errors.New("error")}},
		},
		nil,
		assert.Error,
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, mockAPI := NewMockClient()
			for method, v := range tt.client {
				mockAPI.On(method, v.args...).Return(v.returns...).Once()
			}
			b := c.NewBind("orgID", "spaceID", "bindingID", true)
			got, err := b.BindSpacePolicy()
			tt.wantErr(t, err)
			assert.Equal(t, tt.want, got)
			mockAPI.AssertExpectations(t)
		})
	}
}

func Test_bind_DeleteBindHostPolicy(t *testing.T) {
	type p []interface{}
	type m struct {
		args    p
		returns p
	}
	type mockParams map[string]m
	tests := []struct {
		name    string
		client  mockParams
		want    *Policy
		wantErr assert.ErrorAssertionFunc
	}{{
		"positive",
		mockParams{
			"Resources":    m{args: p{&conjurapi.ResourceFilter{Kind: "host", Search: "bindingID^", Limit: 0, Offset: 0, Role: ""}}, returns: p{[]map[string]interface{}{{"id": "dev:host:cf/orgID/spaceID/bindingID"}}, nil}},
			"RotateAPIKey": m{args: p{"dev:host:cf/orgID/spaceID/bindingID"}, returns: p{nil, nil}},
			"LoadPolicy":   m{args: p{conjurapi.PolicyModePut, "cf/orgID/spaceID", mock.Anything}, returns: p{nil, nil}},
		},
		&Policy{Account: "dev", ApplianceURL: "https://conjur.local", AuthnAPIKey: "secret", AuthnLogin: "host/cf/orgID/spaceID"},
		assert.NoError,
	}, {
		"error from rotate api key",
		mockParams{
			"Resources":    m{args: p{&conjurapi.ResourceFilter{Kind: "host", Search: "bindingID^", Limit: 0, Offset: 0, Role: ""}}, returns: p{[]map[string]interface{}{{"id": "dev:host:cf/orgID/spaceID/bindingID"}}, nil}},
			"RotateAPIKey": m{args: p{"dev:host:cf/orgID/spaceID/bindingID"}, returns: p{nil, errors.New("error")}},
		},
		nil,
		assert.Error,
	}, {
		"error from load policy",
		mockParams{
			"Resources":    m{args: p{&conjurapi.ResourceFilter{Kind: "host", Search: "bindingID^", Limit: 0, Offset: 0, Role: ""}}, returns: p{[]map[string]interface{}{{"id": "dev:host:cf/orgID/spaceID/bindingID"}}, nil}},
			"RotateAPIKey": m{args: p{"dev:host:cf/orgID/spaceID/bindingID"}, returns: p{nil, nil}},
			"LoadPolicy":   m{args: p{conjurapi.PolicyModePut, "cf/orgID/spaceID", mock.Anything}, returns: p{nil, errors.New("error")}},
		},
		nil,
		assert.Error,
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, mockAPI := NewMockClient()
			for method, v := range tt.client {
				mockAPI.On(method, v.args...).Return(v.returns...).Once()
			}
			b, err := c.FromBindingID("bindingID")
			require.NoError(t, err)
			err = b.DeleteBindHostPolicy()
			tt.wantErr(t, err)
			mockAPI.AssertExpectations(t)
		})
	}
}

func Test_bind_HostExists(t *testing.T) {
	type p []interface{}
	type m struct {
		args    p
		returns p
	}
	type mockParams map[string]m
	tests := []struct {
		name    string
		client  mockParams
		want    bool
		wantErr assert.ErrorAssertionFunc
	}{{
		"positive",
		mockParams{
			"ResourceExists": m{args: p{"dev:host:cf/orgID/spaceID"}, returns: p{true, nil}},
		},
		true,
		assert.NoError,
	}, {
		"error from resource exists",
		mockParams{
			"ResourceExists": m{args: p{"dev:host:cf/orgID/spaceID"}, returns: p{false, errors.New("error")}},
		},
		false,
		assert.Error,
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, mockAPI := NewMockClient()
			for method, v := range tt.client {
				mockAPI.On(method, v.args...).Return(v.returns...).Once()
			}
			b := c.NewBind("orgID", "spaceID", "bindingID", true)
			got, err := b.HostExists()
			tt.wantErr(t, err)
			assert.Equal(t, tt.want, got)
			mockAPI.AssertExpectations(t)
		})
	}
}

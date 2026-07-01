//go:build integration

package conjur

import (
	"errors"
	"testing"

	"github.com/cyberark/conjur-api-go/conjurapi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func Test_provision_ProvisionOrgSpacePolicy(t *testing.T) {
	type p []interface{}
	type m struct {
		args    p
		returns p
	}
	type mockParams map[string][]m
	tests := []struct {
		name    string
		client  mockParams
		wantErr assert.ErrorAssertionFunc
	}{{
		"positive",
		mockParams{
			"LoadPolicy": []m{{args: p{conjurapi.PolicyModePost, "cf", mock.Anything}, returns: p{nil, nil}}},
			"ResourceExists": []m{
				{args: p{"dev:policy:cf/11111111-1111-1111-1111-111111111111"}, returns: p{true, nil}},
				{args: p{"dev:policy:cf/11111111-1111-1111-1111-111111111111/22222222-2222-2222-2222-222222222222"}, returns: p{true, nil}},
				{args: p{"dev:group:cf/11111111-1111-1111-1111-111111111111/22222222-2222-2222-2222-222222222222"}, returns: p{true, nil}},
			},
		},
		assert.NoError,
	}, {
		"error in policy",
		mockParams{
			"LoadPolicy": []m{{args: p{conjurapi.PolicyModePost, "cf", mock.Anything}, returns: p{nil, errors.New("error")}}},
		},
		assert.Error,
	}, {
		"error in existence check",
		mockParams{
			"LoadPolicy":     []m{{args: p{conjurapi.PolicyModePost, "cf", mock.Anything}, returns: p{nil, nil}}},
			"ResourceExists": []m{{args: p{"dev:policy:cf/11111111-1111-1111-1111-111111111111"}, returns: p{false, errors.New("error")}}},
		},
		assert.Error,
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, mockAPI := NewMockConjurClient()
			for method, values := range tt.client {
				for _, v := range values {
					mockAPI.On(method, v.args...).Return(v.returns...).Once()
				}
			}
			p := c.NewProvision("11111111-1111-1111-1111-111111111111", "22222222-2222-2222-2222-222222222222", nil, nil)
			err := p.ProvisionOrgSpacePolicy()
			tt.wantErr(t, err)
			mockAPI.AssertExpectations(t)
		})
	}
}

func Test_provision_ProvisionHostPolicy(t *testing.T) {
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
		wantErr assert.ErrorAssertionFunc
	}{{
		"positive",
		mockParams{
			"LoadPolicy": m{
				args:    p{conjurapi.PolicyModePost, "cf/11111111-1111-1111-1111-111111111111/22222222-2222-2222-2222-222222222222", mock.Anything},
				returns: policyResp},
			"AddSecret": m{args: p{"cf/11111111-1111-1111-1111-111111111111/22222222-2222-2222-2222-222222222222/space-host-api-key", "my-api-key"}, returns: p{nil}},
		},
		assert.NoError,
	}, {
		"error in policy",
		mockParams{
			"LoadPolicy": m{args: p{conjurapi.PolicyModePost, "cf/11111111-1111-1111-1111-111111111111/22222222-2222-2222-2222-222222222222", mock.Anything}, returns: p{nil, errors.New("error")}},
		},
		assert.Error,
	}, {
		"empty response from load policy",
		mockParams{
			"LoadPolicy": m{args: p{conjurapi.PolicyModePost, "cf/11111111-1111-1111-1111-111111111111/22222222-2222-2222-2222-222222222222", mock.Anything}, returns: p{nil, nil}},
		},
		assert.Error,
	}, {
		"error in secret add",
		mockParams{
			"LoadPolicy": m{args: p{conjurapi.PolicyModePost, "cf/11111111-1111-1111-1111-111111111111/22222222-2222-2222-2222-222222222222", mock.Anything}, returns: policyResp},
			"AddSecret":  m{args: p{"cf/11111111-1111-1111-1111-111111111111/22222222-2222-2222-2222-222222222222/space-host-api-key", "my-api-key"}, returns: p{errors.New("error")}},
		},
		assert.Error,
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, mockAPI := NewMockConjurClient()
			for method, v := range tt.client {
				mockAPI.On(method, v.args...).Return(v.returns...).Once()
			}
			p := c.NewProvision("11111111-1111-1111-1111-111111111111", "22222222-2222-2222-2222-222222222222", nil, nil)
			err := p.ProvisionHostPolicy()
			tt.wantErr(t, err)
			mockAPI.AssertExpectations(t)
		})
	}
}

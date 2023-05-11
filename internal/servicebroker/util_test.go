//go:build !integration

package servicebroker

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.cyberng.com/Conjur-Enterprise/conjur-service-broker-go/pkg/conjur"
)

func Test_parseContext(t *testing.T) {
	name := "name"
	tests := []struct {
		name string
		ctx  *Context
		want context
	}{{
		"with space id",
		&Context{"space_guid": "space_id"},
		context{
			SpaceID: "space_id",
		}}, {
		"with org name",
		&Context{"organization_name": name},
		context{
			OrgName: &name,
		}}, {
		"nil",
		nil,
		context{},
	}, {
		"invalid field",
		&Context{"invalid": "value"},
		context{},
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := parseContext(tt.ctx)
			require.Equal(t, tt.want, got)
		})
	}
}

func Test_object(t *testing.T) {
	tests := []struct {
		name   string
		policy *conjur.Policy
		want   *Object
	}{{
		"with account",
		&conjur.Policy{
			Account: "test",
		},
		&Object{
			"account":         "test",
			"appliance_url":   "",
			"authn_api_key":   "",
			"authn_login":     "",
			"ssl_certificate": "",
			"version":         float64(0),
		},
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := object(tt.policy)
			require.Equal(t, tt.want, got)
		})
	}
}

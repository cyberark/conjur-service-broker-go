//go:build !integration

package servicebroker

import (
	"testing"

	"github.com/cyberark/conjur-service-broker-go/pkg/conjur"
	"github.com/stretchr/testify/require"
)

const (
	validOrgGUID   = "e027f3f6-80fe-4d22-9374-da23a035ba0b"
	validSpaceGUID = "8c56f85c-c16e-4158-be79-5dac74f970de"
)

func Test_parseContext(t *testing.T) {
	name := "name"
	tests := []struct {
		name    string
		ctx     *Context
		want    context
		wantErr bool
	}{{
		"with valid space guid",
		&Context{"space_guid": validSpaceGUID},
		context{SpaceID: validSpaceGUID},
		false,
	}, {
		"with valid org and space guids",
		&Context{"organization_guid": validOrgGUID, "space_guid": validSpaceGUID},
		context{OrgID: validOrgGUID, SpaceID: validSpaceGUID},
		false,
	}, {
		"with org name only",
		&Context{"organization_name": name},
		context{OrgName: &name},
		false,
	}, {
		"nil context",
		nil,
		context{},
		false,
	}, {
		"invalid field",
		&Context{"invalid": "value"},
		context{},
		false,
	}, {
		"non-UUID space_guid is rejected",
		&Context{"space_guid": "not-a-uuid"},
		context{},
		true,
	}, {
		"non-UUID organization_guid is rejected",
		&Context{"organization_guid": "tenant/../../other"},
		context{},
		true,
	}, {
		"space_guid with path separator is rejected",
		&Context{"space_guid": "abc/def"},
		context{},
		true,
	}, {
		"space_guid with colon is rejected",
		&Context{"space_guid": "myaccount:policy:bad"},
		context{},
		true,
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseContext(tt.ctx)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
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

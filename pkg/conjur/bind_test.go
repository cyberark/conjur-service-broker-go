//go:build !integration

package conjur

import (
	"io"
	"testing"

	"github.com/cyberark/conjur-api-go/conjurapi"
	"github.com/cyberark/conjur-service-broker-go/pkg/conjur/api/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func Test_createBindYAML(t *testing.T) {
	tests := []struct {
		name    string
		args    *bind
		want    string
		wantErr assert.ErrorAssertionFunc
		setup   func(*mocks.MockClient)
	}{{
		"simple bind",
		&bind{bindingID: "test", client: &client{roClient: mocks.NewMockClient(t), config: &Config{ConjurAuthNLogin: "host/test"}}},
		`- !host
  id: test
  annotations:
    authn/api-key: true
`,
		assert.NoError,
		func(m *mocks.MockClient) {
			m.On("Resource", mock.Anything).Return(map[string]interface{}{}, nil).Maybe()
		},
	}, {
		"advanced bind with group",
		&bind{bindingID: "test", client: &client{roClient: mocks.NewMockClient(t), config: &Config{ConjurAccount: "dev", ConjurPolicy: "cf", ConjurAuthNLogin: "host/test"}}, orgID: "orgID", spaceID: "spaceID"},
		`- !host
  id: test
  annotations:
    authn/api-key: true
- !grant
  role: !group
  member: !host test
`,
		assert.NoError,
		func(m *mocks.MockClient) {
			m.On("Resource", mock.Anything).Return(map[string]interface{}{}, nil).Maybe()
			m.On("ResourceExists", "dev:group:cf/orgID/spaceID").Return(true, nil).Once()
		},
	}, {
		"advanced bind with layer",
		&bind{bindingID: "test", client: &client{roClient: mocks.NewMockClient(t), config: &Config{ConjurAccount: "dev", ConjurPolicy: "cf", ConjurAuthNLogin: "host/test"}}, orgID: "orgID", spaceID: "spaceID"},
		`- !host
  id: test
  annotations:
    authn/api-key: true
- !grant
  role: !layer
  member: !host test
`,
		assert.NoError,
		func(m *mocks.MockClient) {
			m.On("Resource", mock.Anything).Return(map[string]interface{}{}, nil).Maybe()
			m.On("ResourceExists", "dev:group:cf/orgID/spaceID").Return(false, nil).Once()
			m.On("ResourceExists", "dev:layer:cf/orgID/spaceID").Return(true, nil).Once()
		},
	}, {
		"advanced bind with neither group nor layer",
		&bind{bindingID: "test", client: &client{roClient: mocks.NewMockClient(t), config: &Config{ConjurAccount: "dev", ConjurPolicy: "cf", ConjurAuthNLogin: "host/test"}}, orgID: "orgID", spaceID: "spaceID"},
		`- !host
  id: test
  annotations:
    authn/api-key: true
`,
		assert.NoError,
		func(m *mocks.MockClient) {
			m.On("Resource", mock.Anything).Return(map[string]interface{}{}, nil).Maybe()
			m.On("ResourceExists", "dev:group:cf/orgID/spaceID").Return(false, nil).Once()
			m.On("ResourceExists", "dev:layer:cf/orgID/spaceID").Return(false, nil).Once()
		},
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setup != nil {
				mockClient := tt.args.client.roClient.(*mocks.MockClient)
				tt.setup(mockClient)
			}
			got, err := tt.args.createBindYAML()
			tt.wantErr(t, err)
			if err == nil {
				gotBytes, err := io.ReadAll(got)
				assert.NoError(t, err)
				assert.Equal(t, tt.want, string(gotBytes))
			}
		})
	}
}

func Test_deleteBindYAML(t *testing.T) {
	tests := []struct {
		name    string
		args    *bind
		want    string
		wantErr assert.ErrorAssertionFunc
	}{{
		"simple delete",
		&bind{bindingID: "test"},
		`- !delete
  record: !host test
`,
		assert.NoError,
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.args.deleteBindYAML()
			tt.wantErr(t, err)
			gotBytes, err := io.ReadAll(got)
			assert.NoError(t, err)
			assert.Equal(t, tt.want, string(gotBytes))
		})
	}
}

// Test_onlyPolicy_credentialURLs verifies that onlyPolicy uses the follower URL
// (AppURL) rather than always using the leader appliance URL. This is the core
// fix for the bug where apps received the leader URL in their credentials even
// when CONJUR_FOLLOWER_URL was configured.
func Test_onlyPolicy_credentialURLs(t *testing.T) {
	makeResponse := func(roleID, apiKey string) *conjurapi.PolicyResponse {
		return &conjurapi.PolicyResponse{
			CreatedRoles: map[string]conjurapi.CreatedRole{
				roleID: {ID: roleID, APIKey: apiKey},
			},
		}
	}

	tests := []struct {
		name        string
		config      *Config
		wantURL     string
		wantSSLCert string
	}{{
		"uses follower URL and cert when both are configured",
		&Config{
			ConjurAccount:                "dev",
			ConjurApplianceURL:           "https://leader.example.com",
			ConjurFollowerURL:            "https://follower.example.com",
			ConjurSSLCertificate:         "leader-cert",
			ConjurFollowerSSLCertificate: "follower-cert",
		},
		"https://follower.example.com",
		"follower-cert",
	}, {
		"uses follower URL with leader cert when follower cert not set",
		&Config{
			ConjurAccount:        "dev",
			ConjurApplianceURL:   "https://leader.example.com",
			ConjurFollowerURL:    "https://follower.example.com",
			ConjurSSLCertificate: "leader-cert",
		},
		"https://follower.example.com",
		"leader-cert",
	}, {
		"falls back to leader URL when follower not configured",
		&Config{
			ConjurAccount:        "dev",
			ConjurApplianceURL:   "https://leader.example.com",
			ConjurSSLCertificate: "leader-cert",
		},
		"https://leader.example.com",
		"leader-cert",
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &bind{
				bindingID: "binding-id",
				client:    &client{config: tt.config},
			}
			roleID := "dev:host:cf/org/space/binding-id"
			policy, err := b.onlyPolicy(makeResponse(roleID, "secret-key"))
			assert.NoError(t, err)
			assert.Equal(t, tt.wantURL, policy.ApplianceURL,
				"ApplianceURL in credentials should point to follower (not leader) when CONJUR_FOLLOWER_URL is set")
			assert.Equal(t, tt.wantSSLCert, policy.SslCertificate)
		})
	}
}

// Test_BindSpacePolicy_credentialURLs verifies that BindSpacePolicy also uses
// AppURL (follower URL) rather than the leader appliance URL.
func Test_BindSpacePolicy_credentialURLs(t *testing.T) {
	// BindSpacePolicy validates that org/space IDs are UUIDs.
	const (
		testOrgID   = "06133883-eea6-497c-b474-7525f1a3e071"
		testSpaceID = "e7deceaf-5b6b-4214-a084-e1982d2e9a13"
	)
	tests := []struct {
		name        string
		config      *Config
		wantURL     string
		wantSSLCert string
	}{{
		"uses follower URL when configured",
		&Config{
			ConjurAccount:                "dev",
			ConjurPolicy:                 "cf",
			ConjurApplianceURL:           "https://leader.example.com",
			ConjurFollowerURL:            "https://follower.example.com",
			ConjurSSLCertificate:         "leader-cert",
			ConjurFollowerSSLCertificate: "follower-cert",
			ConjurVersion:                5,
		},
		"https://follower.example.com",
		"follower-cert",
	}, {
		"falls back to leader URL when follower not configured",
		&Config{
			ConjurAccount:        "dev",
			ConjurPolicy:         "cf",
			ConjurApplianceURL:   "https://leader.example.com",
			ConjurSSLCertificate: "leader-cert",
			ConjurVersion:        5,
		},
		"https://leader.example.com",
		"leader-cert",
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			roMock := mocks.NewMockClient(t)
			roMock.On("RetrieveSecret", mock.Anything).Return([]byte("mock-api-key"), nil)
			b := &bind{
				orgID:     testOrgID,
				spaceID:   testSpaceID,
				bindingID: "binding-id",
				client: &client{
					roClient: roMock,
					config:   tt.config,
				},
			}
			policy, err := b.BindSpacePolicy()
			if !assert.NoError(t, err) {
				return
			}
			assert.Equal(t, tt.wantURL, policy.ApplianceURL,
				"ApplianceURL in credentials should point to follower (not leader) when CONJUR_FOLLOWER_URL is set")
			assert.Equal(t, tt.wantSSLCert, policy.SslCertificate)
		})
	}
}

func Test_dropAccount(t *testing.T) {
	tests := []struct {
		name string
		id   string
		want string
	}{{
		"empty",
		"",
		"",
	}, {
		"full",
		"dev:host:cf/orgID/spaceID",
		"host/cf/orgID/spaceID",
	}, {
		"invalid kind",
		"dev:invalid:cf/orgID/spaceID",
		"cf/orgID/spaceID",
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, dropAccount(tt.id), "dropAccount(%v)", tt.id)
		})
	}
}

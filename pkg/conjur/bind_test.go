//go:build !integration

package conjur

import (
	"io"
	"testing"

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

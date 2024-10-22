//go:build !integration

package conjur

import (
	"io"
	"testing"

	"github.com/cyberark/conjur-api-go/conjurapi"
	"github.com/stretchr/testify/assert"
)

func Test_createBindYAML(t *testing.T) {
	tests := []struct {
		name    string
		args    *bind
		want    string
		wantErr assert.ErrorAssertionFunc
	}{{
		"simple bind",
		&bind{bindingID: "test", client: &client{roClient: &conjurapi.Client{}, config: &Config{}}},
		`- !host
  id: test
  annotations:
    authn/api-key: true
`,
		assert.NoError,
	}, {
		"advanced bind",
		&bind{bindingID: "test", client: &client{roClient: &conjurapi.Client{}, config: &Config{}}, orgID: "orgID", spaceID: "spaceID"},
		`- !host
  id: test
  annotations:
    authn/api-key: true
- !grant
  role: !group
  member: !host test
`,
		assert.NoError,
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.args.createBindYAML()
			tt.wantErr(t, err)
			gotBytes, err := io.ReadAll(got)
			assert.NoError(t, err)
			assert.Equal(t, tt.want, string(gotBytes))
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

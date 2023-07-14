//go:build !integration

package conjur

import (
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_provisionOrgSpaceYAML(t *testing.T) {
	tests := []struct {
		name string
		args *provision
		want string
	}{
		{
			"just IDs",
			&provision{
				orgID:   "1 # 2",
				spaceID: "2",
			}, `- !policy
  id: '1 # 2'
  body:
    - !layer
    - !policy
      id: "2"
      body:
        - !layer
    - !grant
      role: !layer
      member: !layer 2
`,
		}, {
			"with annotations",
			&provision{
				orgID:     "3",
				orgName:   "my-org",
				spaceID:   "4",
				spaceName: "my-space",
			},
			`- !policy
  id: "3"
  annotations:
    pcf/orgName: my-org
    pcf/type: org
  body:
    - !layer
    - !policy
      id: "4"
      annotations:
        pcf/orgName: my-org
        pcf/spaceName: my-space
        pcf/type: space
      body:
        - !layer
    - !grant
      role: !layer
      member: !layer 4
`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.args.provisionOrgSpaceYAML()
			assert.NoError(t, err)
			bytes, err := io.ReadAll(got)
			assert.NoError(t, err)
			assert.Equal(t, tt.want, string(bytes))
		})
	}
}

func Test_provision_provisionSpaceHostYAML(t *testing.T) {
	tests := []struct {
		name    string
		args    *provision
		want    string
		wantErr assert.ErrorAssertionFunc
	}{{
		"positive",
		&provision{
			orgID:     "org_id",
			orgName:   "org_name",
			spaceID:   "space_id",
			spaceName: "space_name",
			client: &client{
				config: &Config{
					ConjurAccount:    "account",
					ConjurAuthNLogin: "host/my-host",
					ConjurPolicy:     "policy",
				},
			},
		},
		`- !host
- !grant
  role: !layer
  member: !host
- !variable
  id: space-host-api-key
- !permit
  role: !host /my-host
  privileges: [read]
  resource: !variable space-host-api-key
`,
		assert.NoError,
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.args.provisionHostYAML()
			tt.wantErr(t, err)
			bytes, err := io.ReadAll(got)
			assert.NoError(t, err)
			assert.Equal(t, tt.want, string(bytes))
		})
	}
}

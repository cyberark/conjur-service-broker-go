//go:build !integration

package conjur

import (
	"io"
	"testing"

	"github.com/stretchr/testify/require"
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
				orgName:   "org",
				spaceID:   "4",
				spaceName: "space",
			},
			`- !policy
  id: "3"
  annotations:
    pcf/orgName: org
    pcf/type: org
  body:
    - !layer
    - !policy
      id: "4"
      annotations:
        pcf/orgName: org
        pcf/spaceName: space
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
			if err != nil {
				t.Error(err)
			}
			bytes, err := io.ReadAll(got)
			if err != nil {
				t.Error(err)
			}
			if string(bytes) != tt.want {
				t.Errorf("provisionOrgSpaceYAML() = \n%v\n, want \n%v\n", string(bytes), tt.want)
			}
		})
	}
}

func Test_provision_provisionSpaceHostYAML(t *testing.T) {
	tests := []struct {
		name    string
		args    *provision
		want    string
		wantErr bool
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
		false,
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.args.provisionHostYAML()
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			bytes, err := io.ReadAll(got)
			require.NoError(t, err)
			require.Equal(t, tt.want, string(bytes))
		})
	}
}

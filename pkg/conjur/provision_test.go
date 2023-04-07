package conjur

import (
	"io"
	"testing"
)

func Test_createOrgSpace(t *testing.T) {
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
			got, err := tt.args.createOrgSpaceYAML()
			if err != nil {
				t.Error(err)
			}
			bytes, err := io.ReadAll(got)
			if err != nil {
				t.Error(err)
			}
			if string(bytes) != tt.want {
				t.Errorf("createOrgSpaceYAML() = \n%v\n, want \n%v\n", string(bytes), tt.want)
			}
		})
	}
}

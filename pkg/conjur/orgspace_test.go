package conjur

import (
	"io"
	"testing"
)

func Test_createOrgSpace(t *testing.T) {
	tests := []struct {
		name string
		args *orgSpace
		want string
	}{
		{
			"just IDs",
			&orgSpace{
				OrgID:   "1 # 2",
				SpaceID: "2",
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
			&orgSpace{
				OrgID:     "3",
				OrgName:   "org",
				SpaceID:   "4",
				SpaceName: "space",
			},
			`- !policy
  id: "3"
  body:
    - !layer
    - !policy
      id: "4"
      body:
        - !layer
      annotations:
        pcf/orgName: org
        pcf/spaceName: space
        pcf/type: space
    - !grant
      role: !layer
      member: !layer 4
  annotations:
    pcf/orgName: org
    pcf/type: org
`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := createOrgSpaceYAML(tt.args)
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

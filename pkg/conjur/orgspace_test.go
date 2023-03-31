package conjur

import (
	"io"
	"testing"
)

func Test_createOrgSpace(t *testing.T) {
	tests := []struct {
		name string
		args *OrgSpace
		want string
	}{
		{
			"just IDs",
			&OrgSpace{
				OrgID:   "1",
				SpaceID: "2",
			}, `- !policy
  id: "1"
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
			&OrgSpace{
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
			got, err := createOrgSpace(tt.args)
			if err != nil {
				t.Error(err)
			}
			bytes, err := io.ReadAll(got)
			if err != nil {
				t.Error(err)
			}
			if string(bytes) != tt.want {
				t.Errorf("createOrgSpace() = \n%v\n, want \n%v\n", string(bytes), tt.want)
			}
		})
	}
}

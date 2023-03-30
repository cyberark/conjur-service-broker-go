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
			},
			`---
- !policy
  id: !!str 1
  body:
    - !layer
    - !policy
      id: !!str 2
      body:
       - !layer	
    - !grant
      role: !layer
      member: !layer 2`,
		}, {
			"with annotations",
			&OrgSpace{
				OrgID:     "3",
				OrgName:   "org",
				SpaceID:   "4",
				SpaceName: "space",
			},
			`---
- !policy
  id: !!str 3
  annotations:
    pcf/type: org
    pcf/orgName: org
  body:
    - !layer
    - !policy
      id: !!str 4
      annotations:
        pcf/type: space
        pcf/orgName: org
        pcf/spaceName: space
      body:
       - !layer	
    - !grant
      role: !layer
      member: !layer 4`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := createOrgSpace(tt.args)
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

package conjur

import (
	"io"
	"text/template"
)

var (
	orgSpaceTemplate = template.Must(template.New("create_org_space").Parse(`---
- !policy
  id: !!str {{ .OrgID }}
  {{ if .OrgName -}}
  annotations:
    pcf/type: org
    pcf/orgName: {{ .OrgName }}
  {{ end -}}
  body:
    - !layer
    - !policy
      id: !!str {{ .SpaceID }}
      {{ if and .OrgName .SpaceName -}}
      annotations:
        pcf/type: space
        pcf/orgName: {{ .OrgName }}
        pcf/spaceName: {{ .SpaceName }}
      {{ end -}}
      body:
       - !layer	
    - !grant
      role: !layer
      member: !layer {{ .SpaceID }}`))
)

func createOrgSpace(orgSpace *OrgSpace) io.Reader {
	reader, writer := io.Pipe()

	go func() {
		defer writer.Close()
		// TODO: error handling
		orgSpaceTemplate.Execute(writer, orgSpace)
	}()

	return reader
}

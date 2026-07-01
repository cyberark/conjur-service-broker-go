// Package servicebroker provides an implementation of the generated service broker server
package servicebroker

import (
	"encoding/json"
	"fmt"
	"regexp"

	"github.com/cyberark/conjur-service-broker-go/pkg/conjur"
)

var uuidRE = regexp.MustCompile(`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`)

type context struct {
	Platform  string  `json:"platform"`
	OrgID     string  `json:"organization_guid"`
	OrgName   *string `json:"organization_name"`
	SpaceID   string  `json:"space_guid"`
	SpaceName *string `json:"space_name"`
}

// parseContext extracts org/space identity from an OSB context object.
// It returns an error if organization_guid or space_guid is present but not a
// valid RFC-4122 UUID, preventing path-traversal and cross-tenant access.
func parseContext(ctx *Context) (context, error) {
	orgID := str(ctx, "organization_guid")
	spaceID := str(ctx, "space_guid")

	if orgID != "" && !uuidRE.MatchString(orgID) {
		return context{}, fmt.Errorf("organization_guid is not a valid UUID")
	}
	if spaceID != "" && !uuidRE.MatchString(spaceID) {
		return context{}, fmt.Errorf("space_guid is not a valid UUID")
	}

	return context{
		str(ctx, "platform"),
		orgID,
		strOrNil(ctx, "organization_name"),
		spaceID,
		strOrNil(ctx, "space_name"),
	}, nil
}

func str(ctx *Context, name string) string {
	if ctx == nil {
		return ""
	}
	v, found := (*ctx)[name]
	if !found {
		return ""
	}
	s, _ := v.(string)
	return s
}

func strOrNil(ctx *Context, name string) *string {
	if ctx == nil {
		return nil
	}
	v, found := (*ctx)[name]
	if !found {
		return nil
	}
	s, ok := v.(string) //nolint:gocritic,sloppyTypeAssert //this is a false positive
	if !ok {
		return nil
	}
	return &s
}

func object(policy *conjur.Policy) *Object {
	// we convert to map[string]interface{} using json library, it might not be the most readable implementation but on the other hand reflection wouldn't be readable either
	bytes, err := json.Marshal(policy)
	if err != nil {
		panic(err)
	}
	var res Object
	err = json.Unmarshal(bytes, &res)
	if err != nil {
		panic(err)
	}
	return &res
}

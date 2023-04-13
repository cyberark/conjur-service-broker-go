// Package servicebroker provides an implementation of the generated service broker server
package servicebroker

import (
	"encoding/json"

	"github.com/cyberark/conjur-service-broker/pkg/conjur"
)

type context struct {
	Platform  string  `json:"platform"`
	OrgID     string  `json:"organization_guid"`
	OrgName   *string `json:"organization_name"`
	SpaceID   string  `json:"space_guid"`
	SpaceName *string `json:"space_name"`
}

func parseContext(ctx *Context) context {
	return context{
		str(ctx, "platform"),
		str(ctx, "organization_guid"),
		strOrNil(ctx, "organization_name"),
		str(ctx, "space_guid"),
		strOrNil(ctx, "space_name"),
	}
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
	s, ok := v.(string) // nolint:gocritic,sloppyTypeAssert // this is a false positive
	if !ok {
		return nil
	}
	return &s
}

func object(policy *conjur.Policy) *Object {
	// TODO: would it better to use reflection?
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

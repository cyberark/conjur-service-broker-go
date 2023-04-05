package servicebroker

import (
	"encoding/json"
	"fmt"

	"github.com/cyberark/conjur-service-broker/pkg/conjur"
)

const (
	expectedServiceID = "c024e536-6dc4-45c6-8a53-127e7f8275ab"
	expectedPlanID    = "3a116ac2-fc8b-496f-a715-e9a1b205d05c.community"
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
	s, ok := v.(string)
	if !ok {
		return ""
	}
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
	s, ok := v.(string)
	if !ok {
		return nil
	}
	return &s
}

func validateServiceAndPlan(serviceID string, planID *string) error {
	if serviceID != expectedServiceID {
		return fmt.Errorf("invalid serviceID expected %v, got %v", expectedServiceID, serviceID)
	}
	if planID != nil && *planID != expectedPlanID {
		return fmt.Errorf("invalid planID expected %s, got %s", expectedPlanID, *planID)
	}
	return nil
}

func object(policy *conjur.CreatedPolicy) *Object {
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

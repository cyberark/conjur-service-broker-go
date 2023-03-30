package conjur

import (
	"fmt"
)

// OrgSpace allows basic operations on the organization and space
type OrgSpace struct {
	OrgID     string
	OrgName   string
	SpaceID   string
	SpaceName string
}

// NewOrgSpace creates an OrgSpace based on
func NewOrgSpace(orgID, spaceID string, orgName, spaceName *string) *OrgSpace {
	res := &OrgSpace{
		OrgID:   orgID,
		SpaceID: spaceID,
	}
	if orgName != nil {
		res.OrgName = *orgName
	}
	if spaceName != nil {
		res.SpaceName = *spaceName
	}
	return res
}

func (o *OrgSpace) orgPolicyID(client ResourceChecker) string {
	return fmt.Sprintf("%v/%v", client.BasePolicy(), o.OrgID)
}

func (o *OrgSpace) spacePolicyID(client ResourceChecker) string {
	return fmt.Sprintf("%v/%v/%v", client.BasePolicy(), o.OrgID, o.SpaceID)
}

func (o *OrgSpace) spaceLayerID(client ResourceChecker) string {
	return fmt.Sprintf("%v/%v/%v", client.BaseLayer(), o.OrgID, o.SpaceID)
}

// CreatePolicy creates all needed conjur polices for given org and space
func (o *OrgSpace) CreatePolicy(client PolicyLoader) error {
	err := client.UpsertPolicy(createOrgSpace(o))
	return err
}

// Exists checks existence of conjur org and space policies
func (o *OrgSpace) Exists(client ResourceChecker) (bool, error) {
	ok, err := client.CheckResource(o.orgPolicyID(client))
	if err != nil {
		return false, err
	}
	if !ok {
		return ok, nil
	}
	ok, err = client.CheckResource(o.spacePolicyID(client))
	if err != nil {
		return false, err
	}
	if !ok {
		return ok, nil
	}
	ok, err = client.CheckResource(o.spaceLayerID(client))
	if err != nil {
		return false, err
	}
	if !ok {
		return ok, nil
	}
	return true, nil
}

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
	client    Client
}

// NewOrgSpace creates an OrgSpace based on
func NewOrgSpace(client Client, orgID, spaceID string, orgName, spaceName *string) *OrgSpace {
	res := &OrgSpace{
		OrgID:   orgID,
		SpaceID: spaceID,
		client:  client,
	}
	if orgName != nil {
		res.OrgName = *orgName
	}
	if spaceName != nil {
		res.SpaceName = *spaceName
	}
	return res
}

func (o *OrgSpace) orgPolicyID() string {
	return fmt.Sprintf("%v/%v", o.client.BasePolicy(), o.OrgID)
}

func (o *OrgSpace) spacePolicyID() string {
	return fmt.Sprintf("%v/%v/%v", o.client.BasePolicy(), o.OrgID, o.SpaceID)
}

func (o *OrgSpace) spaceLayerID() string {
	return fmt.Sprintf("%v/%v/%v", o.client.BaseLayer(), o.OrgID, o.SpaceID)
}

// CreatePolicy creates all needed conjur polices for given org and space
func (o *OrgSpace) CreatePolicy() error {
	err := o.client.UpsertPolicy(createOrgSpace(o))
	return err
}

// Exists checks existence of conjur org and space policies
func (o *OrgSpace) Exists() (bool, error) {
	ok, err := o.client.CheckResource(o.orgPolicyID())
	if err != nil {
		return false, err
	}
	if !ok {
		return ok, nil
	}
	ok, err = o.client.CheckResource(o.spacePolicyID())
	if err != nil {
		return false, err
	}
	if !ok {
		return ok, nil
	}
	ok, err = o.client.CheckResource(o.spaceLayerID())
	if err != nil {
		return false, err
	}
	if !ok {
		return ok, nil
	}
	return true, nil
}

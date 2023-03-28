package conjur

import (
	"fmt"

	"github.com/cyberark/conjur-api-go/conjurapi"
)

// OrgSpace allows basic operations on the organization and space
type OrgSpace struct {
	// TODO: could it be private?
	OrgID     string
	OrgName   string
	SpaceID   string
	SpaceName string
	c         *Client
}

// NewOrgSpace creates an OrgSpace based on
func (c *Client) NewOrgSpace(orgID, spaceID string, orgName, spaceName *string) *OrgSpace {
	res := &OrgSpace{
		OrgID:   orgID,
		SpaceID: spaceID,
		c:       c,
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
	return fmt.Sprintf("%v:policy:%v%v", o.c.ConjurAccount, o.c.ConjurPolicy, o.OrgID)
}

func (o *OrgSpace) spacePolicyID() string {
	return fmt.Sprintf("%v:policy:%v%v/%v", o.c.ConjurAccount, o.c.ConjurPolicy, o.OrgID, o.SpaceID)
}

func (o *OrgSpace) spaceLayerID() string {
	return fmt.Sprintf("%v:layer:%v%v/%v", o.c.ConjurAccount, o.c.ConjurPolicy, o.OrgID, o.SpaceID)
}

// CreatePolicy creates all needed conjur polices for given org and space
func (o *OrgSpace) CreatePolicy() error {
	_, err := o.c.Client.LoadPolicy(
		conjurapi.PolicyModePut,
		o.spacePolicyID(),
		createOrgSpace(o),
	)
	return err
}

// TODO: rethink DX - only one of two methods should stay

// ValidateExists checks existence of conjur org and space policies
func (o *OrgSpace) ValidateExists() error {
	orgPolicy, err := o.c.Client.Resource(o.orgPolicyID())
	if err != nil {
		// TODO: return wrapped fixed error to be used with errors.Is
		return fmt.Errorf("unable to find org policy %v: %w", o.orgPolicyID(), err)
	}
	if len(orgPolicy) == 0 {
		return fmt.Errorf("unable to find org policy %v", o.orgPolicyID())
	}
	spacePolicy, err := o.c.Client.Resource(o.spacePolicyID())
	if err != nil {
		return fmt.Errorf("unable to find space policy %v: %w", o.spacePolicyID(), err)
	}
	if len(spacePolicy) == 0 {
		return fmt.Errorf("unable to find space policy %v", o.spacePolicyID())
	}
	spaceLayer, err := o.c.Client.Resource(o.spaceLayerID())
	if err != nil {
		return fmt.Errorf("unable to find space layer %v in policy: %w", o.spaceLayerID(), err)
	}
	if len(spaceLayer) == 0 {
		return fmt.Errorf("unable to find space layer %v in policy", o.spaceLayerID())
	}
	return nil
}

// Exists checks existence of conjur org and space policies
func (o *OrgSpace) Exists() (bool, error) {
	// TODO: how to verify? by checking permissions or just resource?
	//_, err := o.c.Client.CheckPermission(o.orgPolicyID(), readPermission)
	//if err != nil {
	//	return false, err
	//}
	//o.c.Client.LoadPolicy()
	return true, nil
}

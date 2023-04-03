package conjur

import (
	"fmt"
	"io"
)

type orgSpace struct {
	OrgID     string
	OrgName   string
	SpaceID   string
	SpaceName string
	client    Client
}

// OrgSpace allows basic operations on the organization and space
type OrgSpace interface {
	CreatePolicy() error
	Exists() (bool, error)
}

// NewOrgSpace creates an OrgSpace based on
func NewOrgSpace(client Client, orgID, spaceID string, orgName, spaceName *string) OrgSpace {
	res := &orgSpace{
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

func (o *orgSpace) orgPolicyID() string {
	return fmt.Sprintf("%v/%v", o.client.BasePolicy(), o.OrgID)
}

func (o *orgSpace) spacePolicyID() string {
	return fmt.Sprintf("%v/%v/%v", o.client.BasePolicy(), o.OrgID, o.SpaceID)
}

func (o *orgSpace) spaceLayerID() string {
	return fmt.Sprintf("%v/%v/%v", o.client.BaseLayer(), o.OrgID, o.SpaceID)
}

// CreatePolicy creates all needed conjur polices for given org and space
func (o *orgSpace) CreatePolicy() error {
	orgSpace, err := createOrgSpaceYAML(o)
	if err != nil {
		return err
	}
	err = o.client.UpsertPolicy(orgSpace)
	return err
}

// Exists checks existence of conjur org and space policies
func (o *orgSpace) Exists() (bool, error) {
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

func createOrgSpaceYAML(o *orgSpace) (io.Reader, error) {
	policy := PolicyDocument{
		NewTag(Policy{
			Id:          o.OrgID,
			Annotations: policyAnnotations(o),
			Body: []Tag{
				NewTag[Layer](""),
				NewTag[Policy](Policy{
					Id: o.SpaceID,
					Body: []Tag{
						NewTag[Layer](""),
					},
					Annotations: subPolicyAnnotations(o),
				}),
				NewTag(Grant{
					Role:   NewRef[Layer](""),
					Member: NewRef[Layer](o.SpaceID),
				}),
			},
		}),
	}
	return policyReader(policy)
}

func policyAnnotations(o *orgSpace) map[string]string {
	if len(o.OrgName) == 0 || len(o.OrgName) == 0 {
		return nil
	}
	return map[string]string{"pcf/type": "org", "pcf/orgName": o.OrgName}
}

func subPolicyAnnotations(o *orgSpace) map[string]string {
	if len(o.OrgName) == 0 || len(o.OrgName) == 0 {
		return nil
	}
	return map[string]string{"pcf/type": "space", "pcf/orgName": o.OrgName, "pcf/spaceName": o.SpaceName}
}

package conjur

import (
	"bytes"
	"fmt"
	"io"

	"gopkg.in/yaml.v3"
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
	orgSpace, err := createOrgSpace(o)
	if err != nil {
		return err
	}
	err = o.client.UpsertPolicy(orgSpace)
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

func createOrgSpace(orgSpace *OrgSpace) (io.Reader, error) {
	policy := PolicyDocument{
		NewTag(Policy{
			Id: orgSpace.OrgID,
			Body: []interface{}{
				NewTag[Layer](""),
				NewTag(Policy{
					Id: orgSpace.SpaceID,
					Body: []interface{}{
						NewTag[Layer](""),
					},
				}),
				NewTag(Grant{
					Role:   NewTag[any](Layer("")),
					Member: NewTag[any](Layer(orgSpace.SpaceID)),
				}),
			},
		}),
	}
	if len(orgSpace.OrgName) > 0 && len(orgSpace.OrgName) > 0 { // TODO: make this better
		policy[0].v.Annotations = map[string]string{"pcf/type": "org", "pcf/orgName": orgSpace.OrgName}
		policy[0].v.Body[1].(*Tag[Policy]).v.Annotations = map[string]string{"pcf/type": "space", "pcf/orgName": orgSpace.OrgName, "pcf/spaceName": orgSpace.SpaceName}
	}
	res := new(bytes.Buffer)
	encoder := yaml.NewEncoder(res)
	err := encoder.Encode(policy)
	if err != nil {
		return nil, err
	}
	return res, err
}

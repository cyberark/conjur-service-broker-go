package conjur

import (
	"bytes"
	"fmt"
	"io"

	"gopkg.in/yaml.v3"
)

// OrgSpace allows basic operations on the organization and space
type orgSpace struct {
	OrgID     string
	OrgName   string
	SpaceID   string
	SpaceName string
	client    Client
}

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
	orgSpace, err := createOrgSpace(o)
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

func createOrgSpace(o *orgSpace) (io.Reader, error) {
	policy := PolicyDocument{
		NewTag(Policy{
			Id: o.OrgID,
			Body: []interface{}{
				NewTag[Layer](""),
				NewTag(Policy{
					Id: o.SpaceID,
					Body: []interface{}{
						NewTag[Layer](""),
					},
				}),
				NewTag(Grant{
					Role:   NewTag[any](Layer("")),
					Member: NewTag[any](Layer(o.SpaceID)),
				}),
			},
		}),
	}
	if len(o.OrgName) > 0 && len(o.OrgName) > 0 { // TODO: make this better
		policy[0].v.Annotations = map[string]string{"pcf/type": "org", "pcf/orgName": o.OrgName}
		policy[0].v.Body[1].(*Tag[Policy]).v.Annotations = map[string]string{"pcf/type": "space", "pcf/orgName": o.OrgName, "pcf/spaceName": o.SpaceName}
	}
	res := new(bytes.Buffer)
	encoder := yaml.NewEncoder(res)
	err := encoder.Encode(policy)
	if err != nil {
		return nil, err
	}
	return res, err
}

package conjur

import (
	"fmt"
	"io"

	"github.com/doodlesbykumbi/conjur-policy-go/pkg/conjurpolicy"
)

type orgSpace struct {
	orgID     string
	orgName   string
	spaceID   string
	spaceName string
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
		orgID:   orgID,
		spaceID: spaceID,
		client:  client,
	}
	if orgName != nil {
		res.orgName = *orgName
	}
	if spaceName != nil {
		res.spaceName = *spaceName
	}
	return res
}

func (o *orgSpace) orgPolicyID() string {
	config := o.client.Config()
	return fmt.Sprintf("%v:policy:%v/%v", config.ConjurAccount, config.ConjurPolicy, o.orgID)
}

func (o *orgSpace) spacePolicyID() string {
	config := o.client.Config()
	return fmt.Sprintf("%v:policy:%v/%v/%v", config.ConjurAccount, config.ConjurPolicy, o.orgID, o.spaceID)
}

func (o *orgSpace) spaceLayerID() string {
	config := o.client.Config()
	return fmt.Sprintf("%v:layer:%v/%v/%v", config.ConjurAccount, config.ConjurPolicy, o.orgID, o.spaceID)
}

// CreatePolicy creates all needed conjur polices for given org and space
func (o *orgSpace) CreatePolicy() error {
	yaml, err := o.createOrgSpaceYAML()
	if err != nil {
		return err
	}
	_, err = o.client.UpsertPolicy(yaml)
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

func (o *orgSpace) createOrgSpaceYAML() (io.Reader, error) {
	policy := conjurpolicy.PolicyStatements{
		conjurpolicy.Policy{
			Id:          o.orgID,
			Annotations: o.orgAnnotations(),
			Body: conjurpolicy.PolicyStatements{
				conjurpolicy.Layer{},
				conjurpolicy.Policy{
					Id: o.spaceID,
					Body: conjurpolicy.PolicyStatements{
						conjurpolicy.Layer{},
					},
					Annotations: o.spaceAnnotations(),
				},
				conjurpolicy.Grant{
					Role:   conjurpolicy.LayerRef(""),
					Member: conjurpolicy.LayerRef(o.spaceID),
				},
			},
		},
	}
	return policyReader(policy)
}

func (o *orgSpace) orgAnnotations() map[string]string {
	if len(o.orgName) == 0 || len(o.orgName) == 0 {
		return nil
	}
	return map[string]string{"pcf/type": "org", "pcf/orgName": o.orgName}
}

func (o *orgSpace) spaceAnnotations() map[string]string {
	if len(o.orgName) == 0 || len(o.orgName) == 0 {
		return nil
	}
	return map[string]string{"pcf/type": "space", "pcf/orgName": o.orgName, "pcf/spaceName": o.spaceName}
}

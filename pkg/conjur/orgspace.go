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
	account   string
	policy    string
	client    Client
}

// OrgSpace allows basic operations on the organization and space
type OrgSpace interface {
	CreatePolicy() error
	Exists() (bool, error)
}

// NewOrgSpace creates an OrgSpace based on
func NewOrgSpace(client Client, orgID, spaceID string, orgName, spaceName *string) OrgSpace {
	cfg := client.Config()
	res := &orgSpace{
		orgID:   orgID,
		spaceID: spaceID,
		client:  client,
		policy:  cfg.ConjurPolicy,
		account: cfg.ConjurAccount,
	}
	if orgName != nil {
		res.orgName = *orgName
	}
	if spaceName != nil {
		res.spaceName = *spaceName
	}
	return res
}

func (o *orgSpace) orgSpacePolicyID() string {
	return fmt.Sprintf("%s/%s/%s", o.policy, o.orgID, o.spaceID)
}

func (o *orgSpace) orgPolicyResourceID() string {
	return composeID(o.account, KindPolicy, fmt.Sprintf("%s/%s", o.policy, o.orgID))
}

func (o *orgSpace) spacePolicyResourceID() string {
	return composeID(o.account, KindPolicy, fmt.Sprintf("%s/%s/%s", o.policy, o.orgID, o.spaceID))
}

func (o *orgSpace) spaceLayerResourceID() string {
	return composeID(o.account, KindLayer, fmt.Sprintf("%s/%s/%s", o.policy, o.orgID, o.spaceID))
}

// CreatePolicy creates all needed conjur polices for given org and space
func (o *orgSpace) CreatePolicy() error {
	yaml, err := o.createOrgSpaceYAML()
	if err != nil {
		return err
	}
	_, err = o.client.UpsertPolicy(yaml, o.policy)
	if exists, err := o.Exists(); err != nil || !exists {
		return fmt.Errorf("failed to validate policy exists: %w", err)
	}
	// TODO: check toggle
	yaml, err = o.createSpaceHostYAML()
	if err != nil {
		return err
	}
	policy, err := o.client.UpsertPolicy(yaml, o.orgSpacePolicyID())
	apiKey, err := apiKey(policy)
	if err != nil {
		return err
	}
	err = o.client.SetVariable(fmt.Sprintf("%s/%s", o.orgSpacePolicyID(), spaceHostAPIKey), apiKey)
	if err != nil {
		return err
	}
	return nil
}

// Exists checks existence of conjur org and space policies
func (o *orgSpace) Exists() (bool, error) {
	ok, err := o.client.CheckResource(o.orgPolicyResourceID())
	if err != nil {
		return false, err
	}
	if !ok {
		return ok, nil
	}
	ok, err = o.client.CheckResource(o.spacePolicyResourceID())
	if err != nil {
		return false, err
	}
	if !ok {
		return ok, nil
	}
	ok, err = o.client.CheckResource(o.spaceLayerResourceID())
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

func (o *orgSpace) createSpaceHostYAML() (io.Reader, error) {
	policy := conjurpolicy.PolicyStatements{
		conjurpolicy.Host{},
		conjurpolicy.Grant{
			Role:   conjurpolicy.LayerRef(""),
			Member: conjurpolicy.HostRef(""),
		},
		conjurpolicy.Variable{
			Id: spaceHostAPIKey,
		},
		conjurpolicy.Permit{
			Role:       conjurpolicy.HostRef("/" + o.client.ClientHostID()),
			Privileges: []conjurpolicy.Privilege{conjurpolicy.PrivilegeRead},
			Resources:  conjurpolicy.VariableRef(spaceHostAPIKey),
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

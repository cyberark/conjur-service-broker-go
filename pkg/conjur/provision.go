// Package conjur provides a wrapper around conjur go SDK
package conjur

import (
	"fmt"
	"io"

	"github.com/doodlesbykumbi/conjur-policy-go/pkg/conjurpolicy"
)

type provision struct {
	orgID     string
	orgName   string
	spaceID   string
	spaceName string
	client    *client
}

// Provision allows basic operations on the organization and space
type Provision interface {
	ProvisionOrgSpacePolicy() error
	ProvisionHostPolicy() error
}

func (p *provision) orgSpacePolicyID() string {
	return slashJoin(p.client.config.ConjurPolicy, p.orgID, p.spaceID)
}

func (p *provision) orgPolicyResourceID() string {
	return composeID(p.client.config.ConjurAccount, KindPolicy, slashJoin(p.client.config.ConjurPolicy, p.orgID))
}

func (p *provision) spacePolicyResourceID() string {
	return composeID(p.client.config.ConjurAccount, KindPolicy, p.orgSpacePolicyID())
}

func (p *provision) spaceGroupResourceID() string {
	return composeID(p.client.config.ConjurAccount, KindGroup, p.orgSpacePolicyID())
}

// ProvisionOrgSpacePolicy creates all needed conjur polices for given org and space
func (p *provision) ProvisionOrgSpacePolicy() error {
	yaml, err := p.provisionOrgSpaceYAML()
	if err != nil {
		return err
	}
	_, err = p.client.upsertPolicy(yaml, p.client.config.ConjurPolicy)
	if err != nil {
		return fmt.Errorf("failed to create policy: %w", err)
	}
	if exists, err := p.orgSpacePolicyExists(); err != nil || !exists {
		return fmt.Errorf("failed to validate policy exists: %w", err)
	}
	return nil
}

// ProvisionHostPolicy creates conjur polices on space host level
func (p *provision) ProvisionHostPolicy() error {
	yaml, err := p.provisionHostYAML()
	if err != nil {
		return fmt.Errorf("failed to create space host level yaml: %w", err)
	}
	policy, err := p.client.upsertPolicy(yaml, p.orgSpacePolicyID())
	if err != nil {
		return fmt.Errorf("failed to create space host level policy: %w", err)
	}
	apiKey, err := apiKey(policy)
	if err != nil {
		return err
	}
	err = p.client.setVariable(fmt.Sprintf("%s/%s", p.orgSpacePolicyID(), spaceHostAPIKey), apiKey)
	if err != nil {
		return err
	}
	return nil
}

func (p *provision) orgSpacePolicyExists() (bool, error) {
	for _, id := range []string{p.orgPolicyResourceID(), p.spacePolicyResourceID(), p.spaceGroupResourceID()} {
		ok, err := p.client.resourceExists(id)
		if err != nil {
			return false, err
		}
		if !ok {
			return ok, nil
		}
	}
	return true, nil
}

func (p *provision) provisionOrgSpaceYAML() (io.Reader, error) {
	policy := conjurpolicy.PolicyStatements{
		conjurpolicy.Policy{
			Id:          p.orgID,
			Annotations: p.orgAnnotations(),
			Body: conjurpolicy.PolicyStatements{
				conjurpolicy.Group{},
				conjurpolicy.Policy{
					Id: p.spaceID,
					Body: conjurpolicy.PolicyStatements{
						conjurpolicy.Group{},
					},
					Annotations: p.spaceAnnotations(),
				},
				conjurpolicy.Grant{
					Role:   conjurpolicy.GroupRef(""),
					Member: conjurpolicy.GroupRef(p.spaceID),
				},
			},
		},
	}
	return policyReader(policy)
}

func (p *provision) provisionHostYAML() (io.Reader, error) {
	policy := conjurpolicy.PolicyStatements{
		conjurpolicy.Host{},
		conjurpolicy.Grant{
			Role:   conjurpolicy.GroupRef(""),
			Member: conjurpolicy.HostRef(""),
		},
		conjurpolicy.Variable{
			Id: spaceHostAPIKey,
		},
		conjurpolicy.Permit{
			Role:       conjurpolicy.HostRef("/" + p.client.clientHostID()),
			Privileges: []conjurpolicy.Privilege{conjurpolicy.PrivilegeRead},
			Resources:  conjurpolicy.VariableRef(spaceHostAPIKey),
		},
	}
	return policyReader(policy)
}

func (p *provision) orgAnnotations() map[string]interface{} {
	if len(p.orgName) == 0 || len(p.spaceName) == 0 {
		return nil
	}
	return map[string]interface{}{"pcf/type": "org", "pcf/orgName": p.orgName}
}

func (p *provision) spaceAnnotations() map[string]interface{} {
	if len(p.orgName) == 0 || len(p.spaceName) == 0 {
		return nil
	}
	return map[string]interface{}{"pcf/type": "space", "pcf/orgName": p.orgName, "pcf/spaceName": p.spaceName}
}

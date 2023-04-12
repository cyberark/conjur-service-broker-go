// Package conjur provides a wrapper around conjur go SDK
package conjur

import (
	"fmt"
	"io"

	"github.com/cyberark/conjur-service-broker/internal/ctxutil"
	"github.com/doodlesbykumbi/conjur-policy-go/pkg/conjurpolicy"
)

type provision struct {
	orgID     string
	orgName   string
	spaceID   string
	spaceName string
	account   string
	policy    string
	client    Client
}

// Provision allows basic operations on the organization and space
type Provision interface {
	CreatePolicy(ctx ctxutil.Context) error
	Exists() (bool, error)
}

// NewProvision creates a Provision based on provided configuration
func NewProvision(client Client, orgID, spaceID string, orgName, spaceName *string) Provision {
	cfg := client.Config()
	res := &provision{
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

func (o *provision) orgSpacePolicyID() string {
	return fmt.Sprintf("%s/%s/%s", o.policy, o.orgID, o.spaceID)
}

func (o *provision) orgPolicyResourceID() string {
	return composeID(o.account, KindPolicy, fmt.Sprintf("%s/%s", o.policy, o.orgID))
}

func (o *provision) spacePolicyResourceID() string {
	return composeID(o.account, KindPolicy, fmt.Sprintf("%s/%s/%s", o.policy, o.orgID, o.spaceID))
}

func (o *provision) spaceLayerResourceID() string {
	return composeID(o.account, KindLayer, fmt.Sprintf("%s/%s/%s", o.policy, o.orgID, o.spaceID))
}

// CreatePolicy creates all needed conjur polices for given org and space
func (o *provision) CreatePolicy(ctx ctxutil.Context) error {
	yaml, err := o.createOrgSpaceYAML()
	if err != nil {
		return err
	}
	_, err = o.client.UpsertPolicy(yaml, o.policy)
	if err != nil {
		return fmt.Errorf("failed to create policy: %w", err)
	}
	if exists, err := o.Exists(); err != nil || !exists {
		return fmt.Errorf("failed to validate policy exists: %w", err)
	}
	if !ctx.IsEnableSpaceIdentity() {
		return nil
	}
	yaml, err = o.createSpaceHostYAML()
	if err != nil {
		return fmt.Errorf("failed to create space host level yaml: %w", err)
	}
	policy, err := o.client.UpsertPolicy(yaml, o.orgSpacePolicyID())
	if err != nil {
		return fmt.Errorf("failed to create space host level policy: %w", err)
	}
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
func (o *provision) Exists() (bool, error) {
	ok, err := o.client.ResourceExists(o.orgPolicyResourceID())
	if err != nil {
		return false, err
	}
	if !ok {
		return ok, nil
	}
	ok, err = o.client.ResourceExists(o.spacePolicyResourceID())
	if err != nil {
		return false, err
	}
	if !ok {
		return ok, nil
	}
	ok, err = o.client.ResourceExists(o.spaceLayerResourceID())
	if err != nil {
		return false, err
	}
	if !ok {
		return ok, nil
	}
	return true, nil
}

func (o *provision) createOrgSpaceYAML() (io.Reader, error) {
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

func (o *provision) createSpaceHostYAML() (io.Reader, error) {
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

func (o *provision) orgAnnotations() map[string]string {
	if len(o.orgName) == 0 || len(o.spaceName) == 0 {
		return nil
	}
	return map[string]string{"pcf/type": "org", "pcf/orgName": o.orgName}
}

func (o *provision) spaceAnnotations() map[string]string {
	if len(o.orgName) == 0 || len(o.spaceName) == 0 {
		return nil
	}
	return map[string]string{"pcf/type": "space", "pcf/orgName": o.orgName, "pcf/spaceName": o.spaceName}
}

// Package conjur provides a wrapper around conjur go SDK
package conjur

import (
	"fmt"
	"io"

	"github.com/cyberark/conjur-api-go/conjurapi"
	"github.com/doodlesbykumbi/conjur-policy-go/pkg/conjurpolicy"
)

type bind struct {
	orgID     string
	spaceID   string
	bindingID string
	hostID    string
	client    *client
}

// Policy is a result of policy creation
type Policy struct {
	Account        string `json:"account"`
	ApplianceURL   string `json:"appliance_url"`
	AuthnLogin     string `json:"authn_login"`
	AuthnAPIKey    string `json:"authn_api_key"`
	SslCertificate string `json:"ssl_certificate"`
	Version        uint32 `json:"version"`
}

// Bind allows operations needed for binding an instance
type Bind interface {
	// BindHostPolicy creates policy needed for binding on host identity level
	BindHostPolicy() (*Policy, error)
	// BindSpacePolicy creates policy needed for binding on space identity level
	BindSpacePolicy() (*Policy, error)
	// DeleteBindHostPolicy deletes policy created for binding on host identity level
	DeleteBindHostPolicy() error

	// HostExists checks for host existence for space and host level identity
	HostExists() (bool, error)
}

func (b *bind) BindHostPolicy() (*Policy, error) {
	yaml, err := b.createBindYAML()
	if err != nil {
		return nil, err
	}
	policy, err := b.client.upsertPolicy(yaml, b.policy())
	if err != nil {
		return nil, err
	}
	return b.onlyPolicy(policy)
}

func (b *bind) BindSpacePolicy() (*Policy, error) {
	config := b.client.config
	orgSpaceID := slashJoin(config.ConjurPolicy, b.orgID, b.spaceID)
	apiKey, err := b.client.getVariable(slashJoin(orgSpaceID, spaceHostAPIKey))
	if err != nil {
		return nil, err
	}
	return &Policy{
		Account:        config.ConjurAccount,
		ApplianceURL:   config.ConjurApplianceURL,
		AuthnLogin:     slashJoin(KindHost.String(), orgSpaceID),
		AuthnAPIKey:    apiKey,
		SslCertificate: config.ConjurSSLCertificate,
		Version:        config.ConjurVersion,
	}, nil
}

func (b *bind) DeleteBindHostPolicy() error {
	err := b.client.rotateAPIKey(b.hostID)
	if err != nil {
		return err
	}
	yaml, err := b.deleteBindYAML()
	if err != nil {
		return err
	}
	_, err = b.client.replacePolicy(yaml, b.policy())
	if err != nil {
		return err
	}
	return nil
}

func (b *bind) onlyPolicy(policy *conjurapi.PolicyResponse) (*Policy, error) {
	if policy == nil || len(policy.CreatedRoles) != 1 {
		return nil, fmt.Errorf("expecting exactly one created role")
	}
	var roleID string
	var role conjurapi.CreatedRole
	for k, v := range policy.CreatedRoles {
		roleID = k
		role = v
	}
	if roleID != role.ID {
		return nil, fmt.Errorf("created role ID does not match %v != %v", roleID, role.ID)
	}
	config := b.client.config
	return &Policy{
		Account:        config.ConjurAccount,
		ApplianceURL:   config.ConjurApplianceURL,
		AuthnLogin:     dropAccount(roleID),
		AuthnAPIKey:    role.APIKey,
		SslCertificate: config.ConjurSSLCertificate,
		Version:        config.ConjurVersion,
	}, nil
}

func dropAccount(id string) string {
	_, kind, identifier := parseID(id)
	if !kind.IsAKind() {
		return identifier
	}
	return slashJoin(kind.String(), identifier)
}

func (b *bind) HostExists() (bool, error) {
	return b.client.roleExists(b.hostID)
}

func (b *bind) useSpace() bool {
	return len(b.orgID) > 0 && len(b.spaceID) > 0
}

func (b *bind) policy() string {
	if b.useSpace() {
		return slashJoin(b.client.config.ConjurPolicy, b.orgID, b.spaceID)
	}
	return b.client.config.ConjurPolicy
}

func (b *bind) createBindYAML() (io.Reader, error) {
	policy := conjurpolicy.PolicyStatements{
		conjurpolicy.Host{
			Id:          b.bindingID,
			Annotations: b.hostAnnotations(),
		},
	}
	if b.useSpace() {
		policy = append(policy, conjurpolicy.Grant{
			Role:   conjurpolicy.GroupRef(""),
			Member: conjurpolicy.HostRef(b.bindingID),
		})
	}
	return policyReader(policy)
}

func (b *bind) deleteBindYAML() (io.Reader, error) {
	policy := conjurpolicy.PolicyStatements{
		conjurpolicy.Delete{
			Record: conjurpolicy.HostRef(b.bindingID),
		},
	}
	return policyReader(policy)
}

func (b *bind) hostAnnotations() map[string]interface{} {
	platform, err := b.client.platformAnnotation()
	if len(platform) == 0 || err != nil {
		return map[string]interface{}{"authn/api-key": true}
	}
	return map[string]interface{}{"authn/api-key": true, platform: true}
}

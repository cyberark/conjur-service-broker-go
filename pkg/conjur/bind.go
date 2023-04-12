// Package conjur provides a wrapper around conjur go SDK
package conjur

import (
	"fmt"
	"io"
	"strings"

	"github.com/cyberark/conjur-api-go/conjurapi"
	"github.com/cyberark/conjur-service-broker/internal/ctxutil"
	"github.com/doodlesbykumbi/conjur-policy-go/pkg/conjurpolicy"
)

type bind struct {
	orgID     string
	spaceID   string
	bindingID string
	client    Client
}

// CreatedPolicy is a result of policy creation
type CreatedPolicy struct {
	Account        string `json:"account"`
	ApplianceURL   string `json:"appliance_url"`
	AuthnLogin     string `json:"authn_login"`
	AuthnAPIKey    string `json:"authn_api_key"`
	SslCertificate string `json:"ssl_certificate"`
	Version        uint32 `json:"version"`
}

// Bind allows operations needed for binding an instance
type Bind interface {
	CreatePolicy(ctxutil.Context) (*CreatedPolicy, error)
	HostExists(ctxutil.Context) (bool, error)
	DeletePolicy(ctxutil.Context) error
}

// NewBind creates new binding service
func NewBind(client Client, orgID, spaceID, bindingID string) Bind {
	res := &bind{
		orgID:     orgID,
		spaceID:   spaceID,
		bindingID: bindingID,
		client:    client,
	}
	return res
}

// FromBindingID creates new binding service based on existing binding by it's ID, org and space IDs would be queried from conjur
func FromBindingID(client Client, bindingID string) (Bind, error) {
	orgID, spaceID, err := client.OrgSpaceFromBindingID(bindingID)
	if err != nil {
		return nil, fmt.Errorf("failed to create binding service from binding id: %w", err)
	}
	res := &bind{
		orgID:     orgID,
		spaceID:   spaceID,
		bindingID: bindingID,
		client:    client,
	}
	return res, nil
}

func (b *bind) hostID(ctx ctxutil.Context) string {
	return b.hostPolicyID(ctx) + "/" + b.bindingID
}

func (b *bind) hostPolicyID(ctx ctxutil.Context) string {
	return composeID(b.client.Config().ConjurAccount, KindHost, b.policy(ctx))
}

func (b *bind) CreatePolicy(ctx ctxutil.Context) (*CreatedPolicy, error) {
	yaml, err := b.createBindYAML()
	if err != nil {
		return nil, err
	}
	policy, err := b.client.UpsertPolicy(yaml, b.policy(ctx))
	if err != nil {
		return nil, err
	}
	return b.onlyPolicy(policy)
}

func (b *bind) DeletePolicy(ctx ctxutil.Context) error {
	err := b.client.RotateAPIKey(b.hostID(ctx))
	if err != nil {
		return err
	}
	yaml, err := b.deleteBindYAML()
	if err != nil {
		return err
	}
	_, err = b.client.ReplacePolicy(yaml, b.policy(ctx))
	if err != nil {
		return err
	}
	return nil
}

func (b *bind) onlyPolicy(policy *conjurapi.PolicyResponse) (*CreatedPolicy, error) {
	if len(policy.CreatedRoles) != 1 {
		return nil, fmt.Errorf("expecting exactly one created role")
	}
	var roleID string
	var role conjurapi.CreatedRole
	for k, v := range policy.CreatedRoles {
		roleID = k
		role = v
	}
	if roleID != role.ID {
		return nil, fmt.Errorf("creatred role ID do not match %v != %v", roleID, role.ID)
	}
	config := b.client.Config()
	return &CreatedPolicy{
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
	return strings.Join([]string{kind.String(), identifier}, "/")
}

func (b *bind) HostExists(ctx ctxutil.Context) (bool, error) {
	return b.client.ResourceExists(b.hostID(ctx))
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
			Role:   conjurpolicy.LayerRef(""),
			Member: conjurpolicy.LayerRef(b.bindingID),
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

func (b *bind) hostAnnotations() map[string]string {
	platform, err := b.client.Platform()
	if len(platform) == 0 || err != nil {
		return nil
	}
	return map[string]string{platform: "true"}
}

func (b *bind) useSpace() bool {
	return len(b.orgID) > 0 && len(b.spaceID) > 0
}

func (b *bind) policy(_ ctxutil.Context) string {
	p := []string{b.client.Config().ConjurPolicy}
	if len(b.orgID) > 0 && len(b.spaceID) > 0 {
		p = append(p, b.orgID, b.spaceID)
	}
	return strings.Join(p, "/")
}

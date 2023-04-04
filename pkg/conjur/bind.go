package conjur

import (
	"fmt"
	"io"
	"strings"

	"github.com/cyberark/conjur-api-go/conjurapi"
	"github.com/doodlesbykumbi/conjur-policy-go/pkg/conjurpolicy"
)

type bind struct {
	orgID     string
	spaceID   string
	bindingID string
	hostID    string
	client    Client
}

type CreatedPolicy struct {
	Account        string `json:"account"`
	ApplianceUrl   string `json:"appliance_url"`
	AuthnLogin     string `json:"authn_login"`
	AuthnApiKey    string `json:"authn_api_key"`
	SslCertificate string `json:"ssl_certificate"`
	Version        uint32 `json:"version"`
}

type Bind interface {
	CreatePolicy() (*CreatedPolicy, error)
	HostExists() (bool, error)
}

func NewBind(client Client, orgID, spaceID, bindingID string) Bind {
	res := &bind{
		orgID:     orgID,
		spaceID:   spaceID,
		bindingID: bindingID,
		client:    client,
	}
	res.hostID = hostID(res)
	return res
}

func hostID(b *bind) string {
	res := []string{"host"}
	if b.useSpace() {
		res = append(res, b.orgID, b.spaceID)
	}
	res = append(res, b.client.Config().ConjurPolicy+"/"+b.bindingID)
	return strings.Join(res, ":")
}

func (b *bind) CreatePolicy() (*CreatedPolicy, error) {
	yaml, err := b.createBindYAML()
	if err != nil {
		return nil, err
	}
	policy, err := b.client.UpsertPolicy(yaml)
	if err != nil {
		return nil, err
	}
	return b.onlyPolicy(policy)
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
		ApplianceUrl:   config.ConjurApplianceURL,
		AuthnLogin:     dropAccount(roleID),
		AuthnApiKey:    role.APIKey,
		SslCertificate: config.ConjurSSLCertificate,
		Version:        config.ConjurVersion,
	}, nil
}

func dropAccount(id string) string {
	_, kind, identifier := parseID(id)
	return strings.Join([]string{kind.String(), identifier}, "/")
}

func (b *bind) HostExists() (bool, error) {
	return b.client.RoleExists(b.hostID)
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

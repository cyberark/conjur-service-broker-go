// Package conjur provides a wrapper around conjur go SDK
package conjur

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/cyberark/conjur-api-go/conjurapi"
	"github.com/cyberark/conjur-api-go/conjurapi/authn"
	"github.com/cyberark/conjur-service-broker-go/pkg/conjur/api"
)

// Client is a wrapper on conjur go sdk allowing creation of bind and provision objects
//
//go:generate mockery --name=Client
type Client interface {
	NewBind(orgID, spaceID, bindingID string, enableSpaceIdentity bool) Bind
	FromBindingID(bindingID string) (Bind, error)
	NewProvision(orgID, spaceID string, orgName, spaceName *string) Provision

	ValidateConnectivity() error
}

type client struct {
	client   api.Client // the writable client - used for operations that requires modification
	roClient api.Client // the read-only client - just for read operations - when the follower URL is configured we use it for the read only client
	config   *Config
}

// NewClient creates new conjur API client wrapper
func (config *Config) NewClient() (Client, error) {
	clientConf, err := conjurapi.LoadConfig()
	clientConf = config.mergeConfig(clientConf)
	if err != nil {
		panic(err)
	}
	loginPair := authn.LoginPair{
		Login:  config.ConjurAuthNLogin,
		APIKey: config.ConjurAuthNAPIKey,
	}
	conjur, err := conjurapi.NewClientFromKey(clientConf, loginPair)
	if err != nil {
		return nil, err
	}
	var roClient *conjurapi.Client
	if len(config.ConjurFollowerURL) > 0 {
		clientConf.ApplianceURL = config.ConjurFollowerURL
		if len(config.ConjurFollowerSSLCertificate) > 0 {
			clientConf.SSLCert = config.ConjurFollowerSSLCertificate
		}
		roClient, err = conjurapi.NewClientFromKey(clientConf, loginPair)
		if err != nil {
			return nil, err
		}
	} else {
		roClient = conjur
	}
	if conjur == nil {
		return nil, fmt.Errorf("failed to create conjur client")
	}

	// Clear credentials from the loginPair var after creating the client
	loginPair = authn.LoginPair{}

	return &client{conjur, roClient, config}, nil
}

// NewBind creates new binding service
func (c *client) NewBind(orgID, spaceID, bindingID string, enableSpaceIdentity bool) Bind {
	res := &bind{
		orgID:     orgID,
		spaceID:   spaceID,
		bindingID: bindingID,
		client:    c,
	}
	if enableSpaceIdentity {
		res.hostID = composeID(c.config.ConjurAccount, KindHost, res.policy())
	} else {
		res.hostID = slashJoin(composeID(c.config.ConjurAccount, KindHost, res.policy()), bindingID)
	}
	return res
}

// FromBindingID creates new binding service based on existing binding by its ID, org and space IDs would be queried from conjur
func (c *client) FromBindingID(bindingID string) (Bind, error) {
	orgID, spaceID, err := c.orgSpaceFromBindingID(bindingID)
	if err != nil {
		return nil, fmt.Errorf("failed to create binding service from binding id: %w", err)
	}
	return c.NewBind(orgID, spaceID, bindingID, false), nil // false is safe since this method is only used in context of disabled space identity
}

func (c *client) orgSpaceFromBindingID(bindingID string) (string, string, error) {
	res, err := c.roClient.Resources(&conjurapi.ResourceFilter{
		Kind:   KindHost.String(),
		Search: bindingID + "^",
	})
	if err != nil {
		return "", "", err
	}
	if len(res) == 0 {
		return "", "", nil
	}
	if len(res) > 1 {
		return "", "", fmt.Errorf("expecting exactly one host ending with a given binding id")
	}
	id, ok := res[0]["id"]
	if !ok {
		return "", "", nil
	}
	_, _, identifier := parseID(fmt.Sprintf("%s", id))
	split := strings.SplitN(identifier, "/", 4)
	if len(split) != 4 {
		return "", "", nil
	}
	// expected identifier contains org id on second position and space id on third
	// dev:host:cf/{orgID}/{spaceID}/{bindingID}
	return split[1], split[2], err
}

// NewProvision creates a Provision based on provided configuration
func (c *client) NewProvision(orgID, spaceID string, orgName, spaceName *string) Provision {
	res := &provision{
		orgID:   orgID,
		spaceID: spaceID,
		client:  c,
	}
	if orgName != nil {
		res.orgName = *orgName
	}
	if spaceName != nil {
		res.spaceName = *spaceName
	}
	return res
}

// ValidateConnectivity validates conjur client configuration by checking read access permission to the policy
func (c *client) ValidateConnectivity() error {
	res, err := c.client.ResourceExists(func() string {
		kind := KindUser
		if strings.HasPrefix(c.config.ConjurAuthNLogin, "host/") {
			kind = KindHost
		}
		return composeID(c.config.ConjurAccount, kind, c.clientHostID())
	}())
	if err != nil {
		return fmt.Errorf("failed checking login resource %v: %w", c.config.ConjurAuthNLogin, err)
	}
	if !res {
		return fmt.Errorf("validation failed, identity %v not privileged to read itself", c.config.ConjurAuthNLogin)
	}
	res, err = c.client.CheckPermission(c.basePolicyID(), PrivilegeRead.String())
	if err != nil {
		return fmt.Errorf("failed checking read permissions on policy %v: %w", c.basePolicyID(), err)
	}
	if !res {
		return fmt.Errorf("validation failed, missing read permissions on policy %v, verify your credentials are correct", c.basePolicyID())
	}
	res, err = c.roClient.CheckPermission(c.basePolicyID(), PrivilegeRead.String())
	if err != nil {
		return fmt.Errorf("failed checking follower read permissions on policy %v: %w", c.basePolicyID(), err)
	}
	if !res {
		return fmt.Errorf("validation failed follower, missing read permissions on policy %v, verify your credentials are correct", c.basePolicyID())
	}
	return nil
}

// clientHostID returns host ID of the client
func (c *client) clientHostID() string {
	return strings.TrimPrefix(c.config.ConjurAuthNLogin, "host/")
}

// upsertPolicy creates or updates (appends) a policy
func (c *client) upsertPolicy(policy io.Reader, policyID string) (*conjurapi.PolicyResponse, error) {
	res, err := c.client.LoadPolicy(
		conjurapi.PolicyModePost,
		policyID,
		policy,
	)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// replacePolicy completely replaces an existing policy, implicitly deleting data which is not present in the new policy
func (c *client) replacePolicy(policy io.Reader, policyID string) (*conjurapi.PolicyResponse, error) {
	res, err := c.client.LoadPolicy(
		conjurapi.PolicyModePut,
		policyID,
		policy,
	)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// resourceExists checks for an existence of a resource with a given id
func (c *client) resourceExists(resourceID string) (bool, error) {
	exists, err := c.roClient.ResourceExists(resourceID)
	if err != nil {
		return false, fmt.Errorf("unable to check resource existance %v: %w", resourceID, err)
	}
	return exists, nil
}

// roleExists checks for an existence of a role with a given id
func (c *client) roleExists(resourceID string) (bool, error) {
	exists, err := c.roClient.RoleExists(resourceID)
	if err != nil {
		return false, fmt.Errorf("unable to check role existance %v: %w", resourceID, err)
	}
	return exists, nil
}

// setVariable sets a secret variable
func (c *client) setVariable(variableID, secret string) error {
	return c.client.AddSecret(variableID, secret)
}

// getVariable gets a secret variable
func (c *client) getVariable(variableID string) (string, error) {
	bytes, err := c.roClient.RetrieveSecret(variableID)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// rotateAPIKey checks for an existence of a role with a given id
func (c *client) rotateAPIKey(roleID string) error {
	_, err := c.client.RotateAPIKey(roleID)
	if err != nil {
		return fmt.Errorf("unable to rotate API key for role %v: %w", roleID, err)
	}
	return nil
}

func (c *client) basePolicyID() string {
	return composeID(c.config.ConjurAccount, KindPolicy, c.config.ConjurPolicy)
}

// platformAnnotation checks for platform annotation on host used for service broker authentication
func (c *client) platformAnnotation() (string, error) {
	hostID := composeID(c.config.ConjurAccount, KindHost, c.clientHostID())
	clientHost, err := c.roClient.Resource(hostID)
	if err != nil {
		return "", fmt.Errorf("unable to find resource %v: %w", hostID, err)
	}
	resp := struct {
		Annotations []struct {
			Name  string `json:"name"`
			Value string `json:"value"`
		} `json:"annotations"`
	}{}
	bytes, err := json.Marshal(clientHost)
	if err != nil {
		return "", err
	}
	err = json.Unmarshal(bytes, &resp)
	if err != nil {
		return "", err
	}
	for _, a := range resp.Annotations {
		if a.Name == "platform" {
			return a.Value, nil
		}
	}
	return "", nil
}

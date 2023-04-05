package conjur

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/cyberark/conjur-api-go/conjurapi"
	"github.com/cyberark/conjur-api-go/conjurapi/authn"
	"github.com/cyberark/conjur-api-go/conjurapi/logging"
	"github.com/sirupsen/logrus"
)

// TODO: make this better
func init() {
	// this is just to enable debug in conjur SDK
	if debug := os.Getenv("DEBUG"); debug == "true" {
		logging.ApiLog.Level = logrus.DebugLevel
	}
}

// Client allows interactions with conjure instance
type Client interface {
	CheckPermission(resourceID string, privilege ...Privilege) (bool, error)
	CheckVariablePermission(variableID string, privilege ...VariablePrivilege) (bool, error)
	UpsertPolicy(policy io.Reader) (*conjurapi.PolicyResponse, error)
	ReplacePolicy(policy io.Reader) (*conjurapi.PolicyResponse, error)
	CheckResource(resourceID string) (bool, error)
	RoleExists(roleID string) (bool, error)
	RotateAPIKey(roleID string) error
	Platform() (string, error)

	Config() *Config
	ValidateConnectivity() error
}

type client struct {
	client   *conjurapi.Client
	roClient *conjurapi.Client
	config   *Config
}

// NewClient creates new conjur API client wrapper
func NewClient(config *Config) (Client, error) {
	clientConf, err := conjurapi.LoadConfig()
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
		roClient, err = conjurapi.NewClientFromKey(clientConf, loginPair)
		if err != nil {
			return nil, err
		}
	} else {
		roClient = conjur
	}
	if conjur == nil {
		return nil, fmt.Errorf("failed to create conjur conjur")
	}
	return &client{conjur, roClient, config}, nil
}

// ValidateConnectivity validates conjur client configuration by checking read access permission to the policy
func (c *client) ValidateConnectivity() error {
	_, err := c.client.CheckPermission(c.basePolicy(), PrivilegeRead.String())
	if err != nil {
		return fmt.Errorf("validation failed, missing read permissions on policy %v: %w", c.basePolicy(), err)
	}
	_, err = c.roClient.CheckPermission(c.basePolicy(), PrivilegeRead.String())
	if err != nil {
		return fmt.Errorf("validation failed, ro missing read permissions on policy %v: %w", c.basePolicy(), err)
	}
	return nil
}

// ClientHostID returns host ID of the client
func (c *client) ClientHostID() string {
	return strings.TrimPrefix(c.config.ConjurAuthNLogin, "host/")
}

func (c *client) checkPermission(resourceID string, privilege string) (bool, error) {
	ok, err := c.roClient.CheckPermission(resourceID, privilege)
	if err != nil {
		return false, fmt.Errorf("validation failed, missing read permissions on policy: %w", err)
	}
	if !ok {
		return false, nil
	}
	return true, nil
}

// CheckPermission checks permissions for a given resource id
func (c *client) CheckPermission(resourceID string, privilege ...Privilege) (bool, error) {
	// TODO: check if resource is not a variable
	for _, p := range privilege {
		ok, err := c.checkPermission(resourceID, p.String())
		if err != nil || !ok {
			return ok, err
		}
	}
	return true, nil
}

// CheckVariablePermission checks permissions for a given variable id
func (c *client) CheckVariablePermission(variableID string, privilege ...VariablePrivilege) (bool, error) {
	// TODO: check if variableID is actually a variable id
	for _, p := range privilege {
		ok, err := c.checkPermission(variableID, p.String())
		if err != nil || !ok {
			return ok, err
		}
	}
	return true, nil
}

// UpsertPolicy creates or updates (appends) a policy
func (c *client) UpsertPolicy(policy io.Reader) (*conjurapi.PolicyResponse, error) {
	res, err := c.client.LoadPolicy(
		conjurapi.PolicyModePost,
		c.config.ConjurPolicy,
		policy,
	)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// ReplacePolicy completely replaces an existing policy, implicitly deleting data which is not present in the new policy
func (c *client) ReplacePolicy(policy io.Reader) (*conjurapi.PolicyResponse, error) {
	res, err := c.client.LoadPolicy(
		conjurapi.PolicyModePut,
		c.config.ConjurPolicy,
		policy,
	)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// CheckResource checks for an existence of a resource with a given id
func (c *client) CheckResource(resourceID string) (bool, error) {
	spacePolicy, err := c.roClient.Resource(resourceID)
	if err != nil {
		return false, fmt.Errorf("unable to find resource %v: %w", resourceID, err)
	}
	if len(spacePolicy) == 0 {
		return false, nil
	}
	return true, nil
}

// RoleExists checks for an existence of a role with a given id
func (c *client) RoleExists(roleID string) (bool, error) {
	roleExists, err := c.roClient.RoleExists(roleID)
	if err != nil {
		return false, fmt.Errorf("unable to find role %v: %w", roleID, err)
	}
	return roleExists, nil
}

// RotateAPIKey checks for an existence of a role with a given id
func (c *client) RotateAPIKey(roleID string) error {
	_, err := c.client.RotateAPIKey(roleID)
	if err != nil {
		return fmt.Errorf("unable to rotate API key for role %v: %w", roleID, err)
	}
	return nil
}

func (c *client) basePolicy() string {
	return fmt.Sprintf("%v:policy:%v", c.config.ConjurAccount, c.config.ConjurPolicy)
}

// Platform checks for platform annotation on host used for service broker authentication
func (c *client) Platform() (string, error) {
	hostID := composeID(c.config.ConjurAccount, KindHost, c.ClientHostID())
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

// Config returns a conjur client config
func (c *client) Config() *Config {
	return c.config
}

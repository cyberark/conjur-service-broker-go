package conjur

import (
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

// PermissionChecker allows checking permissions in conjur
type PermissionChecker interface {
	CheckPermission(resourceID string, privilege ...Privilege) (bool, error)
	CheckVariablePermission(variableID string, privilege ...VariablePrivilege) (bool, error)
}

// PolicyLoader allows loading a conjur policy
type PolicyLoader interface {
	UpsertPolicy(policy io.Reader) error
}

// ResourceChecker allows checking existence of resource
type ResourceChecker interface {
	CheckResource(resourceID string) (bool, error)
	BasePolicy() string
	BaseLayer() string
}

// Client is a conjur API client wrapper that allows to manipulate on entities needed by service broker
type Client struct {
	client   *conjurapi.Client
	roClient *conjurapi.Client
	config   *Config
}

// NewClient creates new conjur API client wrapper
func NewClient(config *Config) (*Client, error) {
	clientConf, err := conjurapi.LoadConfig()
	if err != nil {
		panic(err)
	}
	loginPair := authn.LoginPair{
		Login:  config.ConjurAuthNLogin,
		APIKey: config.ConjurAuthNAPIKey,
	}
	client, err := conjurapi.NewClientFromKey(clientConf, loginPair)
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
		roClient = client
	}

	if client == nil {
		return nil, fmt.Errorf("failed to create conjur client")
	}
	res := Client{client, roClient, config}
	return &res, nil
}

// Validate validates conjur client configuration by checking read access permission to the policy
func (c *Client) Validate() error {
	_, err := c.client.CheckPermission(c.BasePolicy(), PrivilegeRead.String())
	if err != nil {
		return fmt.Errorf("validation failed, missing read permissions on policy %v: %w", c.BasePolicy(), err)
	}
	_, err = c.roClient.CheckPermission(c.BasePolicy(), PrivilegeRead.String())
	if err != nil {
		return fmt.Errorf("validation failed, ro missing read permissions on policy %v: %w", c.BasePolicy(), err)
	}
	return nil
}

// HostID returns host ID of the client
func (c *Client) HostID() string {
	if !strings.HasPrefix(c.config.ConjurAuthNLogin, "host/") {
		return "" // TODO: should this be an error?
	}
	return string([]rune(c.config.ConjurAuthNLogin)[5:])
}

func (c *Client) checkPermission(resourceID string, privilege string) (bool, error) {
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
func (c *Client) CheckPermission(resourceID string, privilege ...Privilege) (bool, error) {
	for _, p := range privilege {
		ok, err := c.checkPermission(resourceID, p.String())
		if err != nil || !ok {
			return ok, err
		}
	}
	return true, nil
}

// CheckVariablePermission checks permissions for a given variable id
func (c *Client) CheckVariablePermission(variableID string, privilege ...VariablePrivilege) (bool, error) {
	for _, p := range privilege {
		ok, err := c.checkPermission(variableID, p.String())
		if err != nil || !ok {
			return ok, err
		}
	}
	return true, nil
}

// UpsertPolicy creates or updates (appends) a policy
func (c *Client) UpsertPolicy(policy io.Reader) error {
	_, err := c.client.LoadPolicy(
		conjurapi.PolicyModePost,
		c.config.ConjurPolicy,
		policy,
	)
	if err != nil {
		return err
	}
	return nil
}

// CheckResource checks for an existence of a resource with a given id
func (c *Client) CheckResource(resourceID string) (bool, error) {
	spacePolicy, err := c.roClient.Resource(resourceID)
	if err != nil {
		return false, fmt.Errorf("unable to find resource %v: %w", resourceID, err)
	}
	if len(spacePolicy) == 0 {
		return false, fmt.Errorf("unable to find resource %v", resourceID)
	}
	return true, nil
}

// BasePolicy returns a base policy
func (c *Client) BasePolicy() string {
	return fmt.Sprintf("%v:policy:%v", c.config.ConjurAccount, c.config.ConjurPolicy)
}

// BaseLayer returns a base layer
func (c *Client) BaseLayer() string {
	return fmt.Sprintf("%v:layer:%v", c.config.ConjurAccount, c.config.ConjurPolicy)
}

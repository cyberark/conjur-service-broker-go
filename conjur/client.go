package conjur

import (
	"fmt"
	"os"
	"strings"

	"github.com/cyberark/conjur-api-go/conjurapi/logging"
	"github.com/sirupsen/logrus"

	"github.com/cyberark/conjur-api-go/conjurapi"
	"github.com/cyberark/conjur-api-go/conjurapi/authn"
)

const (
	readPermission = "read"
)

// TODO: make this better
func init() {
	// this is just to enable debug in conjur SDK
	if debug := os.Getenv("DEBUG"); debug == "true" {
		logging.ApiLog.Level = logrus.DebugLevel
	}
}

// Client is a conjur API client wrapper that allows to manipulate on entities needed by service broker
type Client struct {
	*conjurapi.Client
	*Config
	// TODO: improve encapsulation
}

// NewClient creates new conjur API client wrapper
func NewClient(config *Config) (*Client, error) {
	clientConf, err := conjurapi.LoadConfig()
	if err != nil {
		panic(err)
	}
	clientConf = mergeConfig(clientConf, config)
	loginPair := authn.LoginPair{
		Login:  config.ConjurAuthNLogin,
		APIKey: config.ConjurAuthNAPIKey,
	}
	client, err := conjurapi.NewClientFromKey(clientConf, loginPair)
	if err != nil {
		return nil, err
	}
	if client == nil {
		return nil, fmt.Errorf("failed to create conjur client")
	}
	res := Client{client, config}
	return &res, nil
}

func mergeConfig(clientConf conjurapi.Config, config *Config) conjurapi.Config {
	res := clientConf // this is a deep copy since all the fields are primitive types
	res.Account = config.ConjurAccount
	res.ApplianceURL = config.ConjurApplianceURL // TODO: handle follower URL?
	res.SSLCert = config.ConjurSSLCertificate
	return res
}

// Validate validates conjur client configuration by checking read access permission to the policy
func (c *Client) Validate() error {
	policyID := fmt.Sprintf("%s:policy:%s", c.ConjurAccount, c.ConjurPolicy)
	_, err := c.Client.CheckPermission(policyID, readPermission)
	if err != nil {
		return fmt.Errorf("validation failed, missing read permissions on policy: %w", err)
	}
	return nil
}

// HostID returns host ID of the client
func (c *Client) HostID() string {
	if !strings.HasPrefix(c.ConjurAuthNLogin, "host/") {
		return "" // TODO: should this be an error?
	}
	return string([]rune(c.ConjurAuthNLogin)[5:])
}

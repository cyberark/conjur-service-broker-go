// Package main provides integration tests for conjur service broker
package main

import (
	"encoding/json"
	"fmt"

	"github.com/cyberark/conjur-api-go/conjurapi"
	"github.com/cyberark/conjur-api-go/conjurapi/authn"
)

type conjur struct {
	api    *httpFeature
	client *conjurapi.Client
}

type creds struct {
	Credentials struct {
		Account        string `json:"account"`
		ApplianceURL   string `json:"appliance_url"`
		AuthnAPIKey    string `json:"authn_api_key"`
		AuthnLogin     string `json:"authn_login"`
		SslCertificate string `json:"ssl_certificate"`
		Version        int    `json:"version"`
	} `json:"credentials"`
}

func (c *conjur) iCreateConjurClient() error {
	var cr creds
	err := json.Unmarshal([]byte(c.api.body), &cr)
	if err != nil {
		return err
	}
	c.client, _ = conjurapi.NewClientFromKey(conjurapi.Config{
		Account:      cr.Credentials.Account,
		ApplianceURL: cr.Credentials.ApplianceURL,
	}, authn.LoginPair{
		Login:  cr.Credentials.AuthnLogin,
		APIKey: cr.Credentials.AuthnAPIKey,
	})
	_, err = c.client.WhoAmI()
	if err != nil {
		return err
	}
	return nil
}

func (c *conjur) conjurCredentialsAreInvalid() error {
	err := c.client.RefreshToken()
	return err
}

func (c *conjur) conjurCredentialsAreValid() error {
	err := c.client.RefreshToken()
	if err != nil {
		return fmt.Errorf("expecting an error")
	}
	return nil
}

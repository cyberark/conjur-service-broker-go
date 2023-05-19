// Package main provides integration tests for conjur service broker
package main

import (
	"github.com/cucumber/godog"
	"github.com/cyberark/conjur-api-go/conjurapi"
	"github.com/cyberark/conjur-api-go/conjurapi/authn"
)

type conjur struct {
	client *conjurapi.Client
}

func newConjur(cfg cfg) (*conjur, error) {
	client, err := conjurapi.NewClientFromKey(conjurapi.Config{
		Account:      cfg.ConjurAccount,
		ApplianceURL: cfg.ConjurApplianceURL,
	}, authn.LoginPair{
		Login:  cfg.ConjurUser,
		APIKey: cfg.ConjurAPIKey,
	})
	if err != nil {
		return nil, err
	}
	return &conjur{client}, nil
}

func (c *conjur) iLoadASecretIntoConjur() error {
	return godog.ErrPending
}

func (c *conjur) iPrivilegeTheAppHostToAccessASecretInConjur() error {
	return godog.ErrPending
}

func (c *conjur) iPrivilegeTheOrgLayerToAccessASecretInConjur() error {
	return godog.ErrPending
}

func (c *conjur) iPrivilegeTheSpaceLayerToAccessASecretInConjur() error {
	return godog.ErrPending
}

func (c *conjur) thePolicyForTheOrgAndSpaceExists() error {
	return godog.ErrPending
}

func (c *conjur) theSpaceHostAPIKeyVariableExists() error {
	return godog.ErrPending
}

func (c *conjur) theSpaceHostExists() error {
	return godog.ErrPending
}

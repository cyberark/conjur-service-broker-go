// Package main provides e2e tests for conjur service broker
package main

import (
	"github.com/cloudfoundry-community/go-cfclient/v3/client"
	"github.com/cloudfoundry-community/go-cfclient/v3/config"
	"github.com/cucumber/godog"
)

type cf struct {
	client *client.Client
}

func newCF(cfg cfg) (*cf, error) {
	cfConfig, err := config.NewUserPassword(cfg.CFURL, cfg.CFUser, cfg.CFPassword)
	if err != nil {
		return nil, err
	}
	cfClient, err := client.New(cfConfig)
	if err != nil {
		return nil, err
	}
	return &cf{cfClient}, nil
}

func (c *cf) ICreateOrgSpace() error {
	// c.client.Organizations.Create()
	return nil
}

func (c *cf) iCreateAServiceInstanceForConjur() error {
	return godog.ErrPending
}

func (c *cf) iInstallTheConjurServiceBroker() error {
	return godog.ErrPending
}

func (c *cf) iPushTheSampleAppToPCF() error {
	return godog.ErrPending
}

func (c *cf) iStartTheApp() error {
	return godog.ErrPending
}

func (c *cf) iCanRetrieveTheSecretValuesFromTheApp() error {
	return godog.ErrPending
}

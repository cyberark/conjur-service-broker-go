package main

import (
	"context"
	"testing"

	"github.com/caarlos0/env/v7"
	"github.com/cucumber/godog"
)

func InitializeScenario(ctx *godog.ScenarioContext) {
	var cfg cfg
	err := env.Parse(&cfg)
	if err != nil {
		panic(err)
	}
	cf, err := newCF(&cfg)
	if err != nil {
		panic(err)
	}
	conjur, err := newConjur(&cfg)
	if err != nil {
		panic(err)
	}
	ctx.Before(func(ctx context.Context, sc *godog.Scenario) (context.Context, error) {
		state := &state{}
		return context.WithValue(ctx, stateKey{}, state), nil
	})
	ctx.Step(`^I create an org and space$`, cf.iCreateOrgSpace)
	ctx.Step(`^I create space developer user and login$`, cf.iCreateSpaceDeveloperUserAndLogin)
	ctx.Step(`^I can retrieve the secret values from the app$`, cf.iCanRetrieveTheSecretValuesFromTheApp)
	ctx.Step(`^I create a service instance for Conjur$`, cf.iCreateAServiceInstanceForConjur)
	ctx.Step(`^I install the Conjur service broker$`, cf.iInstallTheConjurServiceBroker)
	ctx.Step(`^I install the Conjur service broker with space host identity$`, cf.iInstallTheConjurServiceBrokerWithSpaceHostIdentity)
	ctx.Step(`^I load the secrets into Conjur$`, conjur.iLoadTheSecretsIntoConjur)
	ctx.Step(`^I privilege the app host to access a secret in Conjur$`, conjur.iPrivilegeTheAppHostToAccessASecretInConjur)
	ctx.Step(`^I privilege the org group to access a secret in Conjur$`, conjur.iPrivilegeTheOrgGroupToAccessASecretInConjur)
	ctx.Step(`^I privilege the space group to access a secret in Conjur$`, conjur.iPrivilegeTheSpaceGroupToAccessASecretInConjur)
	ctx.Step(`^I push the sample app to PCF$`, cf.iPushTheSampleAppToPCF)
	ctx.Step(`^I start the app$`, cf.iStartTheApp)
	ctx.Step(`^I remove the service instance$`, cf.iRemoveTheServiceInstance)
	ctx.Step(`^the policies for the org and space exists$`, conjur.thePoliciesForTheOrgAndSpaceExists)
	ctx.Step(`^the space host api key variable exists$`, conjur.theSpaceHostAPIKeyVariableExists)
	ctx.Step(`^the space host exists$`, conjur.theSpaceHostExists)
	ctx.Step(`^the binding host exists$`, conjur.theBindingHostExists)
}

func TestIntegration(t *testing.T) {
	suite := godog.TestSuite{
		ScenarioInitializer: InitializeScenario,
		Options: &godog.Options{
			Strict:   true,
			Tags:     "e2e",
			Format:   "pretty,junit:reports/junit.xml",
			Paths:    []string{"features"},
			TestingT: t, // Testing instance that will run subtests.
		},
	}

	if suite.Run() != 0 {
		t.Fatal("non-zero status returned, failed to run feature tests")
	}
}

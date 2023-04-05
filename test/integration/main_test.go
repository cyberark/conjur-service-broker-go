package main

import (
	"context"
	"testing"

	"github.com/caarlos0/env/v7"
	"github.com/cucumber/godog"
)

func InitializeScenario(ctx *godog.ScenarioContext) {
	var cfg config
	err := env.Parse(&cfg)
	if err != nil {
		panic(err)
	}
	api := &httpFeature{cfg: cfg}

	ctx.Before(func(ctx context.Context, sc *godog.Scenario) (context.Context, error) {
		api.resetResponse(sc)
		return ctx, nil
	})
	ctx.Step(`^my basic auth credentials are incorrect$`, api.myBasicAuthCredentialsAreIncorrect)
	ctx.Step(`^my request doesn\'t include the X-Broker-API-Version header$`, api.myRequestDoesntIncludeTheXBrokerAPIVersionHeader)
	ctx.Step(`^I send "(GET|POST|PUT|DELETE)" request to "([^"]*)"$`, api.iSendrequestTo)
	ctx.Step(`^I send "([^"]*)" request to "([^"]*)" with body:$`, api.iSendRequestToWithBody)
	ctx.Step(`^the response code should be (\d+)$`, api.theResponseCodeShouldBe)
	ctx.Step(`^the response should match json:$`, api.theResponseShouldMatchJSONBody)
	ctx.Step(`^the response should match json "([^"]*)"$`, api.theResponseShouldMatchJSON)
}

func TestIntegration(t *testing.T) {
	suite := godog.TestSuite{
		ScenarioInitializer: InitializeScenario,
		Options: &godog.Options{
			Format:   "pretty",
			Paths:    []string{"features"},
			TestingT: t, // Testing instance that will run subtests.
		},
	}

	if suite.Run() != 0 {
		t.Fatal("non-zero status returned, failed to run feature tests")
	}
}

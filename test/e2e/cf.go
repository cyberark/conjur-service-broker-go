// Package main provides e2e tests for conjur service broker
package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type cf struct {
	cfg *cfg
}

var bindingIDRegexp = regexp.MustCompile(`"binding_guid"\s*:\s*"([\da-f]{8}-[\da-f]{4}-[\da-f]{4}-[\da-f]{4}-[\da-f]{12})"`)

func newCF(cfg *cfg) (*cf, error) {
	if _, err := cfCLI("api", cfg.CFURL); err != nil {
		return nil, err
	}
	return &cf{cfg: cfg}, nil
}

func (c *cf) iCreateOrgSpace(ctx context.Context) (context.Context, error) {
	t := &T{}
	err := authCF(c.cfg.CFUser, c.cfg.CFPassword)
	assert.NoError(t, err)
	state := ctx.Value(stateKey{}).(*state)
	state.orgName = "ci-org-" + randomName()
	state.spaceName = "ci-space-" + randomName()
	_, err = cfCLI("create-org", state.orgName)
	assert.NoError(t, err)
	state.orgID, err = cfCLI("org", state.orgName, "--guid")
	assert.NoError(t, err)
	_, err = cfCLI("target", "-o", state.orgName)
	assert.NoError(t, err)
	_, err = cfCLI("create-space", state.spaceName)
	assert.NoError(t, err)
	state.spaceID, err = cfCLI("space", state.spaceName, "--guid")
	assert.NoError(t, err)
	_, err = cfCLI("target", "-o", state.orgName, "-s", state.spaceName)
	assert.NoError(t, err)
	return context.WithValue(ctx, stateKey{}, state), t.Error()
}

func (c *cf) iCreateSpaceDeveloperUserAndLogin(ctx context.Context) (context.Context, error) {
	t := &T{}
	user, password := "ci-user-"+randomName(), randomName()
	_, err := cfCLI("create-user", user, password)
	assert.NoError(t, err)
	state := ctx.Value(stateKey{}).(*state)
	_, err = cfCLI("set-space-role", user, state.orgName, state.spaceName, "SpaceDeveloper")
	assert.NoError(t, err)
	err = authCF(user, password)
	assert.NoError(t, err)
	_, err = cfCLI("target", "-o", state.orgName, "-s", state.spaceName)
	assert.NoError(t, err)
	return ctx, t.Error()
}

func (c *cf) iCreateAServiceInstanceForConjur(ctx context.Context) (context.Context, error) {
	state := ctx.Value(stateKey{}).(*state)
	if _, err := cfCLI("create-service", "--wait", "-b", state.sbName, "cyberark-conjur", "community", "conjur"); err != nil {
		return ctx, err
	}
	return ctx, nil
}

func (c *cf) iInstallTheConjurServiceBrokerWithSpaceHostIdentity(ctx context.Context) (context.Context, error) {
	return c.installServiceBroker(ctx, true)
}

func (c *cf) iInstallTheConjurServiceBroker(ctx context.Context) (context.Context, error) {
	return c.installServiceBroker(ctx, false)
}

func (c *cf) installServiceBroker(ctx context.Context, enableSpaceHostIdentity bool) (context.Context, error) {
	t := &T{}
	state := ctx.Value(stateKey{}).(*state)
	state.enableSpaceHostIdentity = enableSpaceHostIdentity
	_, err := cfCLI("push", "conjur-service-broker", "--no-start", "--path", "./cyberark-conjur-service-broker.zip")
	require.NoError(t, err)
	apiKey, err := conjurAPIKey(c.cfg)
	require.NoError(t, err)
	for variable, value := range c.envs(enableSpaceHostIdentity, apiKey) {
		_, err := cfCLI("set-env", "conjur-service-broker", variable, value)
		assert.NoError(t, err)
	}
	start, err := cfCLI("restage", "conjur-service-broker")
	require.NoError(t, err)
	route, err := parseRoute(start)
	require.NoError(t, err)
	state.sbName = "cyberark-conjur-" + randomName()
	// we don't need to enable-service-access since the service broker is --space-scoped
	_, err = cfCLI("create-service-broker", "--space-scoped", state.sbName, c.cfg.ServiceBrokerUser, c.cfg.ServiceBrokerPassword, route)
	assert.NoError(t, err)
	return context.WithValue(ctx, stateKey{}, state), t.Error()
}

func (c *cf) envs(enableSpaceHostIdentity bool, apiKey string) map[string]string {
	return map[string]string{
		"SECURITY_USER_NAME":     c.cfg.ServiceBrokerUser,
		"SECURITY_USER_PASSWORD": c.cfg.ServiceBrokerPassword,
		"CONJUR_ACCOUNT":         c.cfg.ConjurAccount,
		"CONJUR_APPLIANCE_URL":   c.cfg.ConjurApplianceURL,
		"CONJUR_SSL_CERTIFICATE": sslCert(c.cfg.ConjurApplianceURL),
		"CONJUR_POLICY":          c.cfg.ConjurPolicy,
		"CONJUR_AUTHN_LOGIN":     c.cfg.ConjurServiceBrokerUser,
		"CONJUR_AUTHN_API_KEY":   apiKey,
		"ENABLE_SPACE_IDENTITY":  strconv.FormatBool(enableSpaceHostIdentity),
		"DEBUG":                  "true",
	}
}

func (c *cf) iPushTheSampleAppToPCF(ctx context.Context) (context.Context, error) {
	t := &T{}
	_, err := cfCLI("delete", "sample-app", "-f", "-r")
	assert.NoError(t, err)
	secrets := "ORG_SECRET: !var app/secrets/org\nSPACE_SECRET: !var app/secrets/space\n" // #nosec G101 false positive
	state := ctx.Value(stateKey{}).(*state)
	if !state.enableSpaceHostIdentity {
		secrets += "APP_SECRET: !var app/secrets/app\n" // #nosec G101 false positive
	}
	err = os.WriteFile("./sample-app/secrets.yml", []byte(secrets), 0600)
	assert.NoError(t, err)
	push, err := cfCLI("push", "--no-start", "--path", "./sample-app", "--manifest", "./sample-app/manifest.yml")
	assert.NoError(t, err)
	state.appURL, err = parseRoute(push)
	assert.NoError(t, err)
	env, err := cfCLI("env", "sample-app")
	assert.NoError(t, err)
	state.bindID, err = parseBindingID(env)
	assert.NoError(t, err)
	return context.WithValue(ctx, stateKey{}, state), t.Error()
}

func (c *cf) iRemoveTheServiceInstance(ctx context.Context) (context.Context, error) {
	t := &T{}
	_, err := cfCLI("unbind-service", "sample-app", "conjur")
	assert.NoError(t, err)
	_, err = cfCLI("delete-service", "conjur", "-f")
	assert.NoError(t, err)
	return ctx, nil
}

func (c *cf) iStartTheApp(ctx context.Context) (context.Context, error) {
	t := &T{}
	_, err := cfCLI("start", "sample-app")
	assert.NoError(t, err)
	return ctx, t.Error()
}

func (c *cf) iCanRetrieveTheSecretValuesFromTheApp(ctx context.Context) (context.Context, error) {
	t := &T{}
	state := ctx.Value(stateKey{}).(*state)
	response, err := http.Get(state.appURL)
	assert.NoError(t, err)
	defer func() {
		_ = response.Body.Close()
	}()
	if response.StatusCode != http.StatusOK {
		t.Errorf("unexpected http status from service: %s", response.Status)
	}
	bytes, err := io.ReadAll(response.Body)
	assert.NoError(t, err)

	body := string(bytes)
	if len(state.secrets) == 0 {
		t.Errorf("no secrets to check")
	}
	for name, value := range state.secrets {
		if !strings.Contains(body, name) || !strings.Contains(body, value) {
			t.Errorf("response from service doesn't contain expected secrets %s=%s\n%s", name, value, body)
		}
	}
	return ctx, t.Error()
}

func cfCLI(args ...string) (string, error) {
	cmd := exec.Command("cf", args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("%w: %s", err, string(output))
	}
	return strings.TrimSpace(string(output)), err
}

func parseRoute(res string) (string, error) {
	scanner := bufio.NewScanner(strings.NewReader(res))
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "routes:") {
			route := strings.TrimPrefix(line, "routes:")
			return "http://" + strings.TrimSpace(route), nil
		}
	}
	return "", fmt.Errorf("routes not found in: %s", res)
}

func parseBindingID(env string) (string, error) {
	match := bindingIDRegexp.FindStringSubmatch(env)
	if len(match) != 2 {
		return "", fmt.Errorf("failed to find binding id in: %s", env)
	}
	return match[1], nil
}

func authCF(user, password string) error {
	cmd := exec.Command("cf", "auth", user)
	cmd.Env = append(cmd.Environ(), "CF_PASSWORD="+password)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%w: %s", err, string(output))
	}
	return nil
}

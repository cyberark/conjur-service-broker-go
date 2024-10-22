// Package main provides integration tests for conjur service broker
package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/cyberark/conjur-api-go/conjurapi"
	"github.com/cyberark/conjur-api-go/conjurapi/authn"
	"github.com/doodlesbykumbi/conjur-policy-go/pkg/conjurpolicy"
	"gopkg.in/yaml.v3"
)

type conjur struct {
	client *conjurapi.Client
	cfg    *cfg
}

func newConjur(cfg *cfg) (*conjur, error) {
	client, err := conjurapi.NewClientFromKey(conjurapi.Config{
		Account:      cfg.ConjurAccount,
		ApplianceURL: cfg.ConjurApplianceURL,
		SSLCert:      sslCert(cfg.ConjurApplianceURL),
	}, authn.LoginPair{
		Login:  cfg.ConjurUser,
		APIKey: cfg.ConjurAPIKey,
	})
	if err != nil {
		return nil, err
	}
	return &conjur{client, cfg}, nil
}

func (c *conjur) iLoadTheSecretsIntoConjur(ctx context.Context) (context.Context, error) {
	_, err := c.client.LoadPolicy(conjurapi.PolicyModePost, "app", secretsPolicy())
	if err != nil {
		return ctx, err
	}
	state := ctx.Value(stateKey{}).(*state)
	state.secrets = make(map[string]string)
	secrets := []string{"org", "space"}
	if !state.enableSpaceHostIdentity {
		secrets = append(secrets, "app")
	}
	for _, s := range secrets {
		state.secrets[s] = s + "-" + randomName()
		if err := c.client.AddSecret(c.cfg.ConjurAccount+":variable:app/secrets/"+s, state.secrets[s]); err != nil {
			return ctx, err
		}
	}
	return context.WithValue(ctx, stateKey{}, state), nil
}

func secretsPolicy() io.Reader {
	buf := new(bytes.Buffer)
	encoder := yaml.NewEncoder(buf)
	err := encoder.Encode(conjurpolicy.PolicyStatements{
		conjurpolicy.Variable{
			Id:          "secrets/org",
			Annotations: map[string]string{"description": "Org-wide secret"},
		},
		conjurpolicy.Variable{
			Id:          "secrets/space",
			Annotations: map[string]string{"description": "Space-wide secret"},
		},
		conjurpolicy.Variable{
			Id:          "secrets/app",
			Annotations: map[string]string{"description": "App-specific secret"},
		},
	})
	if err != nil {
		panic(err)
	}
	return buf
}

func (c *conjur) iPrivilegeTheAppHostToAccessASecretInConjur(ctx context.Context) (context.Context, error) {
	state := ctx.Value(stateKey{}).(*state)
	return ctx, c.iPrivilegeRoleToSecretInConjur(
		conjurpolicy.HostRef(slashJoin(c.cfg.ConjurPolicy, state.orgID, state.spaceID, state.bindID)),
		"app/secrets/app")
}

func (c *conjur) iPrivilegeTheOrgGroupToAccessASecretInConjur(ctx context.Context) (context.Context, error) {
	state := ctx.Value(stateKey{}).(*state)
	return ctx, c.iPrivilegeRoleToSecretInConjur(
		conjurpolicy.GroupRef(slashJoin(c.cfg.ConjurPolicy, state.orgID)),
		"app/secrets/org")
}

func (c *conjur) iPrivilegeTheSpaceGroupToAccessASecretInConjur(ctx context.Context) (context.Context, error) {
	state := ctx.Value(stateKey{}).(*state)
	return ctx, c.iPrivilegeRoleToSecretInConjur(
		conjurpolicy.GroupRef(slashJoin(c.cfg.ConjurPolicy, state.orgID, state.spaceID)),
		"app/secrets/space")
}

func (c *conjur) iPrivilegeRoleToSecretInConjur(role conjurpolicy.ResourceRef, secret string) error {
	policy := conjurpolicy.PolicyStatements{conjurpolicy.Permit{
		Resources:  conjurpolicy.VariableRef(secret),
		Role:       role,
		Privileges: []conjurpolicy.Privilege{conjurpolicy.PrivilegeRead, conjurpolicy.PrivilegeExecute},
	}}
	buffer := bytes.NewBuffer(nil)
	err := yaml.NewEncoder(buffer).Encode(policy)
	if err != nil {
		return err
	}
	_, err = c.client.LoadPolicy(conjurapi.PolicyModePost, "root", buffer)
	if err != nil {
		return err
	}
	return nil
}

func (c *conjur) thePoliciesForTheOrgAndSpaceExists(ctx context.Context) (context.Context, error) {
	state := ctx.Value(stateKey{}).(*state)
	orgPolicyID := slashJoin(c.cfg.ConjurAccount+":policy:"+c.cfg.ConjurPolicy, state.orgID)
	if err := c.theResourceExists(orgPolicyID); err != nil {
		return ctx, err
	}
	orgSpacePolicyID := slashJoin(c.cfg.ConjurAccount+":policy:"+c.cfg.ConjurPolicy, state.orgID, state.spaceID)
	return ctx, c.theResourceExists(orgSpacePolicyID)
}

func (c *conjur) theResourceExists(resourceID string) error {
	exists, err := c.client.ResourceExists(resourceID)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("the resource: %s doesn't exist", resourceID)
	}
	return nil
}

func (c *conjur) theSpaceHostAPIKeyVariableExists(ctx context.Context) (context.Context, error) {
	state := ctx.Value(stateKey{}).(*state)
	spaceHostAPIKeyVariableID := slashJoin(c.cfg.ConjurAccount+":variable:"+c.cfg.ConjurPolicy, state.orgID, state.spaceID, "space-host-api-key")
	return ctx, c.theResourceExists(spaceHostAPIKeyVariableID)
}

func (c *conjur) theSpaceHostExists(ctx context.Context) (context.Context, error) {
	state := ctx.Value(stateKey{}).(*state)
	spaceHostID := slashJoin(c.cfg.ConjurAccount+":host:"+c.cfg.ConjurPolicy, state.orgID, state.spaceID)
	return ctx, c.theResourceExists(spaceHostID)
}

func (c *conjur) theBindingHostExists(ctx context.Context) (context.Context, error) {
	state := ctx.Value(stateKey{}).(*state)
	spaceHostID := slashJoin(c.cfg.ConjurAccount+":host:"+c.cfg.ConjurPolicy, state.orgID, state.spaceID, state.bindID)
	return ctx, c.theResourceExists(spaceHostID)
}

func conjurAPIKey(cfg *cfg) (string, error) {
	c, err := newConjur(cfg)
	if err != nil {
		return "", err
	}
	roleID := strings.Join([]string{cfg.ConjurAccount, "host", strings.TrimPrefix(cfg.ConjurServiceBrokerUser, "host/")}, ":")
	apiKey, err := c.client.RotateAPIKey(roleID)
	if err != nil {
		return "", err
	}
	return string(apiKey), nil
}

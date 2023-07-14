// Package main provides e2e tests for conjur service broker
package main

type state struct {
	orgName   string
	orgID     string
	spaceName string
	spaceID   string
	bindID    string

	secrets map[string]string

	enableSpaceHostIdentity bool

	appURL string
	sbName string
}

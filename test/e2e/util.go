// Package main provides e2e tests for conjur service broker
package main

import (
	"crypto/tls"
	"encoding/pem"
	"errors"
	"fmt"
	"math/rand"
	"net/url"
	"strings"

	"github.com/lucasepe/codename"
)

var rng = must(codename.DefaultRNG())

type stateKey struct{}

func must(rand *rand.Rand, err error) *rand.Rand {
	if err != nil {
		panic(err)
	}
	return rand
}

func randomName() string {
	return codename.Generate(rng, 0)
}

func slashJoin(s ...string) string {
	return strings.Join(s, "/")
}

func sslCert(conjurURL string) string {
	uri, err := url.ParseRequestURI(conjurURL)
	if err != nil {
		return ""
	}
	if uri.Scheme != "https" {
		return ""
	}
	conn, err := tls.Dial("tcp", uri.Host+":443", &tls.Config{
		InsecureSkipVerify: true,
	}) // #nosec G402 this needs to get the self-signed certificate to inject it to conjur client configuration
	if err != nil {
		return ""
	}
	res := make([]string, len(conn.ConnectionState().PeerCertificates))
	for i, c := range conn.ConnectionState().PeerCertificates {
		res[i] = strings.TrimSpace(string(pem.EncodeToMemory(&pem.Block{
			Type:  "CERTIFICATE",
			Bytes: c.Raw,
		})))
	}
	return strings.Join(res, "\n")
}

// T is a util helper that implements behaviour of testing.T for easier use with assert/require from testify
type T struct {
	errs []error
}

// FailNow is used to break the test execution
func (t *T) FailNow() {
	panic(t.Error())
}

// Errorf is used to indicate error while testing
func (t *T) Errorf(format string, args ...interface{}) {
	t.errs = append(t.errs, fmt.Errorf(format, args...))
}

func (t *T) Error() error {
	return errors.Join(t.errs...)
}

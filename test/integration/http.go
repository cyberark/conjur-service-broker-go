// Package main provides integration tests for conjur service broker
package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/cucumber/godog"
	"github.com/nsf/jsondiff"
)

func init() {
	http.DefaultClient.Timeout = 5 * time.Minute
}

type httpFeature struct {
	cfg                        config
	resp                       *http.Response
	body                       string
	authnPass                  string
	omitBrokerAPIVersionHeader bool
}

func (a *httpFeature) resetResponse(_ *godog.Scenario) {
	a.resp = nil
	a.body = ""
	a.omitBrokerAPIVersionHeader = false
	a.authnPass = a.cfg.BasicAuthPassword
}

func (a *httpFeature) myBasicAuthCredentialsAreIncorrect() error {
	a.authnPass = ""
	return nil
}

func (a *httpFeature) myRequestDoesntIncludeTheXBrokerAPIVersionHeader() error {
	a.omitBrokerAPIVersionHeader = true
	return nil
}

func (a *httpFeature) iSendRequestToWithBody(method, endpoint string, bodyDocString *godog.DocString) (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	var body io.Reader
	if bodyDocString != nil {
		body = strings.NewReader(bodyDocString.Content)
	}
	req, err := http.NewRequestWithContext(ctx, method, a.cfg.ServiceURL+endpoint, body)
	if err != nil {
		return
	}

	req.Header.Set("Content-Type", "application/json")
	if !a.omitBrokerAPIVersionHeader {
		req.Header.Set("X-Broker-API-Version", "2.17")
	}

	if len(a.cfg.BasicAuthPassword) > 0 {
		req.SetBasicAuth(a.authnPass, a.cfg.BasicAuthPassword)
	}

	// handle panic
	defer func() {
		switch t := recover().(type) {
		case string:
			err = fmt.Errorf(t)
		case error:
			err = t
		}
	}()

	a.resp, err = http.DefaultClient.Do(req)
	if err != nil {
		return
	}

	bytes, err := io.ReadAll(a.resp.Body)
	a.body = string(bytes)
	if err != nil {
		_ = a.resp.Body.Close()
		return
	}
	_ = a.resp.Body.Close()
	return nil
}

func (a *httpFeature) iSendRequestTo(method, endpoint string) (err error) {
	return a.iSendRequestToWithBody(method, endpoint, nil)
}

func (a *httpFeature) theResponseCodeShouldBe(code int) error {
	if code != a.resp.StatusCode {
		return fmt.Errorf("expected response code to be: %d, but actual is: %d %s\n\n%v", code, a.resp.StatusCode, a.resp.Status, a.body)
	}
	return nil
}

func (a *httpFeature) theResponseShouldMatchJSONBody(body *godog.DocString) (err error) {
	return a.theResponseShouldMatchJSON(body.Content)
}

func (a *httpFeature) theResponseShouldMatchJSON(body string) (err error) {
	if len(a.body) == 0 {
		return fmt.Errorf("server response is empty")
	}
	diffOpts := jsondiff.DefaultConsoleOptions()
	res, diff := jsondiff.Compare([]byte(body), []byte(a.body), &diffOpts)

	if res != jsondiff.FullMatch {
		return fmt.Errorf("expected JSON does not match actual: %s\n%v", diff, a.body)
	}
	return nil
}

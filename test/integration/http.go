package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/cucumber/godog"
	"github.com/nsf/jsondiff"
)

func init() {
	http.DefaultClient.Timeout = 5 * time.Minute
}

type httpFeature struct {
	cfg  config
	resp *http.Response
	body string
}

func (a *httpFeature) resetResponse(*godog.Scenario) {
	a.resp = nil
	a.body = ""
}

func (a *httpFeature) iSendrequestTo(method, endpoint string) (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, method, a.cfg.ServiceURL+endpoint, nil)
	if err != nil {
		return
	}

	req.Header.Set("X-Broker-API-Version", "2.17")

	if len(a.cfg.BasicAuthPassword) > 0 {
		req.SetBasicAuth(a.cfg.BasicAuthUser, a.cfg.BasicAuthPassword)
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
	defer a.resp.Body.Close()
	bytes, err := io.ReadAll(a.resp.Body)
	a.body = string(bytes)
	if err != nil {
		return
	}

	return nil
}

func (a *httpFeature) theResponseCodeShouldBe(code int) error {
	if code != a.resp.StatusCode {
		return fmt.Errorf("expected response code to be: %d, but actual is: %d %s\n\n%v\n", code, a.resp.StatusCode, a.resp.Status, a.body)
	}
	return nil
}

func (a *httpFeature) theResponseShouldMatchJSON(body *godog.DocString) (err error) {
	diffOpts := jsondiff.DefaultConsoleOptions()
	res, diff := jsondiff.Compare([]byte(body.Content), []byte(a.body), &diffOpts)

	if res != jsondiff.FullMatch {
		return fmt.Errorf("expected JSON does not match actual: %s", diff)
	}
	return nil
}

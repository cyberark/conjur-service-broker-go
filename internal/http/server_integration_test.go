//go:build integration

package http

import (
	"errors"
	"net/http"
	"os"
	"strings"
	"syscall"
	"testing"
	"time"

	"github.com/cyberark/conjur-service-broker-go/internal/ctxutil"
	"github.com/cyberark/conjur-service-broker-go/internal/servicebroker/mocks"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func Test_startServer_gracefulShutdown(t *testing.T) {
	srv := mocks.ServerInterface{}
	srv.On("CatalogGet", mock.Anything, mock.Anything).Return()
	var exited bool
	go func() {
		logger, err := zap.NewDevelopment()
		require.NoError(t, err)
		err = startServer(ctxutil.NewContext(), &config{Port: "18080"}, &srv, logger) // port 0 should make go choose a random port
		exited = true
		require.NoError(t, err)
	}()
	require.Eventually(t, func() bool { return getCatalog(t) }, time.Second, 50*time.Millisecond)
	err := syscall.Kill(syscall.Getpid(), syscall.SIGINT)
	require.NoError(t, err)
	require.Eventually(t, func() bool { return exited }, 5*time.Second, 50*time.Millisecond)
	srv.AssertExpectations(t)
}

func getCatalog(t *testing.T) bool {
	req, err := http.NewRequest(http.MethodGet, "http://localhost:18080/v2/catalog", nil)
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Broker-API-Version", "2.17")
	client := &http.Client{
		Transport: &http.Transport{},
	}
	response, err := client.Do(req)
	if errors.Is(err, syscall.ECONNREFUSED) {
		return false // the server is not yet up and running
	}
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, response.StatusCode)
	return true
}

func TestStartHTTPServer(t *testing.T) {
	t.Cleanup(cleanupEnv())
	go StartHTTPServer()

	type J map[string]any
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	httpmock.RegisterResponder("POST", "=~https://(conjur|follower).local/authn/dev/host%2Fservice-broker/authenticate",
		httpmock.NewJsonResponderOrPanic(http.StatusOK, J{"protected": "eyJhbGciOiJjb25qdXIub3JnL3Nsb3NpbG8vdjIiLCJraWQiOiI5M2VjNTEwODRmZTM3Zjc3M2I1ODhlNTYyYWVjZGMxMSJ9", "payload": "eyJzdWIiOiJhZG1pbiIsImlhdCI6MTUxMDc1MzI1OX0=", "signature": "raCufKOf7sKzciZInQTphu1mBbLhAdIJM72ChLB4m5wKWxFnNz_7LawQ9iYEI_we1-tdZtTXoopn_T1qoTplR9_Bo3KkpI5Hj3DB7SmBpR3CSRTnnEwkJ0_aJ8bql5Cbst4i4rSftyEmUqX-FDOqJdAztdi9BUJyLfbeKTW9OGg-QJQzPX1ucB7IpvTFCEjMoO8KUxZpbHj-KpwqAMZRooG4ULBkxp5nSfs-LN27JupU58oRgIfaWASaDmA98O2x6o88MFpxK_M0FeFGuDKewNGrRc8lCOtTQ9cULA080M5CSnruCqu1Qd52r72KIOAfyzNIiBCLTkblz2fZyEkdSKQmZ8J3AakxQE2jyHmMT-eXjfsEIzEt-IRPJIirI3Qm"}))
	httpmock.RegisterResponder("GET", "https://conjur.local/resources/dev/host/service-broker",
		httpmock.NewJsonResponderOrPanic(http.StatusOK, J{}))
	httpmock.RegisterResponder("GET", "=~https://follower.local/resources/dev/(policy|group)/cf.*",
		httpmock.NewJsonResponderOrPanic(http.StatusOK, J{}))
	httpmock.RegisterResponder("GET", "https://conjur.local/resources/dev/policy/cf",
		httpmock.NewJsonResponderOrPanic(http.StatusOK, J{}))
	httpmock.RegisterResponder("POST", "https://conjur.local/policies/dev/policy/cf",
		httpmock.NewJsonResponderOrPanic(http.StatusOK, J{}))

	require.Eventually(t, func() bool { return getCatalog(t) }, 5*time.Second, 50*time.Millisecond)

}

func cleanupEnv() func() {
	values := os.Environ()
	os.Clearenv()
	for k, v := range map[string]string{"DEBUG": "true", "CONJUR_ACCOUNT": "dev", "CONJUR_APPLIANCE_URL": "https://conjur.local", "CONJUR_FOLLOWER_URL": "https://follower.local", "CONJUR_AUTHN_API_KEY": "api-key", "CONJUR_AUTHN_LOGIN": "host/service-broker", "CONJUR_POLICY": "cf", "PORT": "18080"} {
		_ = os.Setenv(k, v)
	}
	return func() {
		os.Clearenv()
		for _, v := range values {
			parts := strings.SplitN(v, "=", 2)
			_ = os.Setenv(parts[0], parts[1])
		}
	}
}

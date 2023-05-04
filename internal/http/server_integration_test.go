//go:build integration

package http

import (
	"errors"
	"net/http"
	"syscall"
	"testing"
	"time"

	"github.com/cyberark/conjur-service-broker/internal/ctxutil"
	"github.com/cyberark/conjur-service-broker/internal/servicebroker/mocks"
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
	response, err := http.DefaultClient.Do(req)
	if errors.Is(err, syscall.ECONNREFUSED) {
		return false // the server is not yet up and running
	}
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, response.StatusCode)
	return true
}

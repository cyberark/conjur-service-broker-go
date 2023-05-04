package servicebroker

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_server_ServiceInstanceDeprovision(t *testing.T) {
	s := &server{}
	w, c := ginTestCtx(t, http.MethodDelete, "/v2/service_instances/{instance_id}", "", false)
	s.ServiceInstanceDeprovision(c, "", ServiceInstanceDeprovisionParams{})
	require.Empty(t, c.Errors.Errors())
	require.Equal(t, "{}", w.Body.String())
	require.Equal(t, http.StatusOK, w.Code)
}

func Test_server_ServiceInstanceGet(t *testing.T) {
	s := &server{}
	w, c := ginTestCtx(t, http.MethodGet, "/v2/service_instances/{instance_id}", "", false)
	s.ServiceInstanceGet(c, "", ServiceInstanceGetParams{})
	c.Writer.Flush() // status code needs the flush
	require.Empty(t, c.Errors.Errors())
	require.Equal(t, "", w.Body.String())
	require.Equal(t, http.StatusNotImplemented, w.Code)
}

func Test_server_ServiceInstanceUpdate(t *testing.T) {
	s := &server{}
	w, c := ginTestCtx(t, http.MethodPatch, "/v2/service_instances/{instance_id}", "", false)
	s.ServiceInstanceUpdate(c, "", ServiceInstanceUpdateParams{})
	require.Empty(t, c.Errors.Errors())
	require.Equal(t, "{}", w.Body.String())
	require.Equal(t, http.StatusOK, w.Code)
}

func Test_server_ServiceInstanceLastOperationGet(t *testing.T) {
	s := &server{}
	w, c := ginTestCtx(t, http.MethodGet, "/v2/service_instances/{instance_id}/last_operation", "", false)
	s.ServiceInstanceLastOperationGet(c, "", ServiceInstanceLastOperationGetParams{})
	c.Writer.Flush() // status code needs the flush
	require.Empty(t, c.Errors.Errors())
	require.Equal(t, "", w.Body.String())
	require.Equal(t, http.StatusNotImplemented, w.Code)
}

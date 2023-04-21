package servicebroker

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/cyberark/conjur-service-broker/pkg/conjur"
	"github.com/gin-gonic/gin"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/require"
)

func Test_server_ServiceInstanceProvision(t *testing.T) {
	tests := []struct {
		name string
	}{{
		"test",
	}}
	client, cleanup := testClient(t)
	defer cleanup()
	s := &server{
		client: client,
	}
	for _, tt := range tests {
		w, c := mockRequest(t, "PUT", "/v2/service_instances/9b292a9c-af66-4797-8d98-b30801f32ax8", `{
			"context": {
				"organization_guid": "e027f3f6-80fe-4d22-9374-da23a035ba0b",
				"space_guid":        "8c56f85c-c16e-4158-be79-5dac74f970de",
				"organization_name": "my-organization",
				"space_name":        "my-space"
			},
			"service_id":        "c024e536-6dc4-45c6-8a53-127e7f8275ab",
			"plan_id":           "3a116ac2-fc8b-496f-a715-e9a1b205d05c.community",
			"organization_guid": "e027f3f6-80fe-4d22-9374-da23a035ba0b",
			"space_guid":        "8c56f85c-c16e-4158-be79-5dac74f970de",
			"parameters":        {}
		}`)
		t.Run(tt.name, func(t *testing.T) {
			s.ServiceInstanceProvision(c, "", ServiceInstanceProvisionParams{})
		})
		require.Empty(t, c.Errors.Errors())
		require.Equal(t, "{}", w.Body.String())
		require.Equal(t, http.StatusCreated, w.Code)
	}
}

func mockRequest(t *testing.T, method, url string, body string) (*httptest.ResponseRecorder, *gin.Context) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	var b io.Reader
	if len(body) > 0 {
		b = strings.NewReader(body)
	}
	req, err := http.NewRequest(method, url, b)
	require.NoError(t, err)
	c.Request = req
	return w, c
}

func testClient(t *testing.T) (conjur.Client, func()) {
	httpmock.Activate()
	httpmock.RegisterResponder("POST", "https://conjur.local/authn/dev/test/authenticate",
		httpmock.NewJsonResponderOrPanic(http.StatusOK, gin.H{"protected": "eyJhbGciOiJjb25qdXIub3JnL3Nsb3NpbG8vdjIiLCJraWQiOiI5M2VjNTEwODRmZTM3Zjc3M2I1ODhlNTYyYWVjZGMxMSJ9", "payload": "eyJzdWIiOiJhZG1pbiIsImlhdCI6MTUxMDc1MzI1OX0=", "signature": "raCufKOf7sKzciZInQTphu1mBbLhAdIJM72ChLB4m5wKWxFnNz_7LawQ9iYEI_we1-tdZtTXoopn_T1qoTplR9_Bo3KkpI5Hj3DB7SmBpR3CSRTnnEwkJ0_aJ8bql5Cbst4i4rSftyEmUqX-FDOqJdAztdi9BUJyLfbeKTW9OGg-QJQzPX1ucB7IpvTFCEjMoO8KUxZpbHj-KpwqAMZRooG4ULBkxp5nSfs-LN27JupU58oRgIfaWASaDmA98O2x6o88MFpxK_M0FeFGuDKewNGrRc8lCOtTQ9cULA080M5CSnruCqu1Qd52r72KIOAfyzNIiBCLTkblz2fZyEkdSKQmZ8J3AakxQE2jyHmMT-eXjfsEIzEt-IRPJIirI3Qm"}))
	httpmock.RegisterResponder("GET", "=~https://conjur.local/resources/dev/(policy|layer)/cf.*",
		httpmock.NewJsonResponderOrPanic(http.StatusOK, gin.H{}))

	httpmock.RegisterResponder("POST", "https://conjur.local/policies/dev/policy/cf",
		httpmock.NewJsonResponderOrPanic(http.StatusOK, gin.H{}))
	client, err := (&conjur.Config{
		ConjurAccount:      "dev",
		ConjurApplianceURL: "https://conjur.local",
		ConjurPolicy:       "cf",
		ConjurAuthNLogin:   "test",
		ConjurAuthNAPIKey:  "test-api-key",
	}).NewClient()
	require.NoError(t, err)
	require.NoError(t, client.ValidateConnectivity())
	return client, httpmock.DeactivateAndReset
}

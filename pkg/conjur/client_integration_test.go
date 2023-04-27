//go:build integration

package conjur_test

import (
	"net/http"
	"testing"

	"github.com/cyberark/conjur-service-broker/pkg/conjur"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/mock"
)

func Test_client_ValidateConnectivity(t *testing.T) {
	tests := []struct {
		name string
		// fields  fields
		wantErr bool
	}{
		{"positive", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, cleanup := NewTestClient()
			defer cleanup()
			if err := client.ValidateConnectivity(); (err != nil) != tt.wantErr {
				t.Errorf("ValidateConnectivity() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

type J map[string]any

func NewTestClient() (conjur.Client, func()) {
	httpmock.Activate()
	httpmock.RegisterResponder("POST", "=~https://(conjur|follower).local/authn/dev/test/authenticate",
		httpmock.NewJsonResponderOrPanic(http.StatusOK, J{"protected": "eyJhbGciOiJjb25qdXIub3JnL3Nsb3NpbG8vdjIiLCJraWQiOiI5M2VjNTEwODRmZTM3Zjc3M2I1ODhlNTYyYWVjZGMxMSJ9", "payload": "eyJzdWIiOiJhZG1pbiIsImlhdCI6MTUxMDc1MzI1OX0=", "signature": "raCufKOf7sKzciZInQTphu1mBbLhAdIJM72ChLB4m5wKWxFnNz_7LawQ9iYEI_we1-tdZtTXoopn_T1qoTplR9_Bo3KkpI5Hj3DB7SmBpR3CSRTnnEwkJ0_aJ8bql5Cbst4i4rSftyEmUqX-FDOqJdAztdi9BUJyLfbeKTW9OGg-QJQzPX1ucB7IpvTFCEjMoO8KUxZpbHj-KpwqAMZRooG4ULBkxp5nSfs-LN27JupU58oRgIfaWASaDmA98O2x6o88MFpxK_M0FeFGuDKewNGrRc8lCOtTQ9cULA080M5CSnruCqu1Qd52r72KIOAfyzNIiBCLTkblz2fZyEkdSKQmZ8J3AakxQE2jyHmMT-eXjfsEIzEt-IRPJIirI3Qm"}))
	httpmock.RegisterResponder("GET", "=~https://follower.local/resources/dev/(policy|layer)/cf.*",
		httpmock.NewJsonResponderOrPanic(http.StatusOK, J{}))
	httpmock.RegisterResponder("GET", "https://conjur.local/resources/dev/policy/cf",
		httpmock.NewJsonResponderOrPanic(http.StatusOK, J{}))
	httpmock.RegisterResponder("POST", "https://conjur.local/policies/dev/policy/cf",
		httpmock.NewJsonResponderOrPanic(http.StatusOK, J{}))
	client, err := (&conjur.Config{
		ConjurAccount:      "dev",
		ConjurApplianceURL: "https://conjur.local",
		ConjurFollowerURL:  "https://follower.local",
		ConjurPolicy:       "cf",
		ConjurAuthNLogin:   "test",
		ConjurAuthNAPIKey:  "test-api-key",
	}).NewClient()
	if err != nil {
		panic(err)
	}
	return client, httpmock.DeactivateAndReset
}

func Test_client_ValidateConnectivityMock(t *testing.T) {
	tests := []struct {
		name string
		// fields  fields
		wantErr bool
	}{
		{"positive", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, mockAPI := conjur.NewMockClient()
			mockAPI.On("CheckPermission", mock.Anything, mock.Anything).Return(true, nil)
			if err := client.ValidateConnectivity(); (err != nil) != tt.wantErr {
				t.Errorf("ValidateConnectivity() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

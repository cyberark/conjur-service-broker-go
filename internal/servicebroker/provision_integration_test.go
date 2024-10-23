//go:build integration

package servicebroker

import (
	"net/http"
	"testing"

	"github.com/cyberark/conjur-api-go/conjurapi"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.cyberng.com/Conjur-Enterprise/conjur-service-broker-go/pkg/conjur"
)

func Test_server_ServiceInstanceProvisionOrgSpacePolicy(t *testing.T) {
	client, mockAPI := conjur.NewMockClient()
	mockAPI.On("LoadPolicy", conjurapi.PolicyModePost, "cf", mock.Anything).Return(nil, nil).Once()
	mockAPI.On("ResourceExists", "dev:policy:cf/e027f3f6-80fe-4d22-9374-da23a035ba0b").Return(true, nil).Once()
	mockAPI.On("ResourceExists", "dev:policy:cf/e027f3f6-80fe-4d22-9374-da23a035ba0b/8c56f85c-c16e-4158-be79-5dac74f970de").Return(true, nil).Once()
	mockAPI.On("ResourceExists", "dev:group:cf/e027f3f6-80fe-4d22-9374-da23a035ba0b/8c56f85c-c16e-4158-be79-5dac74f970de").Return(true, nil).Once()
	s := &server{client: client}
	w, c := ginTestCtx(t, "PUT", "/v2/service_instances/9b292a9c-af66-4797-8d98-b30801f32ax8", provisionBody, false)
	s.ServiceInstanceProvision(c, "", ServiceInstanceProvisionParams{})
	require.Empty(t, c.Errors.Errors())
	require.Equal(t, "{}", w.Body.String())
	require.Equal(t, http.StatusCreated, w.Code)
	mockAPI.AssertExpectations(t)
}

var loadPolicyResp = &conjurapi.PolicyResponse{
	CreatedRoles: map[string]conjurapi.CreatedRole{
		"role": {
			ID:     "role",
			APIKey: "api-key",
		},
	},
}

func Test_server_ServiceInstanceProvisionHostPolicy(t *testing.T) {
	client, mockAPI := conjur.NewMockClient()
	mockAPI.On("LoadPolicy", conjurapi.PolicyModePost, "cf", mock.Anything).Return(nil, nil).Once()
	mockAPI.On("ResourceExists", "dev:policy:cf/e027f3f6-80fe-4d22-9374-da23a035ba0b").Return(true, nil).Once()
	mockAPI.On("ResourceExists", "dev:policy:cf/e027f3f6-80fe-4d22-9374-da23a035ba0b/8c56f85c-c16e-4158-be79-5dac74f970de").Return(true, nil).Once()
	mockAPI.On("ResourceExists", "dev:group:cf/e027f3f6-80fe-4d22-9374-da23a035ba0b/8c56f85c-c16e-4158-be79-5dac74f970de").Return(true, nil).Once()

	mockAPI.On("LoadPolicy", conjurapi.PolicyModePost, "cf/e027f3f6-80fe-4d22-9374-da23a035ba0b/8c56f85c-c16e-4158-be79-5dac74f970de", mock.Anything).Return(loadPolicyResp, nil).Once()
	mockAPI.On("AddSecret", "cf/e027f3f6-80fe-4d22-9374-da23a035ba0b/8c56f85c-c16e-4158-be79-5dac74f970de/space-host-api-key", "api-key").Return(nil).Once()

	s := &server{client: client}
	w, c := ginTestCtx(t, "PUT", "/v2/service_instances/9b292a9c-af66-4797-8d98-b30801f32ax8", provisionBody, true)
	s.ServiceInstanceProvision(c, "", ServiceInstanceProvisionParams{})
	require.Empty(t, c.Errors.Errors())
	require.Equal(t, "{}", w.Body.String())
	require.Equal(t, http.StatusCreated, w.Code)
	mockAPI.AssertExpectations(t)
}

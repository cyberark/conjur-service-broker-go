//go:build integration

package servicebroker

import (
	"net/http"
	"testing"

	"github.com/cyberark/conjur-api-go/conjurapi"
	"github.com/cyberark/conjur-service-broker-go/pkg/conjur"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

var policyResp = &conjurapi.PolicyResponse{
	CreatedRoles: map[string]conjurapi.CreatedRole{"dev:host:role": {
		ID:     "dev:host:role",
		APIKey: "my-api-key",
	}},
}

func Test_server_ServiceBindingBinding(t *testing.T) {
	client, mockAPI := conjur.NewMockConjurClient()
	mockAPI.On("LoadPolicy", conjurapi.PolicyModePost, "cf/e027f3f6-80fe-4d22-9374-da23a035ba0b/8c56f85c-c16e-4158-be79-5dac74f970de", mock.Anything).Return(policyResp, nil).Once()
	mockAPI.On("RoleExists", "dev:host:cf/e027f3f6-80fe-4d22-9374-da23a035ba0b/8c56f85c-c16e-4158-be79-5dac74f970de/bb841d2b-8287-47a9-ac8f-eef4c16106f2").Return(false, nil).Once()
	mockAPI.On("Resource", "dev:host:test").Return(nil, nil).Once()
	s := &server{client: client}
	w, c := ginTestCtx(t, http.MethodPut, "v2/service_instances/9b292a9c-af66-4797-8d98-b30801f32ax8/service_bindings/bb841d2b-8287-47a9-ac8f-eef4c16106f2", bindBody, false)
	s.ServiceBindingBinding(c, "", "bb841d2b-8287-47a9-ac8f-eef4c16106f2", ServiceBindingBindingParams{})
	require.Empty(t, c.Errors.Errors())
	require.Equal(t, `{"credentials":{"account":"dev","appliance_url":"https://conjur.local","authn_api_key":"my-api-key","authn_login":"host/role","ssl_certificate":"","version":0}}`, w.Body.String())
	require.Equal(t, http.StatusCreated, w.Code)
	mockAPI.AssertExpectations(t)
}
func Test_server_ServiceBindingBindingSpaceIdentity(t *testing.T) {
	client, mockAPI := conjur.NewMockConjurClient()
	mockAPI.On("RoleExists", "dev:host:cf/e027f3f6-80fe-4d22-9374-da23a035ba0b/8c56f85c-c16e-4158-be79-5dac74f970de").Return(true, nil).Once()
	mockAPI.On("RetrieveSecret", "cf/e027f3f6-80fe-4d22-9374-da23a035ba0b/8c56f85c-c16e-4158-be79-5dac74f970de/space-host-api-key").Return([]byte("my-api-key"), nil).Once()
	s := &server{client: client}
	w, c := ginTestCtx(t, http.MethodPut, "v2/service_instances/9b292a9c-af66-4797-8d98-b30801f32ax8/service_bindings/bb841d2b-8287-47a9-ac8f-eef4c16106f2", bindBody, true)
	s.ServiceBindingBinding(c, "", "bb841d2b-8287-47a9-ac8f-eef4c16106f2", ServiceBindingBindingParams{})
	require.Empty(t, c.Errors.Errors())
	require.Equal(t, `{"credentials":{"account":"dev","appliance_url":"https://conjur.local","authn_api_key":"my-api-key","authn_login":"host/cf/e027f3f6-80fe-4d22-9374-da23a035ba0b/8c56f85c-c16e-4158-be79-5dac74f970de","ssl_certificate":"","version":0}}`, w.Body.String())
	require.Equal(t, http.StatusCreated, w.Code)
	mockAPI.AssertExpectations(t)
}

func Test_server_ServiceBindingUnbinding(t *testing.T) {
	client, mockAPI := conjur.NewMockConjurClient()
	mockAPI.On("Resources", &conjurapi.ResourceFilter{Kind: "host", Search: "bb841d2b-8287-47a9-ac8f-eef4c16106f2^"}).Return([]map[string]interface{}{{"id": "dev:host:cf/{orgID}/{spaceID}/{bindingID}"}}, nil).Once()
	mockAPI.On("RoleExists", "dev:host:cf/{orgID}/{spaceID}/bb841d2b-8287-47a9-ac8f-eef4c16106f2").Return(true, nil).Once()
	mockAPI.On("RotateAPIKey", "dev:host:cf/{orgID}/{spaceID}/bb841d2b-8287-47a9-ac8f-eef4c16106f2").Return(nil, nil).Once()
	mockAPI.On("LoadPolicy", conjurapi.PolicyModePut, "cf/{orgID}/{spaceID}", mock.Anything).Return(policyResp, nil).Once()
	s := &server{client: client}
	w, c := ginTestCtx(t, http.MethodDelete, "v2/service_instances/9b292a9c-af66-4797-8d98-b30801f32ax8/service_bindings/bb841d2b-8287-47a9-ac8f-eef4c16106f2", bindBody, false)
	s.ServiceBindingUnbinding(c, "", "bb841d2b-8287-47a9-ac8f-eef4c16106f2", ServiceBindingUnbindingParams{})
	require.Empty(t, c.Errors.Errors())
	require.Equal(t, `{}`, w.Body.String())
	require.Equal(t, http.StatusOK, w.Code)
	mockAPI.AssertExpectations(t)
}

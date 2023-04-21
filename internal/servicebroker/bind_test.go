package servicebroker

import (
	"net/http"
	"testing"

	"github.com/cyberark/conjur-service-broker/pkg/conjur"
	"github.com/cyberark/conjur-service-broker/pkg/conjur/mocks"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func Test_server_ServiceBindingBinding(t *testing.T) {
	tests := []struct {
		name string
	}{{
		"test",
	}}
	bind := &mocks.Bind{}
	bind.On("HostExists").Return(false, nil)
	bind.On("BindHostPolicy").Return(&conjur.Policy{}, nil)
	client := &mocks.Client{}
	client.On("NewBind", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(bind)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &server{
				client: client,
			}
			w, c := mockRequest(t, "PUT", "/v2/service_instances/9b292a9c-af66-4797-8d98-b30801f32ax8/service_bindings/bb841d2b-8287-47a9-ac8f-eef4c16106f2", `{
    "service_id": "c024e536-6dc4-45c6-8a53-127e7f8275ab",
    "plan_id": "3a116ac2-fc8b-496f-a715-e9a1b205d05c.community",
      "bind_resource": {
        "app_guid": "bb841d2b-8287-47a9-ac8f-eef4c16106f2"
      },
      "context": {
          "organization_guid": "e027f3f6-80fe-4d22-9374-da23a035ba0b",
          "space_guid": "8c56f85c-c16e-4158-be79-5dac74f970de"
      },
      "parameters": {
        "parameter1": 1,
        "parameter2": "foo"
      }
}`)
			s.ServiceBindingBinding(c, "", "bb841d2b-8287-47a9-ac8f-eef4c16106f2", ServiceBindingBindingParams{})
			require.Empty(t, c.Errors.Errors())
			require.Equal(t, `{"credentials":{"account":"","appliance_url":"","authn_api_key":"","authn_login":"","ssl_certificate":"","version":0}}`, w.Body.String())
			require.Equal(t, http.StatusCreated, w.Code)
		})
	}
}

package servicebroker

import (
	"errors"
	"net/http"
	"reflect"
	"testing"

	"github.com/cyberark/conjur-service-broker/pkg/conjur"
	"github.com/cyberark/conjur-service-broker/pkg/conjur/mocks"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

const (
	bindBody = `{
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
}`
	emptyResp = `{"credentials":{"account":"","appliance_url":"","authn_api_key":"","authn_login":"","ssl_certificate":"","version":0}}`
)

func Test_server_ServiceBinding(t *testing.T) {
	type p []interface{}
	type expects struct {
		status int
		body   string
		errors bool
	}
	type m struct {
		args    p
		returns p
	}
	type mockParams map[string]m
	type method struct {
		name   string
		params p
	}
	type args struct {
		body                string
		enableSpaceIdentity bool
		method              method
		bind                mockParams
		client              mockParams
	}
	tests := []struct {
		name    string
		args    args
		expects expects
	}{{
		"bind get",
		args{
			method: method{"ServiceBindingGet", p{"", "", ServiceBindingGetParams{}}},
		},
		expects{
			status: http.StatusNotImplemented,
		},
	}, {
		"bind get last operation",
		args{
			method: method{"ServiceBindingLastOperationGet", p{"", "", ServiceBindingLastOperationGetParams{}}},
		},
		expects{
			status: http.StatusNotImplemented,
		},
	}, {
		"bind unbind space identity",
		args{
			method:              method{"ServiceBindingUnbinding", p{"", "", ServiceBindingUnbindingParams{}}},
			enableSpaceIdentity: true,
		},
		expects{
			status: http.StatusOK,
			body:   "{}",
		},
	}, {
		"bind unbind host identity",
		args{
			method: method{"ServiceBindingUnbinding", p{"", "binding_id", ServiceBindingUnbindingParams{}}},
			bind: mockParams{
				"HostExists": m{
					returns: p{true, nil},
				},
				"DeleteBindHostPolicy": m{
					returns: p{nil},
				},
			},
			client: mockParams{
				"FromBindingID": m{
					args:    p{"binding_id"},
					returns: p{nil},
				},
			},
		},
		expects{
			status: http.StatusOK,
			body:   "{}",
		},
	}, {
		"bind unbind host identity - host not found",
		args{
			method: method{"ServiceBindingUnbinding", p{"", "binding_id", ServiceBindingUnbindingParams{}}},
			bind: mockParams{
				"HostExists": m{
					returns: p{false, nil},
				},
			},
			client: mockParams{
				"FromBindingID": m{
					args:    p{"binding_id"},
					returns: p{nil},
				},
			},
		},
		expects{
			status: http.StatusGone,
			errors: true,
		},
	}, {
		"bind unbind host identity - error form host exists",
		args{
			method: method{"ServiceBindingUnbinding", p{"", "binding_id", ServiceBindingUnbindingParams{}}},
			bind: mockParams{
				"HostExists": m{
					returns: p{false, errors.New("error")},
				},
			},
			client: mockParams{
				"FromBindingID": m{
					args:    p{"binding_id"},
					returns: p{nil},
				},
			},
		},
		expects{
			status: http.StatusInternalServerError,
			errors: true,
		},
	}, {
		"bind unbind host identity - error from bind",
		args{
			method: method{"ServiceBindingUnbinding", p{"", "binding_id", ServiceBindingUnbindingParams{}}},
			client: mockParams{
				"FromBindingID": m{
					args:    p{"binding_id"},
					returns: p{errors.New("error")},
				},
			},
		},
		expects{
			status: http.StatusInternalServerError,
			errors: true,
		},
	}, {
		"bind unbind host identity - error from delete",
		args{
			method: method{"ServiceBindingUnbinding", p{"", "binding_id", ServiceBindingUnbindingParams{}}},
			bind: mockParams{
				"HostExists": m{
					returns: p{true, nil},
				},
				"DeleteBindHostPolicy": m{
					returns: p{errors.New("error")},
				},
			},
			client: mockParams{
				"FromBindingID": m{
					args:    p{"binding_id"},
					returns: p{nil},
				},
			},
		},
		expects{
			status: http.StatusInternalServerError,
			errors: true,
		},
	}, {
		"bind - invalid request",
		args{
			method: method{"ServiceBindingBinding", p{"", "bb841d2b-8287-47a9-ac8f-eef4c16106f2", ServiceBindingBindingParams{}}},
		},
		expects{
			status: http.StatusBadRequest,
			errors: true,
		},
	}, {
		"bind - host exists error",
		args{
			method: method{"ServiceBindingBinding", p{"", "bb841d2b-8287-47a9-ac8f-eef4c16106f2", ServiceBindingBindingParams{}}},
			body:   bindBody,
			bind: mockParams{
				"HostExists": m{
					returns: p{false, errors.New("error")},
				},
			},
			client: mockParams{
				"NewBind": m{
					args: p{"e027f3f6-80fe-4d22-9374-da23a035ba0b", "8c56f85c-c16e-4158-be79-5dac74f970de", "bb841d2b-8287-47a9-ac8f-eef4c16106f2", false},
				},
			},
		},
		expects{
			status: http.StatusInternalServerError,
			errors: true,
		},
	}, {
		"bind - host identity",
		args{
			method: method{"ServiceBindingBinding", p{"", "bb841d2b-8287-47a9-ac8f-eef4c16106f2", ServiceBindingBindingParams{}}},
			body:   bindBody,
			bind: mockParams{
				"HostExists": m{
					returns: p{false, nil},
				},
				"BindHostPolicy": {
					returns: p{&conjur.Policy{}, nil},
				},
			},
			client: mockParams{
				"NewBind": m{
					args: p{"e027f3f6-80fe-4d22-9374-da23a035ba0b", "8c56f85c-c16e-4158-be79-5dac74f970de", "bb841d2b-8287-47a9-ac8f-eef4c16106f2", false},
				},
			},
		},
		expects{
			status: http.StatusCreated,
			body:   emptyResp,
		},
	}, {
		"bind - host identity - host found",
		args{
			method: method{"ServiceBindingBinding", p{"", "bb841d2b-8287-47a9-ac8f-eef4c16106f2", ServiceBindingBindingParams{}}},
			body:   bindBody,
			bind: mockParams{
				"HostExists": m{
					returns: p{true, nil},
				},
			},
			client: mockParams{
				"NewBind": m{
					args: p{"e027f3f6-80fe-4d22-9374-da23a035ba0b", "8c56f85c-c16e-4158-be79-5dac74f970de", "bb841d2b-8287-47a9-ac8f-eef4c16106f2", false},
				},
			},
		},
		expects{
			status: http.StatusConflict,
			errors: true,
		},
	}, {
		"bind - host identity - error on policy",
		args{
			method: method{"ServiceBindingBinding", p{"", "bb841d2b-8287-47a9-ac8f-eef4c16106f2", ServiceBindingBindingParams{}}},
			body:   bindBody,
			bind: mockParams{
				"HostExists": m{
					returns: p{false, nil},
				},
				"BindHostPolicy": {
					returns: p{nil, errors.New("error")},
				},
			},
			client: mockParams{
				"NewBind": m{
					args: p{"e027f3f6-80fe-4d22-9374-da23a035ba0b", "8c56f85c-c16e-4158-be79-5dac74f970de", "bb841d2b-8287-47a9-ac8f-eef4c16106f2", false},
				},
			},
		},
		expects{
			status: http.StatusInternalServerError,
			errors: true,
		},
	}, {
		"bind - space identity",
		args{
			method:              method{"ServiceBindingBinding", p{"", "bb841d2b-8287-47a9-ac8f-eef4c16106f2", ServiceBindingBindingParams{}}},
			body:                bindBody,
			enableSpaceIdentity: true,
			bind: mockParams{
				"HostExists": m{
					returns: p{true, nil},
				},
				"BindSpacePolicy": {
					returns: p{&conjur.Policy{}, nil},
				},
			},
			client: mockParams{
				"NewBind": m{
					args: p{"e027f3f6-80fe-4d22-9374-da23a035ba0b", "8c56f85c-c16e-4158-be79-5dac74f970de", "bb841d2b-8287-47a9-ac8f-eef4c16106f2", true},
				},
			},
		},
		expects{
			status: http.StatusCreated,
			body:   emptyResp,
		},
	}, {
		"bind - space identity - host not found",
		args{
			method:              method{"ServiceBindingBinding", p{"", "bb841d2b-8287-47a9-ac8f-eef4c16106f2", ServiceBindingBindingParams{}}},
			body:                bindBody,
			enableSpaceIdentity: true,
			bind: mockParams{
				"HostExists": m{
					returns: p{false, nil},
				},
			},
			client: mockParams{
				"NewBind": m{
					args: p{"e027f3f6-80fe-4d22-9374-da23a035ba0b", "8c56f85c-c16e-4158-be79-5dac74f970de", "bb841d2b-8287-47a9-ac8f-eef4c16106f2", true},
				},
			},
		},
		expects{
			status: http.StatusGone,
			errors: true,
		},
	}, {
		"bind - space identity",
		args{
			method:              method{"ServiceBindingBinding", p{"", "bb841d2b-8287-47a9-ac8f-eef4c16106f2", ServiceBindingBindingParams{}}},
			body:                bindBody,
			enableSpaceIdentity: true,
			bind: mockParams{
				"HostExists": m{
					returns: p{true, nil},
				},
				"BindSpacePolicy": {
					returns: p{&conjur.Policy{}, nil},
				},
			},
			client: mockParams{
				"NewBind": m{
					args: p{"e027f3f6-80fe-4d22-9374-da23a035ba0b", "8c56f85c-c16e-4158-be79-5dac74f970de", "bb841d2b-8287-47a9-ac8f-eef4c16106f2", true},
				},
			},
		},
		expects{
			status: http.StatusCreated,
			body:   emptyResp,
		},
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bind := &mocks.Bind{}
			for method, v := range tt.args.bind {
				bind.On(method, v.args...).Return(v.returns...).Once()
			}
			client := &mocks.Client{}
			for method, v := range tt.args.client {
				client.On(method, v.args...).Return(append(p{bind}, v.returns...)...).Once()
			}
			s := &server{client: client}
			w, c := ginTestCtx(t, "", "", tt.args.body, tt.args.enableSpaceIdentity)
			reflect.ValueOf(s).MethodByName(tt.args.method.name).Call(toValues(c, tt.args.method.params))
			c.Writer.Flush()
			if tt.expects.errors {
				assert.NotEmpty(t, c.Errors.Errors())
			} else {
				assert.Empty(t, c.Errors.Errors())
			}
			if len(tt.expects.body) > 2 {
				assert.JSONEq(t, tt.expects.body, w.Body.String())
			} else {
				assert.Equal(t, tt.expects.body, w.Body.String())
			}
			assert.Equal(t, tt.expects.status, w.Code)

			client.AssertExpectations(t)
			bind.AssertExpectations(t)
		})
	}
}

func toValues(c *gin.Context, params []interface{}) []reflect.Value {
	v := make([]reflect.Value, len(params)+1)
	v[0] = reflect.ValueOf(c)
	for i, p := range params {
		v[i+1] = reflect.ValueOf(p)
	}
	return v
}

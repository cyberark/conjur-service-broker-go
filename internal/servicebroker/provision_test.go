//go:build !integration

package servicebroker

import (
	"errors"
	"net/http"
	"reflect"
	"testing"

	"github.com/cyberark/conjur-service-broker-go/pkg/conjur/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func Test_server_ServiceInstance(t *testing.T) {
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
		provision           mockParams
		client              mockParams
	}
	tests := []struct {
		name    string
		args    args
		expects expects
	}{{
		"provision get",
		args{
			method: method{"ServiceInstanceGet", p{"", ServiceInstanceGetParams{}}},
		},
		expects{
			status: http.StatusNotImplemented,
		},
	}, {
		"provision get last operation",
		args{
			method: method{"ServiceInstanceLastOperationGet", p{"", ServiceInstanceLastOperationGetParams{}}},
		},
		expects{
			status: http.StatusNotImplemented,
		},
	}, {
		"deprovision",
		args{
			method: method{"ServiceInstanceDeprovision", p{"", ServiceInstanceDeprovisionParams{}}},
		},
		expects{
			status: http.StatusOK,
			body:   "{}",
		},
	}, {
		"provision update",
		args{
			method: method{"ServiceInstanceUpdate", p{"", ServiceInstanceUpdateParams{}}},
		},
		expects{
			status: http.StatusOK,
			body:   "{}",
		},
	}, {
		"provision - bad request",
		args{
			method: method{"ServiceInstanceProvision", p{"", ServiceInstanceProvisionParams{}}},
		},
		expects{
			status: http.StatusBadRequest,
			errors: true,
		},
	}, {
		"provision space identity",
		args{
			method:    method{"ServiceInstanceProvision", p{"", ServiceInstanceProvisionParams{}}},
			provision: mockParams{"ProvisionOrgSpacePolicy": m{returns: p{nil}}},
			client:    mockParams{"NewProvision": m{args: p{"e027f3f6-80fe-4d22-9374-da23a035ba0b", "8c56f85c-c16e-4158-be79-5dac74f970de", mock.Anything, mock.Anything}, returns: p{nil}}},
			body:      provisionBody,
		},
		expects{
			status: http.StatusCreated,
			body:   "{}",
		},
	}, {
		"provision space identity - error on policy",
		args{
			method:    method{"ServiceInstanceProvision", p{"", ServiceInstanceProvisionParams{}}},
			provision: mockParams{"ProvisionOrgSpacePolicy": m{returns: p{errors.New("error")}}},
			client:    mockParams{"NewProvision": m{args: p{"e027f3f6-80fe-4d22-9374-da23a035ba0b", "8c56f85c-c16e-4158-be79-5dac74f970de", mock.Anything, mock.Anything}, returns: p{nil}}},
			body:      provisionBody,
		},
		expects{
			status: http.StatusInternalServerError,
			errors: true,
		},
	}, {
		"provision host identity",
		args{
			method: method{"ServiceInstanceProvision", p{"", ServiceInstanceProvisionParams{}}},
			provision: mockParams{
				"ProvisionOrgSpacePolicy": m{returns: p{nil}},
				"ProvisionHostPolicy":     m{returns: p{nil}},
			},
			client:              mockParams{"NewProvision": m{args: p{"e027f3f6-80fe-4d22-9374-da23a035ba0b", "8c56f85c-c16e-4158-be79-5dac74f970de", mock.Anything, mock.Anything}, returns: p{nil}}},
			body:                provisionBody,
			enableSpaceIdentity: true,
		},
		expects{
			status: http.StatusCreated,
			body:   "{}",
		},
	}, {
		"provision host identity - error in policy",
		args{
			method: method{"ServiceInstanceProvision", p{"", ServiceInstanceProvisionParams{}}},
			provision: mockParams{
				"ProvisionOrgSpacePolicy": m{returns: p{nil}},
				"ProvisionHostPolicy":     m{returns: p{errors.New("error")}},
			},
			client:              mockParams{"NewProvision": m{args: p{"e027f3f6-80fe-4d22-9374-da23a035ba0b", "8c56f85c-c16e-4158-be79-5dac74f970de", mock.Anything, mock.Anything}, returns: p{nil}}},
			body:                provisionBody,
			enableSpaceIdentity: true,
		},
		expects{
			status: http.StatusInternalServerError,
			errors: true,
		},
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			provision := mocks.NewMockProvision(t)
			for method, v := range tt.args.provision {
				provision.On(method, v.args...).Return(v.returns...).Once()
			}
			client := mocks.NewMockClient(t)
			for method, v := range tt.args.client {
				client.On(method, v.args...).Return(append(p{provision}, v.returns...)...).Once()
			}
			s := NewServerImpl(client)
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
			provision.AssertExpectations(t)
		})
	}
}

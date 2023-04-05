package http

import (
	"testing"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/getkin/kin-openapi/openapi3filter"
)

func Test_errMsg(t *testing.T) {
	type args struct {
		err error
	}
	tests := []struct {
		name string
		args args
		want string
	}{{
		"validation error",
		args{openapi3.MultiError{
			&openapi3filter.RequestError{
				Reason: "doesn't match schema #/components/schemas/ServiceInstanceUpdateRequestBody",
				Err: openapi3.MultiError{
					&openapi3.SchemaError{
						Reason: `property "service_id" is missing`,
					},
				},
			},
		}},
		`doesn't match schema #/components/schemas/ServiceInstanceUpdateRequestBody property "service_id" is missing`,
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := errMsg(tt.args.err); got != tt.want {
				t.Errorf("errMsg() = %v, want %v", got, tt.want)
			}
		})
	}
}

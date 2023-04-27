//go:build !integration

package http

import (
	"testing"
)

func Test_camelCasedStatus(t *testing.T) {
	type args struct {
		code int
	}
	tests := []struct {
		name string
		args args
		want string
	}{{
		"bad request",
		args{400},
		"BadRequest",
	}, {
		"unauthorized",
		args{401},
		"Unauthorized",
	}, {
		"teapot",
		args{418},
		"IMATeapot",
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := camelCasedStatus(tt.args.code); got != tt.want {
				t.Errorf("camelCasedStatus() = %v, want %v", got, tt.want)
			}
		})
	}
}

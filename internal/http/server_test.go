//go:build !integration

package http

import (
	"testing"

	"github.com/cyberark/conjur-service-broker-go/internal/ctxutil"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func Test_initCtx(t *testing.T) {
	type args struct {
		logger *zap.Logger
		cfg    *config
	}
	tests := []struct {
		name      string
		args      args
		wantEmpty bool
	}{{
		"not nil",
		args{
			logger: zap.NewNop(),
			cfg:    &config{},
		},
		false,
	}, {
		"nil",
		args{},
		true,
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := initCtx(tt.args.logger, tt.args.cfg)
			if tt.wantEmpty {
				require.Equal(t, ctxutil.NewContext(), got)
			} else {
				require.NotNil(t, got)
				require.NotEqualf(t, ctxutil.NewContext(), got, "expecting non empty context")
			}
		})
	}
}

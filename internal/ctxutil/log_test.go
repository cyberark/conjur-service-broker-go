package ctxutil

import (
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func Test_ctx_Logger(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	tests := []struct {
		name    string
		ctx     Context
		wantNil bool
	}{{
		"empty",
		NewContext(),
		true,
	}, {
		"non empty",
		NewContext().WithLogger(logger.Sugar()),
		false,
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.ctx.Logger()
			if tt.wantNil {
				require.Nil(t, got)
			} else {
				require.NotNil(t, got)
			}
		})
	}
}

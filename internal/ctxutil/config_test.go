package ctxutil

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_ctx_IsEnableSpaceIdentity(t *testing.T) {
	tests := []struct {
		name string
		ctx  Context
		want bool
	}{{
		"empty context",
		NewContext(),
		false,
	}, {
		"non empty context",
		NewContext().WithEnableSpaceIdentity(true),
		true,
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.ctx.IsEnableSpaceIdentity()
			require.Equal(t, tt.want, got)
		})
	}
}

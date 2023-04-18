package conjur

import (
	"io"
	"testing"

	"github.com/doodlesbykumbi/conjur-policy-go/pkg/conjurpolicy"
	"github.com/stretchr/testify/require"
)

func Test_policyReader(t *testing.T) {
	tests := []struct {
		name    string
		policy  conjurpolicy.PolicyStatements
		want    string
		wantErr bool
	}{{
		"empty",
		[]conjurpolicy.Resource{},
		"\n",
		false,
	}, {
		"non empty",
		[]conjurpolicy.Resource{conjurpolicy.Layer{}},
		"- !layer\n",
		false,
	}, {
		"nil",
		nil,
		"\n",
		false,
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := policyReader(tt.policy)
			if tt.wantErr {
				require.Error(t, err)
			}
			bytes, err := io.ReadAll(got)
			require.NoError(t, err)
			require.Equal(t, string(bytes), tt.want)
		})
	}
}

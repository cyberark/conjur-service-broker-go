//go:build !integration

package conjur

import (
	"io"
	"math"
	"testing"

	"github.com/doodlesbykumbi/conjur-policy-go/pkg/conjurpolicy"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

type invalidType struct {
	conjurpolicy.Resource
	Fail yaml.Node
}

func Test_policyReader(t *testing.T) {
	tests := []struct {
		name    string
		policy  conjurpolicy.PolicyStatements
		want    string
		wantErr assert.ErrorAssertionFunc
	}{{
		"empty",
		[]conjurpolicy.Resource{},
		"\n",
		assert.NoError,
	}, {
		"non empty",
		[]conjurpolicy.Resource{conjurpolicy.Layer{}},
		"- !layer\n",
		assert.NoError,
	}, {
		"nil",
		nil,
		"\n",
		assert.NoError,
	}, {
		"invalid",
		[]conjurpolicy.Resource{invalidType{
			Fail: yaml.Node{Kind: yaml.Kind(math.MaxUint32)},
		}},
		"",
		assert.Error,
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := policyReader(tt.policy)
			tt.wantErr(t, err)
			if got == nil {
				require.Empty(t, tt.want)
				return
			}
			bytes, err := io.ReadAll(got)
			assert.NoError(t, err)
			assert.Equal(t, string(bytes), tt.want)
		})
	}
}

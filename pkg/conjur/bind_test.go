//go:build !integration

package conjur

import (
	"io"
	"testing"

	"github.com/cyberark/conjur-api-go/conjurapi"
	"github.com/stretchr/testify/assert"
)

func Test_createBindYAML(t *testing.T) {
	tests := []struct {
		name    string
		args    *bind
		want    string
		wantErr assert.ErrorAssertionFunc
	}{{
		"simple bind",
		&bind{bindingID: "test", client: &client{roClient: &conjurapi.Client{}, config: &Config{}}},
		`- !host
  id: test
`,
		assert.NoError,
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.args.createBindYAML()
			tt.wantErr(t, err)
			gotBytes, err := io.ReadAll(got)
			assert.NoError(t, err)
			assert.Equal(t, tt.want, string(gotBytes))
		})
	}
}

func Test_deleteBindYAML(t *testing.T) {
	tests := []struct {
		name    string
		args    *bind
		want    string
		wantErr assert.ErrorAssertionFunc
	}{{
		"simple delete",
		&bind{bindingID: "test"},
		`- !delete
  record: !host test
`,
		assert.NoError,
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.args.deleteBindYAML()
			tt.wantErr(t, err)
			gotBytes, err := io.ReadAll(got)
			assert.NoError(t, err)
			assert.Equal(t, tt.want, string(gotBytes))
		})
	}
}

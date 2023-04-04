package conjur

import (
	"io"
	"testing"

	"github.com/cyberark/conjur-api-go/conjurapi"
)

func Test_createBindYAML(t *testing.T) {
	tests := []struct {
		name    string
		args    *bind
		want    string
		wantErr bool
	}{{
		"simple bind",
		&bind{bindingID: "test", client: &client{roClient: &conjurapi.Client{}, config: &Config{}}},
		`- !host
  id: test
`,
		false,
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.args.createBindYAML()
			if (err != nil) != tt.wantErr {
				t.Errorf("createBindYAML() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			gotBytes, err := io.ReadAll(got)
			if err != nil {
				t.Errorf("createBindYAML() error = %v", err)
				return
			}
			if string(gotBytes) != tt.want {
				t.Errorf("createBindYAML() got = \n%v\n, want \n%v\n", string(gotBytes), tt.want)
			}
		})
	}
}

func Test_deleteBindYAML(t *testing.T) {
	tests := []struct {
		name    string
		args    *bind
		want    string
		wantErr bool
	}{{
		"simple delete",
		&bind{bindingID: "test"},
		`- !delete
  record: !host test
`,
		false,
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.args.deleteBindYAML()
			if (err != nil) != tt.wantErr {
				t.Errorf("deleteBindYAML() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			gotBytes, err := io.ReadAll(got)
			if err != nil {
				t.Errorf("deleteBindYAML() error = %v", err)
				return
			}
			if string(gotBytes) != tt.want {
				t.Errorf("deleteBindYAML() got = \n%v\n, want \n%v\n", string(gotBytes), tt.want)
			}
		})
	}
}

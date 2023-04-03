package conjur

import (
	"io"
	"testing"
)

func Test_createBindYAML(t *testing.T) {
	tests := []struct {
		name    string
		args    *bind
		want    string
		wantErr bool
	}{{
		"simple bind",
		&bind{BindID: "test"},
		`- !host
  id: test
- !grant
  role: !layer
  member: !layer test
`,
		false,
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := createBindYAML(tt.args)
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
		&bind{BindID: "test"},
		`- !delete
  record: !host test
`,
		false,
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := deleteBindYAML(tt.args)
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

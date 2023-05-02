// Package servicebroker provides an implementation of the generated service broker server
package servicebroker

import (
	"github.com/cyberark/conjur-service-broker/pkg/conjur"
)

//go:generate oapi-codegen --config ./oapi-codegen.yaml ../../api/openapi.yaml
//go:generate sh -c "echo '//lint:file-ignore ST1005 Ignore generated file' >> servicebroker.gen.go"
//go:generate goimports -w servicebroker.gen.go
//go:generate mockery --name=ServerInterface
type server struct {
	client conjur.Client
}

// NewServerImpl creates the webservice implementation
func NewServerImpl(client conjur.Client) ServerInterface {
	return &server{client}
}

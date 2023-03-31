package servicebroker

import (
	"github.com/cyberark/conjur-service-broker/pkg/conjur"
)

//go:generate oapi-codegen --config ./oapi-codegen.yaml ../../api/openapi.yaml

// ServerImpl service broker server implementation
type ServerImpl struct {
	client *conjur.Client
}

// NewServerImpl creates the webservice implementation
func NewServerImpl(client *conjur.Client) ServerImpl {
	return ServerImpl{client}
}

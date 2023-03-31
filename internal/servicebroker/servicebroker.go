package servicebroker

import (
	"github.com/cyberark/conjur-service-broker/pkg/conjur"
)

//go:generate oapi-codegen --config ./oapi-codegen.yaml ../../api/openapi.yaml

// server service broker server implementation
type server struct {
	client conjur.Client
}

// NewServerImpl creates the webservice implementation
func NewServerImpl(client conjur.Client) ServerInterface {
	return &server{client}
}

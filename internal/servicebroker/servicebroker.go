package servicebroker

import (
	"github.com/cyberark/conjur-service-broker/pkg/conjur"
)

//go:generate oapi-codegen --config ./oapi-codegen.yaml ../../api/openapi.yaml

type serverImpl struct {
	client *conjur.Client
}

// NewServerImpl creates the webservice implementation
func NewServerImpl(client *conjur.Client) serverImpl {
	return serverImpl{client}
}

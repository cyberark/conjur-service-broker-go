package main

import (
	"github.com/cyberark/conjur-service-broker/conjur"
)

// ServerImpl is the webservice implementation
type ServerImpl struct {
	client *conjur.Client
}

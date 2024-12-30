// Package main is the conjur service broker main binary
package main

import (
	"github.com/cyberark/conjur-service-broker-go/internal/http"
)

func main() {
	http.StartHTTPServer(nil)
}

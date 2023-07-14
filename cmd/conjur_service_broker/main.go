// Package main is the conjur service broker main binary
package main

import (
	"github.cyberng.com/Conjur-Enterprise/conjur-service-broker-go/internal/http"
)

func main() {
	http.StartHTTPServer()
}

package main

import (
	"log"

	"github.com/cyberark/conjur-service-broker/internal/http"
)

func main() {
	if err := http.StartHTTPServer(); err != nil {
		log.Fatal(err)
	}
}

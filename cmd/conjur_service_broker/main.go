// Package main is the conjur service broker main binary
package main

import (
	"fmt"
	"log"

	"github.com/cyberark/conjur-service-broker/internal/http"
	"go.uber.org/zap"
)

func main() {
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatal(fmt.Errorf("failed to initialize logger: %w", err))
	}
	if err := http.StartHTTPServer(logger); err != nil {
		logger.Sugar().Fatal("failed to start http server: %s", err)
	}
}

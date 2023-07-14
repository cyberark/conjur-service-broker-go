// Package main provides integration tests for conjur service broker
package main

type config struct {
	ServiceURL        string `env:"SERVICE_URL" envDefault:"http://localhost:8080"`
	BasicAuthUser     string `env:"SECURITY_USER_NAME,unset" envDefault:"test"`
	BasicAuthPassword string `env:"SECURITY_USER_PASSWORD,unset" envDefault:"test"`
}

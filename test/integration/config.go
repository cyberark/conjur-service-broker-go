// Package main provides integration tests for conjur service broker
package main

type config struct {
	ServiceURL        string `env:"SERVICE_URL" envDefault:"http://localhost:8080"`
	BasicAuthUser     string `env:"BASIC_AUTH_USER,unset" envDefault:"test"`
	BasicAuthPassword string `env:"BASIC_AUTH_PASSWORD,unset" envDefault:"test"`
}

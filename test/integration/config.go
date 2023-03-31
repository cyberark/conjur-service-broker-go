package main

type config struct {
	ServiceURL        string `env:"SERVICE_URL" envDefault:"http://localhost:8080"`
	BasicAuthUser     string `env:"BASIC_AUTH_USER,unset"`
	BasicAuthPassword string `env:"BASIC_AUTH_PASSWORD,unset"`
}

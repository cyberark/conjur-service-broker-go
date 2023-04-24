package api

//go:generate interfacer -for github.com/cyberark/conjur-api-go/conjurapi.Client -as api.Client -o client.gen.go
//go:generate goimports -w client.gen.go
//go:generate mockery --name=Client

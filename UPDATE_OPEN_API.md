

// TODO: update this doc

wget https://raw.githubusercontent.com/openservicebrokerapi/servicebroker/master/openapi.yaml -P ./api

go install github.com/deepmap/oapi-codegen/cmd/oapi-codegen@latest

go generate ./...

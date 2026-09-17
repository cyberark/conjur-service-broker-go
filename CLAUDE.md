# CLAUDE.md

## Overview

This repo implements the **CyberArk Secrets Manager (Conjur) Service Broker for Cloud Foundry** — a Go service that implements the [Open Service Broker API](https://www.openservicebrokerapi.org/) (v2.13). It runs as a Cloud Foundry app and lets CF apps automatically receive a Conjur/Secrets Manager `host` identity when deployed, so they can securely retrieve secrets at runtime (typically via the companion [cloudfoundry-conjur-buildpack](https://github.com/cyberark/cloudfoundry-conjur-buildpack)). It talks to a Conjur/Secrets Manager (Self-Hosted or OSS) instance to provision/deprovision policy and mint host credentials on bind/unbind. Part of the CyberArk Secrets Manager Open Source Suite.

## Tech Stack

- **Language:** Go (see `go.mod` for exact version — currently 1.27.x)
- **HTTP framework:** [gin-gonic/gin](https://github.com/gin-gonic/gin)
- **Conjur client:** `github.com/cyberark/conjur-api-go` (see the `replace` directive in `go.mod` — pinned to a specific version; there's a commented-out `replace` for local development against a sibling `../conjur-api-go` checkout)
- **API server code generation:** `oapi-codegen` generates `internal/servicebroker/servicebroker.gen.go` from `api/openapi.yaml` (config in `internal/servicebroker/oapi-codegen.yaml`). Regenerate with `go generate ./...` — never hand-edit the `.gen.go` file.
- **Logging:** logrus + zap (`internal/http/logrus_zap_hook.go` bridges them)
- **Mocks:** generated via `mockery` (`.mockery.yaml`), found under `pkg/conjur/mocks`, `pkg/conjur/api/mocks`, and `internal/servicebroker/mocks`
- **Local dev orchestration:** [Tilt](https://tilt.dev) (`Tiltfile`, `dev/`) with Kubernetes/KinD

## Build

```sh
go build ./cmd/conjur_service_broker
```

Container image build is defined in `Dockerfile` (multi-stage, `CGO_ENABLED=0`) and duplicated inline in the `Tiltfile`'s `docker_build` call for local dev.

## Test

- **Unit + integration tests together (CI-style):** `./scripts/test.sh` — runs `go test` (unit, then `-tags=integration`), merges coverage via `go tool covdata`, and produces `coverage/all_no_gen`, JUnit, and Cobertura reports.
- **Plain unit tests:** `go test ./...` (exclude `test/e2e` and `test/integration`, which are separate Go modules)
- **Integration tests** use the `integration` build tag and live alongside unit tests as `*_integration_test.go` files in `internal/http`, `internal/servicebroker`, and `pkg/conjur`. Run with `go test -tags=integration ./...`.
- **Local dev loop (recommended for iterating):** `tilt up` — starts a containerized Conjur + broker + Postgres, volume-mounts source so code changes are picked up live, and auto-runs unit/integration tests via `local_resource(name='tests', cmd='./scripts/test.sh')`.
- **Containerized integration tests (no local Tilt/KinD needed):** `./scripts/test_integration.sh` — builds `Dockerfile.integration` and runs `tilt ci` inside Docker against a KinD cluster.
- **E2E tests** (`test/e2e/`, own Go module) run against Tanzu VMs in the Jenkins pipeline only; not practically runnable locally, and the "End-to-End testing" stage is currently commented out in `Jenkinsfile`.
- `test/integration` and `test/e2e` are each **separate Go modules** (own `go.mod`/`go.sum`, own Go version directive) — `go build ./...` / `go test ./...` from the repo root will not include them, and bumping the Go version has to be done in all four `go.mod` files (root, `test/integration`, `test/e2e`, `test/e2e/sample-app`) plus the Docker base images (`Dockerfile`, `Dockerfile.tools`, `Dockerfile.integration`, `scripts/Dockerfile.test`, `test/e2e/Dockerfile.test`) to stay consistent.

## Configuration / Authentication

Config is parsed from environment variables via `github.com/caarlos0/env` into `pkg/conjur.Config` (embedded in `internal/http`'s `config`) — see `pkg/conjur/config.go` for the full set of `CONJUR_*` vars and their doc comments, and `README.md`'s Installation section for the operator-facing walkthrough.

The broker supports two mutually-exclusive ways to authenticate to Conjur/Secrets Manager, selected by whether `CONJUR_AUTHN_CERT_SERVICE_ID` is set (`Config.IsAuthnCert()`):

- **API Key auth (default):** requires `CONJUR_AUTHN_LOGIN` + `CONJUR_AUTHN_API_KEY`.
- **Certificate (mTLS) auth (`authn-cert`):** set `CONJUR_AUTHN_CERT_SERVICE_ID` to enable; requires `CONJUR_CLIENT_CERT` + `CONJUR_CLIENT_CERT_KEY` (PEM-encoded, inline). `CONJUR_AUTHN_CERT_HOST_ID` is required for `request` host-mode authenticators but not for `spiffe` ones. This mode is validated in `internal/http/config.go#validate` and wired up in `pkg/conjur/client.go#NewClient`.

Other notable config:
- `CONJUR_FOLLOWER_URL` / `CONJUR_FOLLOWER_SSL_CERTIFICATE` (HA only): a read-only client is built against the follower when set (see `NewClient`'s `roClient` branch for both auth modes); otherwise the same client is used for reads and writes.
- `ENABLE_SPACE_IDENTITY`: space-scoped vs. per-app host identity — see `README.md`'s Space-scoped/Application-scoped Identity sections.
- `SECURITY_USER_NAME` / `SECURITY_USER_PASSWORD`: **required** HTTP Basic Auth creds for the Cloud Controller; the broker refuses to start without them.

## Repo Structure

- `cmd/conjur_service_broker/` — `main.go`, the binary entrypoint; just wires up `internal/http.StartHTTPServer`.
- `api/openapi.yaml` — OpenAPI spec defining the service broker HTTP API; source of truth for generated server code.
- `internal/servicebroker/` — generated gin server (`servicebroker.gen.go`) plus hand-written handler implementations (`catalog.go`, `provision.go`, `bind.go`, `servicebroker.go`).
- `internal/http/` — HTTP server setup, config (incl. authn-cert validation), error handling, request validation, logging glue.
- `internal/ctxutil/` — context helpers for config/logging.
- `internal/version/` — build version info.
- `pkg/conjur/` — Conjur/Secrets Manager client wrapper: policy provisioning, binding (host creation/deletion), config (API Key + certificate auth), YAML policy helpers. Depends on `conjur-api-go` and `doodlesbykumbi/conjur-policy-go`.
- `configs/manifest.yml` — Cloud Foundry `cf push` manifest for the broker app itself.
- `dev/` — Tilt dev-environment dependencies (Conjur server/client containers, Postgres).
- `test/integration/` — standalone Go module with Gherkin-style feature-driven integration tests against a real Conjur instance.
- `test/e2e/` — standalone Go module with end-to-end tests against a full CF/Tanzu deployment (pipeline stage currently disabled — see Test section above).
- `scripts/` — build, test, release, and container-image helper scripts used by CI (`Jenkinsfile`) and Tilt. `scripts/utils`' `go_version` helper reads the `go` directive out of a module's `go.mod` (pass the module dir explicitly, e.g. for `test/e2e`) so Docker build args stay in sync with each module's pinned Go version.
- `Dockerfile` / `Dockerfile.integration` / `Dockerfile.tools` — runtime image, integration-test-runner image, and tooling image respectively.

## Conventions & Gotchas

- **Never hand-edit `*.gen.go` files** (`servicebroker.gen.go`, `pkg/conjur/*.gen.go`) — they're regenerated by `go generate ./...` from `api/openapi.yaml`. Edit the OpenAPI spec or the generator config instead.
- The broker **requires** `SECURITY_USER_NAME` and `SECURITY_USER_PASSWORD` env vars (HTTP Basic Auth for Cloud Controller); it refuses to start without them. Other Conjur-related config is documented above and in `README.md`.
- The broker is a **credential-minting endpoint** — README strongly emphasizes deploying it with `--no-route` and only an internal CF route (`apps.internal`), reachable solely by the Cloud Controller. Keep that security posture in mind if touching deployment/manifest config.
- Mocks (`pkg/conjur/mocks`, `pkg/conjur/api/mocks`, `internal/servicebroker/mocks`) are generated by `mockery`; regenerate rather than editing by hand when interfaces change.
- Pre-commit hooks (`.pre-commit-config.yaml`) run `go generate`, `go build`, `go mod tidy`, `go test` (excluding `test/e2e|test/integration`), `go vet`, `gitleaks`, `shellcheck`/`shfmt`, and `remarklint` — expect these checks in CI/local hooks.
- `NOTICES.txt` must be kept in sync with `go.mod` dependency/license changes (see `CONTRIBUTING.md` release checklist using `go-licenses`).
- **Go version is pinned in four `go.mod` files and five Dockerfiles simultaneously** (see the Test section's module list). Jenkins CI runs an internal `updatePrivateGoDependencies` check that hard-fails the whole pipeline before any build/test stage if the pinned Go minor version has fallen out of the Go team's supported window — when bumping, update all of them together, not just the root `go.mod`.

---
name: ai-pr-review
description: Automated audit of code changes before opening a Pull Request against this repo (conjur-service-broker-go, the Go Cloud Foundry Open Service Broker for CyberArk Secrets Manager/Conjur). Analyzes the diff against the base branch, checks compliance with CyberArk standards (architecture, security, quality) and this repo's own conventions from CLAUDE.md, and generates a PR_REVIEW_REPORT.md ready to paste into the PR description. Use before every commit that finalizes a task, or when the user writes "@ai-pr-review run analysis and generate report".
---

# AI PR Review

Automated audit of code changes before opening a Pull Request. The skill analyzes the diff, checks compliance with the CyberArk and repo-specific conventions used in `conjur-service-broker-go`, and generates a report ready to paste into the PR description.

## Workflow

### 1. Context Extraction

- Run `git diff` against the base branch (`main`).
- Identify the changed components (see `CLAUDE.md` for the full layout):
  - `internal/servicebroker/` — Open Service Broker API handlers (`catalog.go`, `provision.go`, `bind.go`) plus the generated `servicebroker.gen.go`
  - `internal/http/` — server setup, config (incl. `authn-cert` validation), error handling, logging glue
  - `pkg/conjur/` — Conjur/Secrets Manager client wrapper (policy provisioning, host bind/unbind, API Key + certificate auth, YAML policy helpers)
  - `api/openapi.yaml` — source of truth for the generated server code
  - `test/integration/`, `test/e2e/` — separate Go modules with their own `go.mod`

### 2. Multi-Persona Analysis

- **Architect**:
  - Any change to `internal/servicebroker/servicebroker.gen.go` or `pkg/conjur/*.gen.go` by hand is a hard blocker — these are generated via `go generate ./...` from `api/openapi.yaml`; the PR must instead change the spec/generator config and re-run generation.
  - New OSB API behavior (provision/deprovision/bind/unbind) should go through the existing handler structure in `internal/servicebroker/`, not bypass it.
- **Security**:
  - No hardcoded Conjur credentials, API keys, or certs. The two supported Conjur auth modes (API Key via `CONJUR_AUTHN_LOGIN`/`CONJUR_AUTHN_API_KEY`, or certificate/mTLS via `CONJUR_AUTHN_CERT_SERVICE_ID` + `CONJUR_CLIENT_CERT`/`CONJUR_CLIENT_CERT_KEY`) must stay mutually exclusive and validated in `internal/http/config.go#validate`, consistent with `pkg/conjur/client.go#NewClient`.
  - `SECURITY_USER_NAME`/`SECURITY_USER_PASSWORD` (HTTP Basic Auth for the Cloud Controller) must remain required — the broker should refuse to start without them; flag any change that weakens this.
  - This service is a credential-minting endpoint; deployment/manifest changes (`configs/manifest.yml`) must preserve the "internal-only route" posture (`--no-route` / `apps.internal`) documented in `README.md`.
- **Quality**:
  - New logic in `internal/` or `pkg/conjur/` needs unit tests (`go test ./...`), and integration coverage (`*_integration_test.go`, `-tags=integration`) where it touches real Conjur interactions.
  - Mocks under `pkg/conjur/mocks`, `pkg/conjur/api/mocks`, `internal/servicebroker/mocks` must be regenerated via `mockery` (`.mockery.yaml`), never hand-edited, when interfaces change.

### 3. Standard Compliance

- Go version is pinned in **four** `go.mod` files (root, `test/integration`, `test/e2e`, `test/e2e/sample-app`) and **five** Dockerfiles (`Dockerfile`, `Dockerfile.tools`, `Dockerfile.integration`, `scripts/Dockerfile.test`, `test/e2e/Dockerfile.test`) — any Go version bump must update all of them together, or Jenkins' `updatePrivateGoDependencies` check will hard-fail CI.
- `NOTICES.txt` must stay in sync with `go.mod` dependency/license changes (see `CONTRIBUTING.md`'s release checklist, `go-licenses`).
- Pre-commit hooks (`.pre-commit-config.yaml`) run `go generate`, `go build`, `go mod tidy`, `go test`, `go vet`, `gitleaks`, `shellcheck`/`shfmt`, `remarklint` — a PR should pass all of these locally before review.

### 4. Report Generation

Generate `PR_REVIEW_REPORT.md` **fully filled in** — never leave a `<placeholder>` in the final file. Derive each field automatically:

- **One-line summary**: one imperative-mood sentence describing what the diff does.
- **Desired Outcome**: if the branch name or commit messages reference a ticket/issue, pull context from it (e.g. `gh issue view <n>` for a `#<n>` reference); otherwise infer the motivating problem directly from what the diff changes.
- **Implemented Changes**: build this directly from the actual `git diff` and the Multi-Persona Analysis findings (step 2) — what changed, why, and how a reviewer should verify it (manual/integration steps if any). State screenshots as "N/A" unless the diff touches UI.
- **Connected Issue/Story**: scan the current branch name and the commit messages on this branch for a ticket pattern like `[A-Z]+-\d+` (e.g. `CP-26179`, `CNJR-1234`) and use the first match as the "CyberArk internal issue ID". If a GitHub issue number is referenced (`#123` in a branch name or commit message), use it for "Resolves #123" — otherwise omit that line entirely rather than leaving a placeholder.
- **Definition of Done** — tick exactly one box per subsection based on what the diff actually contains, never by default:
  - *Changelog*: tick "CHANGELOG has been updated" only if `CHANGELOG.md` (or equivalent) is part of the diff; tick "does not require a CHANGELOG update" only if the change has no user-facing effect.
  - *Test coverage*: tick "includes new tests" only if test files changed in the diff; tick "does not require tests" only if the change is non-functional (docs, config, tooling, skills).
  - *Documentation*: tick "docs were updated" only if README/docs files changed; tick "does not require updating documentation" only if nothing user-facing changed.
  - *Behavior*: tick "no behavior was changed" only if the diff is limited to non-runtime files (skills, docs, CI config, tests); otherwise tick one of the other two options and explain why.
  - *Security*: tick "no security aspects" only if the Security persona (step 2) genuinely found nothing; otherwise tick the appropriate other option and reference the finding.

If a field genuinely cannot be determined, write a short explicit note (e.g. "No linked issue found in branch/commit history") instead of leaving an angle-bracket placeholder.

The output must follow this structure when pasted into a PR description:

```markdown
<filled summary>

### Desired Outcome

<filled>

### Implemented Changes

- <filled bullets>

### Connected Issue/Story

Resolves #<n>  (omit this line if none found)

CyberArk internal issue ID: <filled or note>

### Definition of Done

#### Changelog
- [ ] The CHANGELOG has been updated, or
- [ ] This PR does not include user-facing changes and doesn't require a CHANGELOG update

#### Test coverage
- [ ] This PR includes new unit and integration tests to go with the code changes, or
- [ ] The changes in this PR do not require tests

#### Documentation
- [ ] Docs (e.g. READMEs) were updated in this PR
- [ ] A follow-up issue to update official docs has been filed here: [insert issue ID]
- [ ] This PR does not require updating any documentation

#### Behavior
- [ ] This PR changes product behavior and has been reviewed by a PO, or
- [ ] These changes are part of a larger initiative that will be reviewed later, or
- [ ] No behavior was changed with this PR

#### Security
- [ ] Security architect has reviewed the changes in this PR,
- [ ] These changes are part of a larger initiative with a separate security review, or
- [ ] There are no security aspects to these changes
```

(the checkboxes above are the template shape — in the actual generated report, exactly one box per subsection must be ticked `[x]`, the rest `[ ]`)

## Usage Instructions

Run this command before every commit that finalizes a task:

```
@ai-pr-review run analysis and generate report
```

## Port.io Metadata

- **Maturity Level**: Tier 2 (Repeatable Skill)
- **Service**: Secrets Manager Integrations

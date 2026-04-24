# AGENTS.md

## Project Overview

This repository is the `tfm` project.

- The repository root is the active Go module (`github.com/hashicorp-services/tfm`).
- `tfm` is a Go CLI for Terraform Cloud / Terraform Enterprise migration and administration workflows.
- Run Go-aware tooling (`go test`, `go vet`, `gopls`) from the repository root.

## Repository Layout

- `./cmd` - Cobra command tree
- `./tfclient` - Terraform Cloud / Enterprise client helpers
- `./vcsclients` - VCS integration clients
- `./output` - output writers/formatters
- `./test` - test configs, Terraform fixtures, cleanup helpers
- `./site` - MkDocs documentation site
- `./.github/workflows` - CI, release, docs, and e2e workflows
- `./.github/agents` - custom Copilot agent definitions (SpecKit, SDD, Terraform execution)
- `./.github/prompts` - prompt templates used by agent workflows
- `./.github/skills` - reference skill packs (Go and Git workflow)
- `./.agents/skills` - repo-local agent skills
- `./.mcp.json` - local MCP server definitions for this workspace
- `./tmp` - workspace-local scratch directory for temporary files

## Temporary Files and Scratch Work

All temporary work files, generated notes, downloaded artifacts, scratch outputs, and agent state that should stay local to this workspace must go under `./tmp`.

Rules:

- Always use `./tmp` for temporary files.
- Never use `/tmp` for this repository.
- Prefer task-specific subdirectories such as `./tmp/build`, `./tmp/review`, or `./tmp/artifacts`.
- Do not commit temporary files unless explicitly requested.

## Local MCP Servers

The root `./.mcp.json` currently defines these local MCP servers:

### terraform

- Runs via Docker image `hashicorp/terraform-mcp-server:0.4.0`
- Expects `TFE_TOKEN` in the environment
- Sets `ENABLE_TF_OPERATIONS=true`
- Starts with `--toolsets=all`

Use this for Terraform module/provider documentation and Terraform-aware MCP operations.

### gopls

- Runs as `gopls mcp`

Use this for Go workspace-aware code navigation and semantic operations.

## Repo-Local Skills and Capabilities

### Repo-local skills (`./.agents/skills`)

Current skill directories include:

- `acquire-codebase-knowledge`
- `cli-developer`
- `code-documenter`
- `code-reviewer`
- `conventional-commit`
- `create-architectural-decision-record`
- `create-github-action-workflow-specification`
- `create-github-issue-feature-from-specification`
- `create-github-issues-feature-from-implementation-plan`
- `create-github-issues-for-unmet-specification-requirements`
- `create-github-pull-request-from-specification`
- `create-implementation-plan`
- `create-specification`
- `debugging-wizard`
- `documentation-writer`
- `generate-custom-instructions-from-codebase`
- `gh-cli`
- `git-commit`
- `git-flow-branch-creator`
- `github-issues`
- `golang-pro`
- `repo-story-time`
- `review-and-refactor`
- `test-master`
- `the-fool`

Use these first when they match the task.

### Bundled GitHub skills (`./.github/skills`)

- `go-development` with references under `./.github/skills/go-development/references`
- `git-workflow` with references under `./.github/skills/git-workflow/references`

### Bundled agent workflows (`./.github/agents`)

Available agent families include:

- SpecKit: `speckit.specify`, `speckit.plan`, `speckit.tasks`, `speckit.analyze`, `speckit.implement`, `speckit.clarify`, `speckit.checklist`, `speckit.constitution`, `speckit.taskstoissues`
- SDD: `sdd-specify`, `sdd-clarify`, `sdd-research`, `sdd-plan-draft`, `sdd-tasks`, `sdd-analyze`, `sdd-checklist`
- Terraform execution helpers: `tf-deployer`, `tf-task-executor`, `tf-report-generator`

## Setup and Development Commands

Run these from repository root.

### Basic setup

- Install or refresh module dependencies: `go mod download`
- Tidy module files after dependency changes: `go mod tidy`

### Build

- Local build with version metadata: `make build-local` or `task build`
- Direct run for a quick version check: `go run . -v`

### Formatting

- Format Go code: `make format` or `task fmt`
- Standard Go formatting: `go fmt ./...`

### Dependency updates

- Update Go dependencies and tidy: `make update` or `task update`

### Docs

- Serve docs locally: `mkdocs serve -f site/mkdocs.yml`

## Testing Instructions

### Go tests

- Run all Go tests: `nocorrect go test -race -v ./...` or `task test`
- Run tests for a specific package: `go test -v ./cmd/...`
- Run a specific test name: `go test -v ./... -run TestName`

### E2E and integration context

End-to-end testing is primarily driven by GitHub Actions and external Terraform Cloud / Enterprise resources.

- Main workflow: `./.github/workflows/end-to-end-test.yml`
- Test configs: `./test/configs`
- Terraform fixtures: `./test/terraform`
- Cleanup helpers: `./test/cleanup`

Notes:

- E2E flows require external credentials, org configuration, and GitHub Actions secrets.
- Do not assume e2e tests can run successfully on a local machine without those dependencies.
- Be careful with cleanup and destructive commands in `./test/cleanup`.

## Code Style and Conventions

### Go

- Go version is `1.23`
- Toolchain is `go1.23.3`
- CLI framework: Cobra
- Configuration: Viper
- Main package layout is command-oriented under `./cmd`

General guidance:

- Keep changes idiomatic and consistent with existing Cobra/Viper patterns.
- Prefer small, focused package changes.
- Keep exported identifiers documented when appropriate.
- Reuse existing helpers in `cmd/helper` and shared clients before adding new abstractions.

### Terraform and HCL

- Terraform files exist under `./.github/terraform` and `./test/terraform`.
- Keep Terraform formatted and validated when editing infra files.

### Documentation

- Main docs live under `./site/docs`.
- Command docs are under `./site/docs/commands`.
- When CLI behavior changes, review matching docs pages and `./README.md`.

## Build, Release, and CI

Primary workflows live in `./.github/workflows`:

- `build.yml` - builds artifacts with GoReleaser
- `release.yml` - release automation
- `docs-deploy.yml` - MkDocs publishing
- `end-to-end-test.yml` - scheduled/manual e2e migration testing
- `jira-issues.yml` - Jira issue automation

Current CI uses GoReleaser for builds. If you modify build or release behavior, review both:

- `./.github/workflows/build.yml`
- `./.goreleaser.yaml`

## Pull Request and Commit Guidance

From `./.github/workflows/README.md`, PR titles should follow:

```text
type(scope): Subject
```

Accepted types documented in the repo:

- `fix`
- `feat`
- `docs`
- `ci`
- `chore`

## Security and Secrets

- Never commit tokens, credentials, or Terraform state.
- `tfm` expects sensitive source/destination credentials via environment variables or secure secret stores.
- The Terraform MCP server requires `TFE_TOKEN`.
- GitHub Actions e2e flows depend on repository secrets and external Terraform org setup.
- Treat `test/configs`, workflow env vars, and CI secret names as sensitive operational context.

## Agent Workflow Notes

- Start at repository root for most work (code, tests, tooling, docs, workflows).
- Prefer repo-local skills and instructions before falling back to generic behavior.
- Use `./.github/copilot-instructions.md` as the primary repo instruction source.
- When creating scratch notes, logs, generated outputs, or downloaded artifacts, place them in `./tmp` and never in `/tmp`.

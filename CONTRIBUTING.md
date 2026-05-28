# Contributing to tfm

Thank you for your interest in contributing to `tfm` — the HashiCorp Implementation Services CLI for Terraform Cloud and Terraform Enterprise migrations.

> [!NOTE]
> This CLI is maintained by [`hashicorp-services/team-advanced-services`](https://github.com/orgs/hashicorp-services/teams/team-advanced-services).
> It does not have official HashiCorp support, but the code owners will work with partners and interested parties to provide assistance where possible.

## Contents

- [Prerequisites](#prerequisites)
- [Getting started](#getting-started)
- [Project layout](#project-layout)
- [Development workflow](#development-workflow)
- [Testing](#testing)
- [Documentation](#documentation)
- [Pull request process](#pull-request-process)
- [Architecture Decision Records](#architecture-decision-records)
- [AI-assisted development](#ai-assisted-development)
- [Code of conduct](#code-of-conduct)
- [License](#license)

---

## Prerequisites

| Tool | Minimum version | Install |
|---|---|---|
| Go | 1.25 | https://go.dev/dl/ |
| [Task](https://taskfile.dev) (`task`) | 3.x | `brew install go-task` / https://taskfile.dev/installation/ |
| golangci-lint | latest | `task tools:install` |
| staticcheck | latest | `task tools:install` |
| govulncheck | latest | `task tools:install` |
| goimports | latest | `task tools:install` |
| Terraform CLI | ≥ 1.6 | https://developer.hashicorp.com/terraform/downloads |
| GitHub CLI (`gh`) | latest | https://cli.github.com |

Install all optional Go dev tools in one step:

```bash
task tools:install
```

Verify they are all available:

```bash
task tools:check
```

---

## Getting started

```bash
# 1. Fork and clone
git clone https://github.com/<your-fork>/tfm.git
cd tfm

# 2. Install dependencies
go mod download

# 3. Build a local binary
task build          # outputs ./tfm

# 4. Confirm the binary works
./tfm --version
```

### Configuration

`tfm` reads credentials and org settings from `.tfm.hcl` or environment variables.
Copy the provided template and fill in your values:

```bash
cp .env.example .env
# edit .env, then source it:
set -a && . ./.env && set +a
```

See the [configuration reference](https://hashicorp-services.github.io/tfm/configuration_file/config_file/) for all available keys.
**Never commit `.env` or `.tfm.hcl` files containing real tokens.**

---

## Project layout

```
tfm/
├── cmd/              # Cobra command tree (copy, list, delete, generate, core…)
│   ├── copy/         # Copy-resource commands
│   ├── list/         # List-resource commands
│   ├── generate/     # Config-generation commands
│   ├── logging/      # Logging initialisation (TFM_LOG / --verbose)
│   └── helper/       # Shared Viper helpers and utilities
├── tfclient/         # Terraform Cloud/Enterprise client setup
├── vcsclients/       # GitHub and GitLab VCS client helpers
├── output/           # Output formatting (table, JSON, deferred)
├── version/          # Build-time version metadata
├── test/             # End-to-end test assets (see Testing below)
│   ├── terraform/    # Terraform configs to provision source test resources
│   ├── configs/      # Pre-built .tfm.hcl configs for CI scenarios
│   ├── cleanup/      # Scripts to tear down test resources
│   └── state/        # State-migration test fixtures
├── site/             # MkDocs documentation site
│   └── docs/
│       ├── commands/ # One page per command family
│       └── …
├── ADR/              # Architecture Decision Records
├── .github/
│   ├── workflows/    # CI, release, docs-deploy, end-to-end-test
│   ├── actions/      # Composite actions for e2e test scenarios
│   ├── prompts/      # GitHub Copilot prompt files
│   └── agents/       # Copilot agent (Speckit) workflow files
├── .specify/         # Spec-driven development templates and scripts
├── Taskfile.yml      # Primary developer task runner (see task --list)
├── Makefile          # Legacy make targets (build-local, format, update)
└── .goreleaser.yaml  # GoReleaser release configuration
```

---

## Development workflow

### Task runner

`task` (go-task) is the primary tool for day-to-day development.
Run `task` with no arguments to list all available tasks:

```bash
task              # list all tasks
task help         # show common workflow patterns
```

Key tasks:

| Task | What it does |
|---|---|
| `task build` | Build `./tfm` with dev ldflags |
| `task test` | Full test suite with race detection |
| `task test:fast` | Tests without race detection (faster iteration) |
| `task test:cover` | Tests with HTML coverage report |
| `task lint` | Run golangci-lint |
| `task fmt` | Format code with `go fmt` |
| `task ci` | Full CI pipeline (fmt-check → vet → lint → vuln → test) |
| `task precommit` | Pre-commit checks (fmt → vet → lint → fast tests) |
| `task update` | Update all Go dependencies + `go mod tidy` |
| `task clean` | Remove build artefacts |
| `task changelog` | Show git log since last tag |

Run `task precommit` before pushing any branch.

### Branching

This project uses GitFlow:

```
feature/<issue-number>-short-description   # new features
fix/<issue-number>-short-description        # bug fixes
release/<version>                           # release branches
hotfix/<version>-<description>              # emergency fixes
```

Branch from `main` for features and fixes. Keep commits focused and atomic.

### Commit messages

Follow [Conventional Commits](https://www.conventionalcommits.org/):

```
type(scope): short description

Optional body explaining the why.

Closes #NNN
```

Accepted types: `feat`, `fix`, `docs`, `ci`, `chore`, `refactor`, `perf`, `test`

Breaking changes must include `BREAKING CHANGE:` in the footer or `!` after the type:

```
feat(copy)!: rename --side flag to --context
```

---

## Testing

### Unit tests

```bash
task test             # race-detection enabled (recommended before PR)
task test:fast        # faster, no race detection
task test:run -- TestMyFunction ./cmd/copy   # run a single test
task test:cover       # generate coverage.html
```

### End-to-end (e2e) tests

E2E testing requires live Terraform Cloud organisations. The full e2e suite runs automatically via the scheduled GitHub Actions workflow (`.github/workflows/end-to-end-test.yml`). The workflow:

1. Runs the Terraform configuration in `test/terraform/` against the `gh-actions-ci-master-workspace` workspace in the `tfm-testing-source` TFC organisation to provision all required source resources (workspaces, SSH keys, variable sets, agent pools, VCS connections, teams).
2. Runs all composite actions under `.github/actions/test-*` to exercise `tfm` commands.
3. Copies resources to `tfm-testing-destination` using `tfm`.
4. Cleans up all copied resources from the destination org.
5. Runs `terraform destroy` from `gh-actions-ci-master-workspace` to remove source test resources.

#### Running e2e tests locally

To provision source resources in your own TFC organisation before running local tests:

```bash
# 1. Set workspace variables — see test/terraform/README.md for required variables
cd test/terraform

# 2. Review the Terraform config
terraform init
terraform plan

# 3. Apply (creates workspaces, teams, variable sets, etc. in your source org)
terraform apply

# 4. Build and run tfm against the provisioned resources
cd ../..
task build
./tfm list workspaces --config test/configs/.e2e-all-workspaces-test.hcl
```

> See [`test/terraform/README.md`](./test/terraform/README.md) and [`.github/workflows/README.md`](./.github/workflows/README.md) for full details on required API tokens, workspace variables, and test org setup.

#### API tokens required for e2e

| Secret / Variable | Purpose |
|---|---|
| `SOURCETOKEN` | GitHub Actions secret — source org owner team token |
| `DESTINATIONTOKEN` | GitHub Actions secret — destination org owner team token |
| `GHTOKEN` | GitHub Actions secret — GitHub token for cloning test repos |
| `GHORGANIZATION` | GitHub Actions secret — GitHub org name |
| `GHUSERNAME` | GitHub Actions secret — GitHub username |
| `TFE_TOKEN` | `gh-actions-ci-master-workspace` workspace variable |

Refer to [`.github/workflows/README.md`](./.github/workflows/README.md) when tokens need to be rotated.

---

## Documentation

`tfm` uses [MkDocs](https://www.mkdocs.org/) with Material theme. The docs site is published to [hashicorp-services.github.io/tfm](https://hashicorp-services.github.io/tfm) from the `main` branch via `.github/workflows/docs-deploy.yml`.

### When to update docs

**Update `site/docs/` whenever you:**

- Add, rename, or remove a CLI command or subcommand
- Add, change, or remove a config key in `.tfm.hcl`
- Add or change any `--flag` on an existing command
- Change default behaviour or output format
- Add a new environment variable

**Files typically affected:**

| Change | Docs file(s) |
|---|---|
| New command | `site/docs/commands/<command-family>.md` (new or updated) |
| New config key | `site/docs/configuration_file/config_file.md` |
| New `generate` subcommand | `site/docs/commands/generate/` |
| Logging / verbosity | `site/docs/logging.md` |
| Copy/migration flow | `site/docs/migration/` |

### Serve docs locally

```bash
# Requires: pip install mkdocs-material
mkdocs serve -f site/mkdocs.yml
```

Open http://127.0.0.1:8000 to preview changes.

### Adding a new command page

1. Create `site/docs/commands/<your-command>.md`.
2. Add an entry in `site/mkdocs.yml` under the relevant `nav` section.
3. Follow the structure of an existing command page (synopsis → flags → config keys → examples).

---

## Pull request process

1. **Open an issue first** for anything beyond trivial fixes. This avoids duplicate work and helps the team plan.

2. **Create a branch** from `main` following the naming conventions above.

3. **Run the pre-commit check** before pushing:

   ```bash
   task precommit
   ```

4. **Open a PR** with a title that follows Conventional Commits format:

   ```
   feat(list): add tfm list variable-sets command
   fix(copy): handle 404 when source workspace has no state
   docs(site): update config_file.md with new varsets-map key
   ```

5. **PR description checklist:**
   - [ ] Issue linked (e.g., `Closes #NNN`)
   - [ ] `task ci` passes locally
   - [ ] `site/docs/` updated for any command, flag, or config change
   - [ ] Unit tests added or updated
   - [ ] No credentials or `.tfm.hcl` files committed
   - [ ] ADR linked or created if the change involves a significant architectural decision

6. **One approval** from `hashicorp-services/team-advanced-services` is required before merge.

7. **Merge strategy:** squash merge for feature branches; merge commit for release branches.

---

## Architecture Decision Records

Significant design decisions are captured in the [`ADR/`](./ADR/) directory using a lightweight ADR format. An index is maintained at [`ADR/index.md`](./ADR/index.md).

**Create an ADR when you are:**
- Choosing between two non-obvious implementation approaches
- Changing the config file format or adding new map keys
- Adding a new API client or VCS provider
- Changing the output format in a way that affects machine-readable consumers
- Making a change that deprecates existing behaviour

Use [`ADR/ADR-template.md`](./ADR/ADR-template.md) as the starting point.

---

## AI-assisted development

This repository includes tooling to assist AI-agent-driven development workflows.

### GitHub Copilot

- `.github/copilot-instructions.md` — project-level Copilot instructions (Go standards, GitFlow, testing conventions)

#### Instruction files (`.github/instructions/`)

These files are automatically applied by GitHub Copilot based on the file type you are editing.
They encode project-specific coding standards and must be read before making changes to the relevant file types.

| File | Applied to | Purpose |
|---|---|---|
| `go.instructions.md` | `**/*.go`, `go.mod`, `go.sum` | Idiomatic Go standards: error handling, naming, generics, type safety |
| `go-mcp.instructions.md` | Go workspace sessions | gopls MCP server usage — symbol search, diagnostics, reference finding |
| `github-actions-ci-cd-best-practices.instructions.md` | `.github/workflows/*.yml` | Secure, efficient GitHub Actions patterns |
| `terraform.instructions.md` | `**/*.tf` | Terraform conventions for files in `test/terraform/` and `.github/terraform/` |
| `markdown-accessibility.instructions.md` | `**/*.md` | Accessible markdown following GitHub's five best practices |
| `update-docs-on-code-change.instructions.md` | All code files | Reminder to update `README.md` and `site/docs/` when application code changes |

If you are making a change that touches any of these file types, read the corresponding instruction file first.

#### Prompt files (`.github/prompts/`)

Ready-to-use Copilot agent prompts:

  | Prompt | Purpose |
  |---|---|
  | `dependabot-pr-review.prompt.md` | Review and comment on all open Dependabot PRs |
  | `go-tfe-review.prompt.md` | Analyse a go-tfe version bump for breaking changes |
  | `tf-go-release.prompt.md` | Review changes since last tag and propose/create a semver release |
  | `speckit.*.prompt.md` | Spec-driven feature development workflow (Speckit) |

### Speckit (spec-driven development)

`.specify/` contains templates and scripts for a structured feature development workflow:

```bash
.specify/
├── memory/constitution.md          # Project principles and constraints
├── scripts/bash/create-new-feature.sh
└── templates/                      # Spec, plan, checklist, and task templates
```

`.github/agents/speckit.*.agent.md` provides the matching Copilot agent definitions.

### AGENTS.md

[`AGENTS.md`](./AGENTS.md) at the repository root is the authoritative guide for agent workflows, skill locations, and workspace conventions. Read it before running any agent-assisted task.

---

## Code of conduct

This project follows the [HashiCorp Community Guidelines](https://www.hashicorp.com/community-guidelines). Please be respectful and constructive in all interactions.

---

## License

`tfm` is licensed under the [Mozilla Public License 2.0](./LICENSE).
By contributing, you agree that your contributions will be licensed under the same terms.

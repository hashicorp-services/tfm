---
description: "Review all open Dependabot PRs in the tfm project: checkout each branch, build, run unit tests and a live TFC smoke test, then post a structured review comment with reasoning and acceptance criteria."
mode: agent
tools: ["runCommands", "editFiles", "codebase", "githubRepo"]
---

# Dependabot PR Review — tfm

You are a senior Go developer and release manager with deep CI/CD expertise reviewing dependency upgrade pull requests for the `tfm` project — a Go CLI for Terraform Cloud/Enterprise migrations.

Before doing anything else, read `./AGENTS.md` for workspace conventions (scratch files go in `./tmp`, never `/tmp`) and apply the guidance in:
- `.github/instructions/go.instructions.md`
- `.github/instructions/github-actions-ci-cd-best-practices.instructions.md`

---

## Scope

If a specific PR number has been provided (e.g. `${input:prNumber:optional: single PR number to review}`), operate on that PR only. Otherwise, operate on **all currently open Dependabot pull requests** in `hashicorp-services/tfm`.

---

## Process

Work through each PR in the following order. Complete all steps for one PR before moving to the next.

### Step 1 — Discover PRs

```bash
gh pr list --repo hashicorp-services/tfm \
  --author app/dependabot \
  --state open \
  --json number,title,headRefName,baseRefName,url \
  --jq '.[] | "\(.number) \(.headRefName) \(.title)"'
```

Record each PR number, branch name, and title.

### Step 2 — Check out the branch

Use a git worktree so multiple branches can coexist without disrupting the working tree:

```bash
git fetch origin
git worktree add ./tmp/${PR_NUMBER} origin/${BRANCH_NAME}
```

All subsequent commands for this PR run from `./tmp/${PR_NUMBER}/tfm`.

### Step 3 — Build

```bash
cd ./tmp/${PR_NUMBER}/tfm
go build -buildvcs=false ./...
```

- Pass `-buildvcs=false` to suppress VCS-stamping errors in detached worktrees.
- Record: ✅ PASS or ❌ FAIL with the first error line.

### Step 4 — Unit tests

```bash
cd ./tmp/${PR_NUMBER}/tfm
go test -buildvcs=false ./... 2>&1 | tail -30
```

- The pre-existing `cmd/nuke/nuke.go` failure (missing `package` declaration) is a **known issue on all branches including `main`** — do **not** count it as a regression introduced by this PR.
- Record: number of packages tested, any new failures vs `main`.

### Step 5 — Live TFC smoke test

Build the binary and run a read-only list command against the source TFC organisation using credentials from the environment:

```bash
cd ./tmp/${PR_NUMBER}/tfm
go build -buildvcs=false -o ./tmp/tfm-pr${PR_NUMBER} .
./tmp/tfm-pr${PR_NUMBER} list workspaces \
  --config ./tmp/.tfm.hcl 2>&1 | head -20
```

Expected: command exits 0 and returns a workspace list.

> **Note:** `./tmp/.tfm.hcl` must already exist with valid `src_tfe_*` credentials. If it is absent, skip this step and note it in the review comment. Never create or commit credential files.

### Step 6 — Dependency change analysis

For each PR, check what changed in `go.mod` / `go.sum`:

```bash
git -C ./tmp/${PR_NUMBER}/tfm diff origin/main -- go.mod
```

Identify:
- Package name(s) updated
- Old version → new version
- Whether this is a minor, patch, or pre-release bump
- Any known breaking changes (check the package's CHANGELOG or GitHub releases)

### Step 7 — Post PR review comment

Post a **single structured comment** to the PR using:

```bash
gh pr comment ${PR_NUMBER} --repo hashicorp-services/tfm --body "$(cat <<'COMMENT'
<comment content>
COMMENT
)"
```

The comment **must** follow this template:

````markdown
## 🤖 Dependabot PR Review — Automated

**PR:** #${PR_NUMBER} — ${TITLE}
**Branch:** `${BRANCH_NAME}`
**Reviewed:** $(date -u +"%Y-%m-%dT%H:%M:%SZ")

### Dependency Change
| | Value |
|---|---|
| Package | `<package>` |
| Old version | `<old>` |
| New version | `<new>` |
| Bump type | patch / minor / major |

### Test Results
| Check | Result |
|---|---|
| `go build ./...` | ✅ PASS / ❌ FAIL |
| `go test ./...` | ✅ PASS / ⚠️ PRE-EXISTING FAILURES ONLY / ❌ NEW FAILURES |
| Live TFC smoke test | ✅ PASS / ⚠️ SKIPPED (no config) / ❌ FAIL |

### Analysis
<2–4 sentences explaining what the dependency update contains, any breaking changes
found in the upstream CHANGELOG, and whether the changes affect tfm functionality.>

### Acceptance Criteria
- [ ] `go build ./...` passes with no new errors
- [ ] `go test ./...` introduces no new test failures beyond pre-existing `cmd/nuke` issue
- [ ] Live TFC smoke test exits 0 (or skipped with documented reason)
- [ ] No breaking API changes that affect tfm's usage of this package
- [ ] `go.sum` updated correctly

### Recommendation
**✅ APPROVE — safe to merge** / **⚠️ NEEDS ATTENTION** / **❌ DO NOT MERGE**

<One sentence rationale.>
````

### Step 8 — Record results locally

Append a row to `./tmp/dependabot-review-$(date +%Y-%m-%d).md`, creating the file with a header if it does not yet exist:

```markdown
# Dependabot PR Review — $(date +%Y-%m-%d)

| PR | Title | Build | Tests | Smoke | Recommendation |
|---|---|---|---|---|---|
| #NNN | title | ✅/❌ | ✅/⚠️/❌ | ✅/⚠️/❌ | APPROVE / ATTENTION / REJECT |
```

---

## Constraints

- **Never** use `/tmp` — all scratch work goes in `./tmp/`.
- **Never** commit `.tfm.hcl`, `.env`, or any file containing credentials.
- **Never** run write/copy/delete commands against TFC organisations during this review — read-only `list` commands only.
- Worktrees in `./tmp/${PR_NUMBER}/` are **intentionally preserved** for manual inspection after the review completes.
- If the `gh` CLI is not authenticated, stop and inform the user before proceeding.

---

## Success Criteria

The prompt is complete when:
1. Every in-scope Dependabot PR has a review comment posted.
2. The local summary file `./tmp/dependabot-review-<date>.md` exists and contains one row per PR.
3. No credentials have been written to tracked files.

---
description: "Review all go-tfe releases since the version in tfm's go.mod, analyse breaking changes, new features, deprecations and security fixes against tfm's actual API usage, then create labelled GitHub issues in hashicorp-services/tfm for each actionable finding."
mode: agent
tools: ["fetch", "githubRepo", "runCommands", "editFiles", "codebase"]
---

# go-tfe Release Review — tfm

You are a senior Go developer and release manager with deep expertise in the Terraform Cloud REST API, the `go-tfe` client library, and the `tfm` migration CLI. You understand how `go-tfe` structs, interfaces, and method signatures map to real TFC/TFE API behaviour, and you can assess the downstream impact of upstream changes on a consumer codebase.

Before doing anything else, read `./AGENTS.md` for workspace conventions (scratch files go in `./tmp`, never `/tmp`) and apply the guidance in:
- `.github/instructions/go.instructions.md`
- `.github/instructions/github-actions-ci-cd-best-practices.instructions.md`

---

## Scope

Analyse every `go-tfe` release published **after** the version currently pinned in `tfm/go.mod`, up to and including the latest release.

If the user has provided a specific version range as input (`${input:fromVersion:optional: override start version e.g. v1.50.0}`), use that as the lower bound instead of auto-detecting from `go.mod`.

---

## Step 1 — Detect current go-tfe version

Read the pinned version from `tfm/go.mod`:

```bash
grep 'github.com/hashicorp/go-tfe' ./tfm/go.mod
```

Record the current version (e.g. `v1.103.0`). All releases **after** this version are in scope.

---

## Step 2 — Fetch go-tfe releases

Retrieve the full list of go-tfe releases newer than the current pinned version:

```
GET https://api.github.com/repos/hashicorp/go-tfe/releases?per_page=100
```

Filter to releases with a tag version greater than the current pinned version, ordered oldest-first.

Also fetch the raw CHANGELOG for narrative context:

```
GET https://raw.githubusercontent.com/hashicorp/go-tfe/main/CHANGELOG.md
```

If the CHANGELOG is unavailable, rely on the GitHub release notes from the API response.

---

## Step 3 — Inventory tfm's go-tfe API usage

Search the `tfm` codebase to identify every go-tfe type, method, and field currently in use:

```bash
grep -rh "tfe\." ./tfm --include="*.go" | grep -oE 'tfe\.[A-Za-z]+' | sort -u
```

Also list all go-tfe imports:

```bash
grep -rh '"github.com/hashicorp/go-tfe"' ./tfm --include="*.go" -l
```

Build a reference set of:
- **Types used**: e.g. `tfe.WorkspaceCreateOptions`, `tfe.Project`
- **Interfaces used**: e.g. `tfe.Workspaces`, `tfe.Projects`
- **Methods called**: e.g. `.Workspaces.List(...)`, `.Projects.Create(...)`

This set is used in Step 4 to assess impact.

---

## Step 4 — Analyse each release

For each in-scope release, working oldest-to-newest, classify every changelog entry into one of these categories:

| Category | Label(s) to apply | Create issue? |
|---|---|---|
| **Breaking change** — removed/renamed type, method, field, or changed signature | `breaking-change`, `dependency` | ✅ Always |
| **New API feature** — new resource type, method, or option that tfm could leverage | `enhancement`, `feature request` | ✅ If tfm doesn't already use it |
| **Deprecation** — method or field marked deprecated | `dependency` | ✅ Always |
| **Security fix** — CVE, auth bypass, data exposure | `dependency`, `breaking-change` | ✅ Always |
| **Internal/test-only change** — no public API impact | — | ❌ Skip |
| **Documentation only** | — | ❌ Skip |

For each **Breaking change** or **Deprecation**, cross-reference against the tfm API usage set from Step 3. Note whether tfm is directly affected.

---

## Step 5 — Write local analysis report

Before creating any issues, write a full analysis to:

```
./tmp/go-tfe-review-$(date +%Y-%m-%d).md
```

Structure:

```markdown
# go-tfe Release Review — <date>

**Current tfm version:** <version from go.mod>
**Latest go-tfe version:** <latest release tag>
**Releases analysed:** <count> (<from> → <to>)

## Summary

| Category | Count |
|---|---|
| Breaking changes | N |
| New features (actionable) | N |
| Deprecations | N |
| Security fixes | N |
| Skipped (internal/docs) | N |

## Findings

### [v1.X.Y] — <release date>

#### Breaking Changes
- **`TypeOrMethod`** — <description>. tfm impact: <AFFECTED / NOT AFFECTED — reason>

#### New Features
- **`NewType/Method`** — <description>. tfm opportunity: <what command or capability this could enable>

#### Deprecations
- **`DeprecatedThing`** — <description>. tfm impact: <AFFECTED / NOT AFFECTED>

#### Security Fixes
- <description>

---
(repeat per release)

## Issues to Create

| # | Title | Labels | Duplicate check |
|---|---|---|---|
| 1 | feat: support X via go-tfe vY | enhancement, feature request | Not found |
| 2 | fix: breaking change in tfe.Y.Z — update tfm | breaking-change, dependency | Not found |
```

---

## Step 6 — Deduplicate against existing issues

Before creating any issue, search the tfm repo for an existing open or closed issue covering the same finding:

```bash
gh issue list --repo hashicorp-services/tfm --state all \
  --search "<key term from finding title>" \
  --json number,title,state | head -5
```

If a sufficiently similar issue already exists, note it in the report and skip creation. Err on the side of creating rather than skipping — only skip if there is a clear title/content match.

---

## Step 7 — Create GitHub issues

For each finding not already covered by an existing issue, create a labelled issue in `hashicorp-services/tfm`:

```bash
gh issue create \
  --repo hashicorp-services/tfm \
  --title "<title>" \
  --label "<comma-separated labels>" \
  --body "<body>"
```

### Issue title format

| Category | Title prefix |
|---|---|
| Breaking change | `fix(go-tfe): breaking change — <description>` |
| New feature | `feat(go-tfe): <new capability> (available since go-tfe <version>)` |
| Deprecation | `chore(go-tfe): deprecation — <description>` |
| Security fix | `fix(go-tfe): security fix — <description>` |

### Issue body template

```markdown
## Summary

<1–3 sentence description of the go-tfe change and why it matters to tfm.>

## go-tfe Release
- **Version:** <vX.Y.Z>
- **Release date:** <date>
- **Release notes:** <link to GitHub release>

## Change Details

<Copy or paraphrase the relevant CHANGELOG entry.>

## Impact on tfm

<Describe which tfm commands, files, or types are affected. If not yet affected, describe the opportunity.>

**Affected files (if any):**
- `tfm/<path>`

## Acceptance Criteria

- [ ] <Specific, testable criterion 1>
- [ ] <Specific, testable criterion 2>
- [ ] `go build ./...` passes after changes
- [ ] `go test ./...` passes with no new failures

## References

- [go-tfe release](<release URL>)
- [go-tfe CHANGELOG](https://github.com/hashicorp/go-tfe/blob/main/CHANGELOG.md)
```

---

## Step 8 — Final summary

After all issues are created, print a concise summary to the terminal:

```
go-tfe Review Complete
======================
Releases analysed : <N> (v<from> → v<to>)
Issues created    : <N>
Issues skipped    : <N> (duplicates)
Report            : ./tmp/go-tfe-review-<date>.md
```

---

## Constraints

- **Never** use `/tmp` — all scratch work goes in `./tmp/`.
- **Never** commit `.tfm.hcl`, `.env`, or credential files.
- **Never** modify `go.mod` or `go.sum` — this prompt analyses only; upgrading is a separate task.
- Use `gh` CLI for all GitHub operations. If `gh` is not authenticated, stop and inform the user.
- If the GitHub API rate limit is hit, pause and inform the user before retrying.

---

## Success Criteria

The prompt is complete when:

1. `./tmp/go-tfe-review-<date>.md` exists with a finding for every in-scope release.
2. A GitHub issue exists in `hashicorp-services/tfm` for every actionable finding not already tracked.
3. Each issue has appropriate labels, a clear title, and populated acceptance criteria.
4. The terminal summary is printed.
5. No credentials or secrets have been written to tracked files.

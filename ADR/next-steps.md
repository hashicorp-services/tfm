# ADR Next Steps

Date: 2026-02-24

## ADR Review Summary

- `0001-initial-adr.md`: Foundational implementation ADR for Go/Cobra/Viper and idempotent copy behavior.
- `0002-deprecated-ce-to-tfc.md`: Deprecated and superseded by `0003-ce-to-tfc.md`.
- `0003-ce-to-tfc.md`: Core CE/OSS to TFC migration decision remains `Proposed` and still drives most migration features.
- `0004-ce-migration-cli-driven-workspaces copy.md`: `Accepted`; defines automated `cloud {}` conversion workflow.
- `0005-ce-migration-variables.md`: Variable migration ADR content exists and should be tracked as active scope.
- `0006-gitlab-vcs.md`: GitLab VCS support ADR draft exists and should be finalized for execution tracking.

## Outstanding Pull Requests (Open)

| PR | Title | Updated | Link |
|---|---|---|---|
| #340 | Bump `github.com/hashicorp/go-tfe` from 1.78.0 to 1.95.0 | 2025-11-17 | https://github.com/hashicorp-services/tfm/pull/340 |
| #338 | Bump `golang.org/x/oauth2` from 0.24.0 to 0.32.0 | 2025-11-10 | https://github.com/hashicorp-services/tfm/pull/338 |
| #337 | Bump `github.com/go-git/go-git/v5` from 5.12.0 to 5.16.3 | 2025-11-05 | https://github.com/hashicorp-services/tfm/pull/337 |
| #328 | Bump `github.com/spf13/cobra` from 1.8.1 to 1.10.1 | 2025-10-31 | https://github.com/hashicorp-services/tfm/pull/328 |
| #329 | Bump `github.com/spf13/pflag` from 1.0.5 to 1.0.10 | 2025-10-31 | https://github.com/hashicorp-services/tfm/pull/329 |
| #301 | Bump the `go_modules` group with 3 updates | 2025-05-07 | https://github.com/hashicorp-services/tfm/pull/301 |

## Next Steps

1. Merge the low-risk CLI dependency pair together (`#328` + `#329`) and validate command behavior (`root`, `copy`, `core`, `list`).
2. Merge Git and auth dependency updates (`#337` + `#338`) after smoke tests for clone/link/auth workflows.
3. Merge `go-tfe` update (`#340`) with focused regression testing on workspace/state/variable operations.
4. Triage stale grouped dependency PR (`#301`): close if superseded by newer PRs, or rebase and validate if still needed.
5. Open implementation PRs tied to ADR scope gaps (not currently represented by open PRs):
   - ADR `0003`: monorepo + CLI workspace migration edge cases.
   - ADR `0004`: branch/PR automation around `cloud {}` conversion workflow.
   - ADR `0006`: finalize GitLab ADR status/numbering and begin VCS-agnostic refactors.

## Notes

- Current open PRs are maintenance/dependency focused; there are no open feature PRs directly advancing CE/OSS migration ADR milestones.

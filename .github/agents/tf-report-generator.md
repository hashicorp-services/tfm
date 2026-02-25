---
name: tf-report-generator
description: Generate deployment reports from workspace and run data
model: Claude Sonnet 4.5 (copilot)
tools: ['execute', 'read', 'edit', 'search', 'terraform-mcp-server/get_run_details', 'terraform-mcp-server/get_workspace_details', 'terraform-mcp-server/list_runs']
---

Use below skills:
- tf-report-template

# tf-report-generator

Generate a comprehensive deployment report from template.

## Input

- `plan.md`, `*.tf` files, git log
- Deployment status and run URL from tf-deployer
- HCP Terraform workspace details

## Output

- Deployment report at `specs/<branch>/reports/deployment_<timestamp>.md`

## Execution Steps

1. Read `.foundations/templates/deployment-report-template.md`
2. Consult `.foundations/design_resources/00-communication-foundations.md` for voice, tone, and documentation standards
3. Collect data: architecture (plan.md), modules (*.tf), git stats, HCP details
4. Fetch workspace and run details via MCP tools
5. Parse security tool output (trivy, vault-radar) if available
6. Replace all `{{PLACEHOLDER}}` tokens with collected data, using clear, concise language per communication standards
7. Use "N/A" for unavailable data — no placeholders may remain
8. Write report file and display path to user

## Constraints

- No `{{PLACEHOLDER}}` may remain in final output
- Document ALL workarounds vs proper fixes
- Include security findings with severity ratings
- Module compliance: percentage of private vs public modules
- Follow `tf-report-template` skill patterns

---
name: sdd-tasks
description: Generate dependency-ordered task breakdowns from planning artifacts. Produces tasks.md with phased checklist format.
model: Claude Sonnet 4.5 (copilot)
tools: ['execute', 'read', 'edit', 'search']
---

Use below skills:
- tf-task-patterns
- terraform-style-guide

# Task Breakdown Generator

Generate an actionable, dependency-ordered task breakdown from planning artifacts. Transforms spec.md and plan.md into a checklist-format tasks.md

## Workflow

1. **Load artifacts**: Read `spec.md`, `plan.md`, `data-model.md`, and `contracts/module-interfaces.md` from the feature directory
2. **Extract stories**: Identify user stories from spec.md with priorities (P1, P2, P3)
3. **Extract architecture**: Pull modules, data model entities, and cross-module data flows from plan.md
4. **Map coverage**: Build requirement → task matrix ensuring every requirement has implementation tasks
5. **Assign phases**: Setup → Foundational → User Stories (priority order) → Polish (per `tf-task-patterns`)
6. **Generate tasks**: Write tasks.md with all required sections (per `tf-task-patterns`)
7. **Validate**: Confirm all requirements covered, constitution 3.2 files present, phase ordering correct

## Output

- **Location**: `specs/{FEATURE}/tasks.md`
- **Format**: Phased checklist following `tf-task-patterns` skill structure
- **Template**: `.foundations/templates/tasks-template.md`

### Required Sections

| Section | Content |
|---------|---------|
| Header | Feature name, input path, prerequisites, tests stance |
| Format explanation | Task ID format, sequential execution note |
| Requirements Coverage Matrix | Requirement → Task(s) → Description |
| Phase sections | Grouped tasks with purpose, checkpoints |
| Dependencies & Execution Order | Phase deps, story deps, cross-module data flow |
| Implementation Strategy | MVP first, incremental delivery |
| File Checklist | File → Task → Purpose |
| Task Summary | Phase → Task range → User Story with total count |

## Constraints

- **Checklist format**: Every task follows `- [ ] T### [US#?] Description with file path`
- **Phase structure**: Setup → Foundational → User Stories → Polish
- **Independent phases**: Each phase must be a complete, independently testable increment
- **Constitution coverage**: Cross-reference constitution 
- **Module wiring**: Each cross-module data flow entry from module-interfaces.md produces a task wiring output to input

## Examples

**Good task** (story phase):
```markdown
- [ ] T008 [US1] Implement CloudFront module in main.tf with OAI creation and S3 origin at /main.tf
```

**Good task** (setup phase):
```markdown
- [ ] T001 Create terraform.tf with Terraform >= 1.7 and AWS provider ~> 5.83 at /terraform.tf
```

**Bad task** (missing ID, label, path):
```markdown
- [ ] Create VPC configuration
```

**Good phase header**:
```markdown
## Phase 3: User Story 1 - Access Static Website via Secure CDN (Priority: P1)

**Goal**: Enable website visitors to access static content through globally distributed CDN

**Independent Test**: Deploy sample HTML and verify HTTPS access via CloudFront URL

**Dependency**: User Story 2 must be implemented together (circular OAI ↔ bucket policy dependency)
```

## Context

$ARGUMENTS

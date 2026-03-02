# Working Agreement — IOrder Agentic Scrum

> Single source of truth for the entire team (human + AI agents).

## 1. Sprint Cadence

| Item | Value |
|---|---|
| Sprint length | **1 week** (Mon → Fri) |
| Planning | Monday 09:00 (async issue) |
| Daily sync | Async — each agent posts update in Daily Issue by 10:00 |
| Review | Friday 14:00 (PR demo + screenshots) |
| Retro | Friday 15:00 (async issue template) |

## 2. Definition of Ready (DoR)

An issue is considered **Ready** when it meets all of the following:

- [ ] Clear title with sufficient context (user story or task description)
- [ ] Acceptance Criteria (AC) clearly listed
- [ ] Labels assigned correctly: `type:*`, `prio:*`, `area:*`
- [ ] Estimate ≤ 1 day (if larger → split the task)
- [ ] Dependencies identified and not blocked
- [ ] Assigned to Sprint Milestone

## 3. Definition of Done (DoD)

An issue is considered **Done** when:

- [ ] Code implements all AC
- [ ] Unit tests pass (coverage does not decrease)
- [ ] Lint + typecheck pass (CI green)
- [ ] PR reviewed with ≥ 1 approval
- [ ] PR description complete (using template)
- [ ] Docs/ADR updated if architecture changes occurred
- [ ] Branch up-to-date with `main`
- [ ] Squash merged into `main`

## 4. Branch Naming Convention

```
<type>/<issue-number>-<short-description>

# Examples:
feat/42-add-health-endpoint
fix/57-fix-login-redirect
chore/61-update-ci-config
docs/65-add-adr-003
```

Types: `feat`, `fix`, `chore`, `docs`, `refactor`, `test`, `ci`

## 5. Commit Convention

Follow [Conventional Commits](https://www.conventionalcommits.org/):

```
<type>(<scope>): <description>

# Examples:
feat(backend): add health check endpoint
fix(frontend): resolve hydration mismatch
chore(ci): pin golangci-lint to v1.61
```

Scopes: `backend`, `frontend`, `mobile`, `ci`, `docs`, `infra`

## 6. PR Rules

- **All changes must go through a PR** — no direct pushes to `main`.
- PR title = commit message (squash merge).
- PR body must use the template at `.github/PULL_REQUEST_TEMPLATE.md`.
- CI must be green (all required checks pass) before merging.
- Minimum 1 approval required.
- Squash merge recommended.
- Branch must be up-to-date with `main` before merging.

## 7. WIP (Work In Progress) Limits

| Column | WIP Limit |
|---|---|
| In Progress | **1 per agent** |
| In Review | **2 per agent** |

Each agent picks only **1 task** at a time. Complete or unblock before picking the next.

## 8. Code Review Guidelines

- Review within **4 hours** (business hours).
- Focus: correctness, security, tests, naming, edge cases.
- Use GitHub Suggestions for minor changes.
- Approve or Request Changes — do not leave "Comment only" hanging.
- Cross-platform review: Backend changes affecting API → tag FE/Mobile reviewers.

## 9. Handoff Rules

| From | To | Trigger |
|---|---|---|
| PO | Dev | Issue moved to **Ready** column |
| Dev | QA/Review | PR opened, moved to **In Review** |
| QA | Dev | Request Changes on PR |
| QA | Release | PR approved + CI green → merge → **Done** |
| Release | PO | Release notes published |

## 10. Reporting Format

All agents update using this format:

```markdown
### Status Update — [Agent Role] — [Date]

**Context:** [What were you working on]
**Plan:** [What you planned to do]
**Result:** [What you accomplished]
**Blockers:** [Any blockers, or "None"]
**Next:** [What you will do next]
```

## 11. Labels

| Category | Labels |
|---|---|
| Type | `type:feature`, `type:bug`, `type:chore` |
| Priority | `prio:P0` (critical), `prio:P1` (high), `prio:P2` (normal) |
| Area | `area:backend`, `area:frontend`, `area:mobile`, `area:infra` |
| Status | `status:blocked`, `status:needs-review`, `status:ready` |
| Size | `size:XS`, `size:S`, `size:M` |

## 12. Dependency & Blocking

- If task A is blocked by task B → add label `status:blocked` + comment with link to issue B.
- Cross-platform dependency (API contract) → Tech Lead creates API contract issue first.
- Do not merge FE/Mobile if the API endpoint has not been merged yet (unless mocked).

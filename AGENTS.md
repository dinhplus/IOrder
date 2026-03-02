# IOrder — Copilot Coding Agent Instructions

IOrder is a restaurant ordering platform built as a mono-repo with a Go backend, Next.js frontend, and React Native mobile app.

## Repo Structure

```
IOrder/
├── backend/          # Go 1.24 + Gin + PostgreSQL
├── frontend/         # Next.js 15 + TypeScript + Tailwind CSS
├── mobile/           # React Native (Expo SDK 52) + TypeScript
├── docs/             # Architecture decisions (ADRs), process docs
├── scripts/          # Setup & utility scripts
└── .github/          # CI/CD workflows, issue templates, CODEOWNERS
```

**Per-directory instructions:** Each sub-project has its own `AGENTS.md` with stack-specific rules:
- [`backend/AGENTS.md`](backend/AGENTS.md) — Go, Gin, PostgreSQL
- [`frontend/AGENTS.md`](frontend/AGENTS.md) — Next.js, TypeScript, Tailwind CSS
- [`mobile/AGENTS.md`](mobile/AGENTS.md) — React Native, Expo, TypeScript

## General Conventions

### Commits (Conventional Commits)
```
<type>(<scope>): <description>

# Examples:
feat(backend): add menu endpoint
fix(frontend): resolve hydration error
chore(ci): update golangci-lint version
```
Scopes: `backend`, `frontend`, `mobile`, `ci`, `docs`, `infra`

### Branches
```
<type>/<issue-number>-<short-description>

# Examples:
feat/42-add-menu-endpoint
fix/57-fix-login-redirect
```

### Pull Requests
- All changes go through PRs — no direct pushes to `main`
- PR title = conventional commit message (squash merge)
- Use `.github/PULL_REQUEST_TEMPLATE.md`
- CI must be green; minimum 1 approval required

### API Error Format
All backend errors follow:
```json
{"error": {"code": "...", "message": "...", "details": [...]}}
```

## Cross-Platform Rules

- **API-first:** Define the API contract before implementing on frontend/mobile.
- **Error codes:** Use `handler.ErrValidation`, `handler.ErrNotFound`, `handler.ErrUnauthorized`, `handler.ErrInternal` in the backend.
- **Shared types:** Frontend and mobile share the same `APIError` class pattern (`src/types/api.ts`).

## Quality Gates

Before opening a PR, run the quality checks for the affected area(s):

### Backend
```bash
cd backend
go build ./...
go test -race ./...
golangci-lint run
govulncheck ./...
```

### Frontend
```bash
cd frontend
pnpm run typecheck
pnpm run lint
pnpm run build
```

### Mobile
```bash
cd mobile
pnpm run typecheck
pnpm run lint
pnpm test
```

## CI Monitoring (Required for Every Agent Session)

After pushing changes or opening a PR, **always** proactively check CI results:

1. Use `list_workflow_runs` (GitHub MCP) to find the latest run for your branch.
2. Use `get_job_logs` (GitHub MCP) to read detailed logs for any failed jobs.
3. If CI fails, investigate the logs, fix the root cause, and push again.
4. Repeat until all CI checks pass.

```
# Example flow after pushing:
list_workflow_runs → find run ID for your branch
get_workflow_run (run ID) → check overall status
get_job_logs (run ID, failed_only=true) → read failure details
# fix the issue, push again, re-check
```

Never consider a task complete while CI is red.

## Key Documentation

| Doc | Purpose |
|---|---|
| [`docs/process/working-agreement.md`](docs/process/working-agreement.md) | DoR/DoD, branch/commit conventions, WIP limits |
| [`docs/architecture/adr/`](docs/architecture/adr/) | Architecture Decision Records |
| [`.github/copilot-instructions.md`](.github/copilot-instructions.md) | Copilot Chat context (mirrors this file) |

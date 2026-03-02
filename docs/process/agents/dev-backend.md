# Agent Role: Dev Backend (Go + Gin)

## Mission

Implement backend services using Go + Gin framework. Ensure API endpoints are accurate per contract, test coverage is solid, and code is clean and idiomatic Go.

## Tech Stack

| Component | Technology |
|---|---|
| Language | Go 1.23 |
| Web framework | Gin |
| Database | PostgreSQL |
| Migrations | golang-migrate |
| Linter | golangci-lint |
| Vuln check | govulncheck |

## Inputs

| Source | What to read |
|---|---|
| GitHub Issues | Issues with label `area:backend` + `status:ready` |
| `docs/architecture/adr/` | Relevant ADRs |
| API contract issues | Endpoint specs from Tech Lead |
| `backend/` | Existing codebase |
| `.golangci.yml` | Lint rules |
| CI results | `ci-backend.yml` output |

## Outputs

| Output | Where |
|---|---|
| Go source code | `backend/` |
| Unit/integration tests | `backend/**/*_test.go` |
| DB migrations | `backend/migrations/` |
| Dockerfile updates | `backend/Dockerfile` |
| PR with implementation | GitHub PR |

## Definition of Ready (before pickup)

- [ ] Issue has `area:backend` + `status:ready`
- [ ] AC is clear (endpoint path, request/response schema)
- [ ] API contract defined (if cross-platform)
- [ ] Not blocked by other tasks

## Definition of Done

- [ ] Code compiles (`go build ./...`)
- [ ] Tests pass (`go test ./...`)
- [ ] Lint clean (`golangci-lint run`)
- [ ] Vuln check clean (`govulncheck ./...`)
- [ ] PR opened per template, CI green
- [ ] Reviewed with ≥ 1 approval
- [ ] Migration tested (if DB change)

## Working Agreements

- **WIP:** 1 task at a time.
- **Branch:** `feat/<issue>-<desc>`, `fix/<issue>-<desc>`
- **Commit:** `feat(backend): ...`, `fix(backend): ...`
- **Test:** Each endpoint must have at least 1 happy path + 1 error test.
- **Error format:** Standard JSON `{"error": {"code": "...", "message": "..."}}`.
- **Config:** 12-factor (env vars), use `backend/config/`.
- **Migrations:** Sequentially numbered, reversible when possible.

## Report Format

```markdown
### Dev Backend Update — [Date]

**Context:** Working on #[issue-number] — [brief description]
**Plan:** [Implement endpoint / Write migrations / Add tests]
**Result:** [Endpoint implemented, X tests added, PR #Y opened]
**Blockers:** [None / Waiting for API contract / DB access issue]
**Next:** [Address review comments / Pick next issue]
```

## Handoff Rules

| Action | Handoff to | How |
|---|---|---|
| PR opened | Tech Lead / QA | Move issue to **In Review**, request review |
| API endpoint merged | Dev FE + Mobile | Comment on API contract issue "endpoint live" |
| Review feedback received | Self | Address comments, re-request review |
| Migration added | DevOps | Note in PR description for staging deploy |
| Blocked | Tech Lead | Comment on issue, add `status:blocked` |

## Code Standards

```go
// Project layout (within backend/)
backend/
├── cmd/server/main.go       // Entry point
├── internal/
│   ├── config/               // Configuration
│   ├── handler/              // Gin handlers (HTTP layer)
│   ├── middleware/            // Auth, logging, recovery
│   ├── model/                // Domain models
│   ├── repository/           // DB access layer
│   ├── service/              // Business logic
│   └── router/               // Route registration
├── migrations/               // SQL migrations
├── Dockerfile
├── go.mod
└── go.sum
```

- Follow standard Go project layout.
- Handlers thin, logic in services.
- Repository pattern for DB access.
- Context propagation throughout.
- Graceful shutdown via signal handling.

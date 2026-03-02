# GitHub Copilot Instructions — IOrder Mono-Repo

This is the IOrder mono-repo: a restaurant ordering platform with Go backend, Next.js frontend, and React Native mobile app.

## Repo Structure

```
IOrder/
├── backend/          # Go 1.24 + Gin + PostgreSQL
├── frontend/         # Next.js 15 + TypeScript + Tailwind CSS
├── mobile/           # React Native (Expo SDK 52) + TypeScript
├── docs/             # Architecture decisions, process docs
└── .github/          # CI/CD workflows, issue templates
```

## General Conventions

- **Commits:** Conventional Commits — `<type>(<scope>): <description>`
  - Scopes: `backend`, `frontend`, `mobile`, `ci`, `docs`, `infra`
- **Branches:** `<type>/<issue-number>-<short-description>`
- **PRs:** Squash merge, 1 approval required, CI must pass
- **Error format:** `{"error":{"code":"...","message":"...","details":[...]}}`

## Backend (Go)

- Module: `github.com/dinhplus/IOrder/backend`
- Layout: `cmd/server/`, `internal/{config,db,handler,middleware,router}/`, `migrations/`
- Framework: Gin (`github.com/gin-gonic/gin`)
- Database: PostgreSQL via `pgx/v5/stdlib`
- Migrations: `golang-migrate` SQL files in `backend/migrations/`
- Logging: `log/slog` — JSON in prod, text in dev
- Config: env vars loaded via `godotenv` for `.env` file
- Error helper: `handler.RespondError(c, status, code, message)`
- All handlers must have happy path + error test
- Run `golangci-lint run` before committing

## Frontend (Next.js)

- Directory: `frontend/`
- Framework: Next.js 15 App Router
- Language: TypeScript strict mode
- Styling: Tailwind CSS v4
- Package manager: pnpm
- API client: `frontend/src/lib/api/client.ts` — use `apiClient.get/post/put/delete`
- Error type: `APIError` from `@/types/api`
- Env: `NEXT_PUBLIC_API_URL` for backend URL
- Run `pnpm run typecheck && pnpm run lint && pnpm run build` before committing

## Mobile (React Native)

- Directory: `mobile/`
- Framework: Expo SDK 52 + expo-router
- Language: TypeScript strict mode
- Package manager: pnpm
- API client: `mobile/src/lib/api/client.ts` — same pattern as frontend
- Error type: `APIError` from `@/types/api`
- Env: API URL via `expo-constants` from `app.json` extra config
- Run `pnpm run typecheck && pnpm run lint && pnpm test` before committing

## CI Monitoring (Required for Every Agent Session)

After pushing changes or opening a PR, **always** proactively check CI results:

1. Use `list_workflow_runs` (GitHub MCP) to find the latest run for your branch.
2. Use `get_job_logs` (GitHub MCP) to read detailed logs for any failed jobs.
3. If CI fails, investigate the logs, fix the root cause, and push again.
4. Repeat until all CI checks pass.

Never consider a task complete while CI is red.

## Code Style

- Go: idiomatic, no magic numbers, context propagation, handle all errors
- TypeScript: strict types, no `any` unless necessary, prefer interfaces over types for objects
- React: functional components, hooks, no class components
- CSS: Tailwind utility classes, no inline styles unless dynamic

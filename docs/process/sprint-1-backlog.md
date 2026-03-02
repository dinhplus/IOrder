# Dry-Run Sprint 1 — Ready Issues

> 12 sample issues for Sprint 1. Copy-paste into GitHub Issues or run the script below.
> All issues meet the Definition of Ready.

---

## Automated Issue Creation Script (via `gh` CLI)

```bash
#!/bin/bash
REPO="dinhplus/IOrder"

# --- BACKEND ISSUES ---

gh issue create -R "$REPO" \
  --title "[Feature]: Backend — Health check endpoint (GET /health)" \
  --label "type:feature,area:backend,prio:P1,size:XS,status:ready" \
  --body "## User Story
As a DevOps agent, I want a health check endpoint so that I can monitor if the backend is running.

## Acceptance Criteria
- [ ] GET /health returns 200 with \`{\"status\":\"ok\",\"version\":\"0.1.0\"}\`
- [ ] Includes DB connectivity check (returns \`degraded\` if DB unreachable)
- [ ] Response time < 100ms

## Technical Notes
- Path: \`backend/internal/handler/health.go\`
- Register in router: \`backend/internal/router/router.go\`
- No auth required"

gh issue create -R "$REPO" \
  --title "[Feature]: Backend — Graceful shutdown with signal handling" \
  --label "type:feature,area:backend,prio:P1,size:S,status:ready" \
  --body "## User Story
As a DevOps agent, I want the server to gracefully shutdown so that in-flight requests complete before exit.

## Acceptance Criteria
- [ ] Server listens for SIGINT and SIGTERM
- [ ] Waits up to 30s for in-flight requests to complete
- [ ] Logs shutdown events
- [ ] Closes DB connections cleanly

## Technical Notes
- Implement in \`backend/cmd/server/main.go\`
- Use \`context.WithTimeout\` for shutdown deadline"

gh issue create -R "$REPO" \
  --title "[Feature]: Backend — Config management (env vars + .env)" \
  --label "type:feature,area:backend,prio:P1,size:S,status:ready" \
  --body "## User Story
As a Dev Backend agent, I want centralized config management so that I can easily configure the app per environment.

## Acceptance Criteria
- [ ] Config struct at \`backend/internal/config/config.go\`
- [ ] Reads from env vars (12-factor)
- [ ] Supports \`.env\` file for local dev
- [ ] Required vars: PORT, DATABASE_URL, LOG_LEVEL, ENV (dev/staging/prod)
- [ ] Validation on startup (fail fast if missing required vars)
- [ ] \`.env.example\` provided

## Technical Notes
- Use \`github.com/caarlos0/env\` or \`github.com/kelseyhightower/envconfig\`
- Add \`.env.example\` to \`backend/\`"

gh issue create -R "$REPO" \
  --title "[Feature]: Backend — Structured logging setup" \
  --label "type:feature,area:backend,prio:P1,size:S,status:ready" \
  --body "## User Story
As a Dev Backend agent, I want structured JSON logging so that logs are parseable in production.

## Acceptance Criteria
- [ ] JSON logger in production, pretty logger in dev
- [ ] Log level configurable via LOG_LEVEL env var
- [ ] Request ID in every log line
- [ ] Gin request logging middleware
- [ ] Logger accessible throughout the app

## Technical Notes
- Use \`log/slog\` (stdlib, Go 1.21+) or \`go.uber.org/zap\`
- Middleware at \`backend/internal/middleware/logger.go\`"

gh issue create -R "$REPO" \
  --title "[Feature]: Backend — Standard error response format" \
  --label "type:feature,area:backend,prio:P1,size:XS,status:ready" \
  --body "## User Story
As a Dev Frontend/Mobile agent, I want consistent error responses so that I can handle errors uniformly.

## Acceptance Criteria
- [ ] Error response format: \`{\"error\":{\"code\":\"VALIDATION_ERROR\",\"message\":\"...\",\"details\":[...]}}\`
- [ ] Error codes enum: VALIDATION_ERROR, NOT_FOUND, UNAUTHORIZED, INTERNAL_ERROR
- [ ] Helper function \`RespondError(c *gin.Context, status int, code string, message string)\`
- [ ] Gin recovery middleware returns JSON (not HTML)

## Technical Notes
- \`backend/internal/handler/response.go\` — helpers
- \`backend/internal/middleware/recovery.go\` — recovery middleware"

gh issue create -R "$REPO" \
  --title "[Feature]: Backend — Base router structure + middleware stack" \
  --label "type:feature,area:backend,prio:P1,size:S,status:ready" \
  --body "## User Story
As a Dev Backend agent, I want a well-organized router so that adding new endpoints is straightforward.

## Acceptance Criteria
- [ ] Router setup at \`backend/internal/router/router.go\`
- [ ] Middleware stack: Logger → Recovery → CORS → RequestID
- [ ] API versioning: all routes under \`/api/v1/\`
- [ ] Health endpoint outside versioned group: \`/health\`
- [ ] CORS configured for frontend origins

## Technical Notes
- Use Gin's RouterGroup for versioning
- CORS: allow configurable origins from env"

gh issue create -R "$REPO" \
  --title "[Feature]: Backend — PostgreSQL connection + migration setup" \
  --label "type:feature,area:backend,prio:P1,size:S,status:ready" \
  --body "## User Story
As a Dev Backend agent, I want DB connection and migration tooling set up so that I can start writing data models.

## Acceptance Criteria
- [ ] PostgreSQL connection pool using \`pgx\` or \`database/sql\` + \`lib/pq\`
- [ ] Connection config from DATABASE_URL env var
- [ ] Connection health check (used by /health endpoint)
- [ ] golang-migrate setup with \`backend/migrations/\` directory
- [ ] First migration: empty (just verifies migration tooling works)
- [ ] Makefile or script: \`make migrate-up\`, \`make migrate-down\`, \`make migrate-create\`

## Technical Notes
- Use \`github.com/golang-migrate/migrate/v4\`
- Store migrations in \`backend/migrations/\` (SQL format)"

gh issue create -R "$REPO" \
  --title "[Chore]: Backend — Dockerfile for production build" \
  --label "type:chore,area:backend,prio:P2,size:S,status:ready" \
  --body "## Description
Create a multi-stage Dockerfile for the backend service.

## Done When
- [ ] Multi-stage build: Go builder → Alpine/distroless runner
- [ ] Binary compiled with \`CGO_ENABLED=0\`
- [ ] Non-root user in final image
- [ ] Image size < 30MB
- [ ] Labels: maintainer, version, description
- [ ] \`.dockerignore\` in \`backend/\`"

# --- FRONTEND ISSUES ---

gh issue create -R "$REPO" \
  --title "[Feature]: Frontend — Next.js app skeleton with TypeScript" \
  --label "type:feature,area:frontend,prio:P1,size:S,status:ready" \
  --body "## User Story
As a Dev Frontend agent, I want the Next.js project initialized so that I can start building pages.

## Acceptance Criteria
- [ ] Next.js app created at \`frontend/\` with App Router
- [ ] TypeScript configured strictly (\`strict: true\`)
- [ ] Tailwind CSS installed and configured
- [ ] Basic layout with header placeholder
- [ ] Home page with \"IOrder\" title
- [ ] ESLint configured (next/core-web-vitals)
- [ ] \`pnpm\` as package manager (pnpm-lock.yaml)

## Technical Notes
- \`npx create-next-app@latest frontend --typescript --tailwind --eslint --app --src-dir\`
- Ensure Node 22 compatibility"

gh issue create -R "$REPO" \
  --title "[Feature]: Frontend — Environment config + API client wrapper" \
  --label "type:feature,area:frontend,prio:P1,size:S,status:ready" \
  --body "## User Story
As a Dev Frontend agent, I want a typed API client so that I can call backend endpoints safely.

## Acceptance Criteria
- [ ] Env config: \`NEXT_PUBLIC_API_URL\` for backend URL
- [ ] \`.env.example\` with all env vars documented
- [ ] Typed fetch wrapper at \`frontend/src/lib/api/client.ts\`
- [ ] Error handling: parse error response format from backend
- [ ] Type for API error: \`{error: {code: string, message: string, details?: any[]}}\`

## Technical Notes
- Match backend error format from \"Standard error response format\" issue
- Consider using \`ky\` or native \`fetch\` with typed wrapper"

# --- MOBILE ISSUES ---

gh issue create -R "$REPO" \
  --title "[Feature]: Mobile — React Native project initialization" \
  --label "type:feature,area:mobile,prio:P1,size:S,status:ready" \
  --body "## User Story
As a Dev Mobile agent, I want the React Native project initialized so that I can start building screens.

## Acceptance Criteria
- [ ] React Native (Expo) project at \`mobile/\`
- [ ] TypeScript configured (\`strict: true\`)
- [ ] ESLint configured
- [ ] Basic navigation structure (React Navigation)
- [ ] Home screen with \"IOrder\" title
- [ ] \`pnpm\` as package manager

## Technical Notes
- \`npx create-expo-app mobile --template tabs\`
- Ensure Node 22 compatibility"

gh issue create -R "$REPO" \
  --title "[Feature]: Mobile — Environment config + API client" \
  --label "type:feature,area:mobile,prio:P1,size:S,status:ready" \
  --body "## User Story
As a Dev Mobile agent, I want a typed API client so that I can call backend endpoints safely.

## Acceptance Criteria
- [ ] Env config for API_URL (Expo constants or react-native-config)
- [ ] \`.env.example\` with all env vars documented
- [ ] Typed fetch wrapper at \`mobile/src/lib/api/client.ts\`
- [ ] Error handling matching backend format
- [ ] Works on both iOS and Android

## Technical Notes
- Match error format with frontend implementation
- Consider sharing types between frontend and mobile"

echo "✅ All 12 issues created!"
```

---

## Issue Content (for manual creation)

### Issue 1: Backend — Health check endpoint

**Labels:** `type:feature`, `area:backend`, `prio:P1`, `size:XS`, `status:ready`

**Body:**

As a DevOps agent, I want a health check endpoint so that I can monitor if the backend is running.

**Acceptance Criteria:**
- [ ] GET /health returns 200 with `{"status":"ok","version":"0.1.0"}`
- [ ] Includes DB connectivity check (returns `degraded` if DB unreachable)
- [ ] Response time < 100ms

---

### Issue 2: Backend — Graceful shutdown

**Labels:** `type:feature`, `area:backend`, `prio:P1`, `size:S`, `status:ready`

**Body:**

As a DevOps agent, I want the server to gracefully shutdown so that in-flight requests complete.

**Acceptance Criteria:**
- [ ] Server listens for SIGINT and SIGTERM
- [ ] Waits up to 30s for in-flight requests to complete
- [ ] Logs shutdown events
- [ ] Closes DB connections cleanly

---

### Issue 3: Backend — Config management

**Labels:** `type:feature`, `area:backend`, `prio:P1`, `size:S`, `status:ready`

**Body:**

Centralized config from env vars + `.env` file.

**Acceptance Criteria:**
- [ ] Config struct at `backend/internal/config/config.go`
- [ ] Reads from env vars (12-factor)
- [ ] Required vars: PORT, DATABASE_URL, LOG_LEVEL, ENV
- [ ] Validation on startup
- [ ] `.env.example` provided

---

### Issue 4: Backend — Structured logging

**Labels:** `type:feature`, `area:backend`, `prio:P1`, `size:S`, `status:ready`

**Body:**

JSON logging in production, pretty in dev.

**Acceptance Criteria:**
- [ ] JSON logger in prod, pretty in dev
- [ ] Log level from LOG_LEVEL env var
- [ ] Request ID in every log line
- [ ] Gin request logging middleware

---

### Issue 5: Backend — Standard error response

**Labels:** `type:feature`, `area:backend`, `prio:P1`, `size:XS`, `status:ready`

**Body:**

Consistent error format for FE/Mobile consumption.

**Acceptance Criteria:**
- [ ] Format: `{"error":{"code":"...","message":"...","details":[...]}}`
- [ ] Error codes enum
- [ ] Helper function `RespondError()`
- [ ] Recovery middleware returns JSON

---

### Issue 6: Backend — Base router structure

**Labels:** `type:feature`, `area:backend`, `prio:P1`, `size:S`, `status:ready`

**Body:**

Well-organized Gin router with middleware stack.

**Acceptance Criteria:**
- [ ] Logger → Recovery → CORS → RequestID middleware
- [ ] API versioning: `/api/v1/`
- [ ] Health outside versioned group
- [ ] Configurable CORS origins

---

### Issue 7: Backend — PostgreSQL + migrations

**Labels:** `type:feature`, `area:backend`, `prio:P1`, `size:S`, `status:ready`

**Body:**

DB connection pool + golang-migrate setup.

**Acceptance Criteria:**
- [ ] PostgreSQL connection from DATABASE_URL
- [ ] Connection health check
- [ ] golang-migrate with `backend/migrations/`
- [ ] First empty migration
- [ ] `make migrate-up/down/create`

---

### Issue 8: Backend — Dockerfile

**Labels:** `type:chore`, `area:backend`, `prio:P2`, `size:S`, `status:ready`

**Body:**

Multi-stage Docker build.

**Acceptance Criteria:**
- [ ] Multi-stage: builder → distroless/alpine
- [ ] CGO_ENABLED=0
- [ ] Non-root user
- [ ] Image < 30MB
- [ ] `.dockerignore`

---

### Issue 9: Frontend — Next.js skeleton

**Labels:** `type:feature`, `area:frontend`, `prio:P1`, `size:S`, `status:ready`

**Body:**

Initialize Next.js with TypeScript + Tailwind + ESLint.

**Acceptance Criteria:**
- [ ] App Router + strict TypeScript
- [ ] Tailwind CSS
- [ ] Basic layout + home page
- [ ] ESLint configured
- [ ] pnpm

---

### Issue 10: Frontend — Env config + API client

**Labels:** `type:feature`, `area:frontend`, `prio:P1`, `size:S`, `status:ready`

**Body:**

Typed API client matching backend error format.

**Acceptance Criteria:**
- [ ] `NEXT_PUBLIC_API_URL`
- [ ] `.env.example`
- [ ] Typed fetch wrapper
- [ ] Error handling matching backend format

---

### Issue 11: Mobile — React Native init

**Labels:** `type:feature`, `area:mobile`, `prio:P1`, `size:S`, `status:ready`

**Body:**

Initialize Expo React Native with TypeScript.

**Acceptance Criteria:**
- [ ] Expo project at `mobile/`
- [ ] TypeScript strict
- [ ] ESLint
- [ ] React Navigation
- [ ] Home screen

---

### Issue 12: Mobile — Env config + API client

**Labels:** `type:feature`, `area:mobile`, `prio:P1`, `size:S`, `status:ready`

**Body:**

Typed API client matching backend/frontend pattern.

**Acceptance Criteria:**
- [ ] Env config for API_URL
- [ ] `.env.example`
- [ ] Typed fetch wrapper
- [ ] Cross-platform compatible

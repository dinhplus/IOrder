# Backend — Code Instructions (GitHub Copilot / Claude)

## Stack
- **Go 1.24** + **Gin** framework
- **PostgreSQL** via `pgx/v5/stdlib`
- **golang-migrate** (SQL migrations in `backend/migrations/`)
- **log/slog** for structured JSON logging
- **golangci-lint** for linting

## Project Layout

```
backend/
├── cmd/server/main.go          # Entry: config → logger → DB → router → HTTP server
├── internal/
│   ├── config/config.go        # Env var config (godotenv + os.Getenv)
│   ├── db/db.go                # pgx/stdlib pool, ping check
│   ├── handler/
│   │   ├── health.go           # GET /health → {"status":"ok","db":"ok"}
│   │   └── response.go         # RespondError / RespondSuccess helpers
│   ├── middleware/
│   │   ├── logger.go           # slog request logger
│   │   └── recovery.go         # JSON panic recovery
│   └── router/router.go        # Gin engine: middleware stack + routes
├── migrations/                 # *.up.sql / *.down.sql (golang-migrate format)
├── Makefile                    # build, run, test, lint, migrate-*
├── Dockerfile                  # Multi-stage: golang:1.24-alpine → alpine:3.21
└── .env.example
```

## Rules

### Error Handling
Always use `RespondError()` from `internal/handler/response.go`:
```go
handler.RespondError(c, http.StatusBadRequest, handler.ErrValidation, "invalid input")
```

### Error Codes
```go
handler.ErrValidation  // VALIDATION_ERROR
handler.ErrNotFound    // NOT_FOUND
handler.ErrUnauthorized // UNAUTHORIZED
handler.ErrInternal    // INTERNAL_ERROR
```

### Config Access
```go
cfg, err := config.Load()
// cfg.Port, cfg.DatabaseURL, cfg.LogLevel, cfg.Env
```

### Logging
```go
slog.InfoContext(ctx, "message", "key", value)
slog.ErrorContext(ctx, "error occurred", "err", err)
```

### Testing Pattern
```go
func TestSomething(t *testing.T) {
    r := gin.New()
    r.GET("/path", handler)
    w := httptest.NewRecorder()
    req := httptest.NewRequest(http.MethodGet, "/path", nil)
    r.ServeHTTP(w, req)
    assert.Equal(t, http.StatusOK, w.Code)
}
```

### Migrations
```sql
-- backend/migrations/000002_add_users.up.sql
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- backend/migrations/000002_add_users.down.sql
DROP TABLE IF EXISTS users;
```

### Route Registration
Add new routes to `internal/router/router.go` in the `/api/v1` group:
```go
v1 := r.Group("/api/v1")
{
    v1.GET("/users", userHandler.List)
    v1.POST("/users", userHandler.Create)
}
```

## Quality Gates (before PR)
```bash
cd backend
go build ./...            # must pass
go test -race ./...        # must pass
golangci-lint run         # must pass (0 warnings)
govulncheck ./...         # must pass
```

# IOrder Backend

Go-based REST API server for the IOrder restaurant ordering platform.

## Stack

- **Go 1.24** - Programming language
- **Gin** - Web framework
- **PostgreSQL** - Database (via pgx/v5/stdlib driver)
- **golang-migrate** - Database migrations
- **slog** - Structured logging
- **golangci-lint** - Code linting

## Project Structure

```
backend/
├── cmd/
│   └── server/
│       └── main.go                 # Application entry point
├── internal/
│   ├── config/
│   │   └── config.go              # Environment configuration
│   ├── db/
│   │   └── db.go                  # Database connection setup
│   ├── handler/
│   │   ├── health.go              # Health check endpoint
│   │   ├── menu.go                # Menu & menu item endpoints
│   │   ├── order.go               # Order management endpoints
│   │   ├── staff.go               # Staff management endpoints
│   │   ├── table.go               # Table & floor plan endpoints
│   │   ├── tenant.go              # Tenant management endpoints
│   │   ├── response.go            # Standard API response helpers
│   │   └── *_test.go              # Handler unit tests
│   ├── middleware/
│   │   ├── logger.go              # Request logging middleware
│   │   └── recovery.go            # Panic recovery middleware
│   ├── repository/
│   │   ├── menu.go                # Menu data access
│   │   ├── order.go               # Order data access
│   │   ├── staff.go               # Staff data access
│   │   ├── table.go               # Table & floor plan data access
│   │   └── tenant.go              # Tenant data access
│   └── router/
│       └── router.go              # Route configuration & middleware stack
├── migrations/
│   ├── 000001_init.up.sql         # Initial schema
│   └── 000001_init.down.sql       # Rollback script
├── Dockerfile                      # Multi-stage container build
├── Makefile                        # Development commands
├── go.mod                          # Go module dependencies
└── .env.example                    # Environment variable template
```

## Quick Start

### Prerequisites

- Go 1.24+
- PostgreSQL 16+
- golang-migrate CLI (`brew install golang-migrate` on macOS)

### Setup

1. **Clone and navigate to backend:**
   ```bash
   cd backend
   ```

2. **Install dependencies:**
   ```bash
   go mod download
   ```

3. **Create environment file:**
   ```bash
   cp .env.example .env
   ```

4. **Start PostgreSQL** (example using Docker):
   ```bash
   docker run --name iorder-postgres \
     -e POSTGRES_PASSWORD=postgres \
     -e POSTGRES_DB=iorder \
     -p 5432:5432 \
     -d postgres:16-alpine
   ```

5. **Run database migrations:**
   ```bash
   make migrate-up
   ```

6. **Start the server:**
   ```bash
   make run
   ```

   Server starts on `http://localhost:8080`

## Environment Variables

Create a `.env` file in the `backend/` directory:

```bash
PORT=8080
DATABASE_URL=postgres://postgres:postgres@localhost:5432/iorder?sslmode=disable
LOG_LEVEL=info          # debug, info, warn, error
ENV=dev                 # dev, production
CORS_ORIGINS=http://localhost:3000
```

## API Endpoints

### Health

- `GET /health` - Health check with database status

### Tenants

- `GET /api/v1/tenants` - List all tenants
- `POST /api/v1/tenants` - Create a new tenant
- `GET /api/v1/tenants/:id` - Get tenant by ID
- `PATCH /api/v1/tenants/:id` - Update tenant

### Menu

**Categories:**
- `GET /api/v1/menu/categories` - List categories
- `POST /api/v1/menu/categories` - Create category
- `PUT /api/v1/menu/categories/:id` - Update category
- `DELETE /api/v1/menu/categories/:id` - Delete category

**Items:**
- `GET /api/v1/menu/items` - List menu items
- `POST /api/v1/menu/items` - Create menu item
- `PUT /api/v1/menu/items/:id` - Update menu item
- `PATCH /api/v1/menu/items/:id/availability` - Toggle item availability
- `DELETE /api/v1/menu/items/:id` - Delete menu item

### Floor Plans & Tables

**Floor Plans:**
- `GET /api/v1/floor-plans` - List floor plans
- `POST /api/v1/floor-plans` - Create floor plan
- `PUT /api/v1/floor-plans/:id` - Update floor plan
- `GET /api/v1/floor-plans/:id/tables` - List tables in floor plan

**Tables:**
- `POST /api/v1/tables` - Create table
- `PUT /api/v1/tables/:id` - Update table
- `DELETE /api/v1/tables/:id` - Delete table

### Orders

- `POST /api/v1/orders` - Create order
- `GET /api/v1/orders` - List orders (by tenant)
- `GET /api/v1/orders/:id` - Get order details
- `PATCH /api/v1/orders/:id/items` - Update order items
- `POST /api/v1/orders/:id/submit` - Submit order (customer → staff)
- `POST /api/v1/orders/:id/confirm` - Confirm order (staff accepts)
- `POST /api/v1/orders/:id/reject` - Reject order
- `POST /api/v1/orders/:id/start-preparation` - Start preparing
- `POST /api/v1/orders/:id/ready` - Mark ready for serving
- `POST /api/v1/orders/:id/serve` - Mark as served
- `POST /api/v1/orders/:id/request-payment` - Request payment
- `POST /api/v1/orders/:id/pay` - Process payment
- `POST /api/v1/orders/:id/close` - Close order
- `POST /api/v1/orders/:id/cancel` - Cancel order

### Staff

- `GET /api/v1/staff` - List staff members
- `POST /api/v1/staff` - Create staff member
- `GET /api/v1/staff/:id` - Get staff member
- `PATCH /api/v1/staff/:id` - Update staff member
- `DELETE /api/v1/staff/:id` - Delete staff member

## Error Response Format

All errors return a consistent JSON structure:

```json
{
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "name is required",
    "details": []
  }
}
```

**Error Codes:**
- `VALIDATION_ERROR` - Invalid request data
- `NOT_FOUND` - Resource not found
- `UNAUTHORIZED` - Authentication required
- `INTERNAL_ERROR` - Server error

## Multi-Tenancy

Most endpoints require the `X-Tenant-ID` header:

```bash
curl -H "X-Tenant-ID: <tenant-uuid>" http://localhost:8080/api/v1/menu/items
```

## Database Schema

See [migrations/000001_init.up.sql](migrations/000001_init.up.sql) for the full schema.

**Key Tables:**
- `tenants` - Restaurant tenants
- `staff` - Staff members per tenant
- `floor_plans` - Restaurant floor layouts
- `restaurant_tables` - Physical tables with QR codes
- `menu_categories` - Food/beverage categories
- `menu_items` - Menu items with pricing
- `orders` - Customer orders with state machine
- `order_items` - Line items for each order
- `order_events` - Audit trail of order state changes
- `customers` - Customer profiles (future use)

## Development Workflow

### Run Tests

```bash
make test
# or with race detector
go test -race ./...
```

### Lint Code

```bash
make lint
# or directly
golangci-lint run
```

### Build Binary

```bash
make build
# produces: bin/server
./bin/server
```

### Database Migrations

**Create a new migration:**
```bash
make migrate-create name=add_users
```

**Run migrations:**
```bash
make migrate-up
```

**Rollback last migration:**
```bash
make migrate-down
```

## Docker

Build and run using Docker:

```bash
# Build image
docker build -t iorder-backend .

# Run container
docker run -p 8080:8080 \
  -e DATABASE_URL="postgres://..." \
  iorder-backend
```

## Code Conventions

### Handler Pattern

```go
func (h *MenuHandler) ListItems() gin.HandlerFunc {
    return func(c *gin.Context) {
        tenantID, ok := tenantIDFromHeader(c)
        if !ok {
            return // Error already sent
        }
        
        items, err := h.repo.ListItems(c.Request.Context(), tenantID)
        if err != nil {
            RespondError(c, http.StatusInternalServerError, ErrInternal, "failed to list items")
            return
        }
        
        RespondSuccess(c, http.StatusOK, items)
    }
}
```

### Repository Pattern

```go
func (r *MenuRepository) GetItem(ctx context.Context, id, tenantID string) (*MenuItem, error) {
    item := &MenuItem{}
    err := r.db.QueryRowContext(ctx, 
        `SELECT ... FROM menu_items WHERE id=$1 AND tenant_id=$2`, 
        id, tenantID,
    ).Scan(&item.ID, ...)
    
    if errors.Is(err, sql.ErrNoRows) {
        return nil, nil // Not found returns nil, not error
    }
    return item, err
}
```

### Error Handling

Always use the standard response helpers:

```go
handler.RespondError(c, http.StatusBadRequest, handler.ErrValidation, "invalid input")
handler.RespondSuccess(c, http.StatusOK, data)
```

### Testing

Every handler must have:
- Happy path test
- Error case tests

```go
func TestListItems_Success(t *testing.T) {
    mock := &mockMenuStore{
        items: []*repository.MenuItem{...},
    }
    h := handler.NewMenuHandler(mock)
    
    r := gin.New()
    r.GET("/items", h.ListItems())
    
    w := httptest.NewRecorder()
    req := httptest.NewRequest(http.MethodGet, "/items", nil)
    req.Header.Set("X-Tenant-ID", "tenant-123")
    
    r.ServeHTTP(w, req)
    
    assert.Equal(t, http.StatusOK, w.Code)
}
```

## Quality Gates

Before opening a PR, ensure all checks pass:

```bash
# Compile check
go build ./...

# Run tests with race detector
go test -race ./...

# Lint check (must be zero warnings)
golangci-lint run

# Vulnerability check
govulncheck ./...
```

## Order State Machine

The order lifecycle follows this state machine (see [ADR-0003](../docs/architecture/adr/0003-order-state-machine.md)):

```
DRAFT → SUBMITTED → CONFIRMED → IN_PREPARATION → READY → SERVED 
  ↓                    ↓             ↓               ↓       ↓
CANCELLED          REJECTED    CANCELLED       CANCELLED  PAYMENT_REQUESTED → PAID → CLOSED
```

## Troubleshooting

### Database connection fails

1. Verify PostgreSQL is running:
   ```bash
   psql -h localhost -U postgres -d iorder
   ```

2. Check `DATABASE_URL` in `.env`

3. Ensure migrations have run:
   ```bash
   make migrate-up
   ```

### Tests fail

- Ensure you're using Go 1.24+
- Run `go mod tidy` to sync dependencies
- Check for uncommitted changes in `go.mod`

### Lint errors

Run `golangci-lint run --fix` to auto-fix simple issues.

## Contributing

1. Create a feature branch: `feat/<issue-number>-description`
2. Write tests for new handlers/functions
3. Run quality gates locally before pushing
4. Open PR with conventional commit message
5. Wait for CI to pass and 1 approval

## Additional Documentation

- [Architecture Overview](../docs/architecture/system-architecture.md)
- [Working Agreement](../docs/process/working-agreement.md)
- [ADRs](../docs/architecture/adr/)
- [Backend Agent Instructions](AGENTS.md)

## License

Proprietary - All rights reserved

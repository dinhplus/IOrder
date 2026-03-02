package handler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"

	"github.com/dinhplus/IOrder/backend/internal/handler"
	"github.com/dinhplus/IOrder/backend/internal/repository"
)

// mockTenantStore is a test double for handler.TenantStore.
type mockTenantStore struct {
	tenant *repository.Tenant
	err    error
}

func (m *mockTenantStore) Create(ctx context.Context, t *repository.Tenant) error {
	if m.err != nil {
		return m.err
	}
	t.ID = "tenant-uuid"
	m.tenant = t
	return nil
}
func (m *mockTenantStore) GetByID(ctx context.Context, id string) (*repository.Tenant, error) {
	return m.tenant, m.err
}
func (m *mockTenantStore) GetBySlug(ctx context.Context, slug string) (*repository.Tenant, error) {
	return m.tenant, m.err
}
func (m *mockTenantStore) Update(ctx context.Context, t *repository.Tenant) error {
	return m.err
}
func (m *mockTenantStore) List(ctx context.Context) ([]*repository.Tenant, error) {
	if m.tenant != nil {
		return []*repository.Tenant{m.tenant}, m.err
	}
	return nil, m.err
}

func TestCreateTenant_Success(t *testing.T) {
	r := gin.New()
	h := handler.NewTenantHandler(&mockTenantStore{})
	r.POST("/api/v1/tenants", h.CreateTenant())

	body := `{"slug":"test-slug","name":"Test Restaurant"}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/tenants", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d: %s", w.Code, w.Body.String())
	}
}

func TestCreateTenant_MissingSlug(t *testing.T) {
	r := gin.New()
	h := handler.NewTenantHandler(&mockTenantStore{})
	r.POST("/api/v1/tenants", h.CreateTenant())

	body := `{"name":"Test Restaurant"}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/tenants", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestCreateTenant_MissingName(t *testing.T) {
	r := gin.New()
	h := handler.NewTenantHandler(&mockTenantStore{})
	r.POST("/api/v1/tenants", h.CreateTenant())

	body := `{"slug":"test-slug"}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/tenants", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestGetTenant_Found(t *testing.T) {
	store := &mockTenantStore{tenant: &repository.Tenant{ID: "abc", Slug: "s", Name: "N"}}
	r := gin.New()
	h := handler.NewTenantHandler(store)
	r.GET("/api/v1/tenants/:id", h.GetTenant())

	req := httptest.NewRequest(http.MethodGet, "/api/v1/tenants/abc", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}

func TestGetTenant_NotFound(t *testing.T) {
	r := gin.New()
	h := handler.NewTenantHandler(&mockTenantStore{})
	r.GET("/api/v1/tenants/:id", h.GetTenant())

	req := httptest.NewRequest(http.MethodGet, "/api/v1/tenants/missing", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", w.Code)
	}
}

func TestUpdateTenant_NotFound(t *testing.T) {
	r := gin.New()
	h := handler.NewTenantHandler(&mockTenantStore{})
	r.PATCH("/api/v1/tenants/:id", h.UpdateTenant())

	body := `{"name":"New Name"}`
	req := httptest.NewRequest(http.MethodPatch, "/api/v1/tenants/missing", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", w.Code)
	}
}

func TestUpdateTenant_Success(t *testing.T) {
	existing := &repository.Tenant{ID: "abc", Slug: "s", Name: "N", Timezone: "UTC", Currency: "USD"}
	store := &mockTenantStore{tenant: existing}
	r := gin.New()
	h := handler.NewTenantHandler(store)
	r.PATCH("/api/v1/tenants/:id", h.UpdateTenant())

	body := `{"name":"Updated Name"}`
	req := httptest.NewRequest(http.MethodPatch, "/api/v1/tenants/abc", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
	var resp map[string]any
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("decode: %v", err)
	}
}

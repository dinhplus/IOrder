package handler_test

import (
	"bytes"
	"context"
	"database/sql"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"

	"github.com/dinhplus/IOrder/backend/internal/handler"
	"github.com/dinhplus/IOrder/backend/internal/repository"
)

type mockMenuStore struct {
	category *repository.MenuCategory
	item     *repository.MenuItem
	err      error
}

func (m *mockMenuStore) CreateCategory(ctx context.Context, c *repository.MenuCategory) error {
	if m.err != nil {
		return m.err
	}
	c.ID = "cat-uuid"
	return nil
}
func (m *mockMenuStore) GetCategory(ctx context.Context, id, tenantID string) (*repository.MenuCategory, error) {
	return m.category, m.err
}
func (m *mockMenuStore) UpdateCategory(ctx context.Context, c *repository.MenuCategory) error {
	return m.err
}
func (m *mockMenuStore) DeleteCategory(ctx context.Context, id, tenantID string) error {
	return m.err
}
func (m *mockMenuStore) ListCategories(ctx context.Context, tenantID string) ([]*repository.MenuCategory, error) {
	if m.category != nil {
		return []*repository.MenuCategory{m.category}, m.err
	}
	return nil, m.err
}
func (m *mockMenuStore) CreateItem(ctx context.Context, item *repository.MenuItem) error {
	if m.err != nil {
		return m.err
	}
	item.ID = "item-uuid"
	return nil
}
func (m *mockMenuStore) GetItem(ctx context.Context, id, tenantID string) (*repository.MenuItem, error) {
	return m.item, m.err
}
func (m *mockMenuStore) UpdateItem(ctx context.Context, item *repository.MenuItem) error {
	return m.err
}
func (m *mockMenuStore) DeleteItem(ctx context.Context, id, tenantID string) error {
	return m.err
}
func (m *mockMenuStore) ListItems(ctx context.Context, tenantID string) ([]*repository.MenuItem, error) {
	if m.item != nil {
		return []*repository.MenuItem{m.item}, m.err
	}
	return nil, m.err
}
func (m *mockMenuStore) ToggleItemAvailability(ctx context.Context, id, tenantID string, available bool) error {
	return m.err
}

func TestListCategories_MissingTenantID(t *testing.T) {
	r := gin.New()
	h := handler.NewMenuHandler(&mockMenuStore{})
	r.GET("/menu/categories", h.ListCategories())

	req := httptest.NewRequest(http.MethodGet, "/menu/categories", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestListCategories_Success(t *testing.T) {
	r := gin.New()
	h := handler.NewMenuHandler(&mockMenuStore{})
	r.GET("/menu/categories", h.ListCategories())

	req := httptest.NewRequest(http.MethodGet, "/menu/categories", nil)
	req.Header.Set("X-Tenant-ID", "tenant-1")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}

func TestCreateCategory_Success(t *testing.T) {
	r := gin.New()
	h := handler.NewMenuHandler(&mockMenuStore{})
	r.POST("/menu/categories", h.CreateCategory())

	body := `{"name":"Appetizers","type":"food"}`
	req := httptest.NewRequest(http.MethodPost, "/menu/categories", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Tenant-ID", "tenant-1")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d: %s", w.Code, w.Body.String())
	}
}

func TestCreateCategory_MissingName(t *testing.T) {
	r := gin.New()
	h := handler.NewMenuHandler(&mockMenuStore{})
	r.POST("/menu/categories", h.CreateCategory())

	body := `{"type":"food"}`
	req := httptest.NewRequest(http.MethodPost, "/menu/categories", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Tenant-ID", "tenant-1")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestDeleteCategory_NotFound(t *testing.T) {
	r := gin.New()
	h := handler.NewMenuHandler(&mockMenuStore{err: sql.ErrNoRows})
	r.DELETE("/menu/categories/:id", h.DeleteCategory())

	req := httptest.NewRequest(http.MethodDelete, "/menu/categories/missing", nil)
	req.Header.Set("X-Tenant-ID", "tenant-1")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", w.Code)
	}
}

func TestDeleteCategory_Success(t *testing.T) {
	r := gin.New()
	h := handler.NewMenuHandler(&mockMenuStore{})
	r.DELETE("/menu/categories/:id", h.DeleteCategory())

	req := httptest.NewRequest(http.MethodDelete, "/menu/categories/cat-1", nil)
	req.Header.Set("X-Tenant-ID", "tenant-1")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNoContent {
		t.Fatalf("expected 204, got %d", w.Code)
	}
}

func TestCreateItem_Success(t *testing.T) {
	r := gin.New()
	h := handler.NewMenuHandler(&mockMenuStore{})
	r.POST("/menu/items", h.CreateItem())

	body := `{"category_id":"cat-1","name":"Spring Roll","price":50000}`
	req := httptest.NewRequest(http.MethodPost, "/menu/items", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Tenant-ID", "tenant-1")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d: %s", w.Code, w.Body.String())
	}
}

func TestCreateItem_MissingName(t *testing.T) {
	r := gin.New()
	h := handler.NewMenuHandler(&mockMenuStore{})
	r.POST("/menu/items", h.CreateItem())

	body := `{"category_id":"cat-1","price":50000}`
	req := httptest.NewRequest(http.MethodPost, "/menu/items", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Tenant-ID", "tenant-1")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestToggleAvailability_NotFound(t *testing.T) {
	r := gin.New()
	h := handler.NewMenuHandler(&mockMenuStore{err: sql.ErrNoRows})
	r.PATCH("/menu/items/:id/availability", h.ToggleAvailability())

	body := `{"is_available":false}`
	req := httptest.NewRequest(http.MethodPatch, "/menu/items/item-1/availability", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Tenant-ID", "tenant-1")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", w.Code)
	}
}

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

type mockTableStore struct {
	floorPlan *repository.FloorPlan
	table     *repository.RestaurantTable
	err       error
}

func (m *mockTableStore) CreateFloorPlan(ctx context.Context, fp *repository.FloorPlan) error {
	if m.err != nil {
		return m.err
	}
	fp.ID = "fp-uuid"
	return nil
}
func (m *mockTableStore) GetFloorPlan(ctx context.Context, id, tenantID string) (*repository.FloorPlan, error) {
	return m.floorPlan, m.err
}
func (m *mockTableStore) UpdateFloorPlan(ctx context.Context, fp *repository.FloorPlan) error {
	return m.err
}
func (m *mockTableStore) ListFloorPlans(ctx context.Context, tenantID string) ([]*repository.FloorPlan, error) {
	if m.floorPlan != nil {
		return []*repository.FloorPlan{m.floorPlan}, m.err
	}
	return nil, m.err
}
func (m *mockTableStore) ListTablesByFloorPlan(ctx context.Context, floorPlanID, tenantID string) ([]*repository.RestaurantTable, error) {
	if m.table != nil {
		return []*repository.RestaurantTable{m.table}, m.err
	}
	return nil, m.err
}
func (m *mockTableStore) CreateTable(ctx context.Context, t *repository.RestaurantTable) error {
	if m.err != nil {
		return m.err
	}
	t.ID = "table-uuid"
	return nil
}
func (m *mockTableStore) GetTable(ctx context.Context, id, tenantID string) (*repository.RestaurantTable, error) {
	return m.table, m.err
}
func (m *mockTableStore) UpdateTable(ctx context.Context, t *repository.RestaurantTable) error {
	return m.err
}
func (m *mockTableStore) DeleteTable(ctx context.Context, id, tenantID string) error {
	return m.err
}

func TestListFloorPlans_MissingTenantID(t *testing.T) {
	r := gin.New()
	h := handler.NewTableHandler(&mockTableStore{})
	r.GET("/floor-plans", h.ListFloorPlans())

	req := httptest.NewRequest(http.MethodGet, "/floor-plans", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestListFloorPlans_Success(t *testing.T) {
	r := gin.New()
	h := handler.NewTableHandler(&mockTableStore{})
	r.GET("/floor-plans", h.ListFloorPlans())

	req := httptest.NewRequest(http.MethodGet, "/floor-plans", nil)
	req.Header.Set("X-Tenant-ID", "tenant-1")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}

func TestCreateFloorPlan_Success(t *testing.T) {
	r := gin.New()
	h := handler.NewTableHandler(&mockTableStore{})
	r.POST("/floor-plans", h.CreateFloorPlan())

	body := `{"name":"Ground Floor","floor_level":1}`
	req := httptest.NewRequest(http.MethodPost, "/floor-plans", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Tenant-ID", "tenant-1")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d: %s", w.Code, w.Body.String())
	}
}

func TestCreateFloorPlan_MissingName(t *testing.T) {
	r := gin.New()
	h := handler.NewTableHandler(&mockTableStore{})
	r.POST("/floor-plans", h.CreateFloorPlan())

	body := `{"floor_level":1}`
	req := httptest.NewRequest(http.MethodPost, "/floor-plans", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Tenant-ID", "tenant-1")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestCreateTable_Success(t *testing.T) {
	r := gin.New()
	h := handler.NewTableHandler(&mockTableStore{})
	r.POST("/tables", h.CreateTable())

	body := `{"floor_plan_id":"fp-1","name":"T1","capacity":4}`
	req := httptest.NewRequest(http.MethodPost, "/tables", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Tenant-ID", "tenant-1")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d: %s", w.Code, w.Body.String())
	}
}

func TestCreateTable_MissingFloorPlanID(t *testing.T) {
	r := gin.New()
	h := handler.NewTableHandler(&mockTableStore{})
	r.POST("/tables", h.CreateTable())

	body := `{"name":"T1","capacity":4}`
	req := httptest.NewRequest(http.MethodPost, "/tables", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Tenant-ID", "tenant-1")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestDeleteTable_NotFound(t *testing.T) {
	r := gin.New()
	h := handler.NewTableHandler(&mockTableStore{err: sql.ErrNoRows})
	r.DELETE("/tables/:id", h.DeleteTable())

	req := httptest.NewRequest(http.MethodDelete, "/tables/missing", nil)
	req.Header.Set("X-Tenant-ID", "tenant-1")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", w.Code)
	}
}

func TestDeleteTable_Success(t *testing.T) {
	r := gin.New()
	h := handler.NewTableHandler(&mockTableStore{})
	r.DELETE("/tables/:id", h.DeleteTable())

	req := httptest.NewRequest(http.MethodDelete, "/tables/table-1", nil)
	req.Header.Set("X-Tenant-ID", "tenant-1")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNoContent {
		t.Fatalf("expected 204, got %d", w.Code)
	}
}

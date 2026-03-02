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

type mockStaffStore struct {
	staff *repository.Staff
	err   error
}

func (m *mockStaffStore) Create(ctx context.Context, s *repository.Staff) error {
	if m.err != nil {
		return m.err
	}
	s.ID = "staff-uuid"
	m.staff = s
	return nil
}
func (m *mockStaffStore) GetByID(ctx context.Context, id, tenantID string) (*repository.Staff, error) {
	return m.staff, m.err
}
func (m *mockStaffStore) ListByTenant(ctx context.Context, tenantID string) ([]*repository.Staff, error) {
	if m.staff != nil {
		return []*repository.Staff{m.staff}, m.err
	}
	return nil, m.err
}
func (m *mockStaffStore) Update(ctx context.Context, s *repository.Staff) error {
	return m.err
}
func (m *mockStaffStore) Delete(ctx context.Context, id, tenantID string) error {
	return m.err
}

func TestCreateStaff_Success(t *testing.T) {
	r := gin.New()
	h := handler.NewStaffHandler(&mockStaffStore{})
	r.POST("/staff", h.CreateStaff())

	body := `{"full_name":"Alice","role":"waiter"}`
	req := httptest.NewRequest(http.MethodPost, "/staff", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Tenant-ID", "tenant-1")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d: %s", w.Code, w.Body.String())
	}
}

func TestCreateStaff_MissingTenantID(t *testing.T) {
	r := gin.New()
	h := handler.NewStaffHandler(&mockStaffStore{})
	r.POST("/staff", h.CreateStaff())

	body := `{"full_name":"Alice","role":"waiter"}`
	req := httptest.NewRequest(http.MethodPost, "/staff", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestCreateStaff_MissingFullName(t *testing.T) {
	r := gin.New()
	h := handler.NewStaffHandler(&mockStaffStore{})
	r.POST("/staff", h.CreateStaff())

	body := `{"role":"waiter"}`
	req := httptest.NewRequest(http.MethodPost, "/staff", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Tenant-ID", "tenant-1")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestListStaff_Success(t *testing.T) {
	store := &mockStaffStore{staff: &repository.Staff{ID: "s1", TenantID: "t1", FullName: "Alice", Role: "waiter", IsActive: true}}
	r := gin.New()
	h := handler.NewStaffHandler(store)
	r.GET("/staff", h.ListStaff())

	req := httptest.NewRequest(http.MethodGet, "/staff", nil)
	req.Header.Set("X-Tenant-ID", "t1")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
}

func TestGetStaff_NotFound(t *testing.T) {
	r := gin.New()
	h := handler.NewStaffHandler(&mockStaffStore{})
	r.GET("/staff/:id", h.GetStaff())

	req := httptest.NewRequest(http.MethodGet, "/staff/missing", nil)
	req.Header.Set("X-Tenant-ID", "t1")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", w.Code)
	}
}

func TestDeleteStaff_NotFound(t *testing.T) {
	r := gin.New()
	h := handler.NewStaffHandler(&mockStaffStore{err: sql.ErrNoRows})
	r.DELETE("/staff/:id", h.DeleteStaff())

	req := httptest.NewRequest(http.MethodDelete, "/staff/missing", nil)
	req.Header.Set("X-Tenant-ID", "t1")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", w.Code)
	}
}

func TestDeleteStaff_Success(t *testing.T) {
	r := gin.New()
	h := handler.NewStaffHandler(&mockStaffStore{})
	r.DELETE("/staff/:id", h.DeleteStaff())

	req := httptest.NewRequest(http.MethodDelete, "/staff/s1", nil)
	req.Header.Set("X-Tenant-ID", "t1")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
}

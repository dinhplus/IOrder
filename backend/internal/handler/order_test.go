package handler_test

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"

	"github.com/dinhplus/IOrder/backend/internal/handler"
	"github.com/dinhplus/IOrder/backend/internal/repository"
)

type mockOrderStore struct {
	order *repository.Order
	items []*repository.OrderItem
	err   error
}

func (m *mockOrderStore) Create(ctx context.Context, o *repository.Order) error {
	if m.err != nil {
		return m.err
	}
	o.ID = "order-uuid"
	return nil
}
func (m *mockOrderStore) GetByID(ctx context.Context, id string) (*repository.Order, error) {
	return m.order, m.err
}
func (m *mockOrderStore) UpdateStatus(ctx context.Context, id, tenantID, status string) error {
	return m.err
}
func (m *mockOrderStore) UpdateItems(ctx context.Context, orderID, tenantID string, items []*repository.OrderItem) error {
	return m.err
}
func (m *mockOrderStore) ListByTenant(ctx context.Context, tenantID string) ([]*repository.Order, error) {
	if m.order != nil {
		return []*repository.Order{m.order}, m.err
	}
	return nil, m.err
}
func (m *mockOrderStore) CreateEvent(ctx context.Context, e *repository.OrderEvent) error {
	return nil
}
func (m *mockOrderStore) GetItems(ctx context.Context, orderID string) ([]*repository.OrderItem, error) {
	return m.items, m.err
}
func (m *mockOrderStore) UpdateItem(ctx context.Context, item *repository.OrderItem) error {
	return m.err
}

func TestCreateOrder_Success(t *testing.T) {
	r := gin.New()
	h := handler.NewOrderHandler(&mockOrderStore{})
	r.POST("/orders", h.CreateOrder())

	body := `{"tenant_id":"t1","table_id":"tbl1","session_id":"sess1"}`
	req := httptest.NewRequest(http.MethodPost, "/orders", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d: %s", w.Code, w.Body.String())
	}
}

func TestCreateOrder_MissingTenantID(t *testing.T) {
	r := gin.New()
	h := handler.NewOrderHandler(&mockOrderStore{})
	r.POST("/orders", h.CreateOrder())

	body := `{"table_id":"tbl1","session_id":"sess1"}`
	req := httptest.NewRequest(http.MethodPost, "/orders", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestListOrders_MissingTenantID(t *testing.T) {
	r := gin.New()
	h := handler.NewOrderHandler(&mockOrderStore{})
	r.GET("/orders", h.ListOrders())

	req := httptest.NewRequest(http.MethodGet, "/orders", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestGetOrder_NotFound(t *testing.T) {
	r := gin.New()
	h := handler.NewOrderHandler(&mockOrderStore{})
	r.GET("/orders/:id", h.GetOrder())

	req := httptest.NewRequest(http.MethodGet, "/orders/missing", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", w.Code)
	}
}

func TestGetOrder_Found(t *testing.T) {
	store := &mockOrderStore{order: &repository.Order{ID: "o1", TenantID: "t1", Status: repository.OrderStatusDraft}}
	r := gin.New()
	h := handler.NewOrderHandler(store)
	r.GET("/orders/:id", h.GetOrder())

	req := httptest.NewRequest(http.MethodGet, "/orders/o1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}

func TestSubmitOrder_InvalidStatus(t *testing.T) {
	store := &mockOrderStore{order: &repository.Order{ID: "o1", TenantID: "t1", Status: repository.OrderStatusConfirmed}}
	r := gin.New()
	h := handler.NewOrderHandler(store)
	r.POST("/orders/:id/submit", h.SubmitOrder())

	req := httptest.NewRequest(http.MethodPost, "/orders/o1/submit", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestSubmitOrder_Success(t *testing.T) {
	store := &mockOrderStore{order: &repository.Order{ID: "o1", TenantID: "t1", Status: repository.OrderStatusDraft}}
	r := gin.New()
	h := handler.NewOrderHandler(store)
	r.POST("/orders/:id/submit", h.SubmitOrder())

	req := httptest.NewRequest(http.MethodPost, "/orders/o1/submit", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
}

func TestCancelOrder_Success(t *testing.T) {
	store := &mockOrderStore{order: &repository.Order{ID: "o1", TenantID: "t1", Status: repository.OrderStatusDraft}}
	r := gin.New()
	h := handler.NewOrderHandler(store)
	r.POST("/orders/:id/cancel", h.CancelOrder())

	req := httptest.NewRequest(http.MethodPost, "/orders/o1/cancel", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}

func TestPayOrder_Success(t *testing.T) {
	store := &mockOrderStore{order: &repository.Order{ID: "o1", TenantID: "t1", Status: repository.OrderStatusPaymentRequested}}
	r := gin.New()
	h := handler.NewOrderHandler(store)
	r.POST("/orders/:id/pay", h.PayOrder())

	req := httptest.NewRequest(http.MethodPost, "/orders/o1/pay", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
}

func TestPayOrder_InvalidStatus(t *testing.T) {
	store := &mockOrderStore{order: &repository.Order{ID: "o1", TenantID: "t1", Status: repository.OrderStatusDraft}}
	r := gin.New()
	h := handler.NewOrderHandler(store)
	r.POST("/orders/:id/pay", h.PayOrder())

	req := httptest.NewRequest(http.MethodPost, "/orders/o1/pay", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestCloseOrder_Success(t *testing.T) {
	store := &mockOrderStore{order: &repository.Order{ID: "o1", TenantID: "t1", Status: repository.OrderStatusPaid}}
	r := gin.New()
	h := handler.NewOrderHandler(store)
	r.POST("/orders/:id/close", h.CloseOrder())

	req := httptest.NewRequest(http.MethodPost, "/orders/o1/close", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
}

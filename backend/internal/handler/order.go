package handler

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/dinhplus/IOrder/backend/internal/repository"
)

// OrderStore is the interface used by OrderHandler.
type OrderStore interface {
	Create(ctx context.Context, o *repository.Order) error
	GetByID(ctx context.Context, id string) (*repository.Order, error)
	UpdateStatus(ctx context.Context, id, tenantID, status string) error
	UpdateItems(ctx context.Context, orderID, tenantID string, items []*repository.OrderItem) error
	ListByTenant(ctx context.Context, tenantID string) ([]*repository.Order, error)
	CreateEvent(ctx context.Context, e *repository.OrderEvent) error
	GetItems(ctx context.Context, orderID string) ([]*repository.OrderItem, error)
	UpdateItem(ctx context.Context, item *repository.OrderItem) error
}

// OrderHandler handles order-related HTTP requests.
type OrderHandler struct {
	repo OrderStore
}

// NewOrderHandler creates an OrderHandler.
func NewOrderHandler(repo OrderStore) *OrderHandler {
	return &OrderHandler{repo: repo}
}

type createOrderRequest struct {
	TenantID  string  `json:"tenant_id"`
	TableID   string  `json:"table_id"`
	SessionID string  `json:"session_id"`
	Notes     *string `json:"notes"`
}

// CreateOrder handles POST /api/v1/orders.
func (h *OrderHandler) CreateOrder() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req createOrderRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			RespondError(c, http.StatusBadRequest, ErrValidation, "invalid request body")
			return
		}
		if req.TenantID == "" {
			RespondError(c, http.StatusBadRequest, ErrValidation, "tenant_id is required")
			return
		}
		if req.TableID == "" {
			RespondError(c, http.StatusBadRequest, ErrValidation, "table_id is required")
			return
		}
		if req.SessionID == "" {
			RespondError(c, http.StatusBadRequest, ErrValidation, "session_id is required")
			return
		}
		o := &repository.Order{
			TenantID:  req.TenantID,
			TableID:   req.TableID,
			SessionID: req.SessionID,
			Status:    repository.OrderStatusDraft,
			Notes:     req.Notes,
		}
		if err := h.repo.Create(c.Request.Context(), o); err != nil {
			RespondError(c, http.StatusInternalServerError, ErrInternal, "failed to create order")
			return
		}
		RespondSuccess(c, http.StatusCreated, o)
	}
}

// ListOrders handles GET /api/v1/orders.
func (h *OrderHandler) ListOrders() gin.HandlerFunc {
	return func(c *gin.Context) {
		tenantID, ok := tenantIDFromHeader(c)
		if !ok {
			return
		}
		orders, err := h.repo.ListByTenant(c.Request.Context(), tenantID)
		if err != nil {
			RespondError(c, http.StatusInternalServerError, ErrInternal, "failed to list orders")
			return
		}
		if orders == nil {
			orders = []*repository.Order{}
		}
		RespondSuccess(c, http.StatusOK, orders)
	}
}

// GetOrder handles GET /api/v1/orders/:id.
func (h *OrderHandler) GetOrder() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		o, err := h.repo.GetByID(c.Request.Context(), id)
		if err != nil {
			RespondError(c, http.StatusInternalServerError, ErrInternal, "failed to get order")
			return
		}
		if o == nil {
			RespondError(c, http.StatusNotFound, ErrNotFound, "order not found")
			return
		}
		RespondSuccess(c, http.StatusOK, o)
	}
}

type orderItemRequest struct {
	ItemID    string          `json:"item_id"`
	ItemName  string          `json:"item_name"`
	ItemPrice int64           `json:"item_price"`
	Quantity  int             `json:"quantity"`
	Modifiers json.RawMessage `json:"modifiers"`
	Notes     *string         `json:"notes"`
	RoutedTo  string          `json:"routed_to"`
}

type updateOrderItemsRequest struct {
	Items []orderItemRequest `json:"items"`
}

// UpdateOrderItems handles PATCH /api/v1/orders/:id/items.
func (h *OrderHandler) UpdateOrderItems() gin.HandlerFunc {
	return func(c *gin.Context) {
		tenantID, ok := tenantIDFromHeader(c)
		if !ok {
			return
		}
		id := c.Param("id")
		var req updateOrderItemsRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			RespondError(c, http.StatusBadRequest, ErrValidation, "invalid request body")
			return
		}
		items := make([]*repository.OrderItem, 0, len(req.Items))
		for _, ri := range req.Items {
			routedTo := ri.RoutedTo
			if routedTo == "" {
				routedTo = "kitchen"
			}
			qty := ri.Quantity
			if qty <= 0 {
				qty = 1
			}
			items = append(items, &repository.OrderItem{
				ItemID:    ri.ItemID,
				ItemName:  ri.ItemName,
				ItemPrice: ri.ItemPrice,
				Quantity:  qty,
				Modifiers: ri.Modifiers,
				Notes:     ri.Notes,
				Status:    "pending",
				RoutedTo:  routedTo,
			})
		}
		if err := h.repo.UpdateItems(c.Request.Context(), id, tenantID, items); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				RespondError(c, http.StatusNotFound, ErrNotFound, "order not found")
				return
			}
			RespondError(c, http.StatusInternalServerError, ErrInternal, "failed to update order items")
			return
		}
		RespondSuccess(c, http.StatusOK, gin.H{"id": id, "updated": true})
	}
}

// transition handles generic order status transitions.
func (h *OrderHandler) transition(fromStatus, toStatus string) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		o, err := h.repo.GetByID(c.Request.Context(), id)
		if err != nil {
			RespondError(c, http.StatusInternalServerError, ErrInternal, "failed to get order")
			return
		}
		if o == nil {
			RespondError(c, http.StatusNotFound, ErrNotFound, "order not found")
			return
		}
		if fromStatus != "" && o.Status != fromStatus {
			RespondError(c, http.StatusBadRequest, ErrValidation, "invalid order status for this transition")
			return
		}
		if err := h.repo.UpdateStatus(c.Request.Context(), id, o.TenantID, toStatus); err != nil {
			RespondError(c, http.StatusInternalServerError, ErrInternal, "failed to update order status")
			return
		}
		_ = h.repo.CreateEvent(c.Request.Context(), &repository.OrderEvent{
			OrderID:    id,
			TenantID:   o.TenantID,
			FromStatus: &o.Status,
			ToStatus:   toStatus,
		})
		o.Status = toStatus
		RespondSuccess(c, http.StatusOK, o)
	}
}

// SubmitOrder handles POST /api/v1/orders/:id/submit.
func (h *OrderHandler) SubmitOrder() gin.HandlerFunc {
	return h.transition(repository.OrderStatusDraft, repository.OrderStatusSubmitted)
}

// ConfirmOrder handles POST /api/v1/orders/:id/confirm.
func (h *OrderHandler) ConfirmOrder() gin.HandlerFunc {
	return h.transition(repository.OrderStatusSubmitted, repository.OrderStatusConfirmed)
}

// RejectOrder handles POST /api/v1/orders/:id/reject.
func (h *OrderHandler) RejectOrder() gin.HandlerFunc {
	return h.transition(repository.OrderStatusSubmitted, repository.OrderStatusRejected)
}

// StartPreparation handles POST /api/v1/orders/:id/start-preparation.
func (h *OrderHandler) StartPreparation() gin.HandlerFunc {
	return h.transition(repository.OrderStatusConfirmed, repository.OrderStatusInPreparation)
}

// MarkReady handles POST /api/v1/orders/:id/ready.
func (h *OrderHandler) MarkReady() gin.HandlerFunc {
	return h.transition(repository.OrderStatusInPreparation, repository.OrderStatusReady)
}

// MarkServed handles POST /api/v1/orders/:id/serve.
func (h *OrderHandler) MarkServed() gin.HandlerFunc {
	return h.transition(repository.OrderStatusReady, repository.OrderStatusServed)
}

// RequestPayment handles POST /api/v1/orders/:id/request-payment.
func (h *OrderHandler) RequestPayment() gin.HandlerFunc {
	return h.transition(repository.OrderStatusServed, repository.OrderStatusPaymentRequested)
}

// CancelOrder handles POST /api/v1/orders/:id/cancel.
func (h *OrderHandler) CancelOrder() gin.HandlerFunc {
	return h.transition("", repository.OrderStatusCancelled)
}

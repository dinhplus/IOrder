package handler

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/dinhplus/IOrder/backend/internal/repository"
)

// TenantStore is the interface used by TenantHandler (enables mocking in tests).
type TenantStore interface {
	Create(ctx context.Context, t *repository.Tenant) error
	GetByID(ctx context.Context, id string) (*repository.Tenant, error)
	GetBySlug(ctx context.Context, slug string) (*repository.Tenant, error)
	Update(ctx context.Context, t *repository.Tenant) error
	List(ctx context.Context) ([]*repository.Tenant, error)
}

// TenantHandler handles tenant-related HTTP requests.
type TenantHandler struct {
	repo TenantStore
}

// NewTenantHandler creates a TenantHandler.
func NewTenantHandler(repo TenantStore) *TenantHandler {
	return &TenantHandler{repo: repo}
}

type createTenantRequest struct {
	Slug     string          `json:"slug"`
	Name     string          `json:"name"`
	LogoURL  *string         `json:"logo_url"`
	Timezone string          `json:"timezone"`
	Currency string          `json:"currency"`
	Settings json.RawMessage `json:"settings"`
}

// CreateTenant handles POST /api/v1/tenants.
func (h *TenantHandler) CreateTenant() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req createTenantRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			RespondError(c, http.StatusBadRequest, ErrValidation, "invalid request body")
			return
		}
		if req.Slug == "" {
			RespondError(c, http.StatusBadRequest, ErrValidation, "slug is required")
			return
		}
		if req.Name == "" {
			RespondError(c, http.StatusBadRequest, ErrValidation, "name is required")
			return
		}

		timezone := req.Timezone
		if timezone == "" {
			timezone = "Asia/Ho_Chi_Minh"
		}
		currency := req.Currency
		if currency == "" {
			currency = "VND"
		}
		settings := req.Settings
		if settings == nil {
			settings = json.RawMessage("{}")
		}

		t := &repository.Tenant{
			Slug:     req.Slug,
			Name:     req.Name,
			LogoURL:  req.LogoURL,
			Timezone: timezone,
			Currency: currency,
			IsActive: true,
			Settings: settings,
		}
		if err := h.repo.Create(c.Request.Context(), t); err != nil {
			RespondError(c, http.StatusInternalServerError, ErrInternal, "failed to create tenant")
			return
		}
		RespondSuccess(c, http.StatusCreated, t)
	}
}

// GetTenant handles GET /api/v1/tenants/:id.
func (h *TenantHandler) GetTenant() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		t, err := h.repo.GetByID(c.Request.Context(), id)
		if err != nil {
			RespondError(c, http.StatusInternalServerError, ErrInternal, "failed to get tenant")
			return
		}
		if t == nil {
			RespondError(c, http.StatusNotFound, ErrNotFound, "tenant not found")
			return
		}
		RespondSuccess(c, http.StatusOK, t)
	}
}

type updateTenantRequest struct {
	Slug     string          `json:"slug"`
	Name     string          `json:"name"`
	LogoURL  *string         `json:"logo_url"`
	Timezone string          `json:"timezone"`
	Currency string          `json:"currency"`
	IsActive *bool           `json:"is_active"`
	Settings json.RawMessage `json:"settings"`
}

// UpdateTenant handles PATCH /api/v1/tenants/:id.
func (h *TenantHandler) UpdateTenant() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		existing, err := h.repo.GetByID(c.Request.Context(), id)
		if err != nil {
			RespondError(c, http.StatusInternalServerError, ErrInternal, "failed to get tenant")
			return
		}
		if existing == nil {
			RespondError(c, http.StatusNotFound, ErrNotFound, "tenant not found")
			return
		}

		var req updateTenantRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			RespondError(c, http.StatusBadRequest, ErrValidation, "invalid request body")
			return
		}

		if req.Slug != "" {
			existing.Slug = req.Slug
		}
		if req.Name != "" {
			existing.Name = req.Name
		}
		if req.LogoURL != nil {
			existing.LogoURL = req.LogoURL
		}
		if req.Timezone != "" {
			existing.Timezone = req.Timezone
		}
		if req.Currency != "" {
			existing.Currency = req.Currency
		}
		if req.IsActive != nil {
			existing.IsActive = *req.IsActive
		}
		if req.Settings != nil {
			existing.Settings = req.Settings
		}

		if err := h.repo.Update(c.Request.Context(), existing); err != nil {
			if err == sql.ErrNoRows {
				RespondError(c, http.StatusNotFound, ErrNotFound, "tenant not found")
				return
			}
			RespondError(c, http.StatusInternalServerError, ErrInternal, "failed to update tenant")
			return
		}
		RespondSuccess(c, http.StatusOK, existing)
	}
}

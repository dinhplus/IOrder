package handler

import (
	"context"
	"database/sql"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/dinhplus/IOrder/backend/internal/repository"
)

// MenuStore is the interface used by MenuHandler.
type MenuStore interface {
	CreateCategory(ctx context.Context, c *repository.MenuCategory) error
	GetCategory(ctx context.Context, id, tenantID string) (*repository.MenuCategory, error)
	UpdateCategory(ctx context.Context, c *repository.MenuCategory) error
	DeleteCategory(ctx context.Context, id, tenantID string) error
	ListCategories(ctx context.Context, tenantID string) ([]*repository.MenuCategory, error)
	CreateItem(ctx context.Context, item *repository.MenuItem) error
	GetItem(ctx context.Context, id, tenantID string) (*repository.MenuItem, error)
	UpdateItem(ctx context.Context, item *repository.MenuItem) error
	DeleteItem(ctx context.Context, id, tenantID string) error
	ListItems(ctx context.Context, tenantID string) ([]*repository.MenuItem, error)
	ToggleItemAvailability(ctx context.Context, id, tenantID string, available bool) error
}

// MenuHandler handles menu-related HTTP requests.
type MenuHandler struct {
	repo MenuStore
}

// NewMenuHandler creates a MenuHandler.
func NewMenuHandler(repo MenuStore) *MenuHandler {
	return &MenuHandler{repo: repo}
}

// tenantIDFromHeader extracts X-Tenant-ID from the request header.
func tenantIDFromHeader(c *gin.Context) (string, bool) {
	id := c.GetHeader("X-Tenant-ID")
	if id == "" {
		RespondError(c, http.StatusBadRequest, ErrValidation, "X-Tenant-ID header is required")
		return "", false
	}
	return id, true
}

// ListCategories handles GET /api/v1/menu/categories.
func (h *MenuHandler) ListCategories() gin.HandlerFunc {
	return func(c *gin.Context) {
		tenantID, ok := tenantIDFromHeader(c)
		if !ok {
			return
		}
		cats, err := h.repo.ListCategories(c.Request.Context(), tenantID)
		if err != nil {
			RespondError(c, http.StatusInternalServerError, ErrInternal, "failed to list categories")
			return
		}
		if cats == nil {
			cats = []*repository.MenuCategory{}
		}
		RespondSuccess(c, http.StatusOK, cats)
	}
}

type createCategoryRequest struct {
	Name        string  `json:"name"`
	Type        string  `json:"type"`
	Description *string `json:"description"`
	SortOrder   int     `json:"sort_order"`
}

// CreateCategory handles POST /api/v1/menu/categories.
func (h *MenuHandler) CreateCategory() gin.HandlerFunc {
	return func(c *gin.Context) {
		tenantID, ok := tenantIDFromHeader(c)
		if !ok {
			return
		}
		var req createCategoryRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			RespondError(c, http.StatusBadRequest, ErrValidation, "invalid request body")
			return
		}
		if req.Name == "" {
			RespondError(c, http.StatusBadRequest, ErrValidation, "name is required")
			return
		}
		catType := req.Type
		if catType == "" {
			catType = "food"
		}
		cat := &repository.MenuCategory{
			TenantID:    tenantID,
			Name:        req.Name,
			Type:        catType,
			Description: req.Description,
			SortOrder:   req.SortOrder,
			IsActive:    true,
		}
		if err := h.repo.CreateCategory(c.Request.Context(), cat); err != nil {
			RespondError(c, http.StatusInternalServerError, ErrInternal, "failed to create category")
			return
		}
		RespondSuccess(c, http.StatusCreated, cat)
	}
}

type updateCategoryRequest struct {
	Name        string  `json:"name"`
	Type        string  `json:"type"`
	Description *string `json:"description"`
	SortOrder   *int    `json:"sort_order"`
	IsActive    *bool   `json:"is_active"`
}

// UpdateCategory handles PUT /api/v1/menu/categories/:id.
func (h *MenuHandler) UpdateCategory() gin.HandlerFunc {
	return func(c *gin.Context) {
		tenantID, ok := tenantIDFromHeader(c)
		if !ok {
			return
		}
		id := c.Param("id")
		existing, err := h.repo.GetCategory(c.Request.Context(), id, tenantID)
		if err != nil {
			RespondError(c, http.StatusInternalServerError, ErrInternal, "failed to get category")
			return
		}
		if existing == nil {
			RespondError(c, http.StatusNotFound, ErrNotFound, "category not found")
			return
		}
		var req updateCategoryRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			RespondError(c, http.StatusBadRequest, ErrValidation, "invalid request body")
			return
		}
		if req.Name != "" {
			existing.Name = req.Name
		}
		if req.Type != "" {
			existing.Type = req.Type
		}
		if req.Description != nil {
			existing.Description = req.Description
		}
		if req.SortOrder != nil {
			existing.SortOrder = *req.SortOrder
		}
		if req.IsActive != nil {
			existing.IsActive = *req.IsActive
		}
		if err := h.repo.UpdateCategory(c.Request.Context(), existing); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				RespondError(c, http.StatusNotFound, ErrNotFound, "category not found")
				return
			}
			RespondError(c, http.StatusInternalServerError, ErrInternal, "failed to update category")
			return
		}
		RespondSuccess(c, http.StatusOK, existing)
	}
}

// DeleteCategory handles DELETE /api/v1/menu/categories/:id.
func (h *MenuHandler) DeleteCategory() gin.HandlerFunc {
	return func(c *gin.Context) {
		tenantID, ok := tenantIDFromHeader(c)
		if !ok {
			return
		}
		id := c.Param("id")
		if err := h.repo.DeleteCategory(c.Request.Context(), id, tenantID); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				RespondError(c, http.StatusNotFound, ErrNotFound, "category not found")
				return
			}
			RespondError(c, http.StatusInternalServerError, ErrInternal, "failed to delete category")
			return
		}
		c.Status(http.StatusNoContent)
	}
}

// ListItems handles GET /api/v1/menu/items.
func (h *MenuHandler) ListItems() gin.HandlerFunc {
	return func(c *gin.Context) {
		tenantID, ok := tenantIDFromHeader(c)
		if !ok {
			return
		}
		items, err := h.repo.ListItems(c.Request.Context(), tenantID)
		if err != nil {
			RespondError(c, http.StatusInternalServerError, ErrInternal, "failed to list items")
			return
		}
		if items == nil {
			items = []*repository.MenuItem{}
		}
		RespondSuccess(c, http.StatusOK, items)
	}
}

type createItemRequest struct {
	CategoryID  string  `json:"category_id"`
	Name        string  `json:"name"`
	Description *string `json:"description"`
	Price       int64   `json:"price"`
	ImageURL    *string `json:"image_url"`
	SortOrder   int     `json:"sort_order"`
}

// CreateItem handles POST /api/v1/menu/items.
func (h *MenuHandler) CreateItem() gin.HandlerFunc {
	return func(c *gin.Context) {
		tenantID, ok := tenantIDFromHeader(c)
		if !ok {
			return
		}
		var req createItemRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			RespondError(c, http.StatusBadRequest, ErrValidation, "invalid request body")
			return
		}
		if req.Name == "" {
			RespondError(c, http.StatusBadRequest, ErrValidation, "name is required")
			return
		}
		if req.CategoryID == "" {
			RespondError(c, http.StatusBadRequest, ErrValidation, "category_id is required")
			return
		}
		if req.Price < 0 {
			RespondError(c, http.StatusBadRequest, ErrValidation, "price must be non-negative")
			return
		}
		item := &repository.MenuItem{
			TenantID:    tenantID,
			CategoryID:  req.CategoryID,
			Name:        req.Name,
			Description: req.Description,
			Price:       req.Price,
			ImageURL:    req.ImageURL,
			IsAvailable: true,
			SortOrder:   req.SortOrder,
		}
		if err := h.repo.CreateItem(c.Request.Context(), item); err != nil {
			RespondError(c, http.StatusInternalServerError, ErrInternal, "failed to create item")
			return
		}
		RespondSuccess(c, http.StatusCreated, item)
	}
}

type updateItemRequest struct {
	CategoryID  string  `json:"category_id"`
	Name        string  `json:"name"`
	Description *string `json:"description"`
	Price       *int64  `json:"price"`
	ImageURL    *string `json:"image_url"`
	IsAvailable *bool   `json:"is_available"`
	SortOrder   *int    `json:"sort_order"`
}

// UpdateItem handles PUT /api/v1/menu/items/:id.
func (h *MenuHandler) UpdateItem() gin.HandlerFunc {
	return func(c *gin.Context) {
		tenantID, ok := tenantIDFromHeader(c)
		if !ok {
			return
		}
		id := c.Param("id")
		existing, err := h.repo.GetItem(c.Request.Context(), id, tenantID)
		if err != nil {
			RespondError(c, http.StatusInternalServerError, ErrInternal, "failed to get item")
			return
		}
		if existing == nil {
			RespondError(c, http.StatusNotFound, ErrNotFound, "item not found")
			return
		}
		var req updateItemRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			RespondError(c, http.StatusBadRequest, ErrValidation, "invalid request body")
			return
		}
		if req.CategoryID != "" {
			existing.CategoryID = req.CategoryID
		}
		if req.Name != "" {
			existing.Name = req.Name
		}
		if req.Description != nil {
			existing.Description = req.Description
		}
		if req.Price != nil {
			existing.Price = *req.Price
		}
		if req.ImageURL != nil {
			existing.ImageURL = req.ImageURL
		}
		if req.IsAvailable != nil {
			existing.IsAvailable = *req.IsAvailable
		}
		if req.SortOrder != nil {
			existing.SortOrder = *req.SortOrder
		}
		if err := h.repo.UpdateItem(c.Request.Context(), existing); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				RespondError(c, http.StatusNotFound, ErrNotFound, "item not found")
				return
			}
			RespondError(c, http.StatusInternalServerError, ErrInternal, "failed to update item")
			return
		}
		RespondSuccess(c, http.StatusOK, existing)
	}
}

type toggleAvailabilityRequest struct {
	IsAvailable bool `json:"is_available"`
}

// ToggleAvailability handles PATCH /api/v1/menu/items/:id/availability.
func (h *MenuHandler) ToggleAvailability() gin.HandlerFunc {
	return func(c *gin.Context) {
		tenantID, ok := tenantIDFromHeader(c)
		if !ok {
			return
		}
		id := c.Param("id")
		var req toggleAvailabilityRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			RespondError(c, http.StatusBadRequest, ErrValidation, "invalid request body")
			return
		}
		if err := h.repo.ToggleItemAvailability(c.Request.Context(), id, tenantID, req.IsAvailable); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				RespondError(c, http.StatusNotFound, ErrNotFound, "item not found")
				return
			}
			RespondError(c, http.StatusInternalServerError, ErrInternal, "failed to update item availability")
			return
		}
		RespondSuccess(c, http.StatusOK, gin.H{"id": id, "is_available": req.IsAvailable})
	}
}

// DeleteItem handles DELETE /api/v1/menu/items/:id.
func (h *MenuHandler) DeleteItem() gin.HandlerFunc {
	return func(c *gin.Context) {
		tenantID, ok := tenantIDFromHeader(c)
		if !ok {
			return
		}
		id := c.Param("id")
		if err := h.repo.DeleteItem(c.Request.Context(), id, tenantID); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				RespondError(c, http.StatusNotFound, ErrNotFound, "item not found")
				return
			}
			RespondError(c, http.StatusInternalServerError, ErrInternal, "failed to delete item")
			return
		}
		c.Status(http.StatusNoContent)
	}
}

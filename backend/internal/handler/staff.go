package handler

import (
	"context"
	"database/sql"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/dinhplus/IOrder/backend/internal/repository"
)

// StaffStore is the interface used by StaffHandler (enables mocking in tests).
type StaffStore interface {
	Create(ctx context.Context, s *repository.Staff) error
	GetByID(ctx context.Context, id, tenantID string) (*repository.Staff, error)
	ListByTenant(ctx context.Context, tenantID string) ([]*repository.Staff, error)
	Update(ctx context.Context, s *repository.Staff) error
	Delete(ctx context.Context, id, tenantID string) error
}

// StaffHandler handles staff-related HTTP requests.
type StaffHandler struct {
	repo StaffStore
}

// NewStaffHandler creates a StaffHandler.
func NewStaffHandler(repo StaffStore) *StaffHandler {
	return &StaffHandler{repo: repo}
}

type createStaffRequest struct {
	FullName string  `json:"full_name"`
	Email    *string `json:"email"`
	Role     string  `json:"role"`
}

// CreateStaff handles POST /api/v1/staff.
func (h *StaffHandler) CreateStaff() gin.HandlerFunc {
	return func(c *gin.Context) {
		tenantID, ok := tenantIDFromHeader(c)
		if !ok {
			return
		}
		var req createStaffRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			RespondError(c, http.StatusBadRequest, ErrValidation, "invalid request body")
			return
		}
		if req.FullName == "" {
			RespondError(c, http.StatusBadRequest, ErrValidation, "full_name is required")
			return
		}
		if req.Role == "" {
			RespondError(c, http.StatusBadRequest, ErrValidation, "role is required")
			return
		}
		s := &repository.Staff{
			TenantID: tenantID,
			FullName: req.FullName,
			Email:    req.Email,
			Role:     req.Role,
			IsActive: true,
		}
		if err := h.repo.Create(c.Request.Context(), s); err != nil {
			RespondError(c, http.StatusInternalServerError, ErrInternal, "failed to create staff")
			return
		}
		RespondSuccess(c, http.StatusCreated, s)
	}
}

// ListStaff handles GET /api/v1/staff.
func (h *StaffHandler) ListStaff() gin.HandlerFunc {
	return func(c *gin.Context) {
		tenantID, ok := tenantIDFromHeader(c)
		if !ok {
			return
		}
		staff, err := h.repo.ListByTenant(c.Request.Context(), tenantID)
		if err != nil {
			RespondError(c, http.StatusInternalServerError, ErrInternal, "failed to list staff")
			return
		}
		if staff == nil {
			staff = []*repository.Staff{}
		}
		RespondSuccess(c, http.StatusOK, staff)
	}
}

// GetStaff handles GET /api/v1/staff/:id.
func (h *StaffHandler) GetStaff() gin.HandlerFunc {
	return func(c *gin.Context) {
		tenantID, ok := tenantIDFromHeader(c)
		if !ok {
			return
		}
		id := c.Param("id")
		s, err := h.repo.GetByID(c.Request.Context(), id, tenantID)
		if err != nil {
			RespondError(c, http.StatusInternalServerError, ErrInternal, "failed to get staff")
			return
		}
		if s == nil {
			RespondError(c, http.StatusNotFound, ErrNotFound, "staff not found")
			return
		}
		RespondSuccess(c, http.StatusOK, s)
	}
}

type updateStaffRequest struct {
	FullName string  `json:"full_name"`
	Email    *string `json:"email"`
	Role     string  `json:"role"`
	IsActive *bool   `json:"is_active"`
}

// UpdateStaff handles PATCH /api/v1/staff/:id.
func (h *StaffHandler) UpdateStaff() gin.HandlerFunc {
	return func(c *gin.Context) {
		tenantID, ok := tenantIDFromHeader(c)
		if !ok {
			return
		}
		id := c.Param("id")
		existing, err := h.repo.GetByID(c.Request.Context(), id, tenantID)
		if err != nil {
			RespondError(c, http.StatusInternalServerError, ErrInternal, "failed to get staff")
			return
		}
		if existing == nil {
			RespondError(c, http.StatusNotFound, ErrNotFound, "staff not found")
			return
		}
		var req updateStaffRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			RespondError(c, http.StatusBadRequest, ErrValidation, "invalid request body")
			return
		}
		if req.FullName != "" {
			existing.FullName = req.FullName
		}
		if req.Email != nil {
			existing.Email = req.Email
		}
		if req.Role != "" {
			existing.Role = req.Role
		}
		if req.IsActive != nil {
			existing.IsActive = *req.IsActive
		}
		if err := h.repo.Update(c.Request.Context(), existing); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				RespondError(c, http.StatusNotFound, ErrNotFound, "staff not found")
				return
			}
			RespondError(c, http.StatusInternalServerError, ErrInternal, "failed to update staff")
			return
		}
		RespondSuccess(c, http.StatusOK, existing)
	}
}

// DeleteStaff handles DELETE /api/v1/staff/:id.
func (h *StaffHandler) DeleteStaff() gin.HandlerFunc {
	return func(c *gin.Context) {
		tenantID, ok := tenantIDFromHeader(c)
		if !ok {
			return
		}
		id := c.Param("id")
		if err := h.repo.Delete(c.Request.Context(), id, tenantID); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				RespondError(c, http.StatusNotFound, ErrNotFound, "staff not found")
				return
			}
			RespondError(c, http.StatusInternalServerError, ErrInternal, "failed to delete staff")
			return
		}
		RespondSuccess(c, http.StatusOK, gin.H{"id": id, "deleted": true})
	}
}

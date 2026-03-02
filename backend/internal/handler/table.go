package handler

import (
	"context"
	"database/sql"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/dinhplus/IOrder/backend/internal/repository"
)

// TableStore is the interface used by TableHandler.
type TableStore interface {
	CreateFloorPlan(ctx context.Context, fp *repository.FloorPlan) error
	GetFloorPlan(ctx context.Context, id, tenantID string) (*repository.FloorPlan, error)
	UpdateFloorPlan(ctx context.Context, fp *repository.FloorPlan) error
	ListFloorPlans(ctx context.Context, tenantID string) ([]*repository.FloorPlan, error)
	ListTablesByFloorPlan(ctx context.Context, floorPlanID, tenantID string) ([]*repository.RestaurantTable, error)
	CreateTable(ctx context.Context, t *repository.RestaurantTable) error
	GetTable(ctx context.Context, id, tenantID string) (*repository.RestaurantTable, error)
	UpdateTable(ctx context.Context, t *repository.RestaurantTable) error
	DeleteTable(ctx context.Context, id, tenantID string) error
}

// TableHandler handles floor-plan and table HTTP requests.
type TableHandler struct {
	repo TableStore
}

// NewTableHandler creates a TableHandler.
func NewTableHandler(repo TableStore) *TableHandler {
	return &TableHandler{repo: repo}
}

// ListFloorPlans handles GET /api/v1/floor-plans.
func (h *TableHandler) ListFloorPlans() gin.HandlerFunc {
	return func(c *gin.Context) {
		tenantID, ok := tenantIDFromHeader(c)
		if !ok {
			return
		}
		fps, err := h.repo.ListFloorPlans(c.Request.Context(), tenantID)
		if err != nil {
			RespondError(c, http.StatusInternalServerError, ErrInternal, "failed to list floor plans")
			return
		}
		if fps == nil {
			fps = []*repository.FloorPlan{}
		}
		RespondSuccess(c, http.StatusOK, fps)
	}
}

type createFloorPlanRequest struct {
	Name       string `json:"name"`
	FloorLevel int    `json:"floor_level"`
}

// CreateFloorPlan handles POST /api/v1/floor-plans.
func (h *TableHandler) CreateFloorPlan() gin.HandlerFunc {
	return func(c *gin.Context) {
		tenantID, ok := tenantIDFromHeader(c)
		if !ok {
			return
		}
		var req createFloorPlanRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			RespondError(c, http.StatusBadRequest, ErrValidation, "invalid request body")
			return
		}
		if req.Name == "" {
			RespondError(c, http.StatusBadRequest, ErrValidation, "name is required")
			return
		}
		level := req.FloorLevel
		if level == 0 {
			level = 1
		}
		fp := &repository.FloorPlan{
			TenantID:   tenantID,
			Name:       req.Name,
			FloorLevel: level,
			IsActive:   true,
		}
		if err := h.repo.CreateFloorPlan(c.Request.Context(), fp); err != nil {
			RespondError(c, http.StatusInternalServerError, ErrInternal, "failed to create floor plan")
			return
		}
		RespondSuccess(c, http.StatusCreated, fp)
	}
}

type updateFloorPlanRequest struct {
	Name       string `json:"name"`
	FloorLevel int    `json:"floor_level"`
	IsActive   *bool  `json:"is_active"`
}

// UpdateFloorPlan handles PUT /api/v1/floor-plans/:id.
func (h *TableHandler) UpdateFloorPlan() gin.HandlerFunc {
	return func(c *gin.Context) {
		tenantID, ok := tenantIDFromHeader(c)
		if !ok {
			return
		}
		id := c.Param("id")
		existing, err := h.repo.GetFloorPlan(c.Request.Context(), id, tenantID)
		if err != nil {
			RespondError(c, http.StatusInternalServerError, ErrInternal, "failed to get floor plan")
			return
		}
		if existing == nil {
			RespondError(c, http.StatusNotFound, ErrNotFound, "floor plan not found")
			return
		}
		var req updateFloorPlanRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			RespondError(c, http.StatusBadRequest, ErrValidation, "invalid request body")
			return
		}
		if req.Name != "" {
			existing.Name = req.Name
		}
		if req.FloorLevel != 0 {
			existing.FloorLevel = req.FloorLevel
		}
		if req.IsActive != nil {
			existing.IsActive = *req.IsActive
		}
		if err := h.repo.UpdateFloorPlan(c.Request.Context(), existing); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				RespondError(c, http.StatusNotFound, ErrNotFound, "floor plan not found")
				return
			}
			RespondError(c, http.StatusInternalServerError, ErrInternal, "failed to update floor plan")
			return
		}
		RespondSuccess(c, http.StatusOK, existing)
	}
}

// ListTablesByFloorPlan handles GET /api/v1/floor-plans/:id/tables.
func (h *TableHandler) ListTablesByFloorPlan() gin.HandlerFunc {
	return func(c *gin.Context) {
		tenantID, ok := tenantIDFromHeader(c)
		if !ok {
			return
		}
		id := c.Param("id")
		tables, err := h.repo.ListTablesByFloorPlan(c.Request.Context(), id, tenantID)
		if err != nil {
			RespondError(c, http.StatusInternalServerError, ErrInternal, "failed to list tables")
			return
		}
		if tables == nil {
			tables = []*repository.RestaurantTable{}
		}
		RespondSuccess(c, http.StatusOK, tables)
	}
}

type createTableRequest struct {
	FloorPlanID string  `json:"floor_plan_id"`
	Name        string  `json:"name"`
	Capacity    int     `json:"capacity"`
	PosX        float64 `json:"pos_x"`
	PosY        float64 `json:"pos_y"`
	Shape       string  `json:"shape"`
}

// CreateTable handles POST /api/v1/tables.
func (h *TableHandler) CreateTable() gin.HandlerFunc {
	return func(c *gin.Context) {
		tenantID, ok := tenantIDFromHeader(c)
		if !ok {
			return
		}
		var req createTableRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			RespondError(c, http.StatusBadRequest, ErrValidation, "invalid request body")
			return
		}
		if req.Name == "" {
			RespondError(c, http.StatusBadRequest, ErrValidation, "name is required")
			return
		}
		if req.FloorPlanID == "" {
			RespondError(c, http.StatusBadRequest, ErrValidation, "floor_plan_id is required")
			return
		}
		capacity := req.Capacity
		if capacity == 0 {
			capacity = 4
		}
		shape := req.Shape
		if shape == "" {
			shape = "rectangle"
		}
		t := &repository.RestaurantTable{
			TenantID:    tenantID,
			FloorPlanID: req.FloorPlanID,
			Name:        req.Name,
			Capacity:    capacity,
			PosX:        req.PosX,
			PosY:        req.PosY,
			Shape:       shape,
			Status:      "available",
		}
		if err := h.repo.CreateTable(c.Request.Context(), t); err != nil {
			RespondError(c, http.StatusInternalServerError, ErrInternal, "failed to create table")
			return
		}
		RespondSuccess(c, http.StatusCreated, t)
	}
}

type updateTableRequest struct {
	FloorPlanID string   `json:"floor_plan_id"`
	Name        string   `json:"name"`
	Capacity    int      `json:"capacity"`
	PosX        *float64 `json:"pos_x"`
	PosY        *float64 `json:"pos_y"`
	Shape       string   `json:"shape"`
	Status      string   `json:"status"`
}

// UpdateTable handles PUT /api/v1/tables/:id.
func (h *TableHandler) UpdateTable() gin.HandlerFunc {
	return func(c *gin.Context) {
		tenantID, ok := tenantIDFromHeader(c)
		if !ok {
			return
		}
		id := c.Param("id")
		existing, err := h.repo.GetTable(c.Request.Context(), id, tenantID)
		if err != nil {
			RespondError(c, http.StatusInternalServerError, ErrInternal, "failed to get table")
			return
		}
		if existing == nil {
			RespondError(c, http.StatusNotFound, ErrNotFound, "table not found")
			return
		}
		var req updateTableRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			RespondError(c, http.StatusBadRequest, ErrValidation, "invalid request body")
			return
		}
		if req.FloorPlanID != "" {
			existing.FloorPlanID = req.FloorPlanID
		}
		if req.Name != "" {
			existing.Name = req.Name
		}
		if req.Capacity != 0 {
			existing.Capacity = req.Capacity
		}
		if req.PosX != nil {
			existing.PosX = *req.PosX
		}
		if req.PosY != nil {
			existing.PosY = *req.PosY
		}
		if req.Shape != "" {
			existing.Shape = req.Shape
		}
		if req.Status != "" {
			existing.Status = req.Status
		}
		if err := h.repo.UpdateTable(c.Request.Context(), existing); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				RespondError(c, http.StatusNotFound, ErrNotFound, "table not found")
				return
			}
			RespondError(c, http.StatusInternalServerError, ErrInternal, "failed to update table")
			return
		}
		RespondSuccess(c, http.StatusOK, existing)
	}
}

// DeleteTable handles DELETE /api/v1/tables/:id.
func (h *TableHandler) DeleteTable() gin.HandlerFunc {
	return func(c *gin.Context) {
		tenantID, ok := tenantIDFromHeader(c)
		if !ok {
			return
		}
		id := c.Param("id")
		if err := h.repo.DeleteTable(c.Request.Context(), id, tenantID); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				RespondError(c, http.StatusNotFound, ErrNotFound, "table not found")
				return
			}
			RespondError(c, http.StatusInternalServerError, ErrInternal, "failed to delete table")
			return
		}
		c.Status(http.StatusNoContent)
	}
}

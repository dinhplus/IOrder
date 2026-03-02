package handler

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
)

const appVersion = "0.1.0"

// HealthHandler serves the /health endpoint.
type HealthHandler struct {
	db *sql.DB
}

// NewHealthHandler creates a HealthHandler. db may be nil when the database is
// unavailable at startup; the handler degrades gracefully in that case.
func NewHealthHandler(db *sql.DB) *HealthHandler {
	return &HealthHandler{db: db}
}

// Check returns the health status of the service.
func (h *HealthHandler) Check() gin.HandlerFunc {
	return func(c *gin.Context) {
		dbStatus := "ok"
		if h.db == nil {
			dbStatus = "unavailable"
		} else if err := h.db.PingContext(c.Request.Context()); err != nil {
			dbStatus = "unavailable"
		}

		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"version": appVersion,
			"db":      dbStatus,
		})
	}
}

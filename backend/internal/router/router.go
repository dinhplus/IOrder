// Package router wires together the Gin engine, middleware, and route groups.
package router

import (
	"database/sql"
	"os"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"

	"github.com/dinhplus/IOrder/backend/internal/config"
	"github.com/dinhplus/IOrder/backend/internal/handler"
	"github.com/dinhplus/IOrder/backend/internal/middleware"
)

// New constructs and returns a fully configured *gin.Engine.
func New(cfg *config.Config, db *sql.DB) *gin.Engine {
	if cfg.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()

	// Core middleware stack.
	r.Use(requestid.New())
	r.Use(middleware.RequestLogger())
	r.Use(middleware.JSONRecovery())
	r.Use(buildCORS())

	// Health endpoint.
	healthHandler := handler.NewHealthHandler(db)
	r.GET("/health", healthHandler.Check())

	// API v1 group — routes registered by feature packages in future PRs.
	_ = r.Group("/api/v1")

	return r
}

func buildCORS() gin.HandlerFunc {
	originsEnv := os.Getenv("CORS_ORIGINS")
	var origins []string
	if originsEnv == "" || originsEnv == "*" {
		origins = []string{"*"}
	} else {
		for _, o := range strings.Split(originsEnv, ",") {
			if trimmed := strings.TrimSpace(o); trimmed != "" {
				origins = append(origins, trimmed)
			}
		}
	}

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = origins
	if len(origins) == 1 && origins[0] == "*" {
		corsConfig.AllowAllOrigins = true
		corsConfig.AllowOrigins = nil
	}
	return cors.New(corsConfig)
}

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
	"github.com/dinhplus/IOrder/backend/internal/repository"
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

	// API v1 group.
	tenantRepo := repository.NewTenantRepository(db)
	menuRepo := repository.NewMenuRepository(db)
	tableRepo := repository.NewTableRepository(db)
	orderRepo := repository.NewOrderRepository(db)

	tenantH := handler.NewTenantHandler(tenantRepo)
	menuH := handler.NewMenuHandler(menuRepo)
	tableH := handler.NewTableHandler(tableRepo)
	orderH := handler.NewOrderHandler(orderRepo)

	v1 := r.Group("/api/v1")
	{
		// Tenant routes
		v1.POST("/tenants", tenantH.CreateTenant())
		v1.GET("/tenants/:id", tenantH.GetTenant())
		v1.PATCH("/tenants/:id", tenantH.UpdateTenant())

		// Menu routes
		v1.GET("/menu/categories", menuH.ListCategories())
		v1.POST("/menu/categories", menuH.CreateCategory())
		v1.PUT("/menu/categories/:id", menuH.UpdateCategory())
		v1.DELETE("/menu/categories/:id", menuH.DeleteCategory())
		v1.GET("/menu/items", menuH.ListItems())
		v1.POST("/menu/items", menuH.CreateItem())
		v1.PUT("/menu/items/:id", menuH.UpdateItem())
		v1.PATCH("/menu/items/:id/availability", menuH.ToggleAvailability())
		v1.DELETE("/menu/items/:id", menuH.DeleteItem())

		// Floor plan & table routes
		v1.GET("/floor-plans", tableH.ListFloorPlans())
		v1.POST("/floor-plans", tableH.CreateFloorPlan())
		v1.PUT("/floor-plans/:id", tableH.UpdateFloorPlan())
		v1.GET("/floor-plans/:id/tables", tableH.ListTablesByFloorPlan())
		v1.POST("/tables", tableH.CreateTable())
		v1.PUT("/tables/:id", tableH.UpdateTable())
		v1.DELETE("/tables/:id", tableH.DeleteTable())

		// Order routes
		v1.POST("/orders", orderH.CreateOrder())
		v1.GET("/orders", orderH.ListOrders())
		v1.GET("/orders/:id", orderH.GetOrder())
		v1.PATCH("/orders/:id/items", orderH.UpdateOrderItems())
		v1.POST("/orders/:id/submit", orderH.SubmitOrder())
		v1.POST("/orders/:id/confirm", orderH.ConfirmOrder())
		v1.POST("/orders/:id/reject", orderH.RejectOrder())
		v1.POST("/orders/:id/start-preparation", orderH.StartPreparation())
		v1.POST("/orders/:id/ready", orderH.MarkReady())
		v1.POST("/orders/:id/serve", orderH.MarkServed())
		v1.POST("/orders/:id/request-payment", orderH.RequestPayment())
		v1.POST("/orders/:id/cancel", orderH.CancelOrder())
	}

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

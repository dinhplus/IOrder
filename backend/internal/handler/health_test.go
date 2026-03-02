package handler_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"

	"github.com/dinhplus/IOrder/backend/internal/handler"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func TestHealthCheck_NoDB_Returns200(t *testing.T) {
	r := gin.New()
	h := handler.NewHealthHandler(nil)
	r.GET("/health", h.Check())

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}

func TestHealthCheck_NoDB_StatusOK(t *testing.T) {
	r := gin.New()
	h := handler.NewHealthHandler(nil)
	r.GET("/health", h.Check())

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	var body map[string]string
	if err := json.NewDecoder(w.Body).Decode(&body); err != nil {
		t.Fatalf("decode response: %v", err)
	}

	if body["status"] != "ok" {
		t.Errorf("expected status 'ok', got %q", body["status"])
	}
}

func TestHealthCheck_NoDB_DBUnavailable(t *testing.T) {
	r := gin.New()
	h := handler.NewHealthHandler(nil)
	r.GET("/health", h.Check())

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	var body map[string]string
	if err := json.NewDecoder(w.Body).Decode(&body); err != nil {
		t.Fatalf("decode response: %v", err)
	}

	if body["db"] != "unavailable" {
		t.Errorf("expected db 'unavailable', got %q", body["db"])
	}
}

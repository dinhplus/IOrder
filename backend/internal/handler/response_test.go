package handler_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"

	"github.com/dinhplus/IOrder/backend/internal/handler"
)

func TestRespondError_Structure(t *testing.T) {
	r := gin.New()
	r.GET("/test", func(c *gin.Context) {
		handler.RespondError(c, http.StatusBadRequest, handler.ErrValidation, "invalid input",
			handler.ErrorDetail{Field: "name", Message: "required"},
		)
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}

	var body handler.APIError
	if err := json.NewDecoder(w.Body).Decode(&body); err != nil {
		t.Fatalf("decode response: %v", err)
	}

	if body.Error.Code != handler.ErrValidation {
		t.Errorf("expected code %q, got %q", handler.ErrValidation, body.Error.Code)
	}
	if body.Error.Message != "invalid input" {
		t.Errorf("unexpected message: %q", body.Error.Message)
	}
	if len(body.Error.Details) != 1 || body.Error.Details[0].Field != "name" {
		t.Errorf("unexpected details: %+v", body.Error.Details)
	}
}

func TestRespondSuccess_Structure(t *testing.T) {
	r := gin.New()
	r.GET("/test", func(c *gin.Context) {
		handler.RespondSuccess(c, http.StatusOK, gin.H{"id": 1})
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	var body map[string]any
	if err := json.NewDecoder(w.Body).Decode(&body); err != nil {
		t.Fatalf("decode response: %v", err)
	}

	if _, ok := body["data"]; !ok {
		t.Errorf("expected 'data' key in response, got: %+v", body)
	}
}

func TestRespondError_NoDetails(t *testing.T) {
	r := gin.New()
	r.GET("/test", func(c *gin.Context) {
		handler.RespondError(c, http.StatusNotFound, handler.ErrNotFound, "not found")
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", w.Code)
	}

	var body handler.APIError
	if err := json.NewDecoder(w.Body).Decode(&body); err != nil {
		t.Fatalf("decode response: %v", err)
	}

	if body.Error.Details != nil {
		t.Errorf("expected nil details, got: %+v", body.Error.Details)
	}
}

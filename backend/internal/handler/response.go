// Package handler contains HTTP handlers and shared response helpers.
package handler

import (
	"github.com/gin-gonic/gin"
)

// ErrorCode is a machine-readable error identifier returned in API error responses.
type ErrorCode string

const (
	ErrValidation   ErrorCode = "VALIDATION_ERROR"
	ErrNotFound     ErrorCode = "NOT_FOUND"
	ErrUnauthorized ErrorCode = "UNAUTHORIZED"
	ErrInternal     ErrorCode = "INTERNAL_ERROR"
)

// ErrorDetail describes a single field-level validation error.
type ErrorDetail struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// ErrorBody is the structured error payload.
type ErrorBody struct {
	Code    ErrorCode     `json:"code"`
	Message string        `json:"message"`
	Details []ErrorDetail `json:"details,omitempty"`
}

// APIError wraps ErrorBody as the top-level JSON object for error responses.
type APIError struct {
	Error ErrorBody `json:"error"`
}

// successBody wraps arbitrary data as the top-level JSON object for success responses.
type successBody struct {
	Data any `json:"data"`
}

// RespondError writes a JSON error response.
func RespondError(c *gin.Context, status int, code ErrorCode, message string, details ...ErrorDetail) {
	body := APIError{
		Error: ErrorBody{
			Code:    code,
			Message: message,
		},
	}
	if len(details) > 0 {
		body.Error.Details = details
	}
	c.JSON(status, body)
}

// RespondSuccess writes a JSON success response wrapping data.
func RespondSuccess(c *gin.Context, status int, data any) {
	c.JSON(status, successBody{Data: data})
}

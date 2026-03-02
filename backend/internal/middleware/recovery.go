package middleware

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/dinhplus/IOrder/backend/internal/handler"
)

// JSONRecovery recovers from panics and returns a structured JSON error response.
func JSONRecovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				slog.Error("panic recovered",
					"error", fmt.Sprintf("%v", r),
					"method", c.Request.Method,
					"path", c.Request.URL.Path,
				)
				handler.RespondError(c, http.StatusInternalServerError, handler.ErrInternal, "an unexpected error occurred")
				c.Abort()
			}
		}()
		c.Next()
	}
}

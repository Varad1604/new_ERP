package middleware

import (
	"time"

	"github.com/enterprise-erp/core/pkg/logger"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// RequestIDMiddleware injects a correlation ID into every request for distributed tracing
func RequestIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		reqID := c.GetHeader("X-Request-ID")
		if reqID == "" {
			reqID = uuid.New().String()
		}
		c.Set("RequestID", reqID)
		c.Header("X-Request-ID", reqID)
		c.Next()
	}
}

// RequestLogger logs structured details about every HTTP request
func RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		c.Next() // Process request

		end := time.Now()
		latency := end.Sub(start)

		reqID, _ := c.Get("RequestID")
		tenantID := c.GetString("tenant_id") // Extracted from auth/domain
		userID := c.GetString("user_id")

		logger.Get().Info("HTTP Request",
			zap.Int("status", c.Writer.Status()),
			zap.String("method", c.Request.Method),
			zap.String("path", path),
			zap.String("query", query),
			zap.String("ip", c.ClientIP()),
			zap.String("user-agent", c.Request.UserAgent()),
			zap.Duration("latency", latency),
			zap.Any("request_id", reqID),
			zap.String("tenant_id", tenantID),
			zap.String("user_id", userID),
			zap.String("errors", c.Errors.ByType(gin.ErrorTypePrivate).String()),
		)
	}
}

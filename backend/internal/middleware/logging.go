package middleware

import (
	"math/rand"
	"time"

	"ASMO-site-backend/pkg/logger"

	"github.com/gin-gonic/gin"
)

func generateRequestID() string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, 16)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

func LoggingMiddleware(logger *logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		requestID := generateRequestID()

		// Создаём logger с request ID
		requestLogger := logger.WithRequestID(requestID)

		// Устанавливаем request ID и logger в контекст
		c.Set("requestID", requestID)
		c.Set("logger", requestLogger)

		// Log запрос
		requestLogger.Info("Request started", map[string]interface{}{
			"method":      c.Request.Method,
			"path":        c.Request.URL.Path,
			"remote_addr": c.Request.RemoteAddr,
			"user_agent":  c.Request.UserAgent(),
		})

		// Процесс запроса
		c.Next()

		// Log ответ
		duration := time.Since(start)
		requestLogger.Info("Request completed", map[string]interface{}{
			"method":      c.Request.Method,
			"path":        c.Request.URL.Path,
			"status_code": c.Writer.Status(),
			"duration_ms": duration.Milliseconds(),
		})
	}
}
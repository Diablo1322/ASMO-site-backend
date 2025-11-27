package handlers

import (
	"database/sql"
	"net/http"
	"time"

	"ASMO-site-backend/internal/metrics"
	"ASMO-site-backend/internal/models"
	"ASMO-site-backend/pkg/logger"

	"github.com/gin-gonic/gin"
)

type HealthHandler struct {
	db     *sql.DB
	logger *logger.Logger
}

func NewHealthHandler(db *sql.DB) *HealthHandler {
    return &HealthHandler{
        db:     db,
        logger: logger.New("health", logger.INFO),
    }
}

func NewHealthHandlerWithLogger(db *sql.DB, logger *logger.Logger) *HealthHandler {
	return &HealthHandler{
		db:     db,
		logger: logger,
	}
}

func (h *HealthHandler) HealthCheck(c *gin.Context) {
	// Получаем logger из контекста (если есть) для request ID
	if ctxLogger, exists := c.Get("logger"); exists {
		if log, ok := ctxLogger.(*logger.Logger); ok {
			h.logger = log
		}
	}

	dbStatus := "connected"
	dbError := ""

	// Проверяем подключение к базе данных
	if h.db != nil {
		start := time.Now()
		if err := h.db.Ping(); err != nil {
			dbStatus = "disconnected"
			dbError = err.Error()

			// Логируем ошибку с дополнительными деталями
			h.logger.Error("Database connection failed", map[string]interface{}{
				"error":   err.Error(),
				"handler": "health",
				"time":    time.Now().Format(time.RFC3339),
			})
		}
		// Записываем метрику для health check
		metrics.RecordDatabaseQuery("ping", "health", time.Since(start))
	} else {
		dbStatus = "no_connection"
		dbError = "database instance is nil"
		h.logger.Error("Database instance is nil", map[string]interface{}{
			"handler": "health",
			"time":    time.Now().Format(time.RFC3339),
		})
	}

	// Формируем ответ
	response := models.HealthResponse{
		Status:  "ok",
		Message: "Service is healthy",
		Timestamp: map[string]interface{}{
			"server": "backend",
			"unix":   time.Now().Unix(),
			"iso":    time.Now().Format(time.RFC3339),
		},
		Database: dbStatus,
		Version:  "1.0.0",
	}

	// Если есть проблемы с БД, меняем статус ответа
	if dbStatus != "connected" {
		response.Status = "degraded"
		response.Message = "Service is running but database is unavailable"

		// Добавляем информацию об ошибке в ответ для дебага (только в development)
		env := c.GetString("ENVIRONMENT")
		if env == "development" || env == "" {
			response.Timestamp["db_error"] = dbError
		}

		h.logger.Warn("Health check: degraded mode", map[string]interface{}{
			"db_status": dbStatus,
			"db_error":  dbError,
			"environment": env,
		})
	} else {
		h.logger.Debug("Health check: all systems operational", map[string]interface{}{
			"db_status": dbStatus,
		})
	}

	c.JSON(http.StatusOK, response)
}
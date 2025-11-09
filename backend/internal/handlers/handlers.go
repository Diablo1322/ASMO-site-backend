package handlers

import (
	"database/sql"
	"net/http"
	"time"

	"ASMO-site-backend/internal/models"
	"ASMO-site-backend/internal/validation"
	"ASMO-site-backend/pkg/logger"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	db     *sql.DB
	logger *logger.Logger
}

func NewHandler(db *sql.DB, logger *logger.Logger) *Handler {
	return &Handler{db: db, logger: logger}
}

func (h *Handler) HealthCheck(c *gin.Context) {
	requestLogger := c.MustGet("logger").(*logger.Logger)
	requestLogger.Debug("Health check requested", nil)

	// Check database connection
	err := h.db.Ping()
	dbStatus := "connected"
	if err != nil {
		dbStatus = "disconnected"
		requestLogger.Error("Database health check failed", map[string]interface{}{
			"error": err.Error(),
		})
	}

	response := models.HealthResponse{
		Status:  "ok",
		Message: "Service is healthy",
		Timestamp: map[string]interface{}{
			"server": "backend",
			"unix":   time.Now().Unix(),
		},
		Database: dbStatus,
	}

	c.JSON(http.StatusOK, response)
}

func (h *Handler) CreateItem(c *gin.Context) {
	requestLogger := c.MustGet("logger").(*logger.Logger)
	
	var req models.CreateItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		requestLogger.Warn("Invalid request body", map[string]interface{}{
			"error": err.Error(),
		})
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
		})
		return
	}

	// Validate request
	if errs := validation.ValidateStruct(req); errs != nil {
		requestLogger.Warn("Validation failed", map[string]interface{}{
			"errors": errs,
		})
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Validation failed",
			"details": errs,
		})
		return
	}

	requestLogger.Info("Creating new item", map[string]interface{}{
		"name":  req.Name,
		"email": req.Email,
	})

	// Database operation would go here
	// For now, just return success
	c.JSON(http.StatusCreated, gin.H{
		"message": "Item created successfully",
		"item":    req,
	})
}

// GetItems - example handler with database query
func (h *Handler) GetItems(c *gin.Context) {
	requestLogger := c.MustGet("logger").(*logger.Logger)
	
	rows, err := h.db.Query("SELECT id, name, email, created_at, updated_at FROM items ORDER BY created_at DESC")
	if err != nil {
		requestLogger.Error("Database query failed", map[string]interface{}{
			"error": err.Error(),
		})
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch items",
		})
		return
	}
	defer rows.Close()

	var items []models.Item
	for rows.Next() {
		var item models.Item
		err := rows.Scan(&item.ID, &item.Name, &item.Email, &item.CreatedAt, &item.UpdatedAt)
		if err != nil {
			requestLogger.Error("Failed to scan row", map[string]interface{}{
				"error": err.Error(),
			})
			continue
		}
		items = append(items, item)
	}

	requestLogger.Debug("Items fetched successfully", map[string]interface{}{
		"count": len(items),
	})

	c.JSON(http.StatusOK, gin.H{
		"items": items,
		"count": len(items),
	})
}
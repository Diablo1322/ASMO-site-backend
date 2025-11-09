package integration

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"ASMO-site-backend/internal/handlers"
	"ASMO-site-backend/internal/models"
	"ASMO-site-backend/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupTestRouter() *gin.Engine {
	logger := logger.New("test", logger.INFO)
	handler := handlers.NewHandler(nil, logger)

	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("logger", logger)
		c.Next()
	})

	router.GET("/api/health", handler.HealthCheck)
	router.POST("/api/items", handler.CreateItem)

	return router
}

func TestHealthCheck(t *testing.T) {
	router := setupTestRouter()

	req := httptest.NewRequest("GET", "/api/health", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response models.HealthResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "ok", response.Status)
	assert.Equal(t, "Service is healthy", response.Message)
}

func TestCreateItem(t *testing.T) {
	router := setupTestRouter()

	t.Run("Valid item creation", func(t *testing.T) {
		item := models.CreateItemRequest{
			Name:  "Test Item",
			Email: "test@example.com",
		}

		body, _ := json.Marshal(item)
		req := httptest.NewRequest("POST", "/api/items", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
	})

	t.Run("Invalid item creation - validation error", func(t *testing.T) {
		item := models.CreateItemRequest{
			Name:  "", // Invalid empty name
			Email: "invalid-email",
		}

		body, _ := json.Marshal(item)
		req := httptest.NewRequest("POST", "/api/items", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}
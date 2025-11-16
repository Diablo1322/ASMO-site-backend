package integration

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"ASMO-site-backend/internal/database"
	"ASMO-site-backend/internal/handlers"
	"ASMO-site-backend/internal/models"
	"ASMO-site-backend/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupTestRouter() *gin.Engine {	
	// Use test database
	testDBURL := "postgres://user:password@localhost:5433/asmo_test_db"
	
	db, err := database.NewPostgresDB(testDBURL)
	if err != nil {
		panic("Failed to connect to test database")
	}

	logger := logger.New("test", logger.INFO)
	handler := handlers.NewHandler(db, logger)

	router := gin.Default()
	
	// Test routes
	api := router.Group("/api")
	{
		api.GET("/health", handler.HealthCheck)
		
		web := api.Group("/WebApplications")
		{
			web.GET("/:id", handler.GetWebProject)
			web.POST("/", handler.CreateWebProject)
		}
		
		mobile := api.Group("/MobileApplications")
		{
			mobile.GET("/:id", handler.GetMobileProject)
			mobile.POST("/", handler.CreateMobileProject)
		}
		
		bots := api.Group("/Bots")
		{
			bots.GET("/:id", handler.GetBotProject)
			bots.POST("/", handler.CreateBotProject)
		}
	}

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
}

func TestCreateWebProject(t *testing.T) {
	router := setupTestRouter()

	project := models.CreateWebProjectRequest{
		Name:        "Test Web Application Project",
		Description: "This is a comprehensive test description for a web application project that meets the minimum length requirements.",
		Img:         "https://example.com/image.jpg",
		Price:       1500.50,
		TimeDevelop: 30,
	}

	body, _ := json.Marshal(project)
	req := httptest.NewRequest("POST", "/api/WebApplications", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestCreateMobileProject(t *testing.T) {
	router := setupTestRouter()

	project := models.CreateMobileProjectRequest{
		Name:        "Test Mobile Application Project",
		Description: "This is a comprehensive test description for a mobile application project that meets the minimum length requirements.",
		Img:         "https://example.com/mobile.jpg",
		Price:       2000.75,
		TimeDevelop: 45,
	}

	body, _ := json.Marshal(project)
	req := httptest.NewRequest("POST", "/api/MobileApplications", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestCreateBotProject(t *testing.T) {
	router := setupTestRouter()

	project := models.CreateBotsProjectRequest{
		Name:        "Test Bot Development Project",
		Description: "This is a comprehensive test description for a bot development project that meets the minimum length requirements.",
		Img:         "https://example.com/bot.jpg",
		Price:       800.25,
		TimeDevelop: 15,
	}

	body, _ := json.Marshal(project)
	req := httptest.NewRequest("POST", "/api/Bots", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestGetNonExistentProject(t *testing.T) {
	router := setupTestRouter()

	req := httptest.NewRequest("GET", "/api/WebApplications/999", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}
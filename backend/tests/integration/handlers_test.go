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
	testutils "ASMO-site-backend/tests/testutils"
	"ASMO-site-backend/internal/cache"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupTestRouter() *gin.Engine {
	// Используем testutils для настройки базы данных
	db, err := testutils.SetupTestDB()
	if err != nil {
		panic("Failed to setup test database: " + err.Error())
	}

	logger.New("test", logger.INFO)

	// Создаем мок Redis для тестов
	redisMock := testutils.NewRedisMock()

	// Приводим к интерфейсу Cache (мок уже реализует интерфейс)
	var cacheInterface cache.Cache = redisMock

	// Инициализируем хэндлеры с Redis моком
	healthHandler := handlers.NewHealthHandlerWithLogger(db, logger.New("test", logger.INFO))
	webHandler := handlers.NewWebProjectsHandler(db, cacheInterface)        // ✅ ИНТЕРФЕЙС
	mobileHandler := handlers.NewMobileProjectsHandler(db, cacheInterface)  // ✅ ИНТЕРФЕЙС
	botHandler := handlers.NewBotProjectsHandler(db, cacheInterface)        // ✅ ИНТЕРФЕЙС
	staffHandler := handlers.NewStaffHandler(db, cacheInterface)

	router := gin.Default()

	// CORS middleware for tests
	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// Test routes
	api := router.Group("/api")
	{
		api.GET("/health", healthHandler.HealthCheck)

		web := api.Group("/WebApplications")
		{
			web.GET("/:id", webHandler.GetWebProject)
			web.GET("/", webHandler.GetWebProjects)
			web.POST("/", webHandler.CreateWebProject)
		}

		mobile := api.Group("/MobileApplications")
		{
			mobile.GET("/:id", mobileHandler.GetMobileProject)
			mobile.GET("/", mobileHandler.GetMobileProjects)
			mobile.POST("/", mobileHandler.CreateMobileProject)
		}

		bots := api.Group("/Bots")
		{
			bots.GET("/:id", botHandler.GetBotProject)
			bots.GET("/", botHandler.GetBotProjects)
			bots.POST("/", botHandler.CreateBotProject)
		}

		staff := api.Group("/Staff")
		{
			staff.GET("/:id", staffHandler.GetStaffMember)
			staff.GET("/", staffHandler.GetStaff)
			staff.POST("/", staffHandler.CreateStaff)
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
	assert.Equal(t, "1.0.0", response.Version) // ✅ Теперь поле есть в модели
	assert.Contains(t, response.Timestamp, "server")
	assert.Contains(t, response.Timestamp, "unix")
	assert.Contains(t, response.Timestamp, "iso")
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

	// Проверяем ответ
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "Web project created successfully", response["message"])
	assert.NotNil(t, response["id"])
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

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "Mobile project created successfully", response["message"])
	assert.NotNil(t, response["id"])
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

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "Bot project created successfully", response["message"])
	assert.NotNil(t, response["id"])
}

func TestCreateStaff(t *testing.T) {
	router := setupTestRouter()

	staff := models.CreateStaffRequest{
		Name:        "Test Staff Member Full Name",
		Description: "This is a comprehensive test description for a staff member that meets the minimum length requirements and provides detailed information about their role and responsibilities.",
		Img:         "https://example.com/staff.jpg",
		Role:        "Senior Developer",
	}

	body, _ := json.Marshal(staff)
	req := httptest.NewRequest("POST", "/api/Staff", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "Staff member created successfully", response["message"])
	assert.NotNil(t, response["id"])
}

func TestGetStaff(t *testing.T) {
	router := setupTestRouter()

	// First create a staff member
	staff := models.CreateStaffRequest{
		Name:        "Test Staff Member For Get",
		Description: "This is a test description for staff member retrieval testing purposes.",
		Img:         "https://example.com/staff-get.jpg",
		Role:        "Test Developer",
	}

	body, _ := json.Marshal(staff)
	req := httptest.NewRequest("POST", "/api/Staff", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)

	// Now get all staff - первый запрос должен быть из БД
	req = httptest.NewRequest("GET", "/api/Staff", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Contains(t, response, "staff")
	assert.Contains(t, response, "count")
	assert.Equal(t, false, response["cached"])

	staffList := response["staff"].([]interface{})
	assert.Greater(t, len(staffList), 0)

	// Второй запрос - должен быть из кэша
	req = httptest.NewRequest("GET", "/api/Staff", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, true, response["cached"])
}

func TestGetNonExistentProject(t *testing.T) {
	router := setupTestRouter()

	req := httptest.NewRequest("GET", "/api/WebApplications/999", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "Web project not found", response["error"])
}

func TestGetNonExistentStaff(t *testing.T) {
	router := setupTestRouter()

	req := httptest.NewRequest("GET", "/api/Staff/999", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "Staff member not found", response["error"])
}

func TestHealthCheckDatabaseStatus(t *testing.T) {
	router := setupTestRouter()

	req := httptest.NewRequest("GET", "/api/health", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response models.HealthResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	// Проверяем все поля HealthResponse
	assert.Equal(t, "ok", response.Status)
	assert.Equal(t, "Service is healthy", response.Message)
	assert.Equal(t, "1.0.0", response.Version)
	assert.Equal(t, "connected", response.Database) // ✅ Должен быть connected в тестах

	// Проверяем timestamp
	assert.Contains(t, response.Timestamp, "server")
	assert.Contains(t, response.Timestamp, "unix")
	assert.Contains(t, response.Timestamp, "iso")
}
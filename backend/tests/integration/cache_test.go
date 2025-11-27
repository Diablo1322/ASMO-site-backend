package integration

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"ASMO-site-backend/internal/models"

	"github.com/stretchr/testify/assert"
)

func TestWebProjectsCache(t *testing.T) {
	router := setupTestRouter()

	// Создаем тестовый проект
	project := models.CreateWebProjectRequest{
		Name:        "Cache Test Web Project",
		Description: "This is a test project for cache functionality with sufficient length to pass validation.",
		Img:         "https://example.com/cache-test.jpg",
		Price:       1000.00,
		TimeDevelop: 20,
	}

	body, _ := json.Marshal(project)
	req := httptest.NewRequest("POST", "/api/WebApplications", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)

	// Первый запрос - должен быть из БД (cached: false)
	req = httptest.NewRequest("GET", "/api/WebApplications", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, false, response["cached"])

	// Второй запрос - должен быть из кэша (cached: true)
	req = httptest.NewRequest("GET", "/api/WebApplications", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, true, response["cached"])
}

func TestStaffCache(t *testing.T) {
	router := setupTestRouter()

	// Создаем тестового сотрудника
	staff := models.CreateStaffRequest{
		Name:        "Cache Test Staff Member",
		Description: "This is a test staff member for cache functionality with sufficient length to pass validation requirements.",
		Img:         "https://example.com/staff-cache.jpg",
		Role:        "Cache Tester",
	}

	body, _ := json.Marshal(staff)
	req := httptest.NewRequest("POST", "/api/Staff", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)

	// Первый запрос - должен быть из БД (cached: false)
	req = httptest.NewRequest("GET", "/api/Staff", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, false, response["cached"])

	// Второй запрос - должен быть из кэша (cached: true)
	req = httptest.NewRequest("GET", "/api/Staff", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, true, response["cached"])
}

func TestCacheInvalidationOnCreate(t *testing.T) {
	router := setupTestRouter()

	// Первый запрос - получаем пустой список (кэшируется)
	req := httptest.NewRequest("GET", "/api/WebApplications", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// Создаем новый проект (должен инвалидировать кэш)
	project := models.CreateWebProjectRequest{
		Name:        "Cache Invalidation Test",
		Description: "Testing cache invalidation when creating new projects with sufficient description length.",
		Img:         "https://example.com/invalidation-test.jpg",
		Price:       1500.00,
		TimeDevelop: 25,
	}

	body, _ := json.Marshal(project)
	req = httptest.NewRequest("POST", "/api/WebApplications", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)

	// Запрос после создания - должен быть из БД (cached: false)
	req = httptest.NewRequest("GET", "/api/WebApplications", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, false, response["cached"])
	assert.Greater(t, len(response["projects"].([]interface{})), 0)
}
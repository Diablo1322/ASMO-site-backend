package handlers

import (
	"database/sql"
	"net/http"
	"strconv"
	"time"

	"ASMO-site-backend/internal/cache"
	"ASMO-site-backend/internal/metrics"
	"ASMO-site-backend/internal/models"
	"ASMO-site-backend/internal/validation"

	"github.com/gin-gonic/gin"
)

type WebProjectsHandler struct {
	db    *sql.DB
	cache cache.Cache
}

func NewWebProjectsHandler(db *sql.DB, cache cache.Cache) *WebProjectsHandler {
	return &WebProjectsHandler{
		db:    db,
		cache: cache,
	}
}

func (h *WebProjectsHandler) GetWebProjects(c *gin.Context) {
	start := time.Now()
	cacheKey := "web_projects:all"

	// Пробуем получить из кэша
	var projects []models.WebProjects
	if err := h.cache.Get(cacheKey, &projects); err == nil {
		metrics.RecordDatabaseQuery("cache_hit", "web_projects", time.Since(start))

		c.JSON(http.StatusOK, gin.H{
			"projects": projects,
			"count":    len(projects),
			"cached":   true,
		})
		return
	}

	// Если нет в кэше, получаем из БД
	rows, err := h.db.Query(`
		SELECT id, name, description, img, price, time_develop, created_at, update_at
		FROM web_projects
		ORDER BY created_at DESC
	`)

	metrics.RecordDatabaseQuery("select", "web_projects", time.Since(start))

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch web projects",
		})
		return
	}
	defer rows.Close()

	projects = []models.WebProjects{}
	for rows.Next() {
		var project models.WebProjects
		err := rows.Scan(
			&project.ID, &project.Name, &project.Description, &project.Img,
			&project.Price, &project.TimeDevelop, &project.CreatedAt, &project.UpdateAt,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to process web projects",
			})
			return
		}
		projects = append(projects, project)
	}

	if len(projects) == 0 {
		projects = []models.WebProjects{}
	}

	// Сохраняем в кэш на 5 минут
	h.cache.Set(cacheKey, projects, 5*time.Minute)

	c.JSON(http.StatusOK, gin.H{
		"projects": projects,
		"count":    len(projects),
		"cached":   false,
	})
}

func (h *WebProjectsHandler) GetWebProject(c *gin.Context) {
	start := time.Now()
	var req models.GetProjectRequest
	if err := c.ShouldBindUri(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid project ID",
		})
		return
	}

	if errs := validation.ValidateStruct(req); errs != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Validation failed",
			"details": errs,
		})
		return
	}

	cacheKey := "web_project:" + strconv.Itoa(req.ID)

	// Пробуем получить из кэша
	var project models.WebProjects
	if err := h.cache.Get(cacheKey, &project); err == nil {
		metrics.RecordDatabaseQuery("cache_hit", "web_projects", time.Since(start))
		c.JSON(http.StatusOK, project)
		return
	}

	err := h.db.QueryRow(`
		SELECT id, name, description, img, price, time_develop, created_at, update_at
		FROM web_projects WHERE id = $1
	`, req.ID).Scan(
		&project.ID, &project.Name, &project.Description, &project.Img,
		&project.Price, &project.TimeDevelop, &project.CreatedAt, &project.UpdateAt,
	)

	metrics.RecordDatabaseQuery("select", "web_projects", time.Since(start))

	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Web project not found",
		})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch web project",
		})
		return
	}

	// Сохраняем в кэш на 10 минут
	h.cache.Set(cacheKey, project, 10*time.Minute)

	c.JSON(http.StatusOK, project)
}

func (h *WebProjectsHandler) CreateWebProject(c *gin.Context) {
	start := time.Now()
	var req models.CreateWebProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
		})
		return
	}

	if errs := validation.ValidateStruct(req); errs != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Validation failed",
			"details": errs,
		})
		return
	}

	var id int
	err := h.db.QueryRow(`
		INSERT INTO web_projects (name, description, img, price, time_develop, created_at, update_at)
		VALUES ($1, $2, $3, $4, $5, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
		RETURNING id
	`, req.Name, req.Description, req.Img, req.Price, req.TimeDevelop).Scan(&id)

	metrics.RecordDatabaseQuery("insert", "web_projects", time.Since(start))

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create web project",
		})
		return
	}

	// Инвалидируем кэш при создании нового проекта
	h.cache.Delete("web_projects:all")

	c.JSON(http.StatusCreated, gin.H{
		"message": "Web project created successfully",
		"id":      id,
	})
}
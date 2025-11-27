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

type BotProjectsHandler struct {
	db    *sql.DB
	cache cache.Cache
}

func NewBotProjectsHandler(db *sql.DB, cache cache.Cache) *BotProjectsHandler {
	return &BotProjectsHandler{
		db:    db,
		cache: cache,
	}
}

func (h *BotProjectsHandler) GetBotProjects(c *gin.Context) {
	start := time.Now()
	cacheKey := "bot_projects:all"

	// Пробуем получить из кэша
	var projects []models.BotsProjects
	if err := h.cache.Get(cacheKey, &projects); err == nil {
		metrics.RecordDatabaseQuery("cache_hit", "bots_projects", time.Since(start))
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
		FROM bots_projects
		ORDER BY created_at DESC
	`)

	metrics.RecordDatabaseQuery("select", "bots_projects", time.Since(start))

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch bot projects",
		})
		return
	}
	defer rows.Close()

	projects = []models.BotsProjects{}
	for rows.Next() {
		var project models.BotsProjects
		err := rows.Scan(
			&project.ID, &project.Name, &project.Description, &project.Img,
			&project.Price, &project.TimeDevelop, &project.CreatedAt, &project.UpdateAt,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to process bot projects",
			})
			return
		}
		projects = append(projects, project)
	}

	if len(projects) == 0 {
		projects = []models.BotsProjects{}
	}

	// Сохраняем в кэш на 5 минут
	h.cache.Set(cacheKey, projects, 5*time.Minute)

	c.JSON(http.StatusOK, gin.H{
		"projects": projects,
		"count":    len(projects),
		"cached":   false,
	})
}

func (h *BotProjectsHandler) GetBotProject(c *gin.Context) {
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

	cacheKey := "bot_project:" + strconv.Itoa(req.ID)

	// Пробуем получить из кэша
	var project models.BotsProjects
	if err := h.cache.Get(cacheKey, &project); err == nil {
		metrics.RecordDatabaseQuery("cache_hit", "bots_projects", time.Since(start))
		c.JSON(http.StatusOK, project)
		return
	}

	err := h.db.QueryRow(`
		SELECT id, name, description, img, price, time_develop, created_at, update_at
		FROM bots_projects WHERE id = $1
	`, req.ID).Scan(
		&project.ID, &project.Name, &project.Description, &project.Img,
		&project.Price, &project.TimeDevelop, &project.CreatedAt, &project.UpdateAt,
	)

	metrics.RecordDatabaseQuery("select", "bots_projects", time.Since(start))

	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Bot project not found",
		})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch bot project",
		})
		return
	}

	// Сохраняем в кэш на 10 минут
	h.cache.Set(cacheKey, project, 10*time.Minute)

	c.JSON(http.StatusOK, project)
}

func (h *BotProjectsHandler) CreateBotProject(c *gin.Context) {
	start := time.Now()
	var req models.CreateBotsProjectRequest
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
		INSERT INTO bots_projects (name, description, img, price, time_develop, created_at, update_at)
		VALUES ($1, $2, $3, $4, $5, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
		RETURNING id
	`, req.Name, req.Description, req.Img, req.Price, req.TimeDevelop).Scan(&id)

	metrics.RecordDatabaseQuery("insert", "bots_projects", time.Since(start))

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create bot project",
		})
		return
	}

	// Инвалидируем кэш при создании нового проекта
	h.cache.Delete("bot_projects:all")

	c.JSON(http.StatusCreated, gin.H{
		"message": "Bot project created successfully",
		"id":      id,
	})
}
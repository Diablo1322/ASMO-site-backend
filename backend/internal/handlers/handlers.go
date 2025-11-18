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
	return &Handler{
		db:     db,
		logger: logger,
	}
}

func (h *Handler) HealthCheck(c *gin.Context) {
	// Проверка подключения базы данных
	err := h.db.Ping()
	dbStatus := "connected"
	if err != nil {
		dbStatus = "disconnected"
		h.logger.Error("Database health check failed", map[string]interface{}{
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

// Web Applications Handlers
func (h *Handler) GetWebProject(c *gin.Context) {
	var req models.GetProjectRequest
	if err := c.ShouldBindUri(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid project ID",
		})
		return
	}

	// Validate request
	if errs := validation.ValidateStruct(req); errs != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Validation failed",
			"details": errs,
		})
		return
	}

	var project models.WebProjects
	err := h.db.QueryRow(`
		SELECT id, name, description, img, price, time_develop, created_at, update_at
		FROM web_projects WHERE id = $1
	`, req.ID).Scan(
		&project.ID, &project.Name, &project.Description, &project.Img,
		&project.Price, &project.TimeDevelop, &project.CreatedAt, &project.UpdateAt,
	)

	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Web project not found",
		})
		return
	} else if err != nil {
		h.logger.Error("Failed to fetch web project", map[string]interface{}{
			"error": err.Error(),
			"id":    req.ID,
		})
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch web project",
		})
		return
	}

	c.JSON(http.StatusOK, project)
}

func (h *Handler) CreateWebProject(c *gin.Context) {
	var req models.CreateWebProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
		})
		return
	}

	// Validate request
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

	if err != nil {
		h.logger.Error("Failed to create web project", map[string]interface{}{
			"error": err.Error(),
			"data":  req,
		})
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create web project",
		})
		return
	}

	h.logger.Info("Web project created successfully", map[string]interface{}{
		"id":   id,
		"name": req.Name,
	})

	c.JSON(http.StatusCreated, gin.H{
		"message": "Web project created successfully",
		"id":      id,
	})
}

// Mobile Applications Handlers
func (h *Handler) GetMobileProject(c *gin.Context) {
	var req models.GetProjectRequest
	if err := c.ShouldBindUri(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid project ID",
		})
		return
	}

	// Validate request
	if errs := validation.ValidateStruct(req); errs != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Validation failed",
			"details": errs,
		})
		return
	}

	var project models.MobileProjects
	err := h.db.QueryRow(`
		SELECT id, name, description, img, price, time_develop, created_at, update_at
		FROM mobile_projects WHERE id = $1
	`, req.ID).Scan(
		&project.ID, &project.Name, &project.Description, &project.Img,
		&project.Price, &project.TimeDevelop, &project.CreatedAt, &project.UpdateAt,
	)

	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Mobile project not found",
		})
		return
	} else if err != nil {
		h.logger.Error("Failed to fetch mobile project", map[string]interface{}{
			"error": err.Error(),
			"id":    req.ID,
		})
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch mobile project",
		})
		return
	}

	c.JSON(http.StatusOK, project)
}

func (h *Handler) CreateMobileProject(c *gin.Context) {
	var req models.CreateMobileProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
		})
		return
	}

	// Validate request
	if errs := validation.ValidateStruct(req); errs != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Validation failed",
			"details": errs,
		})
		return
	}

	var id int
	err := h.db.QueryRow(`
		INSERT INTO mobile_projects (name, description, img, price, time_develop, created_at, update_at)
		VALUES ($1, $2, $3, $4, $5, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
		RETURNING id
	`, req.Name, req.Description, req.Img, req.Price, req.TimeDevelop).Scan(&id)

	if err != nil {
		h.logger.Error("Failed to create mobile project", map[string]interface{}{
			"error": err.Error(),
			"data":  req,
		})
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create mobile project",
		})
		return
	}

	h.logger.Info("Mobile project created successfully", map[string]interface{}{
		"id":   id,
		"name": req.Name,
	})

	c.JSON(http.StatusCreated, gin.H{
		"message": "Mobile project created successfully",
		"id":      id,
	})
}

// Bots Handlers
func (h *Handler) GetBotProject(c *gin.Context) {
	var req models.GetProjectRequest
	if err := c.ShouldBindUri(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid project ID",
		})
		return
	}

	// Validate request
	if errs := validation.ValidateStruct(req); errs != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Validation failed",
			"details": errs,
		})
		return
	}

	var project models.BotsProjects
	err := h.db.QueryRow(`
		SELECT id, name, description, img, price, time_develop, created_at, update_at
		FROM bots_projects WHERE id = $1
	`, req.ID).Scan(
		&project.ID, &project.Name, &project.Description, &project.Img,
		&project.Price, &project.TimeDevelop, &project.CreatedAt, &project.UpdateAt,
	)

	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Bot project not found",
		})
		return
	} else if err != nil {
		h.logger.Error("Failed to fetch bot project", map[string]interface{}{
			"error": err.Error(),
			"id":    req.ID,
		})
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch bot project",
		})
		return
	}

	c.JSON(http.StatusOK, project)
}

func (h *Handler) CreateBotProject(c *gin.Context) {
	var req models.CreateBotsProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
		})
		return
	}

	// Validate request
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

	if err != nil {
		h.logger.Error("Failed to create bot project", map[string]interface{}{
			"error": err.Error(),
			"data":  req,
		})
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create bot project",
		})
		return
	}

	h.logger.Info("Bot project created successfully", map[string]interface{}{
		"id":   id,
		"name": req.Name,
	})

	c.JSON(http.StatusCreated, gin.H{
		"message": "Bot project created successfully",
		"id":      id,
	})
}
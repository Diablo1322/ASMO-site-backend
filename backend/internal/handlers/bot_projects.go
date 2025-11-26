package handlers

import (
	"database/sql"
	"net/http"

	"ASMO-site-backend/internal/models"
	"ASMO-site-backend/internal/validation"

	"github.com/gin-gonic/gin"
)

type BotProjectsHandler struct {
	db *sql.DB
}

func NewBotProjectsHandler(db *sql.DB) *BotProjectsHandler {
	return &BotProjectsHandler{db: db}
}

func (h *BotProjectsHandler) GetBotProjects(c *gin.Context) {
	rows, err := h.db.Query(`
		SELECT id, name, description, img, price, time_develop, created_at, update_at
		FROM bots_projects
		ORDER BY created_at DESC
	`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch bot projects",
		})
		return
	}
	defer rows.Close()

	var projects []models.BotsProjects
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

	c.JSON(http.StatusOK, gin.H{
		"projects": projects,
		"count":    len(projects),
	})
}

func (h *BotProjectsHandler) GetBotProject(c *gin.Context) {
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
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch bot project",
		})
		return
	}

	c.JSON(http.StatusOK, project)
}

func (h *BotProjectsHandler) CreateBotProject(c *gin.Context) {
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

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create bot project",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Bot project created successfully",
		"id":      id,
	})
}
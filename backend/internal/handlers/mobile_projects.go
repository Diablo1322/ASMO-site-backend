package handlers

import (
	"database/sql"
	"net/http"

	"ASMO-site-backend/internal/models"
	"ASMO-site-backend/internal/validation"

	"github.com/gin-gonic/gin"
)

type MobileProjectsHandler struct {
	db *sql.DB
}

func NewMobileProjectsHandler(db *sql.DB) *MobileProjectsHandler {
	return &MobileProjectsHandler{db: db}
}

func (h *MobileProjectsHandler) GetMobileProjects(c *gin.Context) {
	rows, err := h.db.Query(`
		SELECT id, name, description, img, price, time_develop, created_at, update_at
		FROM mobile_projects
		ORDER BY created_at DESC
	`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch mobile projects",
		})
		return
	}
	defer rows.Close()

	var projects []models.MobileProjects
	for rows.Next() {
		var project models.MobileProjects
		err := rows.Scan(
			&project.ID, &project.Name, &project.Description, &project.Img,
			&project.Price, &project.TimeDevelop, &project.CreatedAt, &project.UpdateAt,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to process mobile projects",
			})
			return
		}
		projects = append(projects, project)
	}

	if len(projects) == 0 {
		projects = []models.MobileProjects{}
	}

	c.JSON(http.StatusOK, gin.H{
		"projects": projects,
		"count":    len(projects),
	})
}

func (h *MobileProjectsHandler) GetMobileProject(c *gin.Context) {
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
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch mobile project",
		})
		return
	}

	c.JSON(http.StatusOK, project)
}

func (h *MobileProjectsHandler) CreateMobileProject(c *gin.Context) {
	var req models.CreateMobileProjectRequest
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
		INSERT INTO mobile_projects (name, description, img, price, time_develop, created_at, update_at)
		VALUES ($1, $2, $3, $4, $5, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
		RETURNING id
	`, req.Name, req.Description, req.Img, req.Price, req.TimeDevelop).Scan(&id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create mobile project",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Mobile project created successfully",
		"id":      id,
	})
}
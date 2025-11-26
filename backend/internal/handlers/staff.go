package handlers

import (
	"database/sql"
	"net/http"

	"ASMO-site-backend/internal/models"
	"ASMO-site-backend/internal/validation"

	"github.com/gin-gonic/gin"
)

type StaffHandler struct {
	db *sql.DB
}

func NewStaffHandler(db *sql.DB) *StaffHandler {
	return &StaffHandler{db: db}
}

func (h *StaffHandler) GetStaff(c *gin.Context) {
	rows, err := h.db.Query(`
		SELECT id, name, description, img, role, created_at, update_at
		FROM staff
		ORDER BY created_at DESC
	`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch staff",
		})
		return
	}
	defer rows.Close()

	var staff []models.Staff
	for rows.Next() {
		var member models.Staff
		err := rows.Scan(
			&member.ID, &member.Name, &member.Description, &member.Img,
			&member.Role, &member.CreatedAt, &member.UpdateAt,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to process staff",
			})
			return
		}
		staff = append(staff, member)
	}

	if len(staff) == 0 {
		staff = []models.Staff{}
	}

	c.JSON(http.StatusOK, gin.H{
		"staff": staff,
		"count": len(staff),
	})
}

func (h *StaffHandler) GetStaffMember(c *gin.Context) {
	var req models.GetProjectRequest
	if err := c.ShouldBindUri(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid staff ID",
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

	var member models.Staff
	err := h.db.QueryRow(`
		SELECT id, name, description, img, role, created_at, update_at
		FROM staff WHERE id = $1
	`, req.ID).Scan(
		&member.ID, &member.Name, &member.Description, &member.Img,
		&member.Role, &member.CreatedAt, &member.UpdateAt,
	)

	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Staff member not found",
		})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch staff member",
		})
		return
	}

	c.JSON(http.StatusOK, member)
}

func (h *StaffHandler) CreateStaff(c *gin.Context) {
	var req models.CreateStaffRequest
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
		INSERT INTO staff (name, description, img, role, created_at, update_at)
		VALUES ($1, $2, $3, $4, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
		RETURNING id
	`, req.Name, req.Description, req.Img, req.Role).Scan(&id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create staff member",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Staff member created successfully",
		"id":      id,
	})
}
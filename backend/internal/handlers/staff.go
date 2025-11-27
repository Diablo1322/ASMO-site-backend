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

type StaffHandler struct {
	db    *sql.DB
	cache cache.Cache
}

func NewStaffHandler(db *sql.DB, cache cache.Cache) *StaffHandler {
	return &StaffHandler{
		db:    db,
		cache: cache,
	}
}

func (h *StaffHandler) GetStaff(c *gin.Context) {
	start := time.Now()
	cacheKey := "staff:all"

	// Пробуем получить из кэша
	var staff []models.Staff
	if err := h.cache.Get(cacheKey, &staff); err == nil {
		metrics.RecordDatabaseQuery("cache_hit", "staff", time.Since(start))
		c.JSON(http.StatusOK, gin.H{
			"staff":  staff,
			"count":  len(staff),
			"cached": true,
		})
		return
	}

	// Если нет в кэше, получаем из БД
	rows, err := h.db.Query(`
		SELECT id, name, description, img, role, created_at, update_at
		FROM staff
		ORDER BY created_at DESC
	`)

	metrics.RecordDatabaseQuery("select", "staff", time.Since(start))

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch staff",
		})
		return
	}
	defer rows.Close()

	staff = []models.Staff{}
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

	// Сохраняем в кэш на 5 минут
	h.cache.Set(cacheKey, staff, 5*time.Minute)

	c.JSON(http.StatusOK, gin.H{
		"staff":  staff,
		"count":  len(staff),
		"cached": false,
	})
}

func (h *StaffHandler) GetStaffMember(c *gin.Context) {
	start := time.Now()
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

	cacheKey := "staff:" + strconv.Itoa(req.ID)

	// Пробуем получить из кэша
	var member models.Staff
	if err := h.cache.Get(cacheKey, &member); err == nil {
		metrics.RecordDatabaseQuery("cache_hit", "staff", time.Since(start))
		c.JSON(http.StatusOK, member)
		return
	}

	err := h.db.QueryRow(`
		SELECT id, name, description, img, role, created_at, update_at
		FROM staff WHERE id = $1
	`, req.ID).Scan(
		&member.ID, &member.Name, &member.Description, &member.Img,
		&member.Role, &member.CreatedAt, &member.UpdateAt,
	)

	metrics.RecordDatabaseQuery("select", "staff", time.Since(start))

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

	// Сохраняем в кэш на 10 минут
	h.cache.Set(cacheKey, member, 10*time.Minute)

	c.JSON(http.StatusOK, member)
}

func (h *StaffHandler) CreateStaff(c *gin.Context) {
	start := time.Now()
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

	metrics.RecordDatabaseQuery("insert", "staff", time.Since(start))

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create staff member",
		})
		return
	}

	// Инвалидируем кэш при создании нового сотрудника
	h.cache.Delete("staff:all")

	c.JSON(http.StatusCreated, gin.H{
		"message": "Staff member created successfully",
		"id":      id,
	})
}
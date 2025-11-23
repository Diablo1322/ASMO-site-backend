package main

import (
	"log"
	"strings"

	"ASMO-site-backend/internal/config"
	"ASMO-site-backend/internal/database"
	"ASMO-site-backend/internal/handlers"
	"ASMO-site-backend/internal/validation"
	"ASMO-site-backend/pkg/logger"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Set Gin mode based on environment
	if cfg.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	// Initialize logger
	appLogger := logger.New("backend", logger.INFO)
	appLogger.Info("Application starting", map[string]interface{}{
		"port":        cfg.Port,
		"environment": cfg.Environment,
		"gin_mode":    gin.Mode(),
	})

	// Initialize database
	db, err := database.NewPostgresDB(cfg.DatabaseURL)
	if err != nil {
		appLogger.Error("Failed to connect to database", map[string]interface{}{
			"error": err.Error(),
		})
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	// Initialize validation
	validation.Init()

	// Initialize router
	router := gin.Default()

	// CORS configuration for Next.js frontend on port 3001
	router.Use(cors.New(cors.Config{
		AllowOriginFunc: func(origin string) bool {
			// Разрешаем все локальные адреса на порту 3001 (Next.js)
			return strings.Contains(origin, "://localhost:3001") ||
				strings.Contains(origin, "://127.0.0.1:3001") ||
				strings.Contains(origin, "://0.0.0.0:3001") ||
				// Для продакшна - добавьте ваши домены здесь
				strings.Contains(origin, "://your-production-domain.com")
		},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization", "Accept", "X-Requested-With"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * 60 * 60, // 12 hours
	}))

	// Initialize handlers
	handler := handlers.NewHandler(db, appLogger)

	// Fix favicon
	router.GET("/favicon.ico", func(c *gin.Context) {
		c.Status(204)
	})

	// Health check
	router.GET("/api/health", handler.HealthCheck)

	// Web Applications routes
	web := router.Group("/api/WebApplications")
	{
		web.GET("/:id", handler.GetWebProject)
		web.GET("/", handler.GetWebProjects)
		web.POST("/", handler.CreateWebProject)
	}

	// Mobile Applications routes
	mobile := router.Group("/api/MobileApplications")
	{
		mobile.GET("/:id", handler.GetMobileProject)
		mobile.GET("/", handler.GetMobileProjects)
		mobile.POST("/", handler.CreateMobileProject)
	}

	// Bots routes
	bots := router.Group("/api/Bots")
	{
		bots.GET("/:id", handler.GetBotProject)
		bots.GET("/", handler.GetBotProjects)
		bots.POST("/", handler.CreateBotProject)
	}

	staff := router.Group("/api/Staff")
	{
		staff.GET("/:id", handler.GetStaffMember)
		staff.GET("/", handler.GetStaff)
		staff.POST("/", handler.CreateStaff)
	}

	// Development-only debug endpoint
	if cfg.Environment == "development" {
		router.GET("/api/debug", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"mode":        "development",
				"frontend":    "Next.js on :3001",
				"cors":        "enabled for localhost:3001",
				"environment": cfg.Environment,
			})
		})
	}

	// Root endpoint
	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message":     "ASMO Backend API",
			"version":     "1.0.0",
			"environment": cfg.Environment,
			"frontend":    "Next.js on port 3001",
			"cors":        "configured for localhost:3001",
		})
	})

	// Start server
	appLogger.Info("Server starting", map[string]interface{}{
		"port":        cfg.Port,
		"environment": cfg.Environment,
		"frontend":    "Next.js :3001",
	})
	log.Printf("🚀 Server running in %s mode on http://localhost:%s", cfg.Environment, cfg.Port)
	log.Printf("📡 Frontend (Next.js) should connect from http://localhost:3001")
	log.Fatal(router.Run(":" + cfg.Port))
}
package main

import (
	"log"

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

	// Initialize logger
	appLogger := logger.New("backend", logger.INFO)
	appLogger.Info("Application starting", map[string]interface{}{
		"port":        cfg.Port,
		"environment": "development",
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

	// Настройка CORS
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3001", "http://127.0.0.1:3001", "http://localhost:3001/portgolio/", "http://127.0.0.1:3001/portfolio/"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization", "Accept", "X-Requested-With"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * 60 * 60, // 12 hours
	}))

	// ИЛИ более простая настройка CORS:
	// router.Use(cors.Default())

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

	// Root endpoint
	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "ASMO Backend API",
			"version": "1.0.0",
			"cors":    "enabled for ALL paths on localhost:3001 and 127.0.0.1:3001",
		})
	})

	// Start server
	appLogger.Info("Server starting", map[string]interface{}{
		"port": cfg.Port,
	})
	log.Printf("Server running on http://localhost:%s", cfg.Port)
	log.Fatal(router.Run(":" + cfg.Port))
}

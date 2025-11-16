package main

import (
	"log"

	"ASMO-site-backend/internal/config"
	"ASMO-site-backend/internal/database"
	"ASMO-site-backend/internal/handlers"
	"ASMO-site-backend/internal/middleware"
	"ASMO-site-backend/internal/validation"
	"ASMO-site-backend/pkg/logger"

	"github.com/gin-gonic/gin"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Initialize logger
	appLogger := logger.New("backend", logger.INFO)
	appLogger.Info("Application starting", map[string]interface{}{
		"port":        cfg.Port,
		"environment": "production", // Меняем на production
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

	// Trusted proxies for production
	router.SetTrustedProxies([]string{"nginx", "172.0.0.0/8"})

	// Middleware
	router.Use(middleware.LoggingMiddleware(appLogger))

	// CORS middleware - разрешаем HTTPS origins
	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// Security middleware
	router.Use(func(c *gin.Context) {
		// Set security headers
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("X-Frame-Options", "DENY")
		c.Header("X-XSS-Protection", "1; mode=block")

		// Если запрос пришел через HTTPS
		if c.Request.TLS != nil || c.GetHeader("X-Forwarded-Proto") == "https" {
			c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
		}
	})

	// Initialize handlers
	handler := handlers.NewHandler(db, appLogger)

	// Routes
	api := router.Group("/api")
	{
		api.GET("/health", handler.HealthCheck)

		// Web Applications routes
		web := api.Group("/WebApplications")
		{
			web.GET("/:id", handler.GetWebProject)
			web.POST("/", handler.CreateWebProject)
		}

		// Mobile Applications routes
		mobile := api.Group("/MobileApplications")
		{
			mobile.GET("/:id", handler.GetMobileProject)
			mobile.POST("/", handler.CreateMobileProject)
		}

		// Bots routes
		bots := api.Group("/Bots")
		{
			bots.GET("/:id", handler.GetBotProject)
			bots.POST("/", handler.CreateBotProject)
		}
	}

	// Root endpoint
	router.GET("/", func(c *gin.Context) {
		scheme := "http"
		if c.Request.TLS != nil || c.GetHeader("X-Forwarded-Proto") == "https" {
			scheme = "https"
		}

		baseURL := scheme + "://" + c.Request.Host

		c.JSON(200, gin.H{
			"message": "ASMO Backend API",
			"version": "1.0.0",
			"secure":  scheme == "https",
			"endpoints": []string{
				"GET " + baseURL + "/api/health",
				"GET/POST " + baseURL + "/api/WebApplications",
				"GET/POST " + baseURL + "/api/MobileApplications",
				"GET/POST " + baseURL + "/api/Bots",
			},
		})
	})

	// Start server
	appLogger.Info("Server starting", map[string]interface{}{
		"port": cfg.Port,
		"mode": gin.Mode(),
	})
	log.Fatal(router.Run(":" + cfg.Port))
}
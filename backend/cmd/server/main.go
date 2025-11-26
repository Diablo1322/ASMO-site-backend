package main

import (
	"log"
	"strings"

	"ASMO-site-backend/internal/config"
	"ASMO-site-backend/internal/database"
	"ASMO-site-backend/internal/handlers"
	"ASMO-site-backend/internal/middleware"
	"ASMO-site-backend/internal/validation"
	"ASMO-site-backend/pkg/logger"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Validate production requirements
	if cfg.Environment == "production" {
		if strings.Contains(cfg.DatabaseURL, "password@") && strings.Contains(cfg.DatabaseURL, "password") {
			log.Fatal("Default database password detected in production - use DB_PASSWORD environment variable")
		}
	}

	// Initialize logger
	appLogger := logger.New("backend", logger.INFO)
	appLogger.Info("Application starting", map[string]interface{}{
		"port":        cfg.Port,
		"environment": cfg.Environment,
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

	// Set Gin mode based on environment
	if cfg.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	// ✅ ИСПРАВЛЕНО: Используем обычное присваивание (=) вместо := для существующих переменных
	healthHandler := handlers.NewHealthHandlerWithLogger(db, appLogger)
	webHandler := handlers.NewWebProjectsHandler(db)
	mobileHandler := handlers.NewMobileProjectsHandler(db)
	botHandler := handlers.NewBotProjectsHandler(db)
	staffHandler := handlers.NewStaffHandler(db)

	// Initialize router
	router := gin.Default()

	// CORS configuration
	corsConfig := cors.Config{
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization", "Accept", "X-Requested-With"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * 60 * 60,
	}

	if cfg.Environment == "development" {
		// В разработке разрешаем локальные адреса
		corsConfig.AllowOriginFunc = func(origin string) bool {
			return strings.Contains(origin, "://localhost:3001") ||
				strings.Contains(origin, "://127.0.0.1:3001") ||
				strings.Contains(origin, "://0.0.0.0:3001")
		}
	} else {
		// В продакшене только разрешенные домены
		allowedOrigins := strings.Split(cfg.AllowedOrigins, ",")
		for i, origin := range allowedOrigins {
			allowedOrigins[i] = strings.TrimSpace(origin)
		}
		corsConfig.AllowOrigins = allowedOrigins
	}

	router.Use(cors.New(corsConfig))

	// Rate limiting middleware
	var limiter *rate.Limiter
	if cfg.Environment == "production" {
		limiter = rate.NewLimiter(100, 200)
	} else {
		limiter = rate.NewLimiter(1000, 2000)
	}

	router.Use(func(c *gin.Context) {
		if !limiter.Allow() {
			appLogger.Warn("Rate limit exceeded", map[string]interface{}{
				"ip": c.ClientIP(),
			})
			c.JSON(429, gin.H{
				"error":   "Too many requests",
				"message": "Please try again later",
			})
			c.Abort()
			return
		}
		c.Next()
	})

	// Security headers middleware
	router.Use(func(c *gin.Context) {
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("X-Frame-Options", "DENY")
		c.Header("X-XSS-Protection", "1; mode=block")

		if cfg.Environment == "production" {
			c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
		}
		c.Next()
	})

	// Logging middleware
	router.Use(middleware.LoggingMiddleware(appLogger))

	// Fix favicon
	router.GET("/favicon.ico", func(c *gin.Context) {
		c.Status(204)
	})

	// Health check
	router.GET("/api/health", healthHandler.HealthCheck)

	// Web Applications routes
	web := router.Group("/api/WebApplications")
	{
		web.GET("/:id", webHandler.GetWebProject)
		web.GET("/", webHandler.GetWebProjects)
		web.POST("/", webHandler.CreateWebProject)
	}

	// Mobile Applications routes
	mobile := router.Group("/api/MobileApplications")
	{
		mobile.GET("/:id", mobileHandler.GetMobileProject)
		mobile.GET("/", mobileHandler.GetMobileProjects)
		mobile.POST("/", mobileHandler.CreateMobileProject)
	}

	// Bots routes
	bots := router.Group("/api/Bots")
	{
		bots.GET("/:id", botHandler.GetBotProject)
		bots.GET("/", botHandler.GetBotProjects)
		bots.POST("/", botHandler.CreateBotProject)
	}

	// Staff routes
	staff := router.Group("/api/Staff")
	{
		staff.GET("/:id", staffHandler.GetStaffMember)
		staff.GET("/", staffHandler.GetStaff)
		staff.POST("/", staffHandler.CreateStaff)
	}

	// Root endpoint
	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message":     "ASMO Backend API",
			"version":     "1.0.0",
			"environment": cfg.Environment,
		})
	})

	// 404 handler
	router.NoRoute(func(c *gin.Context) {
		c.JSON(404, gin.H{
			"error":   "Endpoint not found",
			"message": "Check the API documentation",
		})
	})

	// Start server
	appLogger.Info("Server starting", map[string]interface{}{
		"port":        cfg.Port,
		"environment": cfg.Environment,
	})
	log.Printf("Server running in %s mode on http://localhost:%s", cfg.Environment, cfg.Port)
	log.Fatal(router.Run(":" + cfg.Port))
}
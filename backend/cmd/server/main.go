package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"

	"ASMO-site-backend/internal/config"
	"ASMO-site-backend/internal/database"
	"ASMO-site-backend/internal/handlers"
	"ASMO-site-backend/internal/middleware"
	"ASMO-site-backend/internal/validation"
	"ASMO-site-backend/pkg/logger"

	"github.com/gin-gonic/gin"
)

// ReverseProxy создает прокси к фронтенду
func ReverseProxy(target string) gin.HandlerFunc {
	return func(c *gin.Context) {
		remote, err := url.Parse(target)
		if err != nil {
			c.JSON(500, gin.H{"error": "Frontend server error"})
			return
		}

		proxy := httputil.NewSingleHostReverseProxy(remote)
		proxy.Director = func(req *http.Request) {
			req.Header = c.Request.Header.Clone()
			req.Host = remote.Host
			req.URL.Scheme = remote.Scheme
			req.URL.Host = remote.Host
			req.URL.Path = c.Request.URL.Path
		}

		proxy.ServeHTTP(c.Writer, c.Request)
	}
}

func main() {
	// Load configuration
	cfg := config.Load()

	// Initialize logger
	appLogger := logger.New("backend", logger.INFO)
	appLogger.Info("Application starting", map[string]interface{}{
		"port":        "3000",
		"environment": "production",
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

	// CORS middleware
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
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("X-Frame-Options", "DENY")
		c.Header("X-XSS-Protection", "1; mode=block")

		if c.Request.TLS != nil || c.GetHeader("X-Forwarded-Proto") == "https" {
			c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
		}
	})

	// Initialize handlers
	handler := handlers.NewHandler(db, appLogger)

	// API Routes
	api := router.Group("/api")
	{
		api.GET("/health", handler.HealthCheck)

		// Frontend API routes
		api.GET("/projects", handler.GetFrontendProjects)
		api.GET("/projects/:type", handler.GetFrontendProjectsByType)

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

	// Frontend proxy - все остальные запросы на фронтенд
	frontendURL := "http://localhost:3001" // фронтенд на 3001
	router.NoRoute(ReverseProxy(frontendURL))

	// Start server on port 3000
	appLogger.Info("Server starting on port 3000", map[string]interface{}{
		"port": "3000",
		"mode": gin.Mode(),
	})
	log.Fatal(router.Run(":3000"))
}

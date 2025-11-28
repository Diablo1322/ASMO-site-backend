package main

import (
	"log"
	"strings"
	"time"

	"ASMO-site-backend/internal/cache"
	"ASMO-site-backend/internal/config"
	"ASMO-site-backend/internal/database"
	"ASMO-site-backend/internal/handlers"
	"ASMO-site-backend/internal/middleware"
	"ASMO-site-backend/internal/validation"
	"ASMO-site-backend/pkg/logger"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"golang.org/x/time/rate"
)

// Prometheus middleware
func prometheusMiddleware() gin.HandlerFunc {
	// Создаем метрики внутри middleware чтобы избежать глобальных переменных
	httpRequestsTotal := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "path", "status"},
	)

	httpRequestDuration := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "HTTP request duration in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "path"},
	)

	httpResponseSize := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_response_size_bytes",
			Help:    "HTTP response size in bytes",
			Buckets: prometheus.ExponentialBuckets(100, 10, 8),
		},
		[]string{"method", "path"},
	)

	// Регистрируем метрики
	prometheus.MustRegister(httpRequestsTotal)
	prometheus.MustRegister(httpRequestDuration)
	prometheus.MustRegister(httpResponseSize)

	return func(c *gin.Context) {
		start := time.Now()
		path := c.FullPath()
		if path == "" {
			path = c.Request.URL.Path
		}

		// Process request
		c.Next()

		// Record metrics after request is processed
		duration := time.Since(start).Seconds()
		status := c.Writer.Status()

		httpRequestsTotal.WithLabelValues(
			c.Request.Method,
			path,
			string(rune(status)),
		).Inc()

		httpRequestDuration.WithLabelValues(
			c.Request.Method,
			path,
		).Observe(duration)

		httpResponseSize.WithLabelValues(
			c.Request.Method,
			path,
		).Observe(float64(c.Writer.Size()))
	}
}

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
		"redis_url":   cfg.RedisURL,
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

	// Initialize Redis cache
	redisCache, err := cache.NewRedisCache(cfg.RedisURL)
	if err != nil {
		appLogger.Error("Failed to connect to Redis", map[string]interface{}{
			"error": err.Error(),
			"url":   cfg.RedisURL,
		})
		log.Fatal("Failed to connect to Redis:", err)
	}
	defer redisCache.Close()

	appLogger.Info("Redis connected successfully", map[string]interface{}{
		"url": cfg.RedisURL,
	})

	// Initialize validation
	validation.Init()

	// Set Gin mode based on environment
	if cfg.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	// Initialize handlers with Redis cache
	healthHandler := handlers.NewHealthHandlerWithLogger(db, appLogger)
	webHandler := handlers.NewWebProjectsHandler(db, redisCache)
	mobileHandler := handlers.NewMobileProjectsHandler(db, redisCache)
	botHandler := handlers.NewBotProjectsHandler(db, redisCache)
	staffHandler := handlers.NewStaffHandler(db, redisCache)

	// Initialize router
	router := gin.Default()

	// Add Prometheus middleware if metrics are enabled
	if cfg.PrometheusMetrics {
		router.Use(prometheusMiddleware())

		// Expose metrics endpoint
		router.GET("/metrics", gin.WrapH(promhttp.Handler()))

		appLogger.Info("Prometheus metrics enabled", map[string]interface{}{
			"endpoint": "/metrics",
		})
	}

	// CORS configuration
	corsConfig := cors.Config{
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization", "Accept", "X-Requested-With"},
		ExposeHeaders:    []string{"Content-Length", "Content-Range"},
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
	staff := router.Group("/api/Members")
	{
		staff.GET("/:id", staffHandler.GetStaffMember)
		staff.GET("/", staffHandler.GetStaff)
		staff.POST("/", staffHandler.CreateStaff)
	}

	// Root endpoint
	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message":       "ASMO Backend API",
			"version":       "1.0.0",
			"environment":   cfg.Environment,
			"metrics":       cfg.PrometheusMetrics,
			"documentation": "/api/health",
		})
	})

	// 404 handler
	router.NoRoute(func(c *gin.Context) {
		c.JSON(404, gin.H{
			"error":   "Endpoint not found",
			"message": "Check the API documentation at /api/health",
			"metrics": cfg.PrometheusMetrics,
		})
	})

	// Start server
	appLogger.Info("Server starting", map[string]interface{}{
		"port":          cfg.Port,
		"environment":   cfg.Environment,
		"metrics":       cfg.PrometheusMetrics,
		"redis_enabled": true,
	})
	log.Printf("Server running in %s mode on http://localhost:%s", cfg.Environment, cfg.Port)
	log.Fatal(router.Run(":" + cfg.Port))
}
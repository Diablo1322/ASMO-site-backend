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
	// Подгружаем настройки
	cfg := config.Load()
	
	// Инициализация логгера
	appLogger := logger.New("backend", logger.INFO)
	appLogger.Info("Application starting", map[string]interface{}{
		"port":        cfg.Port,
		"environment": "development",
	})
	
	// Инициализация базы данных
	db, err := database.NewPostgresDB(cfg.DatabaseURL)
	if err != nil {
		appLogger.Error("Failed to connect to database", map[string]interface{}{
			"error": err.Error(),
		})
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	appLogger.Info("Running database migrations...", nil)
    err = database.RunMigrations(cfg.DatabaseURL)
	if err != nil {
		appLogger.Error("Failed to run database migrations", map[string]interface{}{
			"error": err.Error(),
		})
		log.Fatal("Failed to run database migrations:", err)
	}
	appLogger.Info("Database migrations completed successfully", nil)
	
	// Инициализация валидации
	validation.Init()
	
	// Инициализация роутера
	router := gin.Default()
	
	// Middleware
	router.Use(middleware.LoggingMiddleware(appLogger))
	
	// CORS middleware
	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", cfg.FrontendURL)
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")
		
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		
		c.Next()
	})
	
	// Инициализация handlers
	handler := handlers.NewHandler(db, appLogger)
	
	// Роутеры
	api := router.Group("/api")
	{
		api.GET("/health", handler.HealthCheck)
		
		// роутеры Web Applications
		web := api.Group("/WebApplications")
		{
			web.GET("/:id", handler.GetWebProject)
			web.POST("/", handler.CreateWebProject)
		}
		
		// роутеры Mobile Applications
		mobile := api.Group("/MobileApplications")
		{
			mobile.GET("/:id", handler.GetMobileProject)
			mobile.POST("/", handler.CreateMobileProject)
		}
		
		// роутеры Bots
		bots := api.Group("/Bots")
		{
			bots.GET("/:id", handler.GetBotProject)
			bots.POST("/", handler.CreateBotProject)
		}
	}
	
	// Запуск сервера
	appLogger.Info("Server starting", map[string]interface{}{
		"port": cfg.Port,
	})
	log.Fatal(router.Run(":" + cfg.Port))
}
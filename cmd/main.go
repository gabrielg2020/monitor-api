package main

import (
	"database/sql"
	"fmt"
	"os"
	"strings"

	"github.com/gabrielg2020/monitor-api/internal/api/handlers"
	"github.com/gabrielg2020/monitor-api/internal/repository"
	"github.com/gabrielg2020/monitor-api/internal/services"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "github.com/gabrielg2020/monitor-api/docs"
	_ "github.com/joho/godotenv/autoload"
	_ "github.com/mattn/go-sqlite3"
)

// @title           Monitoring API
// @version         1.0
// @description     REST API for collecting and serving system metrics from homelab clusters
// @termsOfService  http://swagger.io/terms/

// @contact.name   Gabriel G
// @contact.url    https://monitoring.gabrielg.tech
// @contact.email  gabriel.mg04@outlook.com

// @license.name  MIT
// @license.url   https://opensource.org/licenses/MIT

// @host      localhost:8191
// @BasePath  /api/v1

// @schemes   http https

func main() {
	engine := setupEngine()

	db, err := sql.Open("sqlite3", os.Getenv("DB_PATH"))
	if err != nil {
		panic(err)
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			panic(err)
		}
	}(db)

	if err := db.Ping(); err != nil {
		panic(err)
	}

	// Initialize services and handlers
	healthRepo := repository.NewHealthRepository(db)
	healthService := services.NewHealthService(healthRepo)
	healthHandler := handlers.NewHealthHandler(healthService)

	metricRepo := repository.NewMetricRepository(db)
	metricPostService := services.NewMetricService(metricRepo)
	metricPostHandler := handlers.NewMetricPostHandler(metricPostService)
	metricGetHandler := handlers.NewMetricGetHandler(metricPostService)
	metricLatestHandler := handlers.NewMetricLatestHandler(metricPostService)

	hostRepo := repository.NewHostRepository(db)
	hostService := services.NewHostService(hostRepo)
	hostPostHandler := handlers.NewHostPostHandler(hostService)
	hostGetHandler := handlers.NewHostGetHandler(hostService)

	// Set up endpoints
	engine.GET("/health", healthHandler.HandleHealth)
	engine.GET("/health/detailed", healthHandler.HandleDetailedHealth)

	v1 := engine.Group("/api/v1")
	{
		metrics := v1.Group("/metrics")
		{
			metrics.POST("", metricPostHandler.HandleMetricPost)
			metrics.GET("", metricGetHandler.HandleMetricGet)
			metrics.GET("/latest", metricLatestHandler.HandleMetricLatest)
		}
		hosts := v1.Group("/hosts")
		{
			hosts.POST("", hostPostHandler.HandleHostPost)
			hosts.GET("", hostGetHandler.HandleHostGet)
		}
	}

	engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Start the engine
	port := os.Getenv("PORT")
	fmt.Printf("Starting server on port %s\n", port)
	if err := engine.Run(":" + port); err != nil {
		panic(err)
	}
}

func setupEngine() *gin.Engine {
	engine := gin.New()

	// Get allowed origins from environment
	allowedOriginsStr := os.Getenv("ALLOWED_ORIGINS")
	if allowedOriginsStr == "" {
		allowedOriginsStr = "http://localhost" // default fallback
	}

	// Parse comma-separated origins into a map
	allowedOrigins := make(map[string]bool)
	for _, origin := range strings.Split(allowedOriginsStr, ",") {
		trimmed := strings.TrimSpace(origin)
		if trimmed != "" {
			allowedOrigins[trimmed] = true
		}
	}

	engine.Use(func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")
		if allowedOrigins[origin] {
			c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
		}
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		c.Next()
	})

	engine.Use(gin.Logger())
	engine.Use(gin.Recovery())

	return engine
}

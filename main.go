package main

import (
	"database/sql"
	"fmt"
	"os"
	"strings"

	"github.com/gabrielg2020/monitor-page/handlers"
	"github.com/gabrielg2020/monitor-page/services"
	"github.com/gin-gonic/gin"
	_ "github.com/joho/godotenv/autoload"
	_ "github.com/mattn/go-sqlite3"
)

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
	healthService := services.NewHealthService(db)
	healthHandler := handlers.NewHealthHandler(healthService)

	metricPostService := services.NewMetricPostService(db)
	metricPostHandler := handlers.NewMetricPostHandler(metricPostService)

	metricGetService := services.NewMetricGetService(db)
	metricGetHandler := handlers.NewMetricGetHandler(metricGetService)

	metricLatestService := services.NewMetricLatestService(db)
	metricLatestHandler := handlers.NewMetricLatestHandler(metricLatestService)

	hostPostService := services.NewHostPostService(db)
	hostPostHandler := handlers.NewHostPostHandler(hostPostService)

	hostGetService := services.NewHostGetService(db)
	hostGetHandler := handlers.NewHostGetHandler(hostGetService)

	// Set up endpoints
	engine.GET("/health", healthHandler.HandleHealth)

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

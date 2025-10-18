package main

import (
	"database/sql"
	"os"

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

	metricPushService := services.NewMetricPushService(db)
	metricPushHandler := handlers.NewMetricPushHandler(metricPushService)

	metricGetService := services.NewMetricGetService(db)
	metricGetHandler := handlers.NewMetricGetHandler(metricGetService)

	metricLatestService := services.NewMetricLatestService(db)
	metricLatestHandler := handlers.NewMetricLatestHandler(metricLatestService)

	hostPushService := services.NewHostPushService(db)
	hostPushHandler := handlers.NewHostPushHandler(hostPushService)

	hostGetService := services.NewHostGetService(db)
	hostGetHandler := handlers.NewHostGetHandler(hostGetService)

	// Set up endpoints
	engine.GET("/health", healthHandler.HandleHealth)

	v1 := engine.Group("/api/v1")
	{
		metrics := v1.Group("/metrics")
		{
			metrics.POST("", metricPushHandler.HandleMetricPush)
			metrics.GET("", metricGetHandler.HandleMetricGet)
			metrics.GET("/latest", metricLatestHandler.HandleMetricLatest)
		}
		hosts := v1.Group("/hosts")
		{
			hosts.POST("", hostPushHandler.HandleHostPush)
			hosts.GET("", hostGetHandler.HandleHostGet)
		}
	}

	// Start the engine
	port := os.Getenv("PORT")
	if err := engine.Run(":" + port); err != nil {
		panic(err)
	}
}

func setupEngine() *gin.Engine {
	engine := gin.Default()
	err := engine.SetTrustedProxies([]string{"127.0.0.1"})
	if err != nil {
		panic(err)
	}

	return engine
}

package api

import (
	"database/sql"
	"net/http"

	"github.com/gabrielg2020/monitor-api/internal/api/handlers"
	"github.com/gabrielg2020/monitor-api/internal/middleware"
	"github.com/gabrielg2020/monitor-api/internal/repository"
	"github.com/gabrielg2020/monitor-api/internal/services"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// SetupRouter initialises the router with all routes and middleware
func SetupRouter(
	healthHandler handlers.HealthHandlerInterface,
	hostHandler handlers.HostHandlerInterface,
	metricHandler handlers.MetricHandlerInterface,
	allowedOrigins []string,
) *gin.Engine {
	router := gin.New()

	// Middleware
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(middleware.CORS(allowedOrigins))

	// 405 responses for known routes with unsupported methods
	router.HandleMethodNotAllowed = true

	router.NoMethod(func(c *gin.Context) {
		c.JSON(http.StatusMethodNotAllowed, gin.H{
			"error":   "Method Not Allowed",
			"details": "The method is not allowed for the requested URL.",
		})
	})

	// Health endpoints
	router.GET("/health", healthHandler.GetHealth)
	router.GET("/health/detailed", healthHandler.GetDetailedHealth)

	// Swagger documentation
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		// Host routes
		hosts := v1.Group("/hosts")
		{
			hosts.POST("", hostHandler.Create)
			hosts.GET("", hostHandler.Get)
			hosts.PUT("", hostHandler.Update)
			hosts.DELETE("", hostHandler.Delete)
		}

		// Metric routes
		metrics := v1.Group("/metrics")
		{
			metrics.POST("", metricHandler.Create)
			metrics.GET("", metricHandler.Get)
			metrics.GET("/latest", metricHandler.GetLatest)
		}
	}

	return router
}

// SetupRouterWithDB is a convenience function for production use
// that creates handlers from a database connection
func SetupRouterWithDB(db *sql.DB, allowedOrigins []string) *gin.Engine {
	// Initialise repositories
	healthRepo := repository.NewHealthRepository(db)
	hostRepo := repository.NewHostRepository(db)
	metricRepo := repository.NewMetricRepository(db)

	// Initialise services
	healthService := services.NewHealthService(healthRepo)
	hostService := services.NewHostService(hostRepo)
	metricService := services.NewMetricService(metricRepo)

	// Initialise handlers
	healthHandler := handlers.NewHealthHandler(healthService)
	hostHandler := handlers.NewHostHandler(hostService)
	metricHandler := handlers.NewMetricHandler(metricService)

	return SetupRouter(healthHandler, hostHandler, metricHandler, allowedOrigins)
}

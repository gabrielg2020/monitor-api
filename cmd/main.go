package main

import (
	"fmt"
	"log"

	"github.com/gabrielg2020/monitor-api/internal/api"
	"github.com/gabrielg2020/monitor-api/internal/config"
	"github.com/gabrielg2020/monitor-api/pkg/database"
	"github.com/gin-gonic/gin"

	//nolint:typecheck // ignore missing docs package for Swagger UI
	_ "github.com/gabrielg2020/monitor-api/docs"
	_ "github.com/joho/godotenv/autoload"
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
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Set Gin mode
	gin.SetMode(cfg.Server.Mode)

	// Connect to database
	db, err := database.Connect(cfg.Database.Path)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer func() {
		if err := database.Close(db); err != nil {
			log.Printf("Error closing database: %v", err)
		}
	}()

	// Setup router with all routes
	router := api.SetupRouter(db, cfg.CORS.AllowedOrigins)

	// Start server
	addr := ":" + cfg.Server.Port
	fmt.Printf("Starting server on %s\n", addr)
	fmt.Printf("Swagger documentation: http://localhost%s/swagger/index.html\n", addr)

	if err := router.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

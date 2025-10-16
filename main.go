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
	pushService := services.NewPushService(db)
	pushHandler := handlers.NewPushHandler(pushService)

	// Set up endpoints
	engine.GET("/", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"message":   "Hello!",
			"developer": "Gabriel Guimaraes",
		})
	})

	apiGroup := engine.Group("/api")
	{
		apiGroup.GET("/", func(ctx *gin.Context) {
			ctx.JSON(200, gin.H{"message": "you shouldn't be here!"})
		})
		apiGroup.POST("/push", pushHandler.HandlePush)
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

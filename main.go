package main

import (
	"github.com/gin-gonic/gin"
)

func main() {
	engine := setupEngine()

	// Set up endpoints
	engine.GET("/", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"message":   "Hello!",
			"developer": "Gabriel Guimaraes",
		})
	})

	// Start the engine
	if err := engine.Run(":8080"); err != nil {
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

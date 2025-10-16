package handlers

import (
	"time"

	"github.com/gabrielg2020/monitor-page/services"
	"github.com/gin-gonic/gin"
)

type HealthHandler struct {
	service *services.HealthService
}

func NewHealthHandler(healthService *services.HealthService) *HealthHandler {
	return &HealthHandler{
		service: healthService,
	}
}

func (handler *HealthHandler) HandleHealth(ctx *gin.Context) {
	checks := make(map[string]string)
	status := "healthy"
	statusCode := 200

	if err := handler.service.CheckHealth(); err != nil {
		status = "unhealthy"
		checks["database"] = "unhealthy: " + err.Error()
		statusCode = 503
	} else {
		checks["database"] = "healthy"
	}

	ctx.JSON(statusCode, gin.H{
		"status":    status,
		"timestamp": time.Now().Format(time.RFC3339),
		"checks":    checks,
	})
}

package handlers

import (
	"time"

	"github.com/gabrielg2020/monitor-api/internal/models"
	"github.com/gabrielg2020/monitor-api/internal/services"
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

// HandleHealth godoc
// @Summary      Health check
// @Description  Check API and database health status
// @Tags         system
// @Accept       json
// @Produce      json
// @Success      200  {object}  models.HealthResponse
// @Failure      503  {object}  models.ErrorResponse
// @Router       /health [get]
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

	ctx.JSON(statusCode, models.HealthResponse{
		Status:    status,
		Timestamp: time.Now().Format(time.RFC3339),
		Checks:    checks,
	})
}

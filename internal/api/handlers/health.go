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

func NewHealthHandler(service *services.HealthService) *HealthHandler {
	return &HealthHandler{service: service}
}

// GetHealth godoc
// @Summary      Health check
// @Description  Check API and database health status
// @Tags         system
// @Accept       json
// @Produce      json
// @Success      200  {object}  models.HealthResponse
// @Failure      503  {object}  models.ErrorResponse
// @Router       /health [get]
func (handler *HealthHandler) GetHealth(ctx *gin.Context) {
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

// GetDetailedHealth godoc
// @Summary      Detailed health check
// @Description  Get detailed health information including database stats and table counts
// @Tags         system
// @Accept       json
// @Produce      json
// @Success      200  {object}  object{status=string,timestamp=string,database=object,database_stats=object,table_counts=object}
// @Failure      503  {object}  models.ErrorResponse
// @Router       /health/detailed [get]
func (handler *HealthHandler) GetDetailedHealth(ctx *gin.Context) {
	status := "healthy"
	statusCode := 200

	health, err := handler.service.GetDetailedHealth()
	if err != nil {
		status = "unhealthy"
		statusCode = 503
	}

	// Add status and timestamp
	response := gin.H{
		"status":    status,
		"timestamp": time.Now().Format(time.RFC3339),
	}

	// Merge detailed health info
	for key, value := range health {
		response[key] = value
	}

	ctx.JSON(statusCode, response)
}

package handlers

import (
	"github.com/gabrielg2020/monitor-api/internal/entities"
	"github.com/gabrielg2020/monitor-api/internal/models"
	"github.com/gabrielg2020/monitor-api/internal/services"
	"github.com/gin-gonic/gin"
)

type MetricLatestHandler struct {
	service *services.MetricLatestService
}

func NewMetricLatestHandler(latestService *services.MetricLatestService) *MetricLatestHandler {
	return &MetricLatestHandler{
		service: latestService,
	}
}

// HandleMetricLatest godoc
// @Summary      Get latest metrics
// @Description  Retrieve the most recent metrics for each host or a specific host
// @Tags         metrics
// @Accept       json
// @Produce      json
// @Param        hostname  query  string  false  "Filter by hostname"
// @Success      200  {object}  object{metric=models.SystemMetric}
// @Failure      400  {object}  models.ErrorResponse
// @Failure      500  {object}  models.ErrorResponse
// @Router       /metrics/latest [get]
func (handler *MetricLatestHandler) HandleMetricLatest(ctx *gin.Context) {
	var queryParams entities.MetricLatestQueryParams
	if err := ctx.ShouldBindQuery(&queryParams); err != nil {
		ctx.JSON(400, models.ErrorResponse{
			Error:   "Invalid query parameters",
			Details: err.Error(),
		})
		return
	}

	handler.service.SetQueryParams(&queryParams)

	metric, err := handler.service.GetLatestMetrics()
	if err != nil {
		ctx.JSON(500, models.ErrorResponse{
			Error:   "Failed to retrieve latest metric",
			Details: err.Error(),
		})
		return
	}

	ctx.JSON(200, gin.H{
		"metric": metric,
	})
}

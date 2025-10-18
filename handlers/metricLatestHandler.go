package handlers

import (
	"github.com/gabrielg2020/monitor-page/entities"
	"github.com/gabrielg2020/monitor-page/services"
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

func (handler *MetricLatestHandler) HandleMetricLatest(ctx *gin.Context) {
	var queryParams entities.MetricLatestQueryParams
	if err := ctx.ShouldBindQuery(&queryParams); err != nil {
		ctx.JSON(400, gin.H{
			"message": "Invalid query parameters",
			"error":   err.Error(),
		})
		return
	}

	handler.service.SetQueryParams(&queryParams)

	metric, err := handler.service.GetLatestMetrics()
	if err != nil {
		ctx.JSON(500, gin.H{
			"message": "Failed to retrieve latest metric",
			"error":   err.Error(),
		})
		return
	}

	ctx.JSON(200, gin.H{
		"metric": metric,
	})
}

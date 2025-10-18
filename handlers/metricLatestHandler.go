package handlers

import (
	"fmt"

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
	var hostID *int64
	if hostIDStr := ctx.Query("host_id"); hostIDStr != "" {
		var parsedHostID int64
		_, err := fmt.Sscanf(hostIDStr, "%d", &parsedHostID)
		if err != nil {
			ctx.JSON(400, gin.H{
				"message": "Invalid host_id parameter",
				"error":   err.Error(),
			})
			return
		}
		hostID = &parsedHostID
	}

	metric, err := handler.service.GetLatestMetrics(hostID)
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

package handlers

import (
	"github.com/gabrielg2020/monitor-page/entities"
	"github.com/gabrielg2020/monitor-page/services"
	"github.com/gin-gonic/gin"
)

type MetricPushHandler struct {
	service *services.MetricPushService
}

func NewMetricPushHandler(pushService *services.MetricPushService) *MetricPushHandler {
	return &MetricPushHandler{
		service: pushService,
	}
}

func (handler *MetricPushHandler) HandleMetricPush(ctx *gin.Context) {
	var requestBody struct {
		Record entities.SystemMetric `json:"record"`
	}

	if err := ctx.ShouldBindJSON(&requestBody); err != nil {
		ctx.JSON(400, gin.H{
			"message": "Invalid request body",
			"error":   err.Error(),
		})
		return
	} else if err := handler.service.PushMetric(&requestBody.Record); err != nil {
		ctx.JSON(500, gin.H{
			"message": "Failed to push record",
			"error":   err.Error(),
		})
		return
	}

	ctx.JSON(200, gin.H{"message": "Records pushed successfully"})
}

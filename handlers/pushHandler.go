package handlers

import (
	"github.com/gabrielg2020/monitor-page/entities"
	"github.com/gabrielg2020/monitor-page/services"
	"github.com/gin-gonic/gin"
)

type PushRecordHandler struct {
	service *services.PushService
}

func NewPushHandler(pushService *services.PushService) *PushRecordHandler {
	return &PushRecordHandler{
		service: pushService,
	}
}

func (handler *PushRecordHandler) HandlePush(ctx *gin.Context) {
	var requestBody struct {
		Record entities.SystemMetric `json:"record"`
	}

	if err := ctx.ShouldBindJSON(&requestBody); err != nil {
		ctx.JSON(400, gin.H{
			"message": "Invalid request body",
			"error":   err.Error(),
		})
		return
	} else if err := handler.service.PushRecord(&requestBody.Record); err != nil {
		ctx.JSON(500, gin.H{
			"message": "Failed to push record",
			"error":   err.Error(),
		})
		return
	}

	ctx.JSON(200, gin.H{"message": "Records pushed successfully"})
}

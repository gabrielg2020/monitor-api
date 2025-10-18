package handlers

import (
	"github.com/gabrielg2020/monitor-page/entities"
	"github.com/gabrielg2020/monitor-page/services"
	"github.com/gin-gonic/gin"
)

type HostPushHandler struct {
	service *services.HostPushService
}

func NewHostPushHandler(pushService *services.HostPushService) *HostPushHandler {
	return &HostPushHandler{
		service: pushService,
	}
}

func (handler *HostPushHandler) HandleHostPush(ctx *gin.Context) {
	var requestBody struct {
		Host entities.Host `json:"host"`
	}

	if err := ctx.ShouldBindJSON(&requestBody); err != nil {
		ctx.JSON(400, gin.H{
			"message": "Invalid request body",
			"error":   err.Error(),
		})
		return
	} else if err := handler.service.PushHost(&requestBody.Host); err != nil {
		ctx.JSON(500, gin.H{
			"message": "Failed to push host",
			"error":   err.Error(),
		})
		return
	}

	ctx.JSON(201, gin.H{"message": "Host pushed successfully"})
}

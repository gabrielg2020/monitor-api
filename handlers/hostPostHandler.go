package handlers

import (
	"github.com/gabrielg2020/monitor-page/entities"
	"github.com/gabrielg2020/monitor-page/services"
	"github.com/gin-gonic/gin"
)

type HostPostHandler struct {
	service *services.HostPostService
}

func NewHostPostHandler(postService *services.HostPostService) *HostPostHandler {
	return &HostPostHandler{
		service: postService,
	}
}

func (handler *HostPostHandler) HandleHostPost(ctx *gin.Context) {
	var requestBody struct {
		Host entities.Host `json:"host"`
	}

	if err := ctx.ShouldBindJSON(&requestBody); err != nil {
		ctx.JSON(400, gin.H{
			"message": "Invalid request body",
			"error":   err.Error(),
		})
		return
	}

	id, err := handler.service.PostHost(&requestBody.Host)
	if err != nil {
		ctx.JSON(500, gin.H{
			"message": "Failed to post host",
			"error":   err.Error(),
		})
		return
	}

	ctx.JSON(201, gin.H{
		"message": "Host posted successfully",
		"id":      id,
	})
}

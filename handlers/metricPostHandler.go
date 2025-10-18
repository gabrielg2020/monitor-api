package handlers

import (
	"github.com/gabrielg2020/monitor-page/entities"
	"github.com/gabrielg2020/monitor-page/services"
	"github.com/gin-gonic/gin"
)

type MetricPostHandler struct {
	service *services.MetricPostService
}

func NewMetricPostHandler(postService *services.MetricPostService) *MetricPostHandler {
	return &MetricPostHandler{
		service: postService,
	}
}

func (handler *MetricPostHandler) HandleMetricPost(ctx *gin.Context) {
	var requestBody struct {
		Record entities.SystemMetric `json:"record"`
	}

	if err := ctx.ShouldBindJSON(&requestBody); err != nil {
		ctx.JSON(400, gin.H{
			"message": "Invalid request body",
			"error":   err.Error(),
		})
		return
	}

	id, err := handler.service.PostMetric(&requestBody.Record)
	if err != nil {
		ctx.JSON(500, gin.H{
			"message": "Failed to post metric record",
			"error":   err.Error(),
		})
		return
	}

	ctx.JSON(201, gin.H{
		"message": "Host posted successfully",
		"id":      id,
	})
}

package handlers

import (
	"github.com/gabrielg2020/monitor-api/internal/entities"
	"github.com/gabrielg2020/monitor-api/internal/models"
	"github.com/gabrielg2020/monitor-api/internal/services"
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

// HandleMetricPost godoc
// @Summary      Submit system metrics
// @Description  Submit new system metrics from a monitoring agent
// @Tags         metrics
// @Accept       json
// @Produce      json
// @Param        request  body  models.CreateMetricRequest  true  "Metric data"
// @Success      201  {object}  object{message=string,id=int64}
// @Failure      400  {object}  models.ErrorResponse
// @Failure      500  {object}  models.ErrorResponse
// @Router       /metrics [post]
func (handler *MetricPostHandler) HandleMetricPost(ctx *gin.Context) {
	var requestBody struct {
		Record entities.SystemMetric `json:"record"`
	}

	if err := ctx.ShouldBindJSON(&requestBody); err != nil {
		ctx.JSON(400, models.ErrorResponse{
			Error:   "Invalid request body",
			Details: err.Error(),
		})
		return
	}

	id, err := handler.service.PostMetric(&requestBody.Record)
	if err != nil {
		ctx.JSON(500, models.ErrorResponse{
			Error:   "Failed to post metric record",
			Details: err.Error(),
		})
		return
	}

	ctx.JSON(201, gin.H{
		"message": "Metric posted successfully",
		"id":      id,
	})
}

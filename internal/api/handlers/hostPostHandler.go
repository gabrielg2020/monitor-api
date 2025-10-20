package handlers

import (
	"github.com/gabrielg2020/monitor-api/internal/entities"
	"github.com/gabrielg2020/monitor-api/internal/models"
	"github.com/gabrielg2020/monitor-api/internal/services"
	"github.com/gin-gonic/gin"
)

type HostPostHandler struct {
	service *services.HostService
}

func NewHostPostHandler(postService *services.HostService) *HostPostHandler {
	return &HostPostHandler{
		service: postService,
	}
}

// HandleHostPost godoc
// @Summary      Register a new host
// @Description  Register a new Raspberry Pi host in the monitoring system
// @Tags         hosts
// @Accept       json
// @Produce      json
// @Param        request  body  models.CreateHostRequest  true  "Host information"
// @Success      201  {object}  object{message=string,id=int64}
// @Failure      400  {object}  models.ErrorResponse
// @Failure      500  {object}  models.ErrorResponse
// @Router       /hosts [post]
func (handler *HostPostHandler) HandleHostPost(ctx *gin.Context) {
	var requestBody struct {
		Host entities.Host `json:"host"`
	}

	if err := ctx.ShouldBindJSON(&requestBody); err != nil {
		ctx.JSON(400, models.ErrorResponse{
			Error:   "Invalid request body",
			Details: err.Error(),
		})
		return
	}

	id, err := handler.service.CreateOrUpdateHost(&requestBody.Host)
	if err != nil {
		ctx.JSON(500, models.ErrorResponse{
			Error:   "Failed to register host",
			Details: err.Error(),
		})
		return
	}

	ctx.JSON(201, gin.H{
		"message": "Host registered successfully",
		"id":      id,
	})
}

package handlers

import (
	"github.com/gabrielg2020/monitor-page/entities"
	"github.com/gabrielg2020/monitor-page/services"
	"github.com/gin-gonic/gin"
)

type HostGetHandler struct {
	service *services.HostGetService
}

func NewHostGetHandler(getService *services.HostGetService) *HostGetHandler {
	return &HostGetHandler{
		service: getService,
	}
}

func (handler *HostGetHandler) HandleHostGet(ctx *gin.Context) {
	var queryParams entities.HostQueryParams
	if err := ctx.ShouldBindQuery(&queryParams); err != nil {
		ctx.JSON(400, gin.H{
			"message": "Invalid query parameters",
			"error":   err.Error(),
		})
		return
	}

	if queryParams.ID == 0 && queryParams.Hostname == "" && queryParams.IPAddress == "" {
		ctx.JSON(400, gin.H{
			"message": "At least one query parameter (id, hostname, ip_address) must be provided",
		})
		return
	}

	handler.service.SetQueryParams(&queryParams)

	host, err := handler.service.GetHost()
	if err != nil {
		ctx.JSON(500, gin.H{
			"message": "Failed to retrieve host",
			"error":   err.Error(),
		})
		return
	} else if host == nil {
		ctx.JSON(404, gin.H{
			"message": "Host not found",
		})
		return
	}

	ctx.JSON(200, gin.H{
		"host": host,
	})
}

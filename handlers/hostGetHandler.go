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

	handler.service.SetQueryParams(&queryParams)

	hosts, err := handler.service.GetHost()
	if err != nil {
		ctx.JSON(500, gin.H{
			"message": "Failed to retrieve host",
			"error":   err.Error(),
		})
		return
	}

	ctx.JSON(200, gin.H{
		"hosts": hosts,
		"meta": gin.H{
			"count": len(hosts),
		},
	})
}

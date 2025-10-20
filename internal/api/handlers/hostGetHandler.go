package handlers

import (
	"github.com/gabrielg2020/monitor-api/internal/entities"
	"github.com/gabrielg2020/monitor-api/internal/models"
	"github.com/gabrielg2020/monitor-api/internal/services"
	"github.com/gin-gonic/gin"
)

type HostGetHandler struct {
	service *services.HostService
}

func NewHostGetHandler(service *services.HostService) *HostGetHandler {
	return &HostGetHandler{
		service: service,
	}
}

// HandleHostGet godoc
// @Summary      List all hosts
// @Description  Get a list of all registered hosts in the monitoring system
// @Tags         hosts
// @Accept       json
// @Produce      json
// @Param        limit  query  int  false  "Limit number of results"  default(100)
// @Success      200  {object}  models.HostListResponse
// @Failure      400  {object}  models.ErrorResponse
// @Failure      500  {object}  models.ErrorResponse
// @Router       /hosts [get]
func (handler *HostGetHandler) HandleHostGet(ctx *gin.Context) {
	var queryParams entities.HostQueryParams
	if err := ctx.ShouldBindQuery(&queryParams); err != nil {
		ctx.JSON(500, models.ErrorResponse{
			Error:   "Failed to retrieve hosts",
			Details: err.Error(),
		})
		return
	}

	hosts, err := handler.service.GetHosts(&queryParams)
	if err != nil {
		ctx.JSON(500, models.ErrorResponse{
			Error:   "Failed to retrieve hosts",
			Details: err.Error(),
		})
		return
	}

	// Convert entities.Host to models.Host
	modelHosts := make([]models.Host, len(hosts))
	for i, host := range hosts {
		modelHosts[i] = models.Host{
			ID:        host.ID,
			Hostname:  host.Hostname,
			IPAddress: host.IPAddress,
			Role:      host.Role,
		}
	}

	ctx.JSON(200, models.HostListResponse{
		Hosts: modelHosts,
		Meta: models.Meta{
			Count: len(modelHosts),
		},
	})
}

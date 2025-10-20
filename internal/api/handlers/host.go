package handlers

import (
	"fmt"

	"github.com/gabrielg2020/monitor-api/internal/entities"
	"github.com/gabrielg2020/monitor-api/internal/models"
	"github.com/gabrielg2020/monitor-api/internal/services"
	"github.com/gin-gonic/gin"
)

type HostHandler struct {
	service *services.HostService
}

func NewHostHandler(service *services.HostService) *HostHandler {
	return &HostHandler{service: service}
}

// Create godoc
// @Summary      Register a new host
// @Description  Register a new Raspberry Pi host in the monitoring system or update if already exists
// @Tags         hosts
// @Accept       json
// @Produce      json
// @Param        request  body  models.CreateHostRequest  true  "Host information"
// @Success      201  {object}  object{message=string,id=int64}
// @Failure      400  {object}  models.ErrorResponse
// @Failure      500  {object}  models.ErrorResponse
// @Router       /hosts [post]
func (handler *HostHandler) Create(ctx *gin.Context) {
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

	id, err := handler.service.CreateHost(&requestBody.Host)
	if err != nil {
		ctx.JSON(500, models.ErrorResponse{
			Error:   "Failed to create host",
			Details: err.Error(),
		})
		return
	}

	ctx.JSON(201, gin.H{
		"message": "Host created successfully",
		"id":      id,
	})
}

// Get List godoc
// @Summary      List all hosts
// @Description  Get a list of all registered hosts in the monitoring system
// @Tags         hosts
// @Accept       json
// @Produce      json
// @Param        id          query  int     false  "Filter by host ID"
// @Param        hostname    query  string  false  "Filter by hostname"
// @Param        ip_address  query  string  false  "Filter by IP address"
// @Success      200  {object}  models.HostListResponse
// @Failure      400  {object}  models.ErrorResponse
// @Failure      500  {object}  models.ErrorResponse
// @Router       /hosts [get]
func (handler *HostHandler) Get(ctx *gin.Context) {
	var queryParams entities.HostQueryParams
	if err := ctx.ShouldBindQuery(&queryParams); err != nil {
		ctx.JSON(400, models.ErrorResponse{
			Error:   "Invalid query parameters",
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
		modelHosts[i] = toModelHost(host)
	}

	ctx.JSON(200, models.HostListResponse{
		Hosts: modelHosts,
		Meta: models.Meta{
			Count: len(modelHosts),
		},
	})
}

// Update godoc
// @Summary      Update a host
// @Description  Update an existing host's information
// @Tags         hosts
// @Accept       json
// @Produce      json
// @Param        id       path  int                        true  "Host ID"
// @Param        request  body  models.CreateHostRequest   true  "Host information"
// @Success      200  {object}  object{message=string}
// @Failure      400  {object}  models.ErrorResponse
// @Failure      404  {object}  models.ErrorResponse
// @Failure      500  {object}  models.ErrorResponse
// @Router       /hosts/{id} [put]
func (handler *HostHandler) Update(ctx *gin.Context) {
	id := ctx.Param("id")

	var hostID int64
	if _, err := fmt.Sscanf(id, "%d", &hostID); err != nil {
		ctx.JSON(400, models.ErrorResponse{
			Error:   "Invalid host ID",
			Details: err.Error(),
		})
		return
	}

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

	err := handler.service.UpdateHost(hostID, &requestBody.Host)
	if err != nil {
		ctx.JSON(500, models.ErrorResponse{
			Error:   "Failed to update host",
			Details: err.Error(),
		})
		return
	}

	ctx.JSON(200, gin.H{
		"message": "Host updated successfully",
	})
}

// Delete godoc
// @Summary      Delete a host
// @Description  Delete a host and all its associated metrics
// @Tags         hosts
// @Accept       json
// @Produce      json
// @Param        id   path  int  true  "Host ID"
// @Success      200  {object}  object{message=string}
// @Failure      400  {object}  models.ErrorResponse
// @Failure      404  {object}  models.ErrorResponse
// @Failure      500  {object}  models.ErrorResponse
// @Router       /hosts/{id} [delete]
func (handler *HostHandler) Delete(ctx *gin.Context) {
	id := ctx.Param("id")

	var hostID int64
	if _, err := fmt.Sscanf(id, "%d", &hostID); err != nil {
		ctx.JSON(400, models.ErrorResponse{
			Error:   "Invalid host ID",
			Details: err.Error(),
		})
		return
	}

	err := handler.service.DeleteHost(hostID)
	if err != nil {
		ctx.JSON(500, models.ErrorResponse{
			Error:   "Failed to delete host",
			Details: err.Error(),
		})
		return
	}

	ctx.JSON(200, gin.H{
		"message": "Host deleted successfully",
	})
}

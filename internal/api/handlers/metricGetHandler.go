package handlers

import (
	"strings"
	"time"

	"github.com/gabrielg2020/monitor-api/internal/entities"
	"github.com/gabrielg2020/monitor-api/internal/models"
	"github.com/gabrielg2020/monitor-api/internal/services"
	"github.com/gin-gonic/gin"
)

type MetricGetHandler struct {
	service *services.MetricService
}

func NewMetricGetHandler(getService *services.MetricService) *MetricGetHandler {
	return &MetricGetHandler{
		service: getService,
	}
}

// HandleMetricGet godoc
// @Summary      Get system metrics
// @Description  Retrieve system metrics with optional filtering and time range
// @Tags         metrics
// @Accept       json
// @Produce      json
// @Param        host_id     query  int     false  "Filter by host ID"
// @Param        hostname    query  string  false  "Filter by hostname"
// @Param        limit       query  int     false  "Limit results (max 1000)"  default(100)
// @Param        order       query  string  false  "Sort order (ASC or DESC)"  default(DESC)
// @Param        start_time  query  int     false  "Start timestamp (Unix)"
// @Param        end_time    query  int     false  "End timestamp (Unix)"
// @Success      200  {object}  models.MetricListResponse
// @Failure      400  {object}  models.ErrorResponse
// @Failure      500  {object}  models.ErrorResponse
// @Router       /metrics [get]
func (handler *MetricGetHandler) HandleMetricGet(ctx *gin.Context) {
	var queryParams entities.MetricQueryParams
	if err := ctx.ShouldBindQuery(&queryParams); err != nil {
		ctx.JSON(400, models.ErrorResponse{
			Error:   "Invalid query parameters",
			Details: err.Error(),
		})
		return
	}

	// Set defaults and validate
	if queryParams.Limit == 0 {
		queryParams.Limit = 100
	}
	if queryParams.Limit > 1000 {
		queryParams.Limit = 1000
	}
	if queryParams.Order == "" {
		queryParams.Order = "DESC"
	} else {
		queryParams.Order = strings.ToUpper(queryParams.Order)
		if queryParams.Order != "ASC" && queryParams.Order != "DESC" {
			ctx.JSON(400, models.ErrorResponse{
				Error:   "Invalid order parameter",
				Details: "Must be 'ASC' or 'DESC'",
			})
			return
		}
	}

	now := time.Now().Unix()
	if queryParams.EndTime == nil {
		queryParams.EndTime = &now
	}
	if queryParams.StartTime == nil {
		thirtyDaysAgo := *queryParams.EndTime - (86400 * 30) // Default to last 30 days
		queryParams.StartTime = &thirtyDaysAgo
	}

	records, err := handler.service.GetMetrics(&queryParams)
	if err != nil {
		ctx.JSON(500, models.ErrorResponse{
			Error:   "Failed to retrieve metrics",
			Details: err.Error(),
		})
		return
	}

	// Convert entities to models
	modelMetrics := make([]models.SystemMetric, len(records))
	for i, metric := range records {
		modelMetrics[i] = models.SystemMetric{
			ID:                   metric.ID,
			HostID:               metric.HostID,
			Timestamp:            metric.Timestamp,
			CPUUsage:             metric.CPUUsage,
			MemoryUsagePercent:   metric.MemoryUsagePercent,
			MemoryTotalBytes:     metric.MemoryTotalBytes,
			MemoryUsedBytes:      metric.MemoryUsedBytes,
			MemoryAvailableBytes: metric.MemoryAvailableBytes,
			DiskUsagePercent:     metric.DiskUsagePercent,
			DiskTotalBytes:       metric.DiskTotalBytes,
			DiskUsedBytes:        metric.DiskUsedBytes,
			DiskAvailableBytes:   metric.DiskAvailableBytes,
		}
	}

	ctx.JSON(200, models.MetricListResponse{
		Records: modelMetrics,
		Meta: models.Meta{
			Count: len(modelMetrics),
			Limit: queryParams.Limit,
		},
	})
}

package handlers

import (
	"time"

	"github.com/gabrielg2020/monitor-api/internal/entities"
	"github.com/gabrielg2020/monitor-api/internal/models"
	"github.com/gabrielg2020/monitor-api/internal/services"
	"github.com/gin-gonic/gin"
)

type MetricHandler struct {
	service *services.MetricService
}

func NewMetricHandler(service *services.MetricService) *MetricHandler {
	return &MetricHandler{service: service}
}

// Create godoc
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
func (handler *MetricHandler) Create(ctx *gin.Context) {
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

	id, err := handler.service.CreateMetric(&requestBody.Record)
	if err != nil {
		ctx.JSON(500, models.ErrorResponse{
			Error:   "Failed to create metric record",
			Details: err.Error(),
		})
		return
	}

	ctx.JSON(201, gin.H{
		"message": "Metric create successfully",
		"id":      id,
	})
}

// Get godoc
// @Summary      Get system metrics
// @Description  Retrieve system metrics with optional filtering and time range
// @Tags         metrics
// @Accept       json
// @Produce      json
// @Param        host_id     query  int     false  "Filter by host ID"
// @Param        limit       query  int     false  "Limit results (max 1000)"  default(100)
// @Param        order       query  string  false  "Sort order (ASC or DESC)"  default(DESC)
// @Param        start_time  query  int     false  "Start timestamp (Unix)"
// @Param        end_time    query  int     false  "End timestamp (Unix)"
// @Success      200  {object}  models.MetricListResponse
// @Failure      400  {object}  models.ErrorResponse
// @Failure      500  {object}  models.ErrorResponse
// @Router       /metrics [get]
func (handler *MetricHandler) Get(ctx *gin.Context) {
	var queryParams entities.MetricQueryParams
	if err := ctx.ShouldBindQuery(&queryParams); err != nil {
		ctx.JSON(400, models.ErrorResponse{
			Error:   "Invalid query parameters",
			Details: err.Error(),
		})
		return
	}

	// Validate and set defaults
	if errResp := setMetricQueryDefaults(&queryParams); errResp != nil {
		ctx.JSON(400, errResp)
		return
	}

	now := time.Now().Unix()
	if queryParams.EndTime == nil {
		queryParams.EndTime = &now
	}
	if queryParams.StartTime == nil {
		thirtyDaysAgo := *queryParams.EndTime - (86400 * 30)
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
		modelMetrics[i] = toModelMetric(metric)
	}

	ctx.JSON(200, models.MetricListResponse{
		Records: modelMetrics,
		Meta: models.Meta{
			Count: len(modelMetrics),
			Limit: queryParams.Limit,
		},
	})
}

// GetLatest godoc
// @Summary      Get latest metrics
// @Description  Retrieve the most recent metrics for each host or a specific host
// @Tags         metrics
// @Accept       json
// @Produce      json
// @Param        host_id  query  int  false  "Filter by host ID"
// @Success      200  {object}  object{metric=models.SystemMetric}
// @Failure      400  {object}  models.ErrorResponse
// @Failure      404  {object}  models.ErrorResponse
// @Failure      500  {object}  models.ErrorResponse
// @Router       /metrics/latest [get]
func (handler *MetricHandler) GetLatest(ctx *gin.Context) {
	var queryParams entities.MetricLatestQueryParams
	if err := ctx.ShouldBindQuery(&queryParams); err != nil {
		ctx.JSON(400, models.ErrorResponse{
			Error:   "Invalid query parameters",
			Details: err.Error(),
		})
		return
	}

	metric, err := handler.service.GetLatestMetric(queryParams.HostID)
	if err != nil {
		ctx.JSON(500, models.ErrorResponse{
			Error:   "Failed to retrieve latest metric",
			Details: err.Error(),
		})
		return
	}

	if metric == nil {
		ctx.JSON(404, models.ErrorResponse{
			Error:   "Metric not found",
			Details: "No latest metric found for the specified host",
		})
		return
	}

	ctx.JSON(200, gin.H{
		"metric": metric,
	})
}

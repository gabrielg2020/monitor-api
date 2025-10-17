package handlers

import (
	"fmt"
	"strings"
	"time"

	"github.com/gabrielg2020/monitor-page/entities"
	"github.com/gabrielg2020/monitor-page/services"
	"github.com/gin-gonic/gin"
)

type MetricGetHandler struct {
	service *services.MetricGetService
}

func NewMetricGetHandler(getService *services.MetricGetService) *MetricGetHandler {
	return &MetricGetHandler{
		service: getService,
	}
}

func (handler *MetricGetHandler) HandleMetricGet(ctx *gin.Context) {
	var queryParams entities.MetricQueryParams
	if err := ctx.ShouldBindQuery(&queryParams); err != nil {
		ctx.JSON(400, gin.H{
			"message": "Invalid query parameters",
			"error":   err.Error(),
		})
		return
	}

	fmt.Printf("Host ID =  %v\n", queryParams.HostID)

	// Set Defaults
	if queryParams.Limit == 0 {
		queryParams.Limit = 100
	}
	if queryParams.Limit > 1000 {
		queryParams.Limit = 1000
	}
	if queryParams.Order == "" {
		queryParams.Order = "DESC"
	} else {
		// Check if order is valid
		queryParams.Order = strings.ToUpper(queryParams.Order)
		if queryParams.Order != "ASC" && queryParams.Order != "DESC" {
			ctx.JSON(400, gin.H{
				"message": "Invalid order parameter, must be 'ASC' or 'DESC'",
				"error":   "invalid order parameter",
			})
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

	handler.service.SetQueryParams(&queryParams)

	records, err := handler.service.GetMetrics()
	if err != nil {
		ctx.JSON(500, gin.H{
			"message": "Failed to retrieve metrics",
			"error":   err.Error(),
		})
		return
	} else {
		ctx.JSON(200, gin.H{
			"records": records,
			"meta": gin.H{
				"count": len(records),
				"limit": queryParams.Limit,
			},
		})
	}
}

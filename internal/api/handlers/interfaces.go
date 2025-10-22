package handlers

import "github.com/gin-gonic/gin"

// HealthHandlerInterface defines methods for health check handlers
type HealthHandlerInterface interface {
	GetHealth(ctx *gin.Context)
	GetDetailedHealth(ctx *gin.Context)
}

// HostHandlerInterface defines methods for host handlers
type HostHandlerInterface interface {
	Create(ctx *gin.Context)
	Get(ctx *gin.Context)
	Update(ctx *gin.Context)
	Delete(ctx *gin.Context)
}

// MetricHandlerInterface defines methods for metric handlers
type MetricHandlerInterface interface {
	Create(ctx *gin.Context)
	Get(ctx *gin.Context)
	GetLatest(ctx *gin.Context)
}

var _ HealthHandlerInterface = &HealthHandler{}
var _ HostHandlerInterface = &HostHandler{}
var _ MetricHandlerInterface = &MetricHandler{}

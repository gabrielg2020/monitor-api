package handlers

import (
	"strings"

	"github.com/gabrielg2020/monitor-api/internal/entities"
	"github.com/gabrielg2020/monitor-api/internal/models"
)

// toModelHost converts entity to model
func toModelHost(host entities.Host) models.Host {
	return models.Host{
		ID:        host.ID,
		Hostname:  host.Hostname,
		IPAddress: host.IPAddress,
		Role:      host.Role,
	}
}

// toModelMetric converts entity to model
func toModelMetric(metric entities.SystemMetric) models.SystemMetric {
	return models.SystemMetric{
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

// setMetricQueryDefaults validates and sets defaults for metric query params
func setMetricQueryDefaults(params *entities.MetricQueryParams) *models.ErrorResponse {
	// Set defaults
	if params.Limit == 0 {
		params.Limit = 100
	}
	if params.Limit > 1000 {
		params.Limit = 1000
	}
	if params.Order == "" {
		params.Order = "DESC"
	} else {
		params.Order = strings.ToUpper(params.Order)
		if params.Order != "ASC" && params.Order != "DESC" {
			return &models.ErrorResponse{
				Error:   "Invalid order parameter",
				Details: "Must be 'ASC' or 'DESC'",
			}
		}
	}

	return nil
}

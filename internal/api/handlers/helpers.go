package handlers

import (
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

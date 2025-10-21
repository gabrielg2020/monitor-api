package services

import "github.com/gabrielg2020/monitor-api/internal/entities"

// ValidateSystemMetric validates metric data
func ValidateSystemMetric(params *entities.SystemMetric) error {
	// ID
	if params.ID < 0 {
		return ErrInvalidHostID
	}

	// HostID
	if params.HostID <= 0 {
		return ErrInvalidHostID
	}

	// CPU Usage
	if params.CPUUsage < 0 || params.CPUUsage > 100 {
		return ErrInvalidCPUUsage
	}

	// Memory Usage
	if params.MemoryUsagePercent < 0 || params.MemoryUsagePercent > 100 {
		return ErrInvalidMemoryUsage
	}

	// Disk Usage
	if params.DiskUsagePercent < 0 || params.DiskUsagePercent > 100 {
		return ErrInvalidDiskUsage
	}

	return nil
}

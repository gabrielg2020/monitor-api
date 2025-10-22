package services

import "errors"

// Common service errors
var (
	// Host service errors
	ErrHostNotFound    = errors.New("host not found")
	ErrInvalidHostData = errors.New("invalid host data")
	ErrDuplicateHost   = errors.New("host already exists")

	// Metric service errors
	ErrInvalidHostID      = errors.New("invalid host ID")
	ErrInvalidCPUUsage    = errors.New("CPU usage must be between 0 and 100")
	ErrInvalidMemoryUsage = errors.New("memory usage must be between 0 and 100")
	ErrInvalidDiskUsage   = errors.New("disk usage must be between 0 and 100")
	ErrNilQueryParams     = errors.New("query parameters cannot be nil")
	ErrMetricNotFound     = errors.New("metric not found")
	ErrInvalidTimeRange   = errors.New("invalid time range")
)

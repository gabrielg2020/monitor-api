package entities

type SystemMetric struct {
	ID                   int64   `json:"id" db:"id"`
	HostID               int64   `json:"host_id" db:"host_id"`
	Timestamp            int64   `json:"timestamp" db:"timestamp"`
	CPUUsage             float64 `json:"cpu_usage" db:"cpu_usage"`
	MemoryUsagePercent   float64 `json:"memory_usage_percent" db:"memory_usage_percent"`
	MemoryTotalBytes     int64   `json:"memory_total_bytes" db:"memory_total_bytes"`
	MemoryUsedBytes      int64   `json:"memory_used_bytes" db:"memory_used_bytes"`
	MemoryAvailableBytes int64   `json:"memory_available_bytes" db:"memory_available_bytes"`
	DiskUsagePercent     float64 `json:"disk_usage_percent" db:"disk_usage_percent"`
	DiskTotalBytes       int64   `json:"disk_total_bytes" db:"disk_total_bytes"`
	DiskUsedBytes        int64   `json:"disk_used_bytes" db:"disk_used_bytes"`
	DiskAvailableBytes   int64   `json:"disk_available_bytes" db:"disk_available_bytes"`
}

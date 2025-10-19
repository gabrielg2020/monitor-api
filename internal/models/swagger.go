package models

// Host represents a monitored Raspberry Pi
type Host struct {
	ID        int64  `json:"id" example:"1"`
	Hostname  string `json:"hostname" example:"pi-01"`
	IPAddress string `json:"ip_address" example:"192.168.0.24"`
	Role      string `json:"role" example:"server"`
}

// CreateHostRequest for registering a new host
type CreateHostRequest struct {
	Hostname  string `json:"hostname" binding:"required" example:"pi-01"`
	IPAddress string `json:"ip_address" binding:"required" example:"192.168.0.24"`
	Role      string `json:"role" example:"server"`
}

// HostListResponse contains list of hosts
type HostListResponse struct {
	Hosts []Host `json:"hosts"`
	Meta  Meta   `json:"meta"`
}

// SystemMetric represents system resource usage
type SystemMetric struct {
	ID                   int64   `json:"id" example:"1"`
	HostID               int64   `json:"host_id" example:"1"`
	Timestamp            int64   `json:"timestamp" example:"1729350000"`
	CPUUsage             float64 `json:"cpu_usage" example:"45.2"`
	MemoryUsagePercent   float64 `json:"memory_usage_percent" example:"67.8"`
	MemoryTotalBytes     int64   `json:"memory_total_bytes" example:"4294967296"`
	MemoryUsedBytes      int64   `json:"memory_used_bytes" example:"2910765875"`
	MemoryAvailableBytes int64   `json:"memory_available_bytes" example:"1384191421"`
	DiskUsagePercent     float64 `json:"disk_usage_percent" example:"23.4"`
	DiskTotalBytes       int64   `json:"disk_total_bytes" example:"32212254720"`
	DiskUsedBytes        int64   `json:"disk_used_bytes" example:"7537723520"`
	DiskAvailableBytes   int64   `json:"disk_available_bytes" example:"24674531200"`
}

// CreateMetricRequest for submitting new metrics
type CreateMetricRequest struct {
	Hostname             string  `json:"hostname" binding:"required" example:"pi-01"`
	CPUUsage             float64 `json:"cpu_usage" binding:"required" example:"45.2"`
	MemoryUsagePercent   float64 `json:"memory_usage_percent" binding:"required" example:"67.8"`
	MemoryTotalBytes     int64   `json:"memory_total_bytes" binding:"required" example:"4294967296"`
	MemoryUsedBytes      int64   `json:"memory_used_bytes" binding:"required" example:"2910765875"`
	MemoryAvailableBytes int64   `json:"memory_available_bytes" binding:"required" example:"1384191421"`
	DiskUsagePercent     float64 `json:"disk_usage_percent" binding:"required" example:"23.4"`
	DiskTotalBytes       int64   `json:"disk_total_bytes" binding:"required" example:"32212254720"`
	DiskUsedBytes        int64   `json:"disk_used_bytes" binding:"required" example:"7537723520"`
	DiskAvailableBytes   int64   `json:"disk_available_bytes" binding:"required" example:"24674531200"`
}

// MetricListResponse contains list of metrics
type MetricListResponse struct {
	Records []SystemMetric `json:"records"`
	Meta    Meta           `json:"meta"`
}

// MetricResponse for successful metric submission
type MetricResponse struct {
	Message string `json:"message" example:"Metric received successfully"`
	HostID  int64  `json:"host_id" example:"1"`
}

// Meta contains pagination and count information
type Meta struct {
	Count int `json:"count" example:"10"`
	Limit int `json:"limit" example:"100"`
}

// HealthResponse represents API health status
type HealthResponse struct {
	Status    string            `json:"status" example:"healthy"`
	Timestamp string            `json:"timestamp" example:"2025-10-19T18:30:00Z"`
	Checks    map[string]string `json:"checks"`
}

// ErrorResponse represents an error
type ErrorResponse struct {
	Error   string `json:"error" example:"Invalid request"`
	Details string `json:"details,omitempty"`
}

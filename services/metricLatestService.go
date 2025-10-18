package services

import (
	"database/sql"

	"github.com/gabrielg2020/monitor-page/entities"
)

type MetricLatestService struct {
	db *sql.DB
}

func NewMetricLatestService(con *sql.DB) *MetricLatestService {
	return &MetricLatestService{
		db: con,
	}
}

func (service *MetricLatestService) GetLatestMetrics(hostID *int64) (*entities.SystemMetric, error) {
	querySQL := `
        SELECT id, host_id, timestamp, cpu_usage, memory_usage_percent,
               memory_total_bytes, memory_used_bytes, memory_available_bytes,
               disk_usage_percent, disk_total_bytes, disk_used_bytes, disk_available_bytes
        FROM system_metrics`

	var args []interface{}

	// Add WHERE clauses conditionally
	if hostID != nil {
		querySQL += " WHERE host_id = ?"
		args = append(args, *hostID)
	}

	querySQL += " ORDER BY timestamp DESC LIMIT 1"

	var metric entities.SystemMetric

	err := service.db.QueryRow(querySQL, args...).Scan(
		&metric.ID,
		&metric.HostID,
		&metric.Timestamp,
		&metric.CPUUsage,
		&metric.MemoryUsagePercent,
		&metric.MemoryTotalBytes,
		&metric.MemoryUsedBytes,
		&metric.MemoryAvailableBytes,
		&metric.DiskUsagePercent,
		&metric.DiskTotalBytes,
		&metric.DiskUsedBytes,
		&metric.DiskAvailableBytes,
	)

	if err != nil {
		return nil, err
	}

	return &metric, nil
}

package services

import (
	"database/sql"

	"github.com/gabrielg2020/monitor-page/entities"
)

type MetricGetService struct {
	db           *sql.DB
	requestQuery *entities.MetricQueryParams
}

func NewMetricGetService(con *sql.DB) *MetricGetService {
	return &MetricGetService{
		db: con,
	}
}

func (service *MetricGetService) SetQueryParams(params *entities.MetricQueryParams) {
	service.requestQuery = params
}

func (service *MetricGetService) GetMetrics() ([]entities.SystemMetric, error) {
	querySQL := `
		SELECT id, host_id, timestamp, cpu_usage, memory_usage_percent,
			   memory_total_bytes, memory_used_bytes, memory_available_bytes,
			   disk_usage_percent, disk_total_bytes, disk_used_bytes, disk_available_bytes
		FROM system_metrics
		WHERE 1=1` // Dummy WHERE clause for easier appending

	var args []interface{}

	// Add WHERE clauses conditionally
	if service.requestQuery.HostID != nil {
		querySQL += " AND host_id = ?"
		args = append(args, *service.requestQuery.HostID)
	}

	if service.requestQuery.StartTime != nil {
		querySQL += " AND timestamp >= ?"
		args = append(args, *service.requestQuery.StartTime)
	}

	if service.requestQuery.EndTime != nil {
		querySQL += " AND timestamp <= ?"
		args = append(args, *service.requestQuery.EndTime)
	}

	// Order and limit has to be validated before setting it
	querySQL += " ORDER BY timestamp " + service.requestQuery.Order
	querySQL += " LIMIT ?"
	args = append(args, service.requestQuery.Limit)

	rows, err := service.db.Query(querySQL, args...)
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			panic(err)
		}
	}(rows)

	var metrics []entities.SystemMetric
	for rows.Next() {
		var metric entities.SystemMetric
		if err := rows.Scan(
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
		); err != nil {
			return nil, err
		}
		metrics = append(metrics, metric)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return metrics, nil
}

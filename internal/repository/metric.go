package repository

import (
	"database/sql"
	"errors"

	"github.com/gabrielg2020/monitor-api/internal/entities"
)

type MetricRepository struct {
	db *sql.DB
}

func NewMetricRepository(db *sql.DB) *MetricRepository {
	return &MetricRepository{db: db}
}

// FindByFilters retrieves metrics based on query parameters
func (repo *MetricRepository) FindByFilters(params *entities.MetricQueryParams) ([]entities.SystemMetric, error) {
	querySQL := `
		SELECT id, host_id, timestamp, cpu_usage, memory_usage_percent,
			   memory_total_bytes, memory_used_bytes, memory_available_bytes,
			   disk_usage_percent, disk_total_bytes, disk_used_bytes, disk_available_bytes
		FROM system_metrics
		WHERE 1=1`

	var args []interface{}

	if params.HostID != nil {
		querySQL += " AND host_id = ?"
		args = append(args, *params.HostID)
	}

	if params.StartTime != nil {
		querySQL += " AND timestamp >= ?"
		args = append(args, *params.StartTime)
	}

	if params.EndTime != nil {
		querySQL += " AND timestamp <= ?"
		args = append(args, *params.EndTime)
	}

	// Order and limit
	querySQL += " ORDER BY timestamp " + params.Order
	querySQL += " LIMIT ?"
	args = append(args, params.Limit)

	rows, err := repo.db.Query(querySQL, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return repo.scanMetrics(rows)
}

// FindLatest retrieves the most recent metric for a host
func (repo *MetricRepository) FindLatest(hostID *int64) (*entities.SystemMetric, error) {
	querySQL := `
        SELECT id, host_id, timestamp, cpu_usage, memory_usage_percent,
               memory_total_bytes, memory_used_bytes, memory_available_bytes,
               disk_usage_percent, disk_total_bytes, disk_used_bytes, disk_available_bytes
        FROM system_metrics`

	var args []interface{}

	if hostID != nil {
		querySQL += " WHERE host_id = ?"
		args = append(args, *hostID)
	}

	querySQL += " ORDER BY timestamp DESC LIMIT 1"

	var metric entities.SystemMetric
	err := repo.db.QueryRow(querySQL, args...).Scan(
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
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil // Not found
		}
		return nil, err
	}

	return &metric, nil
}

// Create inserts a new metric record
func (repo *MetricRepository) Create(metric *entities.SystemMetric) (int64, error) {
	insertSQL := `
		INSERT INTO system_metrics (
			host_id, timestamp, cpu_usage, memory_usage_percent,
			memory_total_bytes, memory_used_bytes, memory_available_bytes,
			disk_usage_percent, disk_total_bytes, disk_used_bytes, disk_available_bytes
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	result, err := repo.db.Exec(insertSQL,
		metric.HostID,
		metric.Timestamp,
		metric.CPUUsage,
		metric.MemoryUsagePercent,
		metric.MemoryTotalBytes,
		metric.MemoryUsedBytes,
		metric.MemoryAvailableBytes,
		metric.DiskUsagePercent,
		metric.DiskTotalBytes,
		metric.DiskUsedBytes,
		metric.DiskAvailableBytes,
	)

	if err != nil {
		return 0, err
	}

	return result.LastInsertId()
}

// scanMetrics is a helper to scan multiple rows into SystemMetric slice
func (repo *MetricRepository) scanMetrics(rows *sql.Rows) ([]entities.SystemMetric, error) {
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

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

// FindLatestByHost retrieves the most recent metric for each host
func (repo *MetricRepository) FindLatestByHost() ([]entities.SystemMetric, error) {
	querySQL := `
		SELECT m.id, m.host_id, m.timestamp, m.cpu_usage, m.memory_usage_percent,
			   m.memory_total_bytes, m.memory_used_bytes, m.memory_available_bytes,
			   m.disk_usage_percent, m.disk_total_bytes, m.disk_used_bytes, m.disk_available_bytes
		FROM system_metrics m
		INNER JOIN (
			SELECT host_id, MAX(timestamp) as max_timestamp
			FROM system_metrics
			GROUP BY host_id
		) latest ON m.host_id = latest.host_id AND m.timestamp = latest.max_timestamp
		ORDER BY m.host_id`

	rows, err := repo.db.Query(querySQL)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return repo.scanMetrics(rows)
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

// DeleteOlderThan deletes metrics older than the given timestamp
func (repo *MetricRepository) DeleteOlderThan(timestamp int64) (int64, error) {
	deleteSQL := `DELETE FROM system_metrics WHERE timestamp < ?`
	result, err := repo.db.Exec(deleteSQL, timestamp)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
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

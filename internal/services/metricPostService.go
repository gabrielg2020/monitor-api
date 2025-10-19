package services

import (
	"database/sql"

	"github.com/gabrielg2020/monitor-api/internal/entities"
)

type MetricPostService struct {
	db *sql.DB
}

func NewMetricPostService(con *sql.DB) *MetricPostService {
	return &MetricPostService{
		db: con,
	}
}

func (service *MetricPostService) PostMetric(record *entities.SystemMetric) (int64, error) {
	// Insert into database
	insertSQL := `
    INSERT INTO system_metrics (
    	host_id, timestamp, cpu_usage, memory_usage_percent,
        memory_total_bytes, memory_used_bytes, memory_available_bytes,
    	disk_usage_percent, disk_total_bytes, disk_used_bytes, disk_available_bytes
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	result, err := service.db.Exec(insertSQL,
		record.HostID,
		record.Timestamp,
		record.CPUUsage,
		record.MemoryUsagePercent,
		record.MemoryTotalBytes,
		record.MemoryUsedBytes,
		record.MemoryAvailableBytes,
		record.DiskUsagePercent,
		record.DiskTotalBytes,
		record.DiskUsedBytes,
		record.DiskAvailableBytes,
	)

	if err != nil {
		return 0, err
	}

	// Get the last inserted id
	lastInsertID, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return lastInsertID, nil
}

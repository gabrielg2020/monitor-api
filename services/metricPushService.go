package services

import (
	"database/sql"

	"github.com/gabrielg2020/monitor-page/entities"
)

type MetricPushService struct {
	db *sql.DB
}

func NewMetricPushService(con *sql.DB) *MetricPushService {
	return &MetricPushService{
		db: con,
	}
}

func (service *MetricPushService) MetricPushService(record *entities.SystemMetric) error {
	// Insert into database
	insertSQL := `
    INSERT INTO system_metrics (
    	host_id, timestamp, cpu_usage, memory_usage_percent,
        memory_total_bytes, memory_used_bytes, memory_available_bytes,
    	disk_usage_percent, disk_total_bytes, disk_used_bytes, disk_available_bytes
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	_, err := service.db.Exec(insertSQL,
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
		return err
	}

	return nil
}

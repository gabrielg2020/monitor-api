package repository

import "github.com/gabrielg2020/monitor-api/internal/entities"

// HealthRepositoryInterface defines methods for health checks
type HealthRepositoryInterface interface {
	CheckDatabaseConnection() error
	GetDatabaseStats() (map[string]interface{}, error)
	GetTableCounts() (map[string]int, error)
}

// HostRepositoryInterface defines methods for host repository operations
type HostRepositoryInterface interface {
	FindByFilters(params *entities.HostQueryParams) ([]entities.Host, error)
	Create(host *entities.Host) (int64, error)
	Update(id int64, host *entities.Host) error
	Delete(id int64) error
}

// MetricRepositoryInterface defines methods for metric repository operations
type MetricRepositoryInterface interface {
	FindByFilters(params *entities.MetricQueryParams) ([]entities.SystemMetric, error)
	FindLatest(hostID *int64) (*entities.SystemMetric, error)
	Create(metric *entities.SystemMetric) (int64, error)
}

var _ HealthRepositoryInterface = (*HealthRepository)(nil)
var _ HostRepositoryInterface = (*HostRepository)(nil)

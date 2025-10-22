package services

import "github.com/gabrielg2020/monitor-api/internal/entities"

// HealthServiceInterface defines methods for health checks
type HealthServiceInterface interface {
	CheckHealth() error
	GetDetailedHealth() (map[string]interface{}, error)
}

// HostServiceInterface defines methods for host service operations
type HostServiceInterface interface {
	CreateHost(host *entities.Host) (int64, error)
	GetHosts(params *entities.HostQueryParams) ([]entities.Host, error)
	UpdateHost(id int64, host *entities.Host) error
	DeleteHost(id int64) error
}

// MetricServiceInterface defines methods for metric service operations
type MetricServiceInterface interface {
	CreateMetric(metric *entities.SystemMetric) (int64, error)
	GetMetrics(params *entities.MetricQueryParams) ([]entities.SystemMetric, error)
	GetLatestMetric(hostID *int64) (*entities.SystemMetric, error)
}

var _ HealthServiceInterface = (*HealthService)(nil)
var _ HostServiceInterface = (*HostService)(nil)
var _ MetricServiceInterface = (*MetricService)(nil)

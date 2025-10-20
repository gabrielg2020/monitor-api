package services

import (
	"github.com/gabrielg2020/monitor-api/internal/entities"
	"github.com/gabrielg2020/monitor-api/internal/repository"
)

type MetricService struct {
	repo *repository.MetricRepository
}

func NewMetricService(repo *repository.MetricRepository) *MetricService {
	return &MetricService{repo: repo}
}

// GetMetrics retrieves metrics based on query parameters
func (service *MetricService) GetMetrics(params *entities.MetricQueryParams) ([]entities.SystemMetric, error) {
	return service.repo.FindByFilters(params)
}

// GetLatestMetric retrieves the most recent metric for a specific host or all hosts
func (service *MetricService) GetLatestMetric(hostID *int64) (*entities.SystemMetric, error) {
	return service.repo.FindLatest(hostID)
}

// GetLatestMetrics retrieves the most recent metric for each host
func (service *MetricService) GetLatestMetrics() ([]entities.SystemMetric, error) {
	return service.repo.FindLatestByHost()
}

// CreateMetric stores a new metric record
func (service *MetricService) CreateMetric(metric *entities.SystemMetric) (int64, error) {
	return service.repo.Create(metric)
}

// CleanupOldMetrics deletes metrics older than the given timestamp
func (service *MetricService) CleanupOldMetrics(timestamp int64) (int64, error) {
	return service.repo.DeleteOlderThan(timestamp)
}

package services

import (
	"github.com/gabrielg2020/monitor-api/internal/entities"
	"github.com/gabrielg2020/monitor-api/internal/repository"
)

type MetricService struct {
	repo repository.MetricRepositoryInterface
}

func NewMetricService(repo repository.MetricRepositoryInterface) *MetricService {
	return &MetricService{repo: repo}
}

// CreateMetric stores a new metric record
func (service *MetricService) CreateMetric(metric *entities.SystemMetric) (int64, error) {
	if err := ValidateSystemMetric(metric); err != nil {
		return -1, err
	}
	return service.repo.Create(metric)
}

// GetMetrics retrieves metrics based on query parameters
func (service *MetricService) GetMetrics(params *entities.MetricQueryParams) ([]entities.SystemMetric, error) {
	if params == nil {
		return nil, ErrNilQueryParams
	}

	if *params.HostID <= 0 {
		return nil, ErrInvalidHostID
	}

	return service.repo.FindByFilters(params)
}

// GetLatestMetric retrieves the most recent metric for a specific host or all hosts
func (service *MetricService) GetLatestMetric(hostID *int64) (*entities.SystemMetric, error) {
	return service.repo.FindLatest(hostID)
}

package mocks

import (
	"github.com/gabrielg2020/monitor-api/internal/entities"
	"github.com/stretchr/testify/mock"
)

// MockMetricRepository is a mock implementation of MetricRepositoryInterface
type MockMetricRepository struct {
	mock.Mock
}

// FindByFilters mocks finding metrics by filters
func (mock *MockMetricRepository) FindByFilters(params *entities.MetricQueryParams) ([]entities.SystemMetric, error) {
	args := mock.Called(params)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]entities.SystemMetric), args.Error(1)
}

// FindLatest mocks finding the latest metric
func (mock *MockMetricRepository) FindLatest(hostID *int64) (*entities.SystemMetric, error) {
	args := mock.Called(hostID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.SystemMetric), args.Error(1)
}

// Create mocks creating a new metric
func (mock *MockMetricRepository) Create(metric *entities.SystemMetric) (int64, error) {
	args := mock.Called(metric)
	return args.Get(0).(int64), args.Error(1)
}

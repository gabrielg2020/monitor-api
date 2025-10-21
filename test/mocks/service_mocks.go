package mocks

import (
	"github.com/gabrielg2020/monitor-api/internal/entities"
	"github.com/stretchr/testify/mock"
)

// MockHealthService is a mock implementation of HealthServiceInterface
type MockHealthService struct {
	mock.Mock
}

// CheckHealth mocks the health check
func (m *MockHealthService) CheckHealth() error {
	args := m.Called()
	return args.Error(0)
}

// GetDetailedHealth mocks getting detailed health information
func (m *MockHealthService) GetDetailedHealth() (map[string]interface{}, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(map[string]interface{}), args.Error(1)
}

// MockHostService is a mock implementation of HostServiceInterface
type MockHostService struct {
	mock.Mock
}

// CreateHost mocks creating a new host
func (m *MockHostService) CreateHost(host *entities.Host) (int64, error) {
	args := m.Called(host)
	return args.Get(0).(int64), args.Error(1)
}

// GetHosts mocks getting a host by ID
func (m *MockHostService) GetHosts(params *entities.HostQueryParams) ([]entities.Host, error) {
	args := m.Called(params)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]entities.Host), args.Error(1)
}

// UpdateHost mocks updating a host
func (m *MockHostService) UpdateHost(id int64, host *entities.Host) error {
	args := m.Called(id, host)
	return args.Error(0)
}

// DeleteHost mocks deleting a host
func (m *MockHostService) DeleteHost(id int64) error {
	args := m.Called(id)
	return args.Error(0)
}

// MockMetricService is a mock implementation of MetricServiceInterface
type MockMetricService struct {
	mock.Mock
}

// CreateMetric mocks creating a new metric
func (m *MockMetricService) CreateMetric(metric *entities.SystemMetric) (int64, error) {
	args := m.Called(metric)
	return args.Get(0).(int64), args.Error(1)
}

// GetMetrics mocks getting metrics based on query parameters
func (m *MockMetricService) GetMetrics(params *entities.MetricQueryParams) ([]entities.SystemMetric, error) {
	args := m.Called(params)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]entities.SystemMetric), args.Error(1)
}

// GetLatestMetric mocks getting the latest metric for a specific host or all hosts
func (m *MockMetricService) GetLatestMetric(hostID *int64) (*entities.SystemMetric, error) {
	args := m.Called(hostID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.SystemMetric), args.Error(1)
}

package mocks

import (
	"github.com/gabrielg2020/monitor-api/internal/entities"
	"github.com/stretchr/testify/mock"
)

// MockHealthRepository is a mock implementation of HealthRepositoryInterface
type MockHealthRepository struct {
	mock.Mock
}

// CheckDatabaseConnection mocks the database connection check
func (mock *MockHealthRepository) CheckDatabaseConnection() error {
	args := mock.Called()
	return args.Error(0)
}

// GetDatabaseStats mocks getting database statistics
func (mock *MockHealthRepository) GetDatabaseStats() (map[string]interface{}, error) {
	args := mock.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(map[string]interface{}), args.Error(1)
}

// GetTableCounts mocks getting table counts
func (mock *MockHealthRepository) GetTableCounts() (map[string]int, error) {
	args := mock.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(map[string]int), args.Error(1)
}

// MockHostRepository is a mock implementation of HostRepositoryInterface
type MockHostRepository struct {
	mock.Mock
}

// FindByFilters mocks finding hosts by filters
func (mock *MockHostRepository) FindByFilters(params *entities.HostQueryParams) ([]entities.Host, error) {
	args := mock.Called(params)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]entities.Host), args.Error(1)
}

// Create mocks creating a new host
func (mock *MockHostRepository) Create(host *entities.Host) (int64, error) {
	args := mock.Called(host)
	return args.Get(0).(int64), args.Error(1)
}

// Update mocks updating a host
func (mock *MockHostRepository) Update(id int64, host *entities.Host) error {
	args := mock.Called(id, host)
	return args.Error(0)
}

// Delete mocks deleting a host
func (mock *MockHostRepository) Delete(id int64) error {
	args := mock.Called(id)
	return args.Error(0)
}

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

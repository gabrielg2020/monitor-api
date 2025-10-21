package mocks

import (
	"github.com/gabrielg2020/monitor-api/internal/entities"
	"github.com/stretchr/testify/mock"
)

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

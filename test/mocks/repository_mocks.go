package mocks

import (
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

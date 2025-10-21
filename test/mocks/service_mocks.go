package mocks

import (
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

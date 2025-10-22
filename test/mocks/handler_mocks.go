package mocks

import (
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
)

// MockHealthHandler is a mock implementation of HealthHandlerInterface
type MockHealthHandler struct {
	mock.Mock
}

// GetHealth mocks the GetHealth handler method
func (m *MockHealthHandler) GetHealth(ctx *gin.Context) {
	m.Called(ctx)
}

// GetDetailedHealth mocks the GetDetailedHealth handler method
func (m *MockHealthHandler) GetDetailedHealth(ctx *gin.Context) {
	m.Called(ctx)
}

// MockHostHandler is a mock implementation of HostHandlerInterface
type MockHostHandler struct {
	mock.Mock
}

// Create mocks the Create handler method
func (m *MockHostHandler) Create(ctx *gin.Context) {
	m.Called(ctx)
}

// Get mocks the Get handler method
func (m *MockHostHandler) Get(ctx *gin.Context) {
	m.Called(ctx)
}

// Update mocks the Update handler method
func (m *MockHostHandler) Update(ctx *gin.Context) {
	m.Called(ctx)
}

// Delete mocks the Delete handler method
func (m *MockHostHandler) Delete(ctx *gin.Context) {
	m.Called(ctx)
}

// MockMetricHandler is a mock implementation of MetricHandlerInterface
type MockMetricHandler struct {
	mock.Mock
}

// Create mocks the Create handler method
func (m *MockMetricHandler) Create(ctx *gin.Context) {
	m.Called(ctx)
}

// Get mocks the Get handler method
func (m *MockMetricHandler) Get(ctx *gin.Context) {
	m.Called(ctx)
}

// GetLatest mocks the GetLatest handler method
func (m *MockMetricHandler) GetLatest(ctx *gin.Context) {
	m.Called(ctx)
}

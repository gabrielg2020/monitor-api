// nolint
package services

import (
	"errors"
	"testing"

	"github.com/gabrielg2020/monitor-api/test/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// HealthServiceTestSuite is the test suite for HealthService
type HealthServiceTestSuite struct {
	suite.Suite
	mockRepo *mocks.MockHealthRepository
	service  *HealthService
}

// SetupTest runs before each test in the suite
func (suite *HealthServiceTestSuite) SetupTest() {
	suite.mockRepo = new(mocks.MockHealthRepository)
	suite.service = NewHealthService(suite.mockRepo)
}

// TearDownTest runs after each test
func (suite *HealthServiceTestSuite) TearDownTest() {
	suite.mockRepo.AssertExpectations(suite.T())
}

// TestNewHealthService tests the constructor
func (suite *HealthServiceTestSuite) TestNewHealthService() {
	assert.NotNil(suite.T(), suite.service)
	assert.Equal(suite.T(), suite.mockRepo, suite.service.repo)
}

// TestCheckHealth tests the CheckHealth method
func (suite *HealthServiceTestSuite) TestCheckHealth() {
	tests := []struct {
		name          string
		setupMock     func()
		expectedError error
		description   string
	}{
		{
			name: "successful_health_check",
			setupMock: func() {
				suite.mockRepo.On("CheckDatabaseConnection").Return(nil).Once()
			},
			expectedError: nil,
			description:   "Should return nil when database connection is healthy",
		},
		{
			name: "failed_health_check",
			setupMock: func() {
				suite.mockRepo.On("CheckDatabaseConnection").Return(errors.New("connection refused")).Once()
			},
			expectedError: errors.New("connection refused"),
			description:   "Should return error when database connection fails",
		},
	}

	for _, test := range tests {
		suite.Run(test.name, func() {
			test.setupMock()

			err := suite.service.CheckHealth()

			if test.expectedError != nil {
				assert.Error(suite.T(), err)
				assert.Equal(suite.T(), test.expectedError.Error(), err.Error())
			} else {
				assert.NoError(suite.T(), err)
			}
		})

		// Reset mock for next test
		suite.TearDownTest()
		suite.SetupTest()
	}
}

// TestGetDetailedHealth tests the GetDetailedHealth method
func (suite *HealthServiceTestSuite) TestGetDetailedHealth() {
	tests := []struct {
		name           string
		setupMock      func()
		expectedResult map[string]interface{}
		expectedError  error
		validateResult func(result map[string]interface{})
		description    string
	}{
		{
			name: "all_checks_successful",
			setupMock: func() {
				// Database connection check succeeds
				suite.mockRepo.On("CheckDatabaseConnection").Return(nil).Once()

				// Database stats
				stats := map[string]interface{}{
					"open_connections": 5,
					"in_use":           2,
					"idle":             3,
					"max_open":         10,
				}
				suite.mockRepo.On("GetDatabaseStats").Return(stats, nil).Once()

				// Table counts
				counts := map[string]int{
					"hosts":   10,
					"metrics": 1000,
				}
				suite.mockRepo.On("GetTableCounts").Return(counts, nil).Once()
			},
			expectedError: nil,
			validateResult: func(result map[string]interface{}) {
				// Check database status
				dbStatus, ok := result["database"].(map[string]interface{})
				assert.True(suite.T(), ok)
				assert.Equal(suite.T(), "healthy", dbStatus["status"])

				// Check database stats
				dbStats, ok := result["database_stats"].(map[string]interface{})
				assert.True(suite.T(), ok)
				assert.Equal(suite.T(), 5, dbStats["open_connections"])
				assert.Equal(suite.T(), 2, dbStats["in_use"])

				// Check table counts
				tableCounts, ok := result["table_counts"].(map[string]int)
				assert.True(suite.T(), ok)
				assert.Equal(suite.T(), 10, tableCounts["hosts"])
				assert.Equal(suite.T(), 1000, tableCounts["metrics"])
			},
			description: "All checks successful",
		},
		{
			name: "database_connection_failed",
			setupMock: func() {
				suite.mockRepo.On("CheckDatabaseConnection").Return(errors.New("connection timeout")).Once()
			},
			expectedError: errors.New("connection timeout"),
			validateResult: func(result map[string]interface{}) {
				dbStatus, ok := result["database"].(map[string]interface{})
				assert.True(suite.T(), ok)
				assert.Equal(suite.T(), "unhealthy", dbStatus["status"])
				assert.Equal(suite.T(), "connection timeout", dbStatus["error"])
			},
			description: "Database connection fails",
		},
		{
			name: "database_stats_failed",
			setupMock: func() {
				suite.mockRepo.On("CheckDatabaseConnection").Return(nil).Once()
				suite.mockRepo.On("GetDatabaseStats").Return(nil, errors.New("stats error")).Once()

				counts := map[string]int{"hosts": 5, "metrics": 500}
				suite.mockRepo.On("GetTableCounts").Return(counts, nil).Once()
			},
			expectedError: nil,
			validateResult: func(result map[string]interface{}) {
				// Database should be healthy
				dbStatus, ok := result["database"].(map[string]interface{})
				assert.True(suite.T(), ok)
				assert.Equal(suite.T(), "healthy", dbStatus["status"])

				// Stats should not be present
				_, ok = result["database_stats"]
				assert.False(suite.T(), ok)

				// Table counts should still be present
				tableCounts, ok := result["table_counts"].(map[string]int)
				assert.True(suite.T(), ok)
				assert.Equal(suite.T(), 5, tableCounts["hosts"])
			},
			description: "Database stats retrieval fails",
		},
		{
			name: "table_counts_failed",
			setupMock: func() {
				suite.mockRepo.On("CheckDatabaseConnection").Return(nil).Once()

				stats := map[string]interface{}{"open_connections": 3}
				suite.mockRepo.On("GetDatabaseStats").Return(stats, nil).Once()

				suite.mockRepo.On("GetTableCounts").Return(nil, errors.New("count error")).Once()
			},
			expectedError: nil,
			validateResult: func(result map[string]interface{}) {
				// Database should be healthy
				dbStatus, ok := result["database"].(map[string]interface{})
				assert.True(suite.T(), ok)
				assert.Equal(suite.T(), "healthy", dbStatus["status"])

				// Stats should be present
				dbStats, ok := result["database_stats"].(map[string]interface{})
				assert.True(suite.T(), ok)
				assert.Equal(suite.T(), 3, dbStats["open_connections"])

				// Table counts should not be present
				_, ok = result["table_counts"]
				assert.False(suite.T(), ok)
			},
			description: "Table counts retrieval fails",
		},
	}

	for _, test := range tests {
		suite.Run(test.name, func() {
			test.setupMock()

			result, err := suite.service.GetDetailedHealth()

			if test.expectedError != nil {
				assert.Error(suite.T(), err)
				assert.Equal(suite.T(), test.expectedError.Error(), err.Error())
			} else {
				assert.NoError(suite.T(), err)
			}

			if test.validateResult != nil {
				test.validateResult(result)
			}
		})

		// Reset mock for next test
		suite.TearDownTest()
		suite.SetupTest()
	}
}

// TestGetDetailedHealth_PartialFailures tests partial failures in detailed health
func (suite *HealthServiceTestSuite) TestGetDetailedHealth_PartialFailures() {
	// This is a focused test for edge cases
	suite.mockRepo.On("CheckDatabaseConnection").Return(nil).Once()
	suite.mockRepo.On("GetDatabaseStats").Return(nil, errors.New("stats unavailable")).Once()
	suite.mockRepo.On("GetTableCounts").Return(nil, errors.New("counts unavailable")).Once()

	result, err := suite.service.GetDetailedHealth()

	// Should not return error for partial failures
	assert.NoError(suite.T(), err)

	// Database should still be marked as healthy
	dbStatus := result["database"].(map[string]interface{})
	assert.Equal(suite.T(), "healthy", dbStatus["status"])

	// Optional fields should not be present
	_, hasStats := result["database_stats"]
	assert.False(suite.T(), hasStats)

	_, hasCounts := result["table_counts"]
	assert.False(suite.T(), hasCounts)
}

// Run the test suite
func TestHealthServiceTestSuite(test *testing.T) {
	suite.Run(test, new(HealthServiceTestSuite))
}

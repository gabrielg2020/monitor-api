package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gabrielg2020/monitor-api/internal/models"
	"github.com/gabrielg2020/monitor-api/test/mocks"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// HealthHandlerTestSuite is the test suite for HealthHandler
type HealthHandlerTestSuite struct {
	suite.Suite
	router      *gin.Engine
	mockService *mocks.MockHealthService
	handler     *HealthHandler
}

// SetupTest runs before each test in the suite
func (suite *HealthHandlerTestSuite) SetupTest() {
	gin.SetMode(gin.TestMode)
	suite.router = gin.New()
	suite.mockService = new(mocks.MockHealthService)
	suite.handler = NewHealthHandler(suite.mockService)

	// Register routes
	suite.router.GET("/health", suite.handler.GetHealth)
	suite.router.GET("/health/detailed", suite.handler.GetDetailedHealth)
}

// TearDownTest runs after each test
func (suite *HealthHandlerTestSuite) TearDownTest() {
	suite.mockService.AssertExpectations(suite.T())
}

// TestNewHealthHandler tests the constructor
func (suite *HealthHandlerTestSuite) TestNewHealthHandler() {
	assert.NotNil(suite.T(), suite.handler)
	assert.NotNil(suite.T(), suite.handler.service)
}

// TestGetHealth tests the GetHealth endpoint
func (suite *HealthHandlerTestSuite) TestGetHealth() {
	tests := []struct {
		name               string
		setupMock          func()
		expectedStatus     int
		expectedHealthy    bool
		expectedDBCheck    string
		expectDBCheckField bool
	}{
		{
			name: "healthy_database",
			setupMock: func() {
				suite.mockService.On("CheckHealth").Return(nil).Once()
			},
			expectedStatus:     http.StatusOK,
			expectedHealthy:    true,
			expectedDBCheck:    "healthy",
			expectDBCheckField: true,
		},
		{
			name: "unhealthy_database",
			setupMock: func() {
				suite.mockService.On("CheckHealth").Return(errors.New("connection timeout")).Once()
			},
			expectedStatus:     http.StatusServiceUnavailable,
			expectedHealthy:    false,
			expectedDBCheck:    "unhealthy: connection timeout",
			expectDBCheckField: true,
		},
		{
			name: "database_connection_refused",
			setupMock: func() {
				suite.mockService.On("CheckHealth").Return(errors.New("connection refused")).Once()
			},
			expectedStatus:     http.StatusServiceUnavailable,
			expectedHealthy:    false,
			expectedDBCheck:    "unhealthy: connection refused",
			expectDBCheckField: true,
		},
	}

	for _, test := range tests {
		suite.Run(test.name, func() {
			test.setupMock()

			// Create request
			req, err := http.NewRequest(http.MethodGet, "/health", nil)
			assert.NoError(suite.T(), err)

			// Create response recorder
			w := httptest.NewRecorder()

			// Perform request
			suite.router.ServeHTTP(w, req)

			// Assert status code
			assert.Equal(suite.T(), test.expectedStatus, w.Code)

			// Parse response
			var response models.HealthResponse
			err = json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(suite.T(), err)

			// Assert response fields
			if test.expectedHealthy {
				assert.Equal(suite.T(), "healthy", response.Status)
			} else {
				assert.Equal(suite.T(), "unhealthy", response.Status)
			}

			assert.NotEmpty(suite.T(), response.Timestamp)

			if test.expectDBCheckField {
				assert.Contains(suite.T(), response.Checks, "database")
				assert.Equal(suite.T(), test.expectedDBCheck, response.Checks["database"])
			}
		})

		// Reset for next test
		suite.TearDownTest()
		suite.SetupTest()
	}
}

// TestGetDetailedHealth tests the GetDetailedHealth endpoint
func (suite *HealthHandlerTestSuite) TestGetDetailedHealth() {
	tests := []struct {
		name           string
		setupMock      func()
		expectedStatus int
		expectedHealth bool
		checkResponse  func(*testing.T, map[string]interface{})
	}{
		{
			name: "healthy_with_detailed_info",
			setupMock: func() {
				detailedHealth := map[string]interface{}{
					"database": map[string]interface{}{
						"status": "connected",
					},
					"database_stats": map[string]interface{}{
						"open_connections": 1,
						"in_use":           0,
						"idle":             1,
						"max_open":         10,
					},
					"table_counts": map[string]interface{}{
						"hosts":   5,
						"metrics": 100,
					},
				}
				suite.mockService.On("GetDetailedHealth").Return(detailedHealth, nil).Once()
			},
			expectedStatus: http.StatusOK,
			expectedHealth: true,
			checkResponse: func(t *testing.T, response map[string]interface{}) {
				assert.Equal(t, "healthy", response["status"])
				assert.NotEmpty(t, response["timestamp"])

				// Check database info
				database, ok := response["database"].(map[string]interface{})
				assert.True(t, ok)
				assert.Equal(t, "connected", database["status"])

				// Check database stats
				dbStats, ok := response["database_stats"].(map[string]interface{})
				assert.True(t, ok)
				assert.Equal(t, float64(1), dbStats["open_connections"])
				assert.Equal(t, float64(0), dbStats["in_use"])
				assert.Equal(t, float64(1), dbStats["idle"])
				assert.Equal(t, float64(10), dbStats["max_open"])

				// Check table counts
				tableCounts, ok := response["table_counts"].(map[string]interface{})
				assert.True(t, ok)
				assert.Equal(t, float64(5), tableCounts["hosts"])
				assert.Equal(t, float64(100), tableCounts["metrics"])
			},
		},
		{
			name: "unhealthy_with_error",
			setupMock: func() {
				suite.mockService.On("GetDetailedHealth").Return(nil, errors.New("database error")).Once()
			},
			expectedStatus: http.StatusServiceUnavailable,
			expectedHealth: false,
			checkResponse: func(t *testing.T, response map[string]interface{}) {
				assert.Equal(t, "unhealthy", response["status"])
				assert.NotEmpty(t, response["timestamp"])
			},
		},
		{
			name: "healthy_with_empty_tables",
			setupMock: func() {
				detailedHealth := map[string]interface{}{
					"database": map[string]interface{}{
						"status": "connected",
					},
					"database_stats": map[string]interface{}{
						"open_connections": 1,
						"in_use":           0,
						"idle":             1,
						"max_open":         10,
					},
					"table_counts": map[string]interface{}{
						"hosts":   0,
						"metrics": 0,
					},
				}
				suite.mockService.On("GetDetailedHealth").Return(detailedHealth, nil).Once()
			},
			expectedStatus: http.StatusOK,
			expectedHealth: true,
			checkResponse: func(t *testing.T, response map[string]interface{}) {
				assert.Equal(t, "healthy", response["status"])

				tableCounts, ok := response["table_counts"].(map[string]interface{})
				assert.True(t, ok)
				assert.Equal(t, float64(0), tableCounts["hosts"])
				assert.Equal(t, float64(0), tableCounts["metrics"])
			},
		},
		{
			name: "healthy_with_high_connection_count",
			setupMock: func() {
				detailedHealth := map[string]interface{}{
					"database": map[string]interface{}{
						"status": "connected",
					},
					"database_stats": map[string]interface{}{
						"open_connections": 50,
						"in_use":           25,
						"idle":             25,
						"max_open":         100,
					},
					"table_counts": map[string]interface{}{
						"hosts":          1000,
						"system_metrics": 50000,
					},
				}
				suite.mockService.On("GetDetailedHealth").Return(detailedHealth, nil).Once()
			},
			expectedStatus: http.StatusOK,
			expectedHealth: true,
			checkResponse: func(t *testing.T, response map[string]interface{}) {
				assert.Equal(t, "healthy", response["status"])

				dbStats, ok := response["database_stats"].(map[string]interface{})
				assert.True(t, ok)
				assert.Equal(t, float64(50), dbStats["open_connections"])
				assert.Equal(t, float64(25), dbStats["in_use"])
				assert.Equal(t, float64(25), dbStats["idle"])
				assert.Equal(t, float64(100), dbStats["max_open"])
			},
		},
		{
			name: "partial_data_returned",
			setupMock: func() {
				detailedHealth := map[string]interface{}{
					"database": map[string]interface{}{
						"status": "connected",
					},
				}
				suite.mockService.On("GetDetailedHealth").Return(detailedHealth, nil).Once()
			},
			expectedStatus: http.StatusOK,
			expectedHealth: true,
			checkResponse: func(t *testing.T, response map[string]interface{}) {
				assert.Equal(t, "healthy", response["status"])

				database, ok := response["database"].(map[string]interface{})
				assert.True(t, ok)
				assert.Equal(t, "connected", database["status"])
			},
		},
	}

	for _, test := range tests {
		suite.Run(test.name, func() {
			test.setupMock()

			// Create request
			req, err := http.NewRequest(http.MethodGet, "/health/detailed", nil)
			assert.NoError(suite.T(), err)

			// Create response recorder
			w := httptest.NewRecorder()

			// Perform request
			suite.router.ServeHTTP(w, req)

			// Assert status code
			assert.Equal(suite.T(), test.expectedStatus, w.Code)

			// Parse response
			var response map[string]interface{}
			err = json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(suite.T(), err)

			// Run custom response checks
			test.checkResponse(suite.T(), response)
		})

		// Reset for next test
		suite.TearDownTest()
		suite.SetupTest()
	}
}

// TestGetHealthContentType tests that the correct content type is returned
func (suite *HealthHandlerTestSuite) TestGetHealthContentType() {
	suite.mockService.On("CheckHealth").Return(nil).Once()

	req, err := http.NewRequest(http.MethodGet, "/health", nil)
	assert.NoError(suite.T(), err)

	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), "application/json; charset=utf-8", w.Header().Get("Content-Type"))
}

// TestGetDetailedHealthContentType tests that the correct content type is returned
func (suite *HealthHandlerTestSuite) TestGetDetailedHealthContentType() {
	detailedHealth := map[string]interface{}{
		"database": map[string]interface{}{
			"status": "connected",
		},
	}
	suite.mockService.On("GetDetailedHealth").Return(detailedHealth, nil).Once()

	req, err := http.NewRequest(http.MethodGet, "/health/detailed", nil)
	assert.NoError(suite.T(), err)

	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), "application/json; charset=utf-8", w.Header().Get("Content-Type"))
}

// Run the test suite
func TestHealthHandlerTestSuite(test *testing.T) {
	suite.Run(test, new(HealthHandlerTestSuite))
}

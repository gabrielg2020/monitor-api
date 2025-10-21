package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gabrielg2020/monitor-api/internal/entities"
	"github.com/gabrielg2020/monitor-api/internal/models"
	"github.com/gabrielg2020/monitor-api/test/mocks"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

// MetricHandlerTestSuite is the test suite for MetricHandler
type MetricHandlerTestSuite struct {
	suite.Suite
	router      *gin.Engine
	mockService *mocks.MockMetricService
	handler     *MetricHandler
}

// SetupTest runs before each test in the suite
func (suite *MetricHandlerTestSuite) SetupTest() {
	gin.SetMode(gin.TestMode)
	suite.router = gin.New()
	suite.mockService = new(mocks.MockMetricService)
	suite.handler = NewMetricHandler(suite.mockService)

	// Register routes
	suite.router.POST("/metrics", suite.handler.Create)
	suite.router.GET("/metrics", suite.handler.Get)
	suite.router.GET("/metrics/latest", suite.handler.GetLatest)
}

// TearDownTest runs after each test
func (suite *MetricHandlerTestSuite) TearDownTest() {
	suite.mockService.AssertExpectations(suite.T())
}

// TestNewMetricHandler tests the constructor
func (suite *MetricHandlerTestSuite) TestNewMetricHandler() {
	assert.NotNil(suite.T(), suite.handler)
	assert.NotNil(suite.T(), suite.handler.service)
}

// TestCreate tests the Create endpoint
func (suite *MetricHandlerTestSuite) TestCreate() {
	tests := []struct {
		name           string
		requestBody    interface{}
		setupMock      func()
		expectedStatus int
		checkResponse  func(*testing.T, *httptest.ResponseRecorder)
	}{
		{
			name: "successful_creation",
			requestBody: map[string]interface{}{
				"record": map[string]interface{}{
					"host_id":                1,
					"timestamp":              1609459200,
					"cpu_usage":              45.5,
					"memory_usage_percent":   60.0,
					"memory_total_bytes":     16000000000,
					"memory_used_bytes":      9600000000,
					"memory_available_bytes": 6400000000,
					"disk_usage_percent":     75.0,
					"disk_total_bytes":       500000000000,
					"disk_used_bytes":        375000000000,
					"disk_available_bytes":   125000000000,
				},
			},
			setupMock: func() {
				suite.mockService.On("CreateMetric", &entities.SystemMetric{
					HostID:               1,
					Timestamp:            1609459200,
					CPUUsage:             45.5,
					MemoryUsagePercent:   60.0,
					MemoryTotalBytes:     16000000000,
					MemoryUsedBytes:      9600000000,
					MemoryAvailableBytes: 6400000000,
					DiskUsagePercent:     75.0,
					DiskTotalBytes:       500000000000,
					DiskUsedBytes:        375000000000,
					DiskAvailableBytes:   125000000000,
				}).Return(int64(1), nil).Once()
			},
			expectedStatus: http.StatusCreated,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, "Metric create successfully", response["message"])
				assert.Equal(t, float64(1), response["id"])
			},
		},
		{
			name:           "invalid_json_body",
			requestBody:    "invalid json",
			setupMock:      func() {},
			expectedStatus: http.StatusBadRequest,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response models.ErrorResponse
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, "Invalid request body", response.Error)
				assert.NotEmpty(t, response.Details)
			},
		},
		{
			name: "missing_required_fields",
			requestBody: map[string]interface{}{
				"record": map[string]interface{}{
					"host_id": 1,
					// Missing other required fields
				},
			},
			setupMock: func() {
				suite.mockService.On("CreateMetric", &entities.SystemMetric{
					HostID: 1,
				}).Return(int64(0), errors.New("missing required fields")).Once()
			},
			expectedStatus: http.StatusInternalServerError,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response models.ErrorResponse
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, "Failed to create metric record", response.Error)
			},
		},
		{
			name: "foreign_key_constraint_error",
			requestBody: map[string]interface{}{
				"record": map[string]interface{}{
					"host_id":                999,
					"timestamp":              1609459200,
					"cpu_usage":              45.5,
					"memory_usage_percent":   60.0,
					"memory_total_bytes":     16000000000,
					"memory_used_bytes":      9600000000,
					"memory_available_bytes": 6400000000,
					"disk_usage_percent":     75.0,
					"disk_total_bytes":       500000000000,
					"disk_used_bytes":        375000000000,
					"disk_available_bytes":   125000000000,
				},
			},
			setupMock: func() {
				suite.mockService.On("CreateMetric", &entities.SystemMetric{
					HostID:               999,
					Timestamp:            1609459200,
					CPUUsage:             45.5,
					MemoryUsagePercent:   60.0,
					MemoryTotalBytes:     16000000000,
					MemoryUsedBytes:      9600000000,
					MemoryAvailableBytes: 6400000000,
					DiskUsagePercent:     75.0,
					DiskTotalBytes:       500000000000,
					DiskUsedBytes:        375000000000,
					DiskAvailableBytes:   125000000000,
				}).Return(int64(0), errors.New("FOREIGN KEY constraint failed")).Once()
			},
			expectedStatus: http.StatusInternalServerError,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response models.ErrorResponse
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, "Failed to create metric record", response.Error)
				assert.Contains(t, response.Details, "FOREIGN KEY constraint")
			},
		},
		{
			name: "database_connection_error",
			requestBody: map[string]interface{}{
				"record": map[string]interface{}{
					"host_id":                1,
					"timestamp":              1609459200,
					"cpu_usage":              45.5,
					"memory_usage_percent":   60.0,
					"memory_total_bytes":     16000000000,
					"memory_used_bytes":      9600000000,
					"memory_available_bytes": 6400000000,
					"disk_usage_percent":     75.0,
					"disk_total_bytes":       500000000000,
					"disk_used_bytes":        375000000000,
					"disk_available_bytes":   125000000000,
				},
			},
			setupMock: func() {
				suite.mockService.On("CreateMetric", &entities.SystemMetric{
					HostID:               1,
					Timestamp:            1609459200,
					CPUUsage:             45.5,
					MemoryUsagePercent:   60.0,
					MemoryTotalBytes:     16000000000,
					MemoryUsedBytes:      9600000000,
					MemoryAvailableBytes: 6400000000,
					DiskUsagePercent:     75.0,
					DiskTotalBytes:       500000000000,
					DiskUsedBytes:        375000000000,
					DiskAvailableBytes:   125000000000,
				}).Return(int64(0), errors.New("database connection lost")).Once()
			},
			expectedStatus: http.StatusInternalServerError,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response models.ErrorResponse
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, "Failed to create metric record", response.Error)
				assert.Equal(t, "database connection lost", response.Details)
			},
		},
	}

	for _, test := range tests {
		suite.Run(test.name, func() {
			test.setupMock()

			// Create request body
			var bodyBytes []byte
			var err error
			if str, ok := test.requestBody.(string); ok {
				bodyBytes = []byte(str)
			} else {
				bodyBytes, err = json.Marshal(test.requestBody)
				assert.NoError(suite.T(), err)
			}

			// Create request
			req, err := http.NewRequest(http.MethodPost, "/metrics", bytes.NewBuffer(bodyBytes))
			assert.NoError(suite.T(), err)
			req.Header.Set("Content-Type", "application/json")

			// Create response recorder
			w := httptest.NewRecorder()

			// Perform request
			suite.router.ServeHTTP(w, req)

			// Assert status code
			assert.Equal(suite.T(), test.expectedStatus, w.Code)

			// Run custom response checks
			test.checkResponse(suite.T(), w)
		})

		// Reset for next test
		suite.TearDownTest()
		suite.SetupTest()
	}
}

// TestGet tests the Get endpoint
func (suite *MetricHandlerTestSuite) TestGet() {
	tests := []struct {
		name           string
		queryParams    string
		setupMock      func()
		expectedStatus int
		checkResponse  func(*testing.T, *httptest.ResponseRecorder)
	}{
		{
			name:        "get_metrics_with_defaults",
			queryParams: "",
			setupMock: func() {
				metrics := []entities.SystemMetric{
					{
						ID:                   1,
						HostID:               1,
						Timestamp:            1609459200,
						CPUUsage:             45.5,
						MemoryUsagePercent:   60.0,
						MemoryTotalBytes:     16000000000,
						MemoryUsedBytes:      9600000000,
						MemoryAvailableBytes: 6400000000,
						DiskUsagePercent:     75.0,
						DiskTotalBytes:       500000000000,
						DiskUsedBytes:        375000000000,
						DiskAvailableBytes:   125000000000,
					},
				}
				// Use MatchedBy to match any params with Limit=100 and Order=DESC
				suite.mockService.On("GetMetrics", mock.MatchedBy(func(params *entities.MetricQueryParams) bool {
					return params.Limit == 100 && params.Order == "DESC" && params.HostID == nil
				})).Return(metrics, nil).Once()
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response models.MetricListResponse
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.GreaterOrEqual(t, len(response.Records), 0)
				assert.Equal(t, len(response.Records), response.Meta.Count)
			},
		},
		{
			name:        "get_metrics_with_host_id",
			queryParams: "?host_id=1",
			setupMock: func() {
				metrics := []entities.SystemMetric{
					{
						ID:                   1,
						HostID:               1,
						Timestamp:            1609459200,
						CPUUsage:             45.5,
						MemoryUsagePercent:   60.0,
						MemoryTotalBytes:     16000000000,
						MemoryUsedBytes:      9600000000,
						MemoryAvailableBytes: 6400000000,
						DiskUsagePercent:     75.0,
						DiskTotalBytes:       500000000000,
						DiskUsedBytes:        375000000000,
						DiskAvailableBytes:   125000000000,
					},
				}
				suite.mockService.On("GetMetrics", mock.MatchedBy(func(params *entities.MetricQueryParams) bool {
					return params.Limit == 100 && params.Order == "DESC" && params.HostID != nil && *params.HostID == 1
				})).Return(metrics, nil).Once()
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response models.MetricListResponse
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.GreaterOrEqual(t, len(response.Records), 0)
			},
		},
		{
			name:        "get_metrics_with_time_range",
			queryParams: "?start_time=1609459200&end_time=1609545600",
			setupMock: func() {
				metrics := []entities.SystemMetric{
					{
						ID:                   2,
						HostID:               2,
						Timestamp:            1609500000,
						CPUUsage:             30.0,
						MemoryUsagePercent:   50.0,
						MemoryTotalBytes:     8000000000,
						MemoryUsedBytes:      4000000000,
						MemoryAvailableBytes: 4000000000,
						DiskUsagePercent:     60.0,
						DiskTotalBytes:       250000000000,
						DiskUsedBytes:        150000000000,
						DiskAvailableBytes:   100000000000,
					},
				}
				suite.mockService.On("GetMetrics", mock.MatchedBy(func(params *entities.MetricQueryParams) bool {
					return params.Limit == 100 &&
						params.Order == "DESC" &&
						params.StartTime != nil && *params.StartTime == 1609459200 &&
						params.EndTime != nil && *params.EndTime == 1609545600
				})).Return(metrics, nil).Once()
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response models.MetricListResponse
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.GreaterOrEqual(t, len(response.Records), 0)
			},
		},
		{
			name:        "get_metrics_with_limit",
			queryParams: "?limit=50",
			setupMock: func() {
				var metrics []entities.SystemMetric
				for i := 0; i < 50; i++ {
					metrics = append(metrics, entities.SystemMetric{
						ID:                   int64(i + 1),
						HostID:               1,
						Timestamp:            1609459200 + int64(i*60),
						CPUUsage:             45.5,
						MemoryUsagePercent:   60.0,
						MemoryTotalBytes:     16000000000,
						MemoryUsedBytes:      9600000000,
						MemoryAvailableBytes: 6400000000,
						DiskUsagePercent:     75.0,
						DiskTotalBytes:       500000000000,
						DiskUsedBytes:        375000000000,
						DiskAvailableBytes:   125000000000,
					})
				}
				suite.mockService.On("GetMetrics", mock.MatchedBy(func(params *entities.MetricQueryParams) bool {
					return params.Limit == 50 && params.Order == "DESC"
				})).Return(metrics, nil).Once()
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response models.MetricListResponse
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.LessOrEqual(t, len(response.Records), 50)
				assert.Equal(t, 50, response.Meta.Limit)
			},
		},
		{
			name:        "get_metrics_with_asc_order",
			queryParams: "?order=ASC",
			setupMock: func() {
				metrics := []entities.SystemMetric{
					{
						ID:                   1,
						HostID:               1,
						Timestamp:            1609459200,
						CPUUsage:             45.5,
						MemoryUsagePercent:   60.0,
						MemoryTotalBytes:     16000000000,
						MemoryUsedBytes:      9600000000,
						MemoryAvailableBytes: 6400000000,
						DiskUsagePercent:     75.0,
						DiskTotalBytes:       500000000000,
						DiskUsedBytes:        375000000000,
						DiskAvailableBytes:   125000000000,
					},
				}
				suite.mockService.On("GetMetrics", mock.MatchedBy(func(params *entities.MetricQueryParams) bool {
					return params.Limit == 100 && params.Order == "ASC"
				})).Return(metrics, nil).Once()
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response models.MetricListResponse
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.GreaterOrEqual(t, len(response.Records), 0)
			},
		},
		{
			name:        "get_metrics_empty_result",
			queryParams: "?host_id=999",
			setupMock: func() {
				var metrics []entities.SystemMetric
				suite.mockService.On("GetMetrics", mock.MatchedBy(func(params *entities.MetricQueryParams) bool {
					return params.HostID != nil && *params.HostID == 999
				})).Return(metrics, nil).Once()
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response models.MetricListResponse
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, 0, response.Meta.Count)
			},
		},
		{
			name:           "invalid_query_parameter",
			queryParams:    "?host_id=invalid",
			setupMock:      func() {},
			expectedStatus: http.StatusBadRequest,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response models.ErrorResponse
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, "Invalid query parameters", response.Error)
			},
		},
		{
			name:           "invalid_order",
			queryParams:    "?order=INVALID",
			setupMock:      func() {},
			expectedStatus: http.StatusBadRequest,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response models.ErrorResponse
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
			},
		},
		{
			name:        "database_error",
			queryParams: "",
			setupMock: func() {
				suite.mockService.On("GetMetrics", mock.MatchedBy(func(params *entities.MetricQueryParams) bool {
					return params.Limit == 100 && params.Order == "DESC"
				})).Return(nil, errors.New("database connection lost")).Once()
			},
			expectedStatus: http.StatusInternalServerError,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response models.ErrorResponse
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, "Failed to retrieve metrics", response.Error)
			},
		},
	}

	for _, test := range tests {
		suite.Run(test.name, func() {
			test.setupMock()

			// Create request
			req, err := http.NewRequest(http.MethodGet, "/metrics"+test.queryParams, nil)
			assert.NoError(suite.T(), err)

			// Create response recorder
			w := httptest.NewRecorder()

			// Perform request
			suite.router.ServeHTTP(w, req)

			// Assert status code
			assert.Equal(suite.T(), test.expectedStatus, w.Code)

			// Run custom response checks
			test.checkResponse(suite.T(), w)
		})

		// Reset for next test
		suite.TearDownTest()
		suite.SetupTest()
	}
}

// TestGetLatest tests the GetLatest endpoint
func (suite *MetricHandlerTestSuite) TestGetLatest() {
	hostID := int64(1)

	tests := []struct {
		name           string
		queryParams    string
		setupMock      func()
		expectedStatus int
		checkResponse  func(*testing.T, *httptest.ResponseRecorder)
	}{
		{
			name:        "get_latest_with_host_id",
			queryParams: "?host_id=1",
			setupMock: func() {
				metric := &entities.SystemMetric{
					ID:                   1,
					HostID:               1,
					Timestamp:            1609545600,
					CPUUsage:             55.5,
					MemoryUsagePercent:   70.0,
					MemoryTotalBytes:     16000000000,
					MemoryUsedBytes:      11200000000,
					MemoryAvailableBytes: 4800000000,
					DiskUsagePercent:     85.0,
					DiskTotalBytes:       500000000000,
					DiskUsedBytes:        425000000000,
					DiskAvailableBytes:   75000000000,
				}
				suite.mockService.On("GetLatestMetric", &hostID).Return(metric, nil).Once()
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Contains(t, response, "metric")
				metric := response["metric"].(map[string]interface{})
				assert.Equal(t, float64(1), metric["id"])
				assert.Equal(t, float64(1), metric["host_id"])
			},
		},
		{
			name:        "get_latest_without_host_id",
			queryParams: "",
			setupMock: func() {
				metric := &entities.SystemMetric{
					ID:                   2,
					HostID:               3,
					Timestamp:            1609545600,
					CPUUsage:             45.0,
					MemoryUsagePercent:   65.0,
					MemoryTotalBytes:     32000000000,
					MemoryUsedBytes:      20800000000,
					MemoryAvailableBytes: 11200000000,
					DiskUsagePercent:     75.0,
					DiskTotalBytes:       1000000000000,
					DiskUsedBytes:        750000000000,
					DiskAvailableBytes:   250000000000,
				}
				var nilHostID *int64
				suite.mockService.On("GetLatestMetric", nilHostID).Return(metric, nil).Once()
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Contains(t, response, "metric")
			},
		},
		{
			name:        "metric_not_found",
			queryParams: "?host_id=999",
			setupMock: func() {
				nonExistentHostID := int64(999)
				suite.mockService.On("GetLatestMetric", &nonExistentHostID).Return(nil, nil).Once()
			},
			expectedStatus: http.StatusNotFound,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response models.ErrorResponse
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, "Metric not found", response.Error)
				assert.Contains(t, response.Details, "No latest metric found")
			},
		},
		{
			name:           "invalid_query_parameter",
			queryParams:    "?host_id=invalid",
			setupMock:      func() {},
			expectedStatus: http.StatusBadRequest,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response models.ErrorResponse
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, "Invalid query parameters", response.Error)
			},
		},
		{
			name:        "database_error",
			queryParams: "?host_id=1",
			setupMock: func() {
				suite.mockService.On("GetLatestMetric", &hostID).Return(nil, errors.New("database connection lost")).Once()
			},
			expectedStatus: http.StatusInternalServerError,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response models.ErrorResponse
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, "Failed to retrieve latest metric", response.Error)
				assert.Equal(t, "database connection lost", response.Details)
			},
		},
	}

	for _, test := range tests {
		suite.Run(test.name, func() {
			test.setupMock()

			// Create request
			req, err := http.NewRequest(http.MethodGet, "/metrics/latest"+test.queryParams, nil)
			assert.NoError(suite.T(), err)

			// Create response recorder
			w := httptest.NewRecorder()

			// Perform request
			suite.router.ServeHTTP(w, req)

			// Assert status code
			assert.Equal(suite.T(), test.expectedStatus, w.Code)

			// Run custom response checks
			test.checkResponse(suite.T(), w)
		})

		// Reset for next test
		suite.TearDownTest()
		suite.SetupTest()
	}
}

// TestContentType tests that the correct content type is returned
func (suite *MetricHandlerTestSuite) TestContentType() {
	metrics := []entities.SystemMetric{
		{
			ID:                   1,
			HostID:               1,
			Timestamp:            1609459200,
			CPUUsage:             45.5,
			MemoryUsagePercent:   60.0,
			MemoryTotalBytes:     16000000000,
			MemoryUsedBytes:      9600000000,
			MemoryAvailableBytes: 6400000000,
			DiskUsagePercent:     75.0,
			DiskTotalBytes:       500000000000,
			DiskUsedBytes:        375000000000,
			DiskAvailableBytes:   125000000000,
		},
	}

	suite.mockService.On("GetMetrics", mock.MatchedBy(func(params *entities.MetricQueryParams) bool {
		return params.Limit == 100 && params.Order == "DESC"
	})).Return(metrics, nil).Once()

	req, err := http.NewRequest(http.MethodGet, "/metrics", nil)
	assert.NoError(suite.T(), err)

	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), "application/json; charset=utf-8", w.Header().Get("Content-Type"))
}

// Run the test suite
func TestMetricHandlerTestSuite(test *testing.T) {
	suite.Run(test, new(MetricHandlerTestSuite))
}

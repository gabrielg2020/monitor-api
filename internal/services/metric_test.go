// nolint
package services

import (
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/gabrielg2020/monitor-api/internal/entities"
	"github.com/gabrielg2020/monitor-api/test/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// MetricServiceTestSuite is the test suite for MetricService
type MetricServiceTestSuite struct {
	suite.Suite
	mockRepo *mocks.MockMetricRepository
	service  *MetricService
}

// SetupTest runs before each test in the suite
func (suite *MetricServiceTestSuite) SetupTest() {
	suite.mockRepo = new(mocks.MockMetricRepository)
	suite.service = NewMetricService(suite.mockRepo)
}

// TearDownTest runs after each test
func (suite *MetricServiceTestSuite) TearDownTest() {
	suite.mockRepo.AssertExpectations(suite.T())
}

// TestNewMetricService tests the constructor
func (suite *MetricServiceTestSuite) TestNewMetricService() {
	assert.NotNil(suite.T(), suite.service)
	assert.Equal(suite.T(), suite.mockRepo, suite.service.repo)
}

// TestCreateMetric tests the CreateMetric method
func (suite *MetricServiceTestSuite) TestCreateMetric() {
	timestamp := time.Now().Unix()

	tests := []struct {
		name          string
		metric        *entities.SystemMetric
		setupMock     func()
		expectedID    int64
		expectedError error
		description   string
	}{
		{
			name: "successful_metric_creation",
			metric: &entities.SystemMetric{
				HostID:               1,
				Timestamp:            timestamp,
				CPUUsage:             45.5,
				MemoryUsagePercent:   67.8,
				MemoryTotalBytes:     16777216000,
				MemoryUsedBytes:      11387420000,
				MemoryAvailableBytes: 5389796000,
				DiskUsagePercent:     78.2,
				DiskTotalBytes:       500000000000,
				DiskUsedBytes:        391000000000,
				DiskAvailableBytes:   109000000000,
			},
			setupMock: func() {
				suite.mockRepo.On("Create", &entities.SystemMetric{
					HostID:               1,
					Timestamp:            timestamp,
					CPUUsage:             45.5,
					MemoryUsagePercent:   67.8,
					MemoryTotalBytes:     16777216000,
					MemoryUsedBytes:      11387420000,
					MemoryAvailableBytes: 5389796000,
					DiskUsagePercent:     78.2,
					DiskTotalBytes:       500000000000,
					DiskUsedBytes:        391000000000,
					DiskAvailableBytes:   109000000000,
				}).Return(int64(1), nil).Once()
			},
			expectedID:    1,
			expectedError: nil,
			description:   "Should return the new metric ID on successful creation",
		},
		{
			name: "invalid_host_id",
			metric: &entities.SystemMetric{
				HostID:             0,
				Timestamp:          timestamp,
				CPUUsage:           45.5,
				MemoryUsagePercent: 67.8,
				DiskUsagePercent:   78.2,
			},
			setupMock:     func() {}, // No mock call expected due to validation
			expectedID:    -1,
			expectedError: ErrInvalidHostID,
			description:   "Should return an error when trying to create a metric with invalid HostID",
		},
		{
			name: "negative_host_id",
			metric: &entities.SystemMetric{
				HostID:             -1,
				Timestamp:          timestamp,
				CPUUsage:           45.5,
				MemoryUsagePercent: 67.8,
				DiskUsagePercent:   78.2,
			},
			setupMock:     func() {},
			expectedID:    -1,
			expectedError: ErrInvalidHostID,
			description:   "Should return an error when trying to create a metric with invalid HostID",
		},
		{
			name: "invalid_cpu_usage_over_100",
			metric: &entities.SystemMetric{
				HostID:             1,
				Timestamp:          timestamp,
				CPUUsage:           101.5,
				MemoryUsagePercent: 67.8,
				DiskUsagePercent:   78.2,
			},
			setupMock:     func() {},
			expectedID:    -1,
			expectedError: ErrInvalidCPUUsage,
			description:   "Should return an error when trying to create a metric with invalid CPU usage",
		},
		{
			name: "invalid_cpu_usage_negative",
			metric: &entities.SystemMetric{
				HostID:             1,
				Timestamp:          timestamp,
				CPUUsage:           -5.0,
				MemoryUsagePercent: 67.8,
				DiskUsagePercent:   78.2,
			},
			setupMock:     func() {},
			expectedID:    -1,
			expectedError: ErrInvalidCPUUsage,
			description:   "Should return an error when trying to create a metric with invalid CPU usage",
		},
		{
			name: "invalid_memory_usage_over_100",
			metric: &entities.SystemMetric{
				HostID:             1,
				Timestamp:          timestamp,
				CPUUsage:           45.5,
				MemoryUsagePercent: 105.0,
				DiskUsagePercent:   78.2,
			},
			setupMock:     func() {},
			expectedID:    -1,
			expectedError: ErrInvalidMemoryUsage,
			description:   "Should return an error when trying to create a metric with invalid Memory usage",
		},
		{
			name: "invalid_disk_usage_negative",
			metric: &entities.SystemMetric{
				HostID:             1,
				Timestamp:          timestamp,
				CPUUsage:           45.5,
				MemoryUsagePercent: 67.8,
				DiskUsagePercent:   -10.0,
			},
			setupMock:     func() {},
			expectedID:    -1,
			expectedError: ErrInvalidDiskUsage,
			description:   "Should return an error when trying to create a metric with invalid Disk usage",
		},
		{
			name: "database_connection_error",
			metric: &entities.SystemMetric{
				HostID:             1,
				Timestamp:          timestamp,
				CPUUsage:           45.5,
				MemoryUsagePercent: 67.8,
				DiskUsagePercent:   78.2,
			},
			setupMock: func() {
				suite.mockRepo.On("Create", &entities.SystemMetric{
					HostID:             1,
					Timestamp:          timestamp,
					CPUUsage:           45.5,
					MemoryUsagePercent: 67.8,
					DiskUsagePercent:   78.2,
				}).Return(int64(-1), errors.New("database connection lost")).Once()
			},
			expectedID:    -1,
			expectedError: errors.New("database connection lost"),
			description:   "Should return an error when there is a database connection issue",
		},
	}

	for _, test := range tests {
		suite.Run(test.name, func() {
			test.setupMock()

			id, err := suite.service.CreateMetric(test.metric)

			assert.Equal(suite.T(), test.expectedID, id)
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

// TestGetMetrics tests the GetMetrics method
func (suite *MetricServiceTestSuite) TestGetMetrics() {
	timestamp := time.Now().Unix()
	startTime := timestamp - 3600
	endTime := timestamp
	hostID := int64(1)

	tests := []struct {
		name            string
		params          *entities.MetricQueryParams
		setupMock       func()
		expectedMetrics []entities.SystemMetric
		expectedError   error
		description     string
	}{
		{
			name: "get_metrics_with_all_filters",
			params: &entities.MetricQueryParams{
				HostID:    &hostID,
				StartTime: &startTime,
				EndTime:   &endTime,
				Order:     "ASC",
				Limit:     50,
			},
			setupMock: func() {
				metrics := []entities.SystemMetric{
					{
						ID:                 1,
						HostID:             1,
						Timestamp:          startTime + 100,
						CPUUsage:           45.5,
						MemoryUsagePercent: 67.8,
						DiskUsagePercent:   78.2,
					},
					{
						ID:                 2,
						HostID:             1,
						Timestamp:          startTime + 200,
						CPUUsage:           48.2,
						MemoryUsagePercent: 69.1,
						DiskUsagePercent:   78.5,
					},
				}
				suite.mockRepo.On("FindByFilters", &entities.MetricQueryParams{
					HostID:    &hostID,
					StartTime: &startTime,
					EndTime:   &endTime,
					Order:     "ASC",
					Limit:     50,
				}).Return(metrics, nil).Once()
			},
			expectedMetrics: []entities.SystemMetric{
				{
					ID:                 1,
					HostID:             1,
					Timestamp:          startTime + 100,
					CPUUsage:           45.5,
					MemoryUsagePercent: 67.8,
					DiskUsagePercent:   78.2,
				},
				{
					ID:                 2,
					HostID:             1,
					Timestamp:          startTime + 200,
					CPUUsage:           48.2,
					MemoryUsagePercent: 69.1,
					DiskUsagePercent:   78.5,
				},
			},
			expectedError: nil,
			description:   "Should return a list of system metrics",
		},
		{
			name:            "get_metrics_with_nil_params",
			params:          nil,
			setupMock:       func() {},
			expectedMetrics: []entities.SystemMetric(nil),
			expectedError:   ErrNilQueryParams,
			description:     "Should return a list of system metrics",
		},
		{
			name: "get_metrics_invalid_host_id",
			params: &entities.MetricQueryParams{
				HostID: func() *int64 { var id int64 = -5; return &id }(),
				Order:  "DESC",
				Limit:  10,
			},
			setupMock:       func() {},
			expectedMetrics: []entities.SystemMetric(nil),
			expectedError:   ErrInvalidHostID,
			description:     "Should return an error for invalid HostID",
		},
		{
			name: "get_metrics_empty_result",
			params: &entities.MetricQueryParams{
				HostID: &hostID,
				Order:  "DESC",
				Limit:  10,
			},
			setupMock: func() {
				suite.mockRepo.On("FindByFilters", &entities.MetricQueryParams{
					HostID: &hostID,
					Order:  "DESC",
					Limit:  10,
				}).Return([]entities.SystemMetric{}, nil).Once()
			},
			expectedMetrics: []entities.SystemMetric{},
			expectedError:   nil,
			description:     "Should return an empty list of system metrics",
		},
		{
			name: "database_error",
			params: &entities.MetricQueryParams{
				HostID: &hostID,
				Order:  "DESC",
				Limit:  100,
			},
			setupMock: func() {
				suite.mockRepo.On("FindByFilters", &entities.MetricQueryParams{
					HostID: &hostID,
					Order:  "DESC",
					Limit:  100,
				}).Return(nil, errors.New("query timeout")).Once()
			},
			expectedMetrics: nil,
			expectedError:   errors.New("query timeout"),
			description:     "Should return an error when there is a database query issue",
		},
	}

	for _, test := range tests {
		suite.Run(test.name, func() {
			test.setupMock()

			metrics, err := suite.service.GetMetrics(test.params)

			assert.Equal(suite.T(), test.expectedMetrics, metrics)
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

// TestGetLatestMetric tests the GetMetrics method
func (suite *MetricServiceTestSuite) TestGetLatestMetric() {
	hostID := int64(1)
	timestamp := time.Now().Unix()
	latestMetric := entities.SystemMetric{
		ID:                 1,
		HostID:             1,
		Timestamp:          timestamp - 100,
		CPUUsage:           45.5,
		MemoryUsagePercent: 67.8,
		DiskUsagePercent:   78.2,
	}

	tests := []struct {
		name           string
		hostID         *int64
		setupMock      func()
		expectedMetric *entities.SystemMetric
		expectedError  error
		description    string
	}{
		{
			name:   "get_latest_metric",
			hostID: nil,
			setupMock: func() {
				suite.mockRepo.On("FindLatest", (*int64)(nil)).Return(&latestMetric, nil).Once()
			},
			expectedMetric: &latestMetric,
			expectedError:  nil,
			description:    "Should return the latest system metric",
		},
		{
			name:   "get_latest_metric_with_host_id",
			hostID: &hostID,
			setupMock: func() {
				suite.mockRepo.On("FindLatest", &hostID).Return(&latestMetric, nil).Once()
			},
			expectedMetric: &latestMetric,
			expectedError:  nil,
			description:    "Should return the latest system metric for the given host",
		},
		{
			name:   "get_latest_metrics_empty_result",
			hostID: &hostID,
			setupMock: func() {
				suite.mockRepo.On("FindLatest", &hostID).Return(nil, sql.ErrNoRows).Once()
			},
			expectedMetric: nil,
			expectedError:  sql.ErrNoRows,
			description:    "Should return nil when no metrics are found for the given host",
		},
		{
			name:   "database_error",
			hostID: &hostID,
			setupMock: func() {
				suite.mockRepo.On("FindLatest", &hostID).Return(nil, errors.New("query timeout")).Once()
			},
			expectedMetric: nil,
			expectedError:  errors.New("query timeout"),
			description:    "Should return an error when there is a database query issue",
		},
	}

	for _, test := range tests {
		suite.Run(test.name, func() {
			test.setupMock()

			metric, err := suite.service.GetLatestMetric(test.hostID)

			assert.Equal(suite.T(), test.expectedMetric, metric)
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

// Run the test suite
func TestMetricServiceTestSuite(test *testing.T) {
	suite.Run(test, new(MetricServiceTestSuite))
}

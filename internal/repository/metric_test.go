// nolint
package repository

import (
	"database/sql"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gabrielg2020/monitor-api/internal/entities"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// MetricRepositoryTestSuite is the test suite for MetricRepository
type MetricRepositoryTestSuite struct {
	suite.Suite
	db   *sql.DB
	mock sqlmock.Sqlmock
	repo *MetricRepository
}

// SetupTest runs before each test in the suite
func (suite *MetricRepositoryTestSuite) SetupTest() {
	var err error
	suite.db, suite.mock, err = sqlmock.New(
		sqlmock.MonitorPingsOption(true),
		sqlmock.QueryMatcherOption(sqlmock.QueryMatcherRegexp),
	)
	suite.Require().NoError(err)

	suite.repo = NewMetricRepository(suite.db)
}

// TearDownTest runs after each test
func (suite *MetricRepositoryTestSuite) TearDownTest() {
	suite.db.Close()

	// Ensure all expectations were met
	err := suite.mock.ExpectationsWereMet()
	suite.NoError(err)
}

// TestNewMetricRepository tests the constructor
func (suite *MetricRepositoryTestSuite) TestNewMetricRepository() {
	assert.NotNil(suite.T(), suite.repo)
	assert.Equal(suite.T(), suite.db, suite.repo.db)
}

// TestFindByFilters tests the FindByFilters method
func (suite *MetricRepositoryTestSuite) TestFindByFilters() {
	hostID := int64(1)
	startTime := int64(1000)
	endTime := int64(2000)

	tests := []struct {
		name            string
		params          *entities.MetricQueryParams
		setupMock       func()
		expectedMetrics []entities.SystemMetric
		expectedError   error
	}{
		{
			name: "filter_by_host_id",
			params: &entities.MetricQueryParams{
				HostID: &hostID,
				Order:  "DESC",
				Limit:  10,
			},
			setupMock: func() {
				rows := sqlmock.NewRows([]string{
					"id", "host_id", "timestamp", "cpu_usage", "memory_usage_percent",
					"memory_total_bytes", "memory_used_bytes", "memory_available_bytes",
					"disk_usage_percent", "disk_total_bytes", "disk_used_bytes", "disk_available_bytes",
				}).
					AddRow(1, 1, 1500, 45.5, 60.0, 16000000000, 9600000000, 6400000000, 75.0, 500000000000, 375000000000, 125000000000)

				suite.mock.ExpectQuery("SELECT id, host_id, timestamp, cpu_usage, memory_usage_percent, memory_total_bytes, memory_used_bytes, memory_available_bytes, disk_usage_percent, disk_total_bytes, disk_used_bytes, disk_available_bytes FROM system_metrics WHERE 1=1 AND host_id = \\? ORDER BY timestamp DESC LIMIT \\?").
					WithArgs(int64(1), 10).
					WillReturnRows(rows)
			},
			expectedMetrics: []entities.SystemMetric{
				{
					ID:                   1,
					HostID:               1,
					Timestamp:            1500,
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
			},
			expectedError: nil,
		},
		{
			name: "filter_by_time_range",
			params: &entities.MetricQueryParams{
				StartTime: &startTime,
				EndTime:   &endTime,
				Order:     "ASC",
				Limit:     5,
			},
			setupMock: func() {
				rows := sqlmock.NewRows([]string{
					"id", "host_id", "timestamp", "cpu_usage", "memory_usage_percent",
					"memory_total_bytes", "memory_used_bytes", "memory_available_bytes",
					"disk_usage_percent", "disk_total_bytes", "disk_used_bytes", "disk_available_bytes",
				}).
					AddRow(2, 2, 1200, 30.0, 50.0, 8000000000, 4000000000, 4000000000, 60.0, 250000000000, 150000000000, 100000000000).
					AddRow(3, 2, 1800, 35.0, 55.0, 8000000000, 4400000000, 3600000000, 65.0, 250000000000, 162500000000, 87500000000)

				suite.mock.ExpectQuery("SELECT id, host_id, timestamp, cpu_usage, memory_usage_percent, memory_total_bytes, memory_used_bytes, memory_available_bytes, disk_usage_percent, disk_total_bytes, disk_used_bytes, disk_available_bytes FROM system_metrics WHERE 1=1 AND timestamp >= \\? AND timestamp <= \\? ORDER BY timestamp ASC LIMIT \\?").
					WithArgs(int64(1000), int64(2000), 5).
					WillReturnRows(rows)
			},
			expectedMetrics: []entities.SystemMetric{
				{
					ID:                   2,
					HostID:               2,
					Timestamp:            1200,
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
				{
					ID:                   3,
					HostID:               2,
					Timestamp:            1800,
					CPUUsage:             35.0,
					MemoryUsagePercent:   55.0,
					MemoryTotalBytes:     8000000000,
					MemoryUsedBytes:      4400000000,
					MemoryAvailableBytes: 3600000000,
					DiskUsagePercent:     65.0,
					DiskTotalBytes:       250000000000,
					DiskUsedBytes:        162500000000,
					DiskAvailableBytes:   87500000000,
				},
			},
			expectedError: nil,
		},
		{
			name: "filter_by_host_and_time_range",
			params: &entities.MetricQueryParams{
				HostID:    &hostID,
				StartTime: &startTime,
				EndTime:   &endTime,
				Order:     "DESC",
				Limit:     20,
			},
			setupMock: func() {
				rows := sqlmock.NewRows([]string{
					"id", "host_id", "timestamp", "cpu_usage", "memory_usage_percent",
					"memory_total_bytes", "memory_used_bytes", "memory_available_bytes",
					"disk_usage_percent", "disk_total_bytes", "disk_used_bytes", "disk_available_bytes",
				}).
					AddRow(4, 1, 1500, 40.0, 65.0, 16000000000, 10400000000, 5600000000, 70.0, 500000000000, 350000000000, 150000000000)

				suite.mock.ExpectQuery("SELECT id, host_id, timestamp, cpu_usage, memory_usage_percent, memory_total_bytes, memory_used_bytes, memory_available_bytes, disk_usage_percent, disk_total_bytes, disk_used_bytes, disk_available_bytes FROM system_metrics WHERE 1=1 AND host_id = \\? AND timestamp >= \\? AND timestamp <= \\? ORDER BY timestamp DESC LIMIT \\?").
					WithArgs(int64(1), int64(1000), int64(2000), 20).
					WillReturnRows(rows)
			},
			expectedMetrics: []entities.SystemMetric{
				{
					ID:                   4,
					HostID:               1,
					Timestamp:            1500,
					CPUUsage:             40.0,
					MemoryUsagePercent:   65.0,
					MemoryTotalBytes:     16000000000,
					MemoryUsedBytes:      10400000000,
					MemoryAvailableBytes: 5600000000,
					DiskUsagePercent:     70.0,
					DiskTotalBytes:       500000000000,
					DiskUsedBytes:        350000000000,
					DiskAvailableBytes:   150000000000,
				},
			},
			expectedError: nil,
		},
		{
			name: "no_filters",
			params: &entities.MetricQueryParams{
				Order: "DESC",
				Limit: 100,
			},
			setupMock: func() {
				rows := sqlmock.NewRows([]string{
					"id", "host_id", "timestamp", "cpu_usage", "memory_usage_percent",
					"memory_total_bytes", "memory_used_bytes", "memory_available_bytes",
					"disk_usage_percent", "disk_total_bytes", "disk_used_bytes", "disk_available_bytes",
				}).
					AddRow(5, 3, 3000, 50.0, 70.0, 32000000000, 22400000000, 9600000000, 80.0, 1000000000000, 800000000000, 200000000000)

				suite.mock.ExpectQuery("SELECT id, host_id, timestamp, cpu_usage, memory_usage_percent, memory_total_bytes, memory_used_bytes, memory_available_bytes, disk_usage_percent, disk_total_bytes, disk_used_bytes, disk_available_bytes FROM system_metrics WHERE 1=1 ORDER BY timestamp DESC LIMIT \\?").
					WithArgs(100).
					WillReturnRows(rows)
			},
			expectedMetrics: []entities.SystemMetric{
				{
					ID:                   5,
					HostID:               3,
					Timestamp:            3000,
					CPUUsage:             50.0,
					MemoryUsagePercent:   70.0,
					MemoryTotalBytes:     32000000000,
					MemoryUsedBytes:      22400000000,
					MemoryAvailableBytes: 9600000000,
					DiskUsagePercent:     80.0,
					DiskTotalBytes:       1000000000000,
					DiskUsedBytes:        800000000000,
					DiskAvailableBytes:   200000000000,
				},
			},
			expectedError: nil,
		},
		{
			name: "empty_result_set",
			params: &entities.MetricQueryParams{
				HostID: &hostID,
				Order:  "DESC",
				Limit:  10,
			},
			setupMock: func() {
				rows := sqlmock.NewRows([]string{
					"id", "host_id", "timestamp", "cpu_usage", "memory_usage_percent",
					"memory_total_bytes", "memory_used_bytes", "memory_available_bytes",
					"disk_usage_percent", "disk_total_bytes", "disk_used_bytes", "disk_available_bytes",
				})

				suite.mock.ExpectQuery("SELECT id, host_id, timestamp, cpu_usage, memory_usage_percent, memory_total_bytes, memory_used_bytes, memory_available_bytes, disk_usage_percent, disk_total_bytes, disk_used_bytes, disk_available_bytes FROM system_metrics WHERE 1=1 AND host_id = \\? ORDER BY timestamp DESC LIMIT \\?").
					WithArgs(int64(1), 10).
					WillReturnRows(rows)
			},
			expectedMetrics: []entities.SystemMetric(nil),
			expectedError:   nil,
		},
		{
			name: "database_error",
			params: &entities.MetricQueryParams{
				HostID: &hostID,
				Order:  "DESC",
				Limit:  10,
			},
			setupMock: func() {
				suite.mock.ExpectQuery("SELECT id, host_id, timestamp, cpu_usage, memory_usage_percent, memory_total_bytes, memory_used_bytes, memory_available_bytes, disk_usage_percent, disk_total_bytes, disk_used_bytes, disk_available_bytes FROM system_metrics WHERE 1=1 AND host_id = \\? ORDER BY timestamp DESC LIMIT \\?").
					WithArgs(int64(1), 10).
					WillReturnError(errors.New("connection timeout"))
			},
			expectedMetrics: nil,
			expectedError:   errors.New("connection timeout"),
		},
		{
			name: "scan_error",
			params: &entities.MetricQueryParams{
				Order: "DESC",
				Limit: 10,
			},
			setupMock: func() {
				rows := sqlmock.NewRows([]string{
					"id", "host_id", "timestamp", "cpu_usage", "memory_usage_percent",
					"memory_total_bytes", "memory_used_bytes", "memory_available_bytes",
					"disk_usage_percent", "disk_total_bytes", "disk_used_bytes", "disk_available_bytes",
				}).
					AddRow("invalid", 1, 1500, 45.5, 60.0, 16000000000, 9600000000, 6400000000, 75.0, 500000000000, 375000000000, 125000000000)

				suite.mock.ExpectQuery("SELECT id, host_id, timestamp, cpu_usage, memory_usage_percent, memory_total_bytes, memory_used_bytes, memory_available_bytes, disk_usage_percent, disk_total_bytes, disk_used_bytes, disk_available_bytes FROM system_metrics WHERE 1=1 ORDER BY timestamp DESC LIMIT \\?").
					WithArgs(10).
					WillReturnRows(rows)
			},
			expectedMetrics: nil,
			expectedError:   errors.New("sql: Scan error"),
		},
	}

	for _, test := range tests {
		suite.Run(test.name, func() {
			test.setupMock()

			metrics, err := suite.repo.FindByFilters(test.params)

			if test.expectedError != nil {
				assert.Error(suite.T(), err)
				if test.name == "scan_error" {
					assert.Contains(suite.T(), err.Error(), "Scan error")
				} else {
					assert.Equal(suite.T(), test.expectedError.Error(), err.Error())
				}
				assert.Nil(suite.T(), metrics)
			} else {
				assert.NoError(suite.T(), err)
				assert.Equal(suite.T(), test.expectedMetrics, metrics)
			}
		})

		// Reset for next test
		suite.TearDownTest()
		suite.SetupTest()
	}
}

// TestFindLatest tests the FindLatest method
func (suite *MetricRepositoryTestSuite) TestFindLatest() {
	hostID := int64(1)

	tests := []struct {
		name           string
		hostID         *int64
		setupMock      func()
		expectedMetric *entities.SystemMetric
		expectedError  error
	}{
		{
			name:   "find_latest_with_host_id",
			hostID: &hostID,
			setupMock: func() {
				rows := sqlmock.NewRows([]string{
					"id", "host_id", "timestamp", "cpu_usage", "memory_usage_percent",
					"memory_total_bytes", "memory_used_bytes", "memory_available_bytes",
					"disk_usage_percent", "disk_total_bytes", "disk_used_bytes", "disk_available_bytes",
				}).
					AddRow(1, 1, 2000, 55.5, 70.0, 16000000000, 11200000000, 4800000000, 85.0, 500000000000, 425000000000, 75000000000)

				suite.mock.ExpectQuery("SELECT id, host_id, timestamp, cpu_usage, memory_usage_percent, memory_total_bytes, memory_used_bytes, memory_available_bytes, disk_usage_percent, disk_total_bytes, disk_used_bytes, disk_available_bytes FROM system_metrics WHERE host_id = \\? ORDER BY timestamp DESC LIMIT 1").
					WithArgs(int64(1)).
					WillReturnRows(rows)
			},
			expectedMetric: &entities.SystemMetric{
				ID:                   1,
				HostID:               1,
				Timestamp:            2000,
				CPUUsage:             55.5,
				MemoryUsagePercent:   70.0,
				MemoryTotalBytes:     16000000000,
				MemoryUsedBytes:      11200000000,
				MemoryAvailableBytes: 4800000000,
				DiskUsagePercent:     85.0,
				DiskTotalBytes:       500000000000,
				DiskUsedBytes:        425000000000,
				DiskAvailableBytes:   75000000000,
			},
			expectedError: nil,
		},
		{
			name:   "find_latest_without_host_id",
			hostID: nil,
			setupMock: func() {
				rows := sqlmock.NewRows([]string{
					"id", "host_id", "timestamp", "cpu_usage", "memory_usage_percent",
					"memory_total_bytes", "memory_used_bytes", "memory_available_bytes",
					"disk_usage_percent", "disk_total_bytes", "disk_used_bytes", "disk_available_bytes",
				}).
					AddRow(2, 3, 3000, 45.0, 65.0, 32000000000, 20800000000, 11200000000, 75.0, 1000000000000, 750000000000, 250000000000)

				suite.mock.ExpectQuery("SELECT id, host_id, timestamp, cpu_usage, memory_usage_percent, memory_total_bytes, memory_used_bytes, memory_available_bytes, disk_usage_percent, disk_total_bytes, disk_used_bytes, disk_available_bytes FROM system_metrics ORDER BY timestamp DESC LIMIT 1").
					WillReturnRows(rows)
			},
			expectedMetric: &entities.SystemMetric{
				ID:                   2,
				HostID:               3,
				Timestamp:            3000,
				CPUUsage:             45.0,
				MemoryUsagePercent:   65.0,
				MemoryTotalBytes:     32000000000,
				MemoryUsedBytes:      20800000000,
				MemoryAvailableBytes: 11200000000,
				DiskUsagePercent:     75.0,
				DiskTotalBytes:       1000000000000,
				DiskUsedBytes:        750000000000,
				DiskAvailableBytes:   250000000000,
			},
			expectedError: nil,
		},
		{
			name:   "no_metrics_found",
			hostID: &hostID,
			setupMock: func() {
				suite.mock.ExpectQuery("SELECT id, host_id, timestamp, cpu_usage, memory_usage_percent, memory_total_bytes, memory_used_bytes, memory_available_bytes, disk_usage_percent, disk_total_bytes, disk_used_bytes, disk_available_bytes FROM system_metrics WHERE host_id = \\? ORDER BY timestamp DESC LIMIT 1").
					WithArgs(int64(1)).
					WillReturnError(sql.ErrNoRows)
			},
			expectedMetric: nil,
			expectedError:  nil, // Returns nil, nil when not found
		},
		{
			name:   "database_error",
			hostID: &hostID,
			setupMock: func() {
				suite.mock.ExpectQuery("SELECT id, host_id, timestamp, cpu_usage, memory_usage_percent, memory_total_bytes, memory_used_bytes, memory_available_bytes, disk_usage_percent, disk_total_bytes, disk_used_bytes, disk_available_bytes FROM system_metrics WHERE host_id = \\? ORDER BY timestamp DESC LIMIT 1").
					WithArgs(int64(1)).
					WillReturnError(errors.New("database connection lost"))
			},
			expectedMetric: nil,
			expectedError:  errors.New("database connection lost"),
		},
		{
			name:   "scan_error",
			hostID: &hostID,
			setupMock: func() {
				rows := sqlmock.NewRows([]string{
					"id", "host_id", "timestamp", "cpu_usage", "memory_usage_percent",
					"memory_total_bytes", "memory_used_bytes", "memory_available_bytes",
					"disk_usage_percent", "disk_total_bytes", "disk_used_bytes", "disk_available_bytes",
				}).
					AddRow("invalid", 1, 2000, 55.5, 70.0, 16000000000, 11200000000, 4800000000, 85.0, 500000000000, 425000000000, 75000000000)

				suite.mock.ExpectQuery("SELECT id, host_id, timestamp, cpu_usage, memory_usage_percent, memory_total_bytes, memory_used_bytes, memory_available_bytes, disk_usage_percent, disk_total_bytes, disk_used_bytes, disk_available_bytes FROM system_metrics WHERE host_id = \\? ORDER BY timestamp DESC LIMIT 1").
					WithArgs(int64(1)).
					WillReturnRows(rows)
			},
			expectedMetric: nil,
			expectedError:  errors.New("sql: Scan error"),
		},
	}

	for _, test := range tests {
		suite.Run(test.name, func() {
			test.setupMock()

			metric, err := suite.repo.FindLatest(test.hostID)

			if test.expectedError != nil {
				assert.Error(suite.T(), err)
				if test.name == "scan_error" {
					assert.Contains(suite.T(), err.Error(), "Scan error")
				} else {
					assert.Equal(suite.T(), test.expectedError.Error(), err.Error())
				}
				assert.Nil(suite.T(), metric)
			} else {
				assert.NoError(suite.T(), err)
				if test.expectedMetric == nil {
					assert.Nil(suite.T(), metric)
				} else {
					assert.Equal(suite.T(), test.expectedMetric, metric)
				}
			}
		})

		// Reset for next test
		suite.TearDownTest()
		suite.SetupTest()
	}
}

// TestCreate tests the Create method
func (suite *MetricRepositoryTestSuite) TestCreate() {
	tests := []struct {
		name          string
		metric        *entities.SystemMetric
		setupMock     func()
		expectedID    int64
		expectedError error
	}{
		{
			name: "successful_creation",
			metric: &entities.SystemMetric{
				HostID:               1,
				Timestamp:            1500,
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
			setupMock: func() {
				suite.mock.ExpectExec("INSERT INTO system_metrics \\( host_id, timestamp, cpu_usage, memory_usage_percent, memory_total_bytes, memory_used_bytes, memory_available_bytes, disk_usage_percent, disk_total_bytes, disk_used_bytes, disk_available_bytes \\) VALUES \\(\\?, \\?, \\?, \\?, \\?, \\?, \\?, \\?, \\?, \\?, \\?\\)").
					WithArgs(int64(1), int64(1500), 45.5, 60.0, int64(16000000000), int64(9600000000), int64(6400000000), 75.0, int64(500000000000), int64(375000000000), int64(125000000000)).
					WillReturnResult(sqlmock.NewResult(10, 1))
			},
			expectedID:    10,
			expectedError: nil,
		},
		{
			name: "foreign_key_constraint_error",
			metric: &entities.SystemMetric{
				HostID:               999,
				Timestamp:            1500,
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
			setupMock: func() {
				suite.mock.ExpectExec("INSERT INTO system_metrics \\( host_id, timestamp, cpu_usage, memory_usage_percent, memory_total_bytes, memory_used_bytes, memory_available_bytes, disk_usage_percent, disk_total_bytes, disk_used_bytes, disk_available_bytes \\) VALUES \\(\\?, \\?, \\?, \\?, \\?, \\?, \\?, \\?, \\?, \\?, \\?\\)").
					WithArgs(int64(999), int64(1500), 45.5, 60.0, int64(16000000000), int64(9600000000), int64(6400000000), 75.0, int64(500000000000), int64(375000000000), int64(125000000000)).
					WillReturnError(errors.New("FOREIGN KEY constraint failed"))
			},
			expectedID:    -1,
			expectedError: errors.New("FOREIGN KEY constraint failed"),
		},
		{
			name: "database_connection_error",
			metric: &entities.SystemMetric{
				HostID:               1,
				Timestamp:            1500,
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
			setupMock: func() {
				suite.mock.ExpectExec("INSERT INTO system_metrics \\( host_id, timestamp, cpu_usage, memory_usage_percent, memory_total_bytes, memory_used_bytes, memory_available_bytes, disk_usage_percent, disk_total_bytes, disk_used_bytes, disk_available_bytes \\) VALUES \\(\\?, \\?, \\?, \\?, \\?, \\?, \\?, \\?, \\?, \\?, \\?\\)").
					WithArgs(int64(1), int64(1500), 45.5, 60.0, int64(16000000000), int64(9600000000), int64(6400000000), 75.0, int64(500000000000), int64(375000000000), int64(125000000000)).
					WillReturnError(errors.New("database connection lost"))
			},
			expectedID:    -1,
			expectedError: errors.New("database connection lost"),
		},
	}

	for _, test := range tests {
		suite.Run(test.name, func() {
			test.setupMock()

			id, err := suite.repo.Create(test.metric)

			if test.expectedError != nil {
				assert.Error(suite.T(), err)
				assert.Equal(suite.T(), test.expectedError.Error(), err.Error())
				assert.Equal(suite.T(), int64(-1), id)
			} else {
				assert.NoError(suite.T(), err)
				assert.Equal(suite.T(), test.expectedID, id)
			}
		})

		// Reset for next test
		suite.TearDownTest()
		suite.SetupTest()
	}
}

// TestScanMetricsErrorHandling tests error handling in scanMetrics helper
func (suite *MetricRepositoryTestSuite) TestScanMetricsErrorHandling() {
	// Test rows.Err() handling
	suite.mock.ExpectQuery("SELECT id, host_id, timestamp, cpu_usage, memory_usage_percent, memory_total_bytes, memory_used_bytes, memory_available_bytes, disk_usage_percent, disk_total_bytes, disk_used_bytes, disk_available_bytes FROM system_metrics WHERE 1=1 ORDER BY timestamp DESC LIMIT \\?").
		WithArgs(10).
		WillReturnRows(sqlmock.NewRows([]string{
			"id", "host_id", "timestamp", "cpu_usage", "memory_usage_percent",
			"memory_total_bytes", "memory_used_bytes", "memory_available_bytes",
			"disk_usage_percent", "disk_total_bytes", "disk_used_bytes", "disk_available_bytes",
		}).
			AddRow(1, 1, 1500, 45.5, 60.0, 16000000000, 9600000000, 6400000000, 75.0, 500000000000, 375000000000, 125000000000).
			RowError(0, errors.New("row iteration error")))

	params := &entities.MetricQueryParams{
		Order: "DESC",
		Limit: 10,
	}

	metrics, err := suite.repo.FindByFilters(params)

	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), "row iteration error", err.Error())
	assert.Nil(suite.T(), metrics)
}

// Run the test suite
func TestMetricRepositoryTestSuite(test *testing.T) {
	suite.Run(test, new(MetricRepositoryTestSuite))
}

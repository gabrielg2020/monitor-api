package repository

import (
	"database/sql"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// HealthRepositoryTestSuite is the test suite for HealthRepository
type HealthRepositoryTestSuite struct {
	suite.Suite
	db   *sql.DB
	mock sqlmock.Sqlmock
	repo *HealthRepository
}

// SetupTest runs before each test in the suite
func (suite *HealthRepositoryTestSuite) SetupTest() {
	var err error
	suite.db, suite.mock, err = sqlmock.New(
		sqlmock.MonitorPingsOption(true),
		sqlmock.QueryMatcherOption(sqlmock.QueryMatcherRegexp),
	)
	suite.Require().NoError(err)

	suite.repo = NewHealthRepository(suite.db)
}

// TearDownTest runs after each test
func (suite *HealthRepositoryTestSuite) TearDownTest() {
	suite.db.Close()

	// Ensure all expectations were met
	err := suite.mock.ExpectationsWereMet()
	suite.NoError(err)
}

// TestNewHealthRepository tests the constructor
func (suite *HealthRepositoryTestSuite) TestNewHealthRepository() {
	assert.NotNil(suite.T(), suite.repo)
	assert.Equal(suite.T(), suite.db, suite.repo.db)
}

// TestCheckDatabaseConnection tests the CheckDatabaseConnection method
func (suite *HealthRepositoryTestSuite) TestCheckDatabaseConnection() {
	tests := []struct {
		name          string
		setupMock     func()
		expectedError error
	}{
		{
			name: "successful_connection",
			setupMock: func() {
				suite.mock.ExpectPing()
			},
			expectedError: nil,
		},
		{
			name: "connection_failed",
			setupMock: func() {
				suite.mock.ExpectPing().WillReturnError(errors.New("connection refused"))
			},
			expectedError: errors.New("connection refused"),
		},
		{
			name: "timeout_error",
			setupMock: func() {
				suite.mock.ExpectPing().WillReturnError(errors.New("i/o timeout"))
			},
			expectedError: errors.New("i/o timeout"),
		},
		{
			name: "database_not_found",
			setupMock: func() {
				suite.mock.ExpectPing().WillReturnError(errors.New("database does not exist"))
			},
			expectedError: errors.New("database does not exist"),
		},
	}

	for _, test := range tests {
		suite.Run(test.name, func() {
			test.setupMock()

			err := suite.repo.CheckDatabaseConnection()

			if test.expectedError != nil {
				assert.Error(suite.T(), err)
				assert.Equal(suite.T(), test.expectedError.Error(), err.Error())
			} else {
				assert.NoError(suite.T(), err)
			}
		})

		// Reset for next test
		suite.TearDownTest()
		suite.SetupTest()
	}
}

// TestGetDatabaseStats tests the GetDatabaseStats method
func (suite *HealthRepositoryTestSuite) TestGetDatabaseStats() {
	// Note: sql.DB.Stats() doesn't actually query the database,
	// it returns the internal connection pool statistics.
	// So no need to set up specific expectations here.

	stats, err := suite.repo.GetDatabaseStats()

	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), stats)

	// Check that all expected keys are present
	expectedKeys := []string{"open_connections", "in_use", "idle", "max_open"}
	for _, key := range expectedKeys {
		_, exists := stats[key]
		assert.True(suite.T(), exists, "Expected key %s to exist in stats", key)
	}

	// Verify the values are of correct type (int)
	assert.IsType(suite.T(), 0, stats["open_connections"])
	assert.IsType(suite.T(), 0, stats["in_use"])
	assert.IsType(suite.T(), 0, stats["idle"])
	assert.IsType(suite.T(), 0, stats["max_open"])

	// Values should be non-negative
	assert.GreaterOrEqual(suite.T(), stats["open_connections"].(int), 0)
	assert.GreaterOrEqual(suite.T(), stats["in_use"].(int), 0)
	assert.GreaterOrEqual(suite.T(), stats["idle"].(int), 0)
	assert.GreaterOrEqual(suite.T(), stats["max_open"].(int), 0)
}

// TestGetTableCounts tests the GetTableCounts method
func (suite *HealthRepositoryTestSuite) TestGetTableCounts() {
	tests := []struct {
		name           string
		setupMock      func()
		expectedCounts map[string]int
		expectedError  error
	}{
		{
			name: "successful_counts",
			setupMock: func() {
				// Mock host count query
				hostRows := sqlmock.NewRows([]string{"count"}).AddRow(10)
				suite.mock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM hosts").
					WillReturnRows(hostRows)

				// Mock metric count query
				metricRows := sqlmock.NewRows([]string{"count"}).AddRow(1000)
				suite.mock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM system_metrics").
					WillReturnRows(metricRows)
			},
			expectedCounts: map[string]int{
				"hosts":   10,
				"metrics": 1000,
			},
			expectedError: nil,
		},
		{
			name: "zero_counts",
			setupMock: func() {
				// Mock host count query with zero
				hostRows := sqlmock.NewRows([]string{"count"}).AddRow(0)
				suite.mock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM hosts").
					WillReturnRows(hostRows)

				// Mock metric count query with zero
				metricRows := sqlmock.NewRows([]string{"count"}).AddRow(0)
				suite.mock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM system_metrics").
					WillReturnRows(metricRows)
			},
			expectedCounts: map[string]int{
				"hosts":   0,
				"metrics": 0,
			},
			expectedError: nil,
		},
		{
			name: "large_counts",
			setupMock: func() {
				// Mock host count query with large number
				hostRows := sqlmock.NewRows([]string{"count"}).AddRow(100000)
				suite.mock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM hosts").
					WillReturnRows(hostRows)

				// Mock metric count query with large number
				metricRows := sqlmock.NewRows([]string{"count"}).AddRow(5000000)
				suite.mock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM system_metrics").
					WillReturnRows(metricRows)
			},
			expectedCounts: map[string]int{
				"hosts":   100000,
				"metrics": 5000000,
			},
			expectedError: nil,
		},
		{
			name: "host_table_error",
			setupMock: func() {
				// Mock host count query with error
				suite.mock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM hosts").
					WillReturnError(errors.New("table 'hosts' doesn't exist"))
			},
			expectedCounts: nil,
			expectedError:  errors.New("table 'hosts' doesn't exist"),
		},
		{
			name: "metric_table_error",
			setupMock: func() {
				// Mock host count query success
				hostRows := sqlmock.NewRows([]string{"count"}).AddRow(10)
				suite.mock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM hosts").
					WillReturnRows(hostRows)

				// Mock metric count query with error
				suite.mock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM system_metrics").
					WillReturnError(errors.New("table 'system_metrics' doesn't exist"))
			},
			expectedCounts: nil,
			expectedError:  errors.New("table 'system_metrics' doesn't exist"),
		},
		{
			name: "database_connection_error_on_hosts",
			setupMock: func() {
				suite.mock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM hosts").
					WillReturnError(errors.New("database connection lost"))
			},
			expectedCounts: nil,
			expectedError:  errors.New("database connection lost"),
		},
		{
			name: "scan_error_on_hosts",
			setupMock: func() {
				// Return rows with wrong type/structure
				hostRows := sqlmock.NewRows([]string{"invalid_column"}).AddRow("not_a_number")
				suite.mock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM hosts").
					WillReturnRows(hostRows)
			},
			expectedCounts: nil,
			expectedError:  errors.New("sql: Scan error on column index 0, name \"invalid_column\": converting driver.Value type string (\"not_a_number\") to a int: invalid syntax"),
		},
	}

	for _, test := range tests {
		suite.Run(test.name, func() {
			test.setupMock()

			counts, err := suite.repo.GetTableCounts()

			if test.expectedError != nil {
				assert.Error(suite.T(), err)
				// For scan errors, just check that an error occurred
				if test.name == "scan_error_on_hosts" {
					assert.Contains(suite.T(), err.Error(), "Scan error")
				} else {
					assert.Equal(suite.T(), test.expectedError.Error(), err.Error())
				}
				assert.Nil(suite.T(), counts)
			} else {
				assert.NoError(suite.T(), err)
				assert.Equal(suite.T(), test.expectedCounts, counts)
			}
		})

		// Reset for next test
		suite.TearDownTest()
		suite.SetupTest()
	}
}

// TestGetTableCountsPartialSuccess tests partial success scenarios
func (suite *HealthRepositoryTestSuite) TestGetTableCountsPartialSuccess() {
	// This tests that if the first query succeeds but the second fails,
	// it should still return an error and nil counts

	// Mock successful host count
	hostRows := sqlmock.NewRows([]string{"count"}).AddRow(50)
	suite.mock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM hosts").
		WillReturnRows(hostRows)

	// Mock failed metric count
	suite.mock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM system_metrics").
		WillReturnError(errors.New("permission denied"))

	counts, err := suite.repo.GetTableCounts()

	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), "permission denied", err.Error())
	assert.Nil(suite.T(), counts)
}

// TestGetTableCountsWithNullValues tests handling of NULL values
func (suite *HealthRepositoryTestSuite) TestGetTableCountsWithNullValues() {
	// Test that NULL values are handled correctly (though COUNT(*) shouldn't return NULL)

	// Mock host count query with NULL (edge case)
	hostRows := sqlmock.NewRows([]string{"count"}).AddRow(nil)
	suite.mock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM hosts").
		WillReturnRows(hostRows)

	counts, err := suite.repo.GetTableCounts()

	// Should error when trying to scan NULL into int
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), counts)
}

// Run the test suite
func TestHealthRepositoryTestSuite(test *testing.T) {
	suite.Run(test, new(HealthRepositoryTestSuite))
}

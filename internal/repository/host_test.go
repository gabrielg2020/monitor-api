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

// HostRepositoryTestSuite is the test suite for HostRepository
type HostRepositoryTestSuite struct {
	suite.Suite
	db   *sql.DB
	mock sqlmock.Sqlmock
	repo *HostRepository
}

// SetupTest runs before each test in the suite
func (suite *HostRepositoryTestSuite) SetupTest() {
	var err error
	suite.db, suite.mock, err = sqlmock.New(
		sqlmock.MonitorPingsOption(true),
		sqlmock.QueryMatcherOption(sqlmock.QueryMatcherRegexp),
	)
	suite.Require().NoError(err)

	suite.repo = NewHostRepository(suite.db)
}

// TearDownTest runs after each test
func (suite *HostRepositoryTestSuite) TearDownTest() {
	suite.db.Close()

	// Ensure all expectations were met
	err := suite.mock.ExpectationsWereMet()
	suite.NoError(err)
}

// TestNewHostRepository tests the constructor
func (suite *HostRepositoryTestSuite) TestNewHostRepository() {
	assert.NotNil(suite.T(), suite.repo)
	assert.Equal(suite.T(), suite.db, suite.repo.db)
}

// TestFindByFilters tests the FindByFilters method
func (suite *HostRepositoryTestSuite) TestFindByFilters() {
	tests := []struct {
		name          string
		params        *entities.HostQueryParams
		setupMock     func()
		expectedHosts []entities.Host
		expectedError error
	}{
		{
			name: "filter_by_id",
			params: &entities.HostQueryParams{
				ID: 1,
			},
			setupMock: func() {
				rows := sqlmock.NewRows([]string{"id", "hostname", "ip_address", "role"}).
					AddRow(1, "server-01.example.com", "192.168.1.100", "web-server")

				suite.mock.ExpectQuery("SELECT id, hostname, ip_address, role FROM hosts WHERE 1=1 AND id = \\?").
					WithArgs(int64(1)).
					WillReturnRows(rows)
			},
			expectedHosts: []entities.Host{
				{ID: 1, Hostname: "server-01.example.com", IPAddress: "192.168.1.100", Role: "web-server"},
			},
			expectedError: nil,
		},
		{
			name: "filter_by_hostname",
			params: &entities.HostQueryParams{
				Hostname: "server-02.example.com",
			},
			setupMock: func() {
				rows := sqlmock.NewRows([]string{"id", "hostname", "ip_address", "role"}).
					AddRow(2, "server-02.example.com", "192.168.1.101", "database")

				suite.mock.ExpectQuery("SELECT id, hostname, ip_address, role FROM hosts WHERE 1=1 AND hostname = \\?").
					WithArgs("server-02.example.com").
					WillReturnRows(rows)
			},
			expectedHosts: []entities.Host{
				{ID: 2, Hostname: "server-02.example.com", IPAddress: "192.168.1.101", Role: "database"},
			},
			expectedError: nil,
		},
		{
			name: "filter_by_ip_address",
			params: &entities.HostQueryParams{
				IPAddress: "192.168.1.102",
			},
			setupMock: func() {
				rows := sqlmock.NewRows([]string{"id", "hostname", "ip_address", "role"}).
					AddRow(3, "server-03.example.com", "192.168.1.102", "cache")

				suite.mock.ExpectQuery("SELECT id, hostname, ip_address, role FROM hosts WHERE 1=1 AND ip_address = \\?").
					WithArgs("192.168.1.102").
					WillReturnRows(rows)
			},
			expectedHosts: []entities.Host{
				{ID: 3, Hostname: "server-03.example.com", IPAddress: "192.168.1.102", Role: "cache"},
			},
			expectedError: nil,
		},
		{
			name: "filter_by_multiple_fields",
			params: &entities.HostQueryParams{
				ID:        1,
				Hostname:  "server-01.example.com",
				IPAddress: "192.168.1.100",
			},
			setupMock: func() {
				rows := sqlmock.NewRows([]string{"id", "hostname", "ip_address", "role"}).
					AddRow(1, "server-01.example.com", "192.168.1.100", "web-server")

				suite.mock.ExpectQuery("SELECT id, hostname, ip_address, role FROM hosts WHERE 1=1 AND id = \\? AND hostname = \\? AND ip_address = \\?").
					WithArgs(int64(1), "server-01.example.com", "192.168.1.100").
					WillReturnRows(rows)
			},
			expectedHosts: []entities.Host{
				{ID: 1, Hostname: "server-01.example.com", IPAddress: "192.168.1.100", Role: "web-server"},
			},
			expectedError: nil,
		},
		{
			name:   "no_filters",
			params: &entities.HostQueryParams{},
			setupMock: func() {
				rows := sqlmock.NewRows([]string{"id", "hostname", "ip_address", "role"}).
					AddRow(1, "server-01.example.com", "192.168.1.100", "web-server").
					AddRow(2, "server-02.example.com", "192.168.1.101", "database")

				suite.mock.ExpectQuery("SELECT id, hostname, ip_address, role FROM hosts WHERE 1=1").
					WillReturnRows(rows)
			},
			expectedHosts: []entities.Host{
				{ID: 1, Hostname: "server-01.example.com", IPAddress: "192.168.1.100", Role: "web-server"},
				{ID: 2, Hostname: "server-02.example.com", IPAddress: "192.168.1.101", Role: "database"},
			},
			expectedError: nil,
		},
		{
			name: "no_results_found",
			params: &entities.HostQueryParams{
				Hostname: "non-existent.example.com",
			},
			setupMock: func() {
				rows := sqlmock.NewRows([]string{"id", "hostname", "ip_address", "role"})

				suite.mock.ExpectQuery("SELECT id, hostname, ip_address, role FROM hosts WHERE 1=1 AND hostname = \\?").
					WithArgs("non-existent.example.com").
					WillReturnRows(rows)
			},
			expectedHosts: []entities.Host(nil),
			expectedError: nil,
		},
		{
			name: "database_error",
			params: &entities.HostQueryParams{
				ID: 1,
			},
			setupMock: func() {
				suite.mock.ExpectQuery("SELECT id, hostname, ip_address, role FROM hosts WHERE 1=1 AND id = \\?").
					WithArgs(int64(1)).
					WillReturnError(errors.New("query execution failed"))
			},
			expectedHosts: nil,
			expectedError: errors.New("query execution failed"),
		},
		{
			name:   "scan_error",
			params: &entities.HostQueryParams{},
			setupMock: func() {
				rows := sqlmock.NewRows([]string{"id", "hostname", "ip_address", "role"}).
					AddRow("invalid", "server-01.example.com", "192.168.1.100", "web-server")

				suite.mock.ExpectQuery("SELECT id, hostname, ip_address, role FROM hosts WHERE 1=1").
					WillReturnRows(rows)
			},
			expectedHosts: nil,
			expectedError: errors.New("sql: Scan error on column index 0, name \"id\": converting driver.Value type string (\"invalid\") to a int64: invalid syntax"),
		},
	}

	for _, test := range tests {
		suite.Run(test.name, func() {
			test.setupMock()

			hosts, err := suite.repo.FindByFilters(test.params)

			if test.expectedError != nil {
				assert.Error(suite.T(), err)
				if test.name == "scan_error" {
					assert.Contains(suite.T(), err.Error(), "Scan error")
				} else {
					assert.Equal(suite.T(), test.expectedError.Error(), err.Error())
				}
				assert.Nil(suite.T(), hosts)
			} else {
				assert.NoError(suite.T(), err)
				assert.Equal(suite.T(), test.expectedHosts, hosts)
			}
		})

		// Reset for next test
		suite.TearDownTest()
		suite.SetupTest()
	}
}

// TestCreate tests the Create method
func (suite *HostRepositoryTestSuite) TestCreate() {
	tests := []struct {
		name          string
		host          *entities.Host
		setupMock     func()
		expectedID    int64
		expectedError error
	}{
		{
			name: "successful_creation",
			host: &entities.Host{
				Hostname:  "new-server.example.com",
				IPAddress: "192.168.1.200",
				Role:      "application",
			},
			setupMock: func() {
				suite.mock.ExpectExec("INSERT INTO hosts \\(hostname, ip_address, role, created_at, last_seen\\) VALUES \\(\\?, \\?, \\?, \\?, \\?\\)").
					WithArgs("new-server.example.com", "192.168.1.200", "application",
						sqlmock.AnyArg(), sqlmock.AnyArg()).
					WillReturnResult(sqlmock.NewResult(10, 1))
			},
			expectedID:    10,
			expectedError: nil,
		},
		{
			name: "duplicate_hostname_error",
			host: &entities.Host{
				Hostname:  "existing-server.example.com",
				IPAddress: "192.168.1.201",
				Role:      "database",
			},
			setupMock: func() {
				suite.mock.ExpectExec("INSERT INTO hosts \\(hostname, ip_address, role, created_at, last_seen\\) VALUES \\(\\?, \\?, \\?, \\?, \\?\\)").
					WithArgs("existing-server.example.com", "192.168.1.201", "database",
						sqlmock.AnyArg(), sqlmock.AnyArg()).
					WillReturnError(errors.New("UNIQUE constraint failed: hosts.hostname"))
			},
			expectedID:    0,
			expectedError: errors.New("UNIQUE constraint failed: hosts.hostname"),
		},
		{
			name: "duplicate_ip_address_error",
			host: &entities.Host{
				Hostname:  "new-server.example.com",
				IPAddress: "existing_ip_address",
				Role:      "database",
			},
			setupMock: func() {
				suite.mock.ExpectExec("INSERT INTO hosts \\(hostname, ip_address, role, created_at, last_seen\\) VALUES \\(\\?, \\?, \\?, \\?, \\?\\)").
					WithArgs("new-server.example.com", "existing_ip_address", "database",
						sqlmock.AnyArg(), sqlmock.AnyArg()).
					WillReturnError(errors.New("UNIQUE constraint failed: hosts.ip_address"))
			},
			expectedID:    0,
			expectedError: errors.New("UNIQUE constraint failed: hosts.ip_address"),
		},
		{
			name: "database_connection_error",
			host: &entities.Host{
				Hostname:  "server.example.com",
				IPAddress: "192.168.1.202",
				Role:      "cache",
			},
			setupMock: func() {
				suite.mock.ExpectExec("INSERT INTO hosts \\(hostname, ip_address, role, created_at, last_seen\\) VALUES \\(\\?, \\?, \\?, \\?, \\?\\)").
					WithArgs("server.example.com", "192.168.1.202", "cache",
						sqlmock.AnyArg(), sqlmock.AnyArg()).
					WillReturnError(errors.New("database connection lost"))
			},
			expectedID:    0,
			expectedError: errors.New("database connection lost"),
		},
	}

	for _, test := range tests {
		suite.Run(test.name, func() {
			test.setupMock()

			id, err := suite.repo.Create(test.host)

			if test.expectedError != nil {
				assert.Error(suite.T(), err)
				assert.Equal(suite.T(), test.expectedError.Error(), err.Error())
				assert.Equal(suite.T(), int64(0), id)
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

// TestCreateWithLastInsertIdError tests Create when LastInsertId fails
func (suite *HostRepositoryTestSuite) TestCreateWithLastInsertIdError() {
	host := &entities.Host{
		Hostname:  "test-server.example.com",
		IPAddress: "192.168.1.50",
		Role:      "test",
	}

	suite.mock.ExpectExec("INSERT INTO hosts \\(hostname, ip_address, role, created_at, last_seen\\) VALUES \\(\\?, \\?, \\?, \\?, \\?\\)").
		WithArgs("test-server.example.com", "192.168.1.50", "test",
			sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewErrorResult(errors.New("LastInsertId not supported")))

	id, err := suite.repo.Create(host)

	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), "LastInsertId not supported", err.Error())
	assert.Equal(suite.T(), int64(0), id)
}

// TestUpdate tests the Update method
func (suite *HostRepositoryTestSuite) TestUpdate() {
	tests := []struct {
		name          string
		id            int64
		host          *entities.Host
		setupMock     func()
		expectedError error
	}{
		{
			name: "successful_update",
			id:   1,
			host: &entities.Host{
				Role: "updated-role",
			},
			setupMock: func() {
				suite.mock.ExpectExec("UPDATE hosts SET role = \\?, last_seen = \\? WHERE id = \\?").
					WithArgs("updated-role", sqlmock.AnyArg(), int64(1)).
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
			expectedError: nil,
		},
		{
			name: "update_non_existent_host",
			id:   999,
			host: &entities.Host{
				Role: "web-server",
			},
			setupMock: func() {
				suite.mock.ExpectExec("UPDATE hosts SET role = \\?, last_seen = \\? WHERE id = \\?").
					WithArgs("web-server", sqlmock.AnyArg(), int64(999)).
					WillReturnResult(sqlmock.NewResult(0, 0))
			},
			expectedError: nil, // Note: Current implementation doesn't check rows affected
		},
		{
			name: "database_error",
			id:   2,
			host: &entities.Host{
				Role: "cache",
			},
			setupMock: func() {
				suite.mock.ExpectExec("UPDATE hosts SET role = \\?, last_seen = \\? WHERE id = \\?").
					WithArgs("cache", sqlmock.AnyArg(), int64(2)).
					WillReturnError(errors.New("database locked"))
			},
			expectedError: errors.New("database locked"),
		},
	}

	for _, test := range tests {
		suite.Run(test.name, func() {
			test.setupMock()

			err := suite.repo.Update(test.id, test.host)

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

// TestDelete tests the Delete method
func (suite *HostRepositoryTestSuite) TestDelete() {
	tests := []struct {
		name          string
		id            int64
		setupMock     func()
		expectedError error
	}{
		{
			name: "successful_deletion",
			id:   1,
			setupMock: func() {
				suite.mock.ExpectExec("DELETE FROM hosts WHERE id = \\?").
					WithArgs(int64(1)).
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
			expectedError: nil,
		},
		{
			name: "delete_non_existent_host",
			id:   999,
			setupMock: func() {
				suite.mock.ExpectExec("DELETE FROM hosts WHERE id = \\?").
					WithArgs(int64(999)).
					WillReturnResult(sqlmock.NewResult(0, 0))
			},
			expectedError: sql.ErrNoRows,
		},
		{
			name: "database_error",
			id:   2,
			setupMock: func() {
				suite.mock.ExpectExec("DELETE FROM hosts WHERE id = \\?").
					WithArgs(int64(2)).
					WillReturnError(errors.New("foreign key constraint failed"))
			},
			expectedError: errors.New("foreign key constraint failed"),
		},
		{
			name: "rows_affected_error",
			id:   3,
			setupMock: func() {
				suite.mock.ExpectExec("DELETE FROM hosts WHERE id = \\?").
					WithArgs(int64(3)).
					WillReturnResult(sqlmock.NewErrorResult(errors.New("rows affected error")))
			},
			expectedError: errors.New("rows affected error"),
		},
	}

	for _, test := range tests {
		suite.Run(test.name, func() {
			test.setupMock()

			err := suite.repo.Delete(test.id)

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

// TestScanHostsErrorHandling tests error handling in scanHosts helper
func (suite *HostRepositoryTestSuite) TestScanHostsErrorHandling() {
	// Test rows.Err() handling
	suite.mock.ExpectQuery("SELECT id, hostname, ip_address, role FROM hosts WHERE 1=1").
		WillReturnRows(sqlmock.NewRows([]string{"id", "hostname", "ip_address", "role"}).
			AddRow(1, "server-01", "192.168.1.100", "web").
			RowError(0, errors.New("row error")))

	hosts, err := suite.repo.FindByFilters(&entities.HostQueryParams{})

	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), "row error", err.Error())
	assert.Nil(suite.T(), hosts)
}

// Run the test suite
func TestHostRepositoryTestSuite(test *testing.T) {
	suite.Run(test, new(HostRepositoryTestSuite))
}

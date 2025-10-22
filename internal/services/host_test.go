// nolint
package services

import (
	"database/sql"
	"errors"
	"testing"

	"github.com/gabrielg2020/monitor-api/internal/entities"
	"github.com/gabrielg2020/monitor-api/test/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// HostServiceTestSuite is the test suite for HostService
type HostServiceTestSuite struct {
	suite.Suite
	mockRepo *mocks.MockHostRepository
	service  *HostService
}

// SetupTest runs before each test in the suite
func (suite *HostServiceTestSuite) SetupTest() {
	suite.mockRepo = new(mocks.MockHostRepository)
	suite.service = NewHostService(suite.mockRepo)
}

// TearDownTest runs after each test
func (suite *HostServiceTestSuite) TearDownTest() {
	suite.mockRepo.AssertExpectations(suite.T())
}

// TestNewHostService tests the constructor
func (suite *HostServiceTestSuite) TestNewHostService() {
	assert.NotNil(suite.T(), suite.service)
	assert.Equal(suite.T(), suite.mockRepo, suite.service.repo)
}

// TestCreateHost tests the CreateHost method
func (suite *HostServiceTestSuite) TestCreateHost() {
	tests := []struct {
		name          string
		host          *entities.Host
		setupMock     func()
		expectedID    int64
		expectedError error
		description   string
	}{
		{
			name: "successful_host_creation",
			host: &entities.Host{
				Hostname:  "server-01.example.com",
				IPAddress: "192.168.1.100",
				Role:      "web-server",
			},
			setupMock: func() {
				suite.mockRepo.On("Create", &entities.Host{
					Hostname:  "server-01.example.com",
					IPAddress: "192.168.1.100",
					Role:      "web-server",
				}).Return(int64(1), nil).Once()
			},
			expectedID:    1,
			expectedError: nil,
			description:   "Should return the new host ID on successful creation",
		},
		{
			name: "duplicate_hostname_error",
			host: &entities.Host{
				Hostname:  "existing-server.example.com",
				IPAddress: "192.168.1.101",
				Role:      "database",
			},
			setupMock: func() {
				suite.mockRepo.On("Create", &entities.Host{
					Hostname:  "existing-server.example.com",
					IPAddress: "192.168.1.101",
					Role:      "database",
				}).Return(int64(0), errors.New("UNIQUE constraint failed: hosts.hostname")).Once()
			},
			expectedID:    0,
			expectedError: errors.New("UNIQUE constraint failed: hosts.hostname"),
			description:   "Should return an error when trying to create a host with a duplicate hostname",
		},
		{
			name: "duplicate_ip_address_error",
			host: &entities.Host{
				Hostname:  "server-02.example.com",
				IPAddress: "duplicate-ip",
				Role:      "database",
			},
			setupMock: func() {
				suite.mockRepo.On("Create", &entities.Host{
					Hostname:  "server-02.example.com",
					IPAddress: "duplicate-ip",
					Role:      "database",
				}).Return(int64(0), errors.New("UNIQUE constraint failed: hosts.ip_address")).Once()
			},
			expectedID:    0,
			expectedError: errors.New("UNIQUE constraint failed: hosts.ip_address"),
			description:   "Should return an error when trying to create a host with a duplicate IP address",
		},
		{
			name: "invalid_ip_address",
			host: &entities.Host{
				Hostname:  "server-03.example.com",
				IPAddress: "invalid-ip",
				Role:      "cache",
			},
			setupMock: func() {
				suite.mockRepo.On("Create", &entities.Host{
					Hostname:  "server-03.example.com",
					IPAddress: "invalid-ip",
					Role:      "cache",
				}).Return(int64(0), errors.New("invalid IP address format")).Once()
			},
			expectedID:    0,
			expectedError: errors.New("invalid IP address format"),
			description:   "Should return an error when trying to create a host with invalid IP address",
		},
		{
			name: "database_connection_error",
			host: &entities.Host{
				Hostname:  "server-04.example.com",
				IPAddress: "192.168.1.103",
				Role:      "monitoring",
			},
			setupMock: func() {
				suite.mockRepo.On("Create", &entities.Host{
					Hostname:  "server-04.example.com",
					IPAddress: "192.168.1.103",
					Role:      "monitoring",
				}).Return(int64(0), errors.New("database connection lost")).Once()
			},
			expectedID:    0,
			expectedError: errors.New("database connection lost"),
			description:   "Should return an error when there is a database connection issue",
		},
	}

	for _, test := range tests {
		suite.Run(test.name, func() {
			test.setupMock()

			id, err := suite.service.CreateHost(test.host)

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

// TestGetHosts tests the GetHosts method with filters
func (suite *HostServiceTestSuite) TestGetHosts() {
	tests := []struct {
		name          string
		params        *entities.HostQueryParams
		setupMock     func()
		expectedHosts []entities.Host
		expectedError error
		description   string
	}{
		{
			name: "get_hosts_by_hostname",
			params: &entities.HostQueryParams{
				Hostname: "server-01.example.com",
			},
			setupMock: func() {
				hosts := []entities.Host{
					{
						ID:        1,
						Hostname:  "server-01.example.com",
						IPAddress: "192.168.1.100",
						Role:      "web-server",
					},
				}
				suite.mockRepo.On("FindByFilters", &entities.HostQueryParams{
					Hostname: "server-01.example.com",
				}).Return(hosts, nil).Once()
			},
			expectedHosts: []entities.Host{
				{
					ID:        1,
					Hostname:  "server-01.example.com",
					IPAddress: "192.168.1.100",
					Role:      "web-server",
				},
			},
			expectedError: nil,
			description:   "Should return the hosts on successful query",
		},
		{
			name: "get_hosts_by_ip_address",
			params: &entities.HostQueryParams{
				IPAddress: "192.168.1.100",
			},
			setupMock: func() {
				hosts := []entities.Host{
					{
						ID:        1,
						Hostname:  "server-01.example.com",
						IPAddress: "192.168.1.100",
						Role:      "web-server",
					},
				}
				suite.mockRepo.On("FindByFilters", &entities.HostQueryParams{
					IPAddress: "192.168.1.100",
				}).Return(hosts, nil).Once()
			},
			expectedHosts: []entities.Host{
				{
					ID:        1,
					Hostname:  "server-01.example.com",
					IPAddress: "192.168.1.100",
					Role:      "web-server",
				},
			},
			expectedError: nil,
			description:   "Should return the hosts on successful query",
		},
		{
			name: "get_hosts_by_id",
			params: &entities.HostQueryParams{
				ID: 1,
			},
			setupMock: func() {
				hosts := []entities.Host{
					{
						ID:        1,
						Hostname:  "server-01.example.com",
						IPAddress: "192.168.1.100",
						Role:      "web-server",
					},
				}
				suite.mockRepo.On("FindByFilters", &entities.HostQueryParams{
					ID: 1,
				}).Return(hosts, nil).Once()
			},
			expectedHosts: []entities.Host{
				{
					ID:        1,
					Hostname:  "server-01.example.com",
					IPAddress: "192.168.1.100",
					Role:      "web-server",
				},
			},
			expectedError: nil,
			description:   "Should return the hosts on successful query",
		},
		{
			name: "no_hosts_found",
			params: &entities.HostQueryParams{
				Hostname: "non-existent-server.example.com",
			},
			setupMock: func() {
				suite.mockRepo.On("FindByFilters", &entities.HostQueryParams{
					Hostname: "non-existent-server.example.com",
				}).Return([]entities.Host{}, nil).Once()
			},
			expectedHosts: []entities.Host{},
			expectedError: nil,
			description:   "Should return an empty slice when no hosts match the query",
		},
		{
			name: "database_error",
			params: &entities.HostQueryParams{
				ID: 999,
			},
			setupMock: func() {
				suite.mockRepo.On("FindByFilters", &entities.HostQueryParams{
					ID: 999,
				}).Return(nil, errors.New("database query failed")).Once()
			},
			expectedHosts: nil,
			expectedError: errors.New("database query failed"),
			description:   "Should return an error when there is a database query issue",
		},
	}

	for _, test := range tests {
		suite.Run(test.name, func() {
			test.setupMock()

			hosts, err := suite.service.GetHosts(test.params)

			assert.Equal(suite.T(), test.expectedHosts, hosts)
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

// TestUpdateHost tests the UpdateHost method
func (suite *HostServiceTestSuite) TestUpdateHost() {
	tests := []struct {
		name          string
		id            int64
		host          *entities.Host
		setupMock     func()
		expectedError error
		description   string
	}{
		{
			name: "successful_update",
			id:   1,
			host: &entities.Host{
				Role: "database",
			},
			setupMock: func() {
				suite.mockRepo.On("Update", int64(1), &entities.Host{
					Role: "database",
				}).Return(nil).Once()
			},
			expectedError: nil,
			description:   "Should update the host successfully",
		},
		{
			name: "update_non_existent_host",
			id:   999,
			host: &entities.Host{
				Role: "web-server",
			},
			setupMock: func() {
				suite.mockRepo.On("Update", int64(999), &entities.Host{
					Role: "web-server",
				}).Return(sql.ErrNoRows).Once()
			},
			expectedError: sql.ErrNoRows,
			description:   "Should return sql.ErrNoRows when trying to update a non-existent host",
		},
		{
			name: "database_error_during_update",
			id:   2,
			host: &entities.Host{
				Role: "cache",
			},
			setupMock: func() {
				suite.mockRepo.On("Update", int64(2), &entities.Host{
					Role: "cache",
				}).Return(errors.New("database locked")).Once()
			},
			expectedError: errors.New("database locked"),
			description:   "Should return an error when database error occurs",
		},
	}

	for _, test := range tests {
		suite.Run(test.name, func() {
			test.setupMock()

			err := suite.service.UpdateHost(test.id, test.host)

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

// TestDeleteHost tests the DeleteHost method
func (suite *HostServiceTestSuite) TestDeleteHost() {
	tests := []struct {
		name          string
		id            int64
		setupMock     func()
		expectedError error
		description   string
	}{
		{
			name: "successful_deletion",
			id:   1,
			setupMock: func() {
				suite.mockRepo.On("Delete", int64(1)).Return(nil).Once()
			},
			expectedError: nil,
			description:   "Should delete the host successfully",
		},
		{
			name: "deletion_of_non_existent_host",
			id:   999,
			setupMock: func() {
				suite.mockRepo.On("Delete", int64(999)).Return(sql.ErrNoRows).Once()
			},
			expectedError: sql.ErrNoRows,
			description:   "Should return sql.ErrNoRows when trying to delete a non-existent host",
		},
		{
			name: "database_error_during_deletion",
			id:   2,
			setupMock: func() {
				suite.mockRepo.On("Delete", int64(2)).Return(errors.New("foreign key constraint failed")).Once()
			},
			expectedError: errors.New("foreign key constraint failed"),
			description:   "Should return an error when database error occurs",
		},
	}

	for _, test := range tests {
		suite.Run(test.name, func() {
			test.setupMock()

			err := suite.service.DeleteHost(test.id)

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
func TestHostServiceTestSuite(test *testing.T) {
	suite.Run(test, new(HostServiceTestSuite))
}

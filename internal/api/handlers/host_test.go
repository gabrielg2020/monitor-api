// nolint
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
	"github.com/stretchr/testify/suite"
)

// HostHandlerTestSuite is the test suite for HostHandler
type HostHandlerTestSuite struct {
	suite.Suite
	router      *gin.Engine
	mockService *mocks.MockHostService
	handler     *HostHandler
}

// SetupTest runs before each test in the suite
func (suite *HostHandlerTestSuite) SetupTest() {
	gin.SetMode(gin.TestMode)
	suite.router = gin.New()
	suite.mockService = new(mocks.MockHostService)
	suite.handler = NewHostHandler(suite.mockService)

	// Register routes
	suite.router.POST("/hosts", suite.handler.Create)
	suite.router.GET("/hosts", suite.handler.Get)
	suite.router.PUT("/hosts/:id", suite.handler.Update)
	suite.router.DELETE("/hosts/:id", suite.handler.Delete)
}

// TearDownTest runs after each test
func (suite *HostHandlerTestSuite) TearDownTest() {
	suite.mockService.AssertExpectations(suite.T())
}

// TestNewHostHandler tests the constructor
func (suite *HostHandlerTestSuite) TestNewHostHandler() {
	assert.NotNil(suite.T(), suite.handler)
	assert.NotNil(suite.T(), suite.handler.service)
}

// TestCreate tests the Create endpoint
func (suite *HostHandlerTestSuite) TestCreate() {
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
				"host": map[string]interface{}{
					"hostname":   "pi-monitor-01",
					"ip_address": "192.168.1.100",
					"role":       "monitor",
				},
			},
			setupMock: func() {
				suite.mockService.On("CreateHost", &entities.Host{
					Hostname:  "pi-monitor-01",
					IPAddress: "192.168.1.100",
					Role:      "monitor",
				}).Return(int64(1), nil).Once()
			},
			expectedStatus: http.StatusCreated,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, "Host created successfully", response["message"])
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
				"host": map[string]interface{}{
					"hostname": "pi-monitor-01",
					// Missing ip_address and role
				},
			},
			setupMock: func() {
				suite.mockService.On("CreateHost", &entities.Host{
					Hostname: "pi-monitor-01",
				}).Return(int64(0), errors.New("missing required fields")).Once()
			},
			expectedStatus: http.StatusInternalServerError,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response models.ErrorResponse
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, "Failed to create host", response.Error)
			},
		},
		{
			name: "duplicate_hostname_error",
			requestBody: map[string]interface{}{
				"host": map[string]interface{}{
					"hostname":   "existing-host",
					"ip_address": "192.168.1.200",
					"role":       "monitor",
				},
			},
			setupMock: func() {
				suite.mockService.On("CreateHost", &entities.Host{
					Hostname:  "existing-host",
					IPAddress: "192.168.1.200",
					Role:      "monitor",
				}).Return(int64(0), errors.New("UNIQUE constraint failed: hosts.hostname")).Once()
			},
			expectedStatus: http.StatusInternalServerError,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response models.ErrorResponse
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, "Failed to create host", response.Error)
				assert.Contains(t, response.Details, "UNIQUE constraint failed")
			},
		},
		{
			name: "database_connection_error",
			requestBody: map[string]interface{}{
				"host": map[string]interface{}{
					"hostname":   "pi-monitor-02",
					"ip_address": "192.168.1.101",
					"role":       "monitor",
				},
			},
			setupMock: func() {
				suite.mockService.On("CreateHost", &entities.Host{
					Hostname:  "pi-monitor-02",
					IPAddress: "192.168.1.101",
					Role:      "monitor",
				}).Return(int64(0), errors.New("database connection lost")).Once()
			},
			expectedStatus: http.StatusInternalServerError,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response models.ErrorResponse
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, "Failed to create host", response.Error)
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
			req, err := http.NewRequest(http.MethodPost, "/hosts", bytes.NewBuffer(bodyBytes))
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
func (suite *HostHandlerTestSuite) TestGet() {
	tests := []struct {
		name           string
		queryParams    string
		setupMock      func()
		expectedStatus int
		checkResponse  func(*testing.T, *httptest.ResponseRecorder)
	}{
		{
			name:        "get_all_hosts",
			queryParams: "",
			setupMock: func() {
				hosts := []entities.Host{
					{ID: 1, Hostname: "pi-monitor-01", IPAddress: "192.168.1.100", Role: "monitor"},
					{ID: 2, Hostname: "pi-monitor-02", IPAddress: "192.168.1.101", Role: "monitor"},
				}
				suite.mockService.On("GetHosts", &entities.HostQueryParams{}).Return(hosts, nil).Once()
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response models.HostListResponse
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Len(t, response.Hosts, 2)
				assert.Equal(t, 2, response.Meta.Count)
				assert.Equal(t, "pi-monitor-01", response.Hosts[0].Hostname)
				assert.Equal(t, "pi-monitor-02", response.Hosts[1].Hostname)
			},
		},
		{
			name:        "get_hosts_by_id",
			queryParams: "?id=1",
			setupMock: func() {
				hosts := []entities.Host{
					{ID: 1, Hostname: "pi-monitor-01", IPAddress: "192.168.1.100", Role: "monitor"},
				}
				suite.mockService.On("GetHosts", &entities.HostQueryParams{ID: 1}).Return(hosts, nil).Once()
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response models.HostListResponse
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Len(t, response.Hosts, 1)
				assert.Equal(t, int64(1), response.Hosts[0].ID)
			},
		},
		{
			name:        "get_hosts_by_hostname",
			queryParams: "?hostname=pi-monitor-01",
			setupMock: func() {
				hosts := []entities.Host{
					{ID: 1, Hostname: "pi-monitor-01", IPAddress: "192.168.1.100", Role: "monitor"},
				}
				suite.mockService.On("GetHosts", &entities.HostQueryParams{Hostname: "pi-monitor-01"}).Return(hosts, nil).Once()
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response models.HostListResponse
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Len(t, response.Hosts, 1)
				assert.Equal(t, "pi-monitor-01", response.Hosts[0].Hostname)
			},
		},
		{
			name:        "get_hosts_by_ip_address",
			queryParams: "?ip_address=192.168.1.100",
			setupMock: func() {
				hosts := []entities.Host{
					{ID: 1, Hostname: "pi-monitor-01", IPAddress: "192.168.1.100", Role: "monitor"},
				}
				suite.mockService.On("GetHosts", &entities.HostQueryParams{IPAddress: "192.168.1.100"}).Return(hosts, nil).Once()
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response models.HostListResponse
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Len(t, response.Hosts, 1)
				assert.Equal(t, "192.168.1.100", response.Hosts[0].IPAddress)
			},
		},
		{
			name:        "get_hosts_empty_result",
			queryParams: "",
			setupMock: func() {
				var hosts []entities.Host
				suite.mockService.On("GetHosts", &entities.HostQueryParams{}).Return(hosts, nil).Once()
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response models.HostListResponse
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Len(t, response.Hosts, 0)
				assert.Equal(t, 0, response.Meta.Count)
			},
		},
		{
			name:           "invalid_query_parameter",
			queryParams:    "?id=invalid",
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
			queryParams: "",
			setupMock: func() {
				suite.mockService.On("GetHosts", &entities.HostQueryParams{}).Return(nil, errors.New("database connection lost")).Once()
			},
			expectedStatus: http.StatusInternalServerError,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response models.ErrorResponse
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, "Failed to retrieve hosts", response.Error)
				assert.Equal(t, "database connection lost", response.Details)
			},
		},
	}

	for _, test := range tests {
		suite.Run(test.name, func() {
			test.setupMock()

			// Create request
			req, err := http.NewRequest(http.MethodGet, "/hosts"+test.queryParams, nil)
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

// TestUpdate tests the Update endpoint
func (suite *HostHandlerTestSuite) TestUpdate() {
	tests := []struct {
		name           string
		hostID         string
		requestBody    interface{}
		setupMock      func()
		expectedStatus int
		checkResponse  func(*testing.T, *httptest.ResponseRecorder)
	}{
		{
			name:   "successful_update",
			hostID: "1",
			requestBody: map[string]interface{}{
				"host": map[string]interface{}{
					"role": "updated-role",
				},
			},
			setupMock: func() {
				suite.mockService.On("UpdateHost", int64(1), &entities.Host{
					Role: "updated-role",
				}).Return(nil).Once()
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, "Host updated successfully", response["message"])
			},
		},
		{
			name:   "invalid_host_id",
			hostID: "invalid",
			requestBody: map[string]interface{}{
				"host": map[string]interface{}{
					"role": "monitor",
				},
			},
			setupMock:      func() {},
			expectedStatus: http.StatusBadRequest,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response models.ErrorResponse
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, "Invalid host ID", response.Error)
			},
		},
		{
			name:        "invalid_json_body",
			hostID:      "1",
			requestBody: "invalid json",
			setupMock: func() {
				// No mock setup needed as validation happens before service call
			},
			expectedStatus: http.StatusBadRequest,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response models.ErrorResponse
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, "Invalid request body", response.Error)
			},
		},
		{
			name:   "host_not_found",
			hostID: "999",
			requestBody: map[string]interface{}{
				"host": map[string]interface{}{
					"role": "monitor",
				},
			},
			setupMock: func() {
				suite.mockService.On("UpdateHost", int64(999), &entities.Host{
					Role: "monitor",
				}).Return(errors.New("host not found")).Once()
			},
			expectedStatus: http.StatusInternalServerError,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response models.ErrorResponse
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, "Failed to update host", response.Error)
			},
		},
		{
			name:   "database_error",
			hostID: "1",
			requestBody: map[string]interface{}{
				"host": map[string]interface{}{
					"role": "monitor",
				},
			},
			setupMock: func() {
				suite.mockService.On("UpdateHost", int64(1), &entities.Host{
					Role: "monitor",
				}).Return(errors.New("database connection lost")).Once()
			},
			expectedStatus: http.StatusInternalServerError,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response models.ErrorResponse
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, "Failed to update host", response.Error)
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
			req, err := http.NewRequest(http.MethodPut, "/hosts/"+test.hostID, bytes.NewBuffer(bodyBytes))
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

// TestDelete tests the Delete endpoint
func (suite *HostHandlerTestSuite) TestDelete() {
	tests := []struct {
		name           string
		hostID         string
		setupMock      func()
		expectedStatus int
		checkResponse  func(*testing.T, *httptest.ResponseRecorder)
	}{
		{
			name:   "successful_deletion",
			hostID: "1",
			setupMock: func() {
				suite.mockService.On("DeleteHost", int64(1)).Return(nil).Once()
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, "Host deleted successfully", response["message"])
			},
		},
		{
			name:           "invalid_host_id",
			hostID:         "invalid",
			setupMock:      func() {},
			expectedStatus: http.StatusBadRequest,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response models.ErrorResponse
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, "Invalid host ID", response.Error)
			},
		},
		{
			name:   "host_not_found",
			hostID: "999",
			setupMock: func() {
				suite.mockService.On("DeleteHost", int64(999)).Return(errors.New("host not found")).Once()
			},
			expectedStatus: http.StatusInternalServerError,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response models.ErrorResponse
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, "Failed to delete host", response.Error)
			},
		},
		{
			name:   "database_error",
			hostID: "1",
			setupMock: func() {
				suite.mockService.On("DeleteHost", int64(1)).Return(errors.New("database connection lost")).Once()
			},
			expectedStatus: http.StatusInternalServerError,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response models.ErrorResponse
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, "Failed to delete host", response.Error)
				assert.Equal(t, "database connection lost", response.Details)
			},
		},
		{
			name:   "foreign_key_constraint_error",
			hostID: "1",
			setupMock: func() {
				suite.mockService.On("DeleteHost", int64(1)).Return(errors.New("FOREIGN KEY constraint failed")).Once()
			},
			expectedStatus: http.StatusInternalServerError,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response models.ErrorResponse
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, "Failed to delete host", response.Error)
				assert.Contains(t, response.Details, "FOREIGN KEY constraint")
			},
		},
	}

	for _, test := range tests {
		suite.Run(test.name, func() {
			test.setupMock()

			// Create request
			req, err := http.NewRequest(http.MethodDelete, "/hosts/"+test.hostID, nil)
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
func (suite *HostHandlerTestSuite) TestContentType() {
	hosts := []entities.Host{
		{ID: 1, Hostname: "pi-monitor-01", IPAddress: "192.168.1.100", Role: "monitor"},
	}
	suite.mockService.On("GetHosts", &entities.HostQueryParams{}).Return(hosts, nil).Once()

	req, err := http.NewRequest(http.MethodGet, "/hosts", nil)
	assert.NoError(suite.T(), err)

	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), "application/json; charset=utf-8", w.Header().Get("Content-Type"))
}

// Run the test suite
func TestHostHandlerTestSuite(test *testing.T) {
	suite.Run(test, new(HostHandlerTestSuite))
}

// nolint
package api

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gabrielg2020/monitor-api/test/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

// RouterTestSuite is the test suite for router setup
type RouterTestSuite struct {
	suite.Suite
	mockHealthHandler *mocks.MockHealthHandler
	mockHostHandler   *mocks.MockHostHandler
	mockMetricHandler *mocks.MockMetricHandler
}

// SetupTest runs before each test in the suite
func (suite *RouterTestSuite) SetupTest() {
	suite.mockHealthHandler = new(mocks.MockHealthHandler)
	suite.mockHostHandler = new(mocks.MockHostHandler)
	suite.mockMetricHandler = new(mocks.MockMetricHandler)
}

// TearDownTest runs after each test
func (suite *RouterTestSuite) TearDownTest() {
	suite.mockHealthHandler.AssertExpectations(suite.T())
	suite.mockHostHandler.AssertExpectations(suite.T())
	suite.mockMetricHandler.AssertExpectations(suite.T())
}

// TestSetupRouter tests the router initialisation
func (suite *RouterTestSuite) TestSetupRouter() {
	allowedOrigins := []string{"http://localhost:3000"}
	router := SetupRouter(
		suite.mockHealthHandler,
		suite.mockHostHandler,
		suite.mockMetricHandler,
		allowedOrigins,
	)

	assert.NotNil(suite.T(), router)
}

// TestHealthRoutes tests that health routes are registered and call correct handlers
func (suite *RouterTestSuite) TestHealthRoutes() {
	tests := []struct {
		name           string
		method         string
		path           string
		setupMock      func()
		expectedStatus int
	}{
		{
			name:   "health_endpoint_calls_handler",
			method: http.MethodGet,
			path:   "/health",
			setupMock: func() {
				suite.mockHealthHandler.On("GetHealth", mock.AnythingOfType("*gin.Context")).Once()
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:   "detailed_health_endpoint_calls_handler",
			method: http.MethodGet,
			path:   "/health/detailed",
			setupMock: func() {
				suite.mockHealthHandler.On("GetDetailedHealth", mock.AnythingOfType("*gin.Context")).Once()
			},
			expectedStatus: http.StatusOK,
		},
	}

	for _, test := range tests {
		suite.Run(test.name, func() {
			test.setupMock()

			router := SetupRouter(
				suite.mockHealthHandler,
				suite.mockHostHandler,
				suite.mockMetricHandler,
				[]string{"*"},
			)

			req, err := http.NewRequest(test.method, test.path, nil)
			assert.NoError(suite.T(), err)

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(suite.T(), test.expectedStatus, w.Code)
		})

		// Reset for next test
		suite.TearDownTest()
		suite.SetupTest()
	}
}

// TestSwaggerRoute tests that Swagger documentation is accessible
func (suite *RouterTestSuite) TestSwaggerRoute() {
	router := SetupRouter(
		suite.mockHealthHandler,
		suite.mockHostHandler,
		suite.mockMetricHandler,
		[]string{"*"},
	)

	req, err := http.NewRequest(http.MethodGet, "/swagger", nil)
	assert.NoError(suite.T(), err)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Swagger should be accessible (200) or redirect (301/302)
	// We're not testing Swagger itself, just that the route is registered
	assert.NotEqual(suite.T(), http.StatusNotFound, w.Code, "Swagger route should be registered")
}

// TestAPIv1HostRoutes tests that host routes are registered under /api/v1 and call correct handlers
func (suite *RouterTestSuite) TestAPIv1HostRoutes() {
	tests := []struct {
		name      string
		method    string
		path      string
		setupMock func()
	}{
		{
			name:   "post_hosts_calls_create",
			method: http.MethodPost,
			path:   "/api/v1/hosts",
			setupMock: func() {
				suite.mockHostHandler.On("Create", mock.AnythingOfType("*gin.Context")).Once()
			},
		},
		{
			name:   "get_hosts_calls_get",
			method: http.MethodGet,
			path:   "/api/v1/hosts",
			setupMock: func() {
				suite.mockHostHandler.On("Get", mock.AnythingOfType("*gin.Context")).Once()
			},
		},
		{
			name:   "put_hosts_calls_update",
			method: http.MethodPut,
			path:   "/api/v1/hosts",
			setupMock: func() {
				suite.mockHostHandler.On("Update", mock.AnythingOfType("*gin.Context")).Once()
			},
		},
		{
			name:   "delete_hosts_calls_delete",
			method: http.MethodDelete,
			path:   "/api/v1/hosts",
			setupMock: func() {
				suite.mockHostHandler.On("Delete", mock.AnythingOfType("*gin.Context")).Once()
			},
		},
	}

	for _, test := range tests {
		suite.Run(test.name, func() {
			test.setupMock()

			router := SetupRouter(
				suite.mockHealthHandler,
				suite.mockHostHandler,
				suite.mockMetricHandler,
				[]string{"*"},
			)

			req, err := http.NewRequest(test.method, test.path, nil)
			assert.NoError(suite.T(), err)

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			// We're not checking status codes here because that's the handler's responsibility
			// We just verify the route exists (not 404) and the correct handler was called
			assert.NotEqual(suite.T(), http.StatusNotFound, w.Code, "Route should be registered")
		})

		// Reset for next test
		suite.TearDownTest()
		suite.SetupTest()
	}
}

// TestAPIv1MetricRoutes tests that metric routes are registered under /api/v1 and call correct handlers
func (suite *RouterTestSuite) TestAPIv1MetricRoutes() {
	tests := []struct {
		name      string
		method    string
		path      string
		setupMock func()
	}{
		{
			name:   "post_metrics_calls_create",
			method: http.MethodPost,
			path:   "/api/v1/metrics",
			setupMock: func() {
				suite.mockMetricHandler.On("Create", mock.AnythingOfType("*gin.Context")).Once()
			},
		},
		{
			name:   "get_metrics_calls_get",
			method: http.MethodGet,
			path:   "/api/v1/metrics",
			setupMock: func() {
				suite.mockMetricHandler.On("Get", mock.AnythingOfType("*gin.Context")).Once()
			},
		},
		{
			name:   "get_latest_metrics_calls_get_latest",
			method: http.MethodGet,
			path:   "/api/v1/metrics/latest",
			setupMock: func() {
				suite.mockMetricHandler.On("GetLatest", mock.AnythingOfType("*gin.Context")).Once()
			},
		},
	}

	for _, test := range tests {
		suite.Run(test.name, func() {
			test.setupMock()

			router := SetupRouter(
				suite.mockHealthHandler,
				suite.mockHostHandler,
				suite.mockMetricHandler,
				[]string{"*"},
			)

			req, err := http.NewRequest(test.method, test.path, nil)
			assert.NoError(suite.T(), err)

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			// We're not checking status codes here because that's the handler's responsibility
			// We just verify the route exists (not 404) and the correct handler was called
			assert.NotEqual(suite.T(), http.StatusNotFound, w.Code, "Route should be registered")
		})

		// Reset for next test
		suite.TearDownTest()
		suite.SetupTest()
	}
}

// TestRouteNotFound tests that 404 is returned for non-existent routes
func (suite *RouterTestSuite) TestRouteNotFound() {
	tests := []struct {
		name           string
		method         string
		path           string
		expectedStatus int
	}{
		{
			name:           "non_existent_route",
			method:         http.MethodGet,
			path:           "/nonexistent",
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "non_existent_api_route",
			method:         http.MethodGet,
			path:           "/api/v1/nonexistent",
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "misspelled_health_route",
			method:         http.MethodGet,
			path:           "/helth",
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "invalid_api_version",
			method:         http.MethodGet,
			path:           "/api/v99/hosts",
			expectedStatus: http.StatusNotFound,
		},
	}

	for _, test := range tests {
		suite.Run(test.name, func() {
			router := SetupRouter(
				suite.mockHealthHandler,
				suite.mockHostHandler,
				suite.mockMetricHandler,
				[]string{"*"},
			)

			req, err := http.NewRequest(test.method, test.path, nil)
			assert.NoError(suite.T(), err)

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(suite.T(), test.expectedStatus, w.Code)
		})

		// Reset for next test
		suite.TearDownTest()
		suite.SetupTest()
	}
}

// TestMethodNotAllowed tests that 405 is returned when using wrong HTTP method on existing route
func (suite *RouterTestSuite) TestMethodNotAllowed() {
	tests := []struct {
		name           string
		method         string
		path           string
		expectedStatus int
	}{
		{
			name:           "post_on_health_endpoint",
			method:         http.MethodPost,
			path:           "/health",
			expectedStatus: http.StatusMethodNotAllowed,
		},
		{
			name:           "put_on_health_endpoint",
			method:         http.MethodPut,
			path:           "/health",
			expectedStatus: http.StatusMethodNotAllowed,
		},
		{
			name:           "delete_on_health_endpoint",
			method:         http.MethodDelete,
			path:           "/health",
			expectedStatus: http.StatusMethodNotAllowed,
		},
		{
			name:           "patch_on_hosts_endpoint",
			method:         http.MethodPatch,
			path:           "/api/v1/hosts",
			expectedStatus: http.StatusMethodNotAllowed,
		},
		{
			name:           "post_on_metrics_latest",
			method:         http.MethodPost,
			path:           "/api/v1/metrics/latest",
			expectedStatus: http.StatusMethodNotAllowed,
		},
		{
			name:           "delete_on_metrics_latest",
			method:         http.MethodDelete,
			path:           "/api/v1/metrics/latest",
			expectedStatus: http.StatusMethodNotAllowed,
		},
	}

	for _, test := range tests {
		suite.Run(test.name, func() {
			router := SetupRouter(
				suite.mockHealthHandler,
				suite.mockHostHandler,
				suite.mockMetricHandler,
				[]string{"*"},
			)

			req, err := http.NewRequest(test.method, test.path, nil)
			assert.NoError(suite.T(), err)

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(suite.T(), test.expectedStatus, w.Code)
		})

		// Reset for next test
		suite.TearDownTest()
		suite.SetupTest()
	}
}

// TestRouteGrouping tests that routes are properly grouped under /api/v1
func (suite *RouterTestSuite) TestRouteGrouping() {
	tests := []struct {
		name        string
		path        string
		method      string
		shouldExist bool
		setupMock   func()
	}{
		{
			name:        "hosts_under_v1",
			path:        "/api/v1/hosts",
			method:      http.MethodGet,
			shouldExist: true,
			setupMock: func() {
				suite.mockHostHandler.On("Get", mock.AnythingOfType("*gin.Context")).Once()
			},
		},
		{
			name:        "metrics_under_v1",
			path:        "/api/v1/metrics",
			method:      http.MethodGet,
			shouldExist: true,
			setupMock: func() {
				suite.mockMetricHandler.On("Get", mock.AnythingOfType("*gin.Context")).Once()
			},
		},
		{
			name:        "metrics_latest_under_v1",
			path:        "/api/v1/metrics/latest",
			method:      http.MethodGet,
			shouldExist: true,
			setupMock: func() {
				suite.mockMetricHandler.On("GetLatest", mock.AnythingOfType("*gin.Context")).Once()
			},
		},
		{
			name:        "hosts_not_at_root",
			path:        "/hosts",
			method:      http.MethodGet,
			shouldExist: false,
			setupMock:   func() {},
		},
		{
			name:        "metrics_not_at_root",
			path:        "/metrics",
			method:      http.MethodGet,
			shouldExist: false,
			setupMock:   func() {},
		},
		{
			name:        "health_at_root_not_v1",
			path:        "/health",
			method:      http.MethodGet,
			shouldExist: true,
			setupMock: func() {
				suite.mockHealthHandler.On("GetHealth", mock.AnythingOfType("*gin.Context")).Once()
			},
		},
		{
			name:        "health_not_under_v1",
			path:        "/api/v1/health",
			method:      http.MethodGet,
			shouldExist: false,
			setupMock:   func() {},
		},
	}

	for _, test := range tests {
		suite.Run(test.name, func() {
			test.setupMock()

			router := SetupRouter(
				suite.mockHealthHandler,
				suite.mockHostHandler,
				suite.mockMetricHandler,
				[]string{"*"},
			)

			req, err := http.NewRequest(test.method, test.path, nil)
			assert.NoError(suite.T(), err)

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if test.shouldExist {
				assert.NotEqual(suite.T(), http.StatusNotFound, w.Code, "Route should exist: %s", test.path)
			} else {
				assert.Equal(suite.T(), http.StatusNotFound, w.Code, "Route should not exist: %s", test.path)
			}
		})

		// Reset for next test
		suite.TearDownTest()
		suite.SetupTest()
	}
}

// TestAllRoutesRegistered tests that every expected route exists
func (suite *RouterTestSuite) TestAllRoutesRegistered() {
	routes := []struct {
		method    string
		path      string
		setupMock func()
	}{
		{
			method: http.MethodGet,
			path:   "/health",
			setupMock: func() {
				suite.mockHealthHandler.On("GetHealth", mock.AnythingOfType("*gin.Context")).Once()
			},
		},
		{
			method: http.MethodGet,
			path:   "/health/detailed",
			setupMock: func() {
				suite.mockHealthHandler.On("GetDetailedHealth", mock.AnythingOfType("*gin.Context")).Once()
			},
		},
		{
			method: http.MethodPost,
			path:   "/api/v1/hosts",
			setupMock: func() {
				suite.mockHostHandler.On("Create", mock.AnythingOfType("*gin.Context")).Once()
			},
		},
		{
			method: http.MethodGet,
			path:   "/api/v1/hosts",
			setupMock: func() {
				suite.mockHostHandler.On("Get", mock.AnythingOfType("*gin.Context")).Once()
			},
		},
		{
			method: http.MethodPut,
			path:   "/api/v1/hosts",
			setupMock: func() {
				suite.mockHostHandler.On("Update", mock.AnythingOfType("*gin.Context")).Once()
			},
		},
		{
			method: http.MethodDelete,
			path:   "/api/v1/hosts",
			setupMock: func() {
				suite.mockHostHandler.On("Delete", mock.AnythingOfType("*gin.Context")).Once()
			},
		},
		{
			method: http.MethodPost,
			path:   "/api/v1/metrics",
			setupMock: func() {
				suite.mockMetricHandler.On("Create", mock.AnythingOfType("*gin.Context")).Once()
			},
		},
		{
			method: http.MethodGet,
			path:   "/api/v1/metrics",
			setupMock: func() {
				suite.mockMetricHandler.On("Get", mock.AnythingOfType("*gin.Context")).Once()
			},
		},
		{
			method: http.MethodGet,
			path:   "/api/v1/metrics/latest",
			setupMock: func() {
				suite.mockMetricHandler.On("GetLatest", mock.AnythingOfType("*gin.Context")).Once()
			},
		},
	}

	for _, route := range routes {
		suite.Run(route.method+"_"+route.path, func() {
			route.setupMock()

			router := SetupRouter(
				suite.mockHealthHandler,
				suite.mockHostHandler,
				suite.mockMetricHandler,
				[]string{"*"},
			)

			req, err := http.NewRequest(route.method, route.path, nil)
			assert.NoError(suite.T(), err)

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			// Should not be 404 (route should exist and be accessible)
			assert.NotEqual(suite.T(), http.StatusNotFound, w.Code,
				"Route %s %s should be accessible", route.method, route.path)
		})

		// Reset for next test
		suite.TearDownTest()
		suite.SetupTest()
	}
}

// Run the test suite
func TestRouterTestSuite(t *testing.T) {
	suite.Run(t, new(RouterTestSuite))
}

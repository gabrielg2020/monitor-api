// nolint
package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// MiddlewareTestSuite is the test suite for middleware
type MiddlewareTestSuite struct {
	suite.Suite
}

// SetupTest runs before each test in the suite
func (suite *MiddlewareTestSuite) SetupTest() {
	gin.SetMode(gin.TestMode)
}

// TestCORSMiddleware tests that CORS middleware handles origins correctly
func (suite *MiddlewareTestSuite) TestCORSMiddleware() {
	tests := []struct {
		name           string
		allowedOrigins []string
		requestOrigin  string
		expectHeader   bool
		expectedOrigin string
	}{
		{
			name:           "specific_origin_allowed",
			allowedOrigins: []string{"http://localhost:3000"},
			requestOrigin:  "http://localhost:3000",
			expectHeader:   true,
			expectedOrigin: "http://localhost:3000",
		},
		{
			name:           "wildcard_origin",
			allowedOrigins: []string{"*"},
			requestOrigin:  "http://example.com",
			expectHeader:   true,
			expectedOrigin: "*",
		},
		{
			name:           "origin_not_allowed",
			allowedOrigins: []string{"http://localhost:3000"},
			requestOrigin:  "http://malicious.com",
			expectHeader:   false,
			expectedOrigin: "",
		},
		{
			name:           "multiple_origins_first_matches",
			allowedOrigins: []string{"http://localhost:3000", "http://localhost:4000"},
			requestOrigin:  "http://localhost:3000",
			expectHeader:   true,
			expectedOrigin: "http://localhost:3000",
		},
		{
			name:           "multiple_origins_second_matches",
			allowedOrigins: []string{"http://localhost:3000", "http://localhost:4000"},
			requestOrigin:  "http://localhost:4000",
			expectHeader:   true,
			expectedOrigin: "http://localhost:4000",
		},
		{
			name:           "empty_origin_header",
			allowedOrigins: []string{"http://localhost:3000"},
			requestOrigin:  "",
			expectHeader:   false,
			expectedOrigin: "",
		},
	}

	for _, test := range tests {
		suite.Run(test.name, func() {
			// Create a test router with CORS middleware
			router := gin.New()
			router.Use(CORS(test.allowedOrigins))

			// Add a simple test endpoint
			router.GET("/test", func(c *gin.Context) {
				c.Status(http.StatusOK)
			})

			req, err := http.NewRequest(http.MethodGet, "/test", nil)
			assert.NoError(suite.T(), err)
			if test.requestOrigin != "" {
				req.Header.Set("Origin", test.requestOrigin)
			}

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if test.expectHeader {
				assert.Equal(suite.T(), test.expectedOrigin, w.Header().Get("Access-Control-Allow-Origin"))
			} else {
				assert.Empty(suite.T(), w.Header().Get("Access-Control-Allow-Origin"))
			}

			// Verify other CORS headers are always set
			assert.Equal(suite.T(), "true", w.Header().Get("Access-Control-Allow-Credentials"))
			assert.NotEmpty(suite.T(), w.Header().Get("Access-Control-Allow-Methods"))
			assert.NotEmpty(suite.T(), w.Header().Get("Access-Control-Allow-Headers"))
		})
	}
}

// TestCORSPreflightRequest tests OPTIONS requests for CORS preflight
func (suite *MiddlewareTestSuite) TestCORSPreflightRequest() {
	tests := []struct {
		name                  string
		allowedOrigins        []string
		requestOrigin         string
		requestMethod         string
		expectedStatus        int
		expectCORSHeaders     bool
		expectedAllowedOrigin string
	}{
		{
			name:                  "preflight_allowed_origin",
			allowedOrigins:        []string{"http://localhost:3000"},
			requestOrigin:         "http://localhost:3000",
			requestMethod:         "POST",
			expectedStatus:        http.StatusNoContent,
			expectCORSHeaders:     true,
			expectedAllowedOrigin: "http://localhost:3000",
		},
		{
			name:                  "preflight_wildcard_origin",
			allowedOrigins:        []string{"*"},
			requestOrigin:         "http://example.com",
			requestMethod:         "POST",
			expectedStatus:        http.StatusNoContent,
			expectCORSHeaders:     true,
			expectedAllowedOrigin: "*",
		},
		{
			name:                  "preflight_not_allowed_origin",
			allowedOrigins:        []string{"http://localhost:3000"},
			requestOrigin:         "http://malicious.com",
			requestMethod:         "POST",
			expectedStatus:        http.StatusNoContent,
			expectCORSHeaders:     false,
			expectedAllowedOrigin: "",
		},
		{
			name:                  "preflight_put_request",
			allowedOrigins:        []string{"http://localhost:3000"},
			requestOrigin:         "http://localhost:3000",
			requestMethod:         "PUT",
			expectedStatus:        http.StatusNoContent,
			expectCORSHeaders:     true,
			expectedAllowedOrigin: "http://localhost:3000",
		},
		{
			name:                  "preflight_delete_request",
			allowedOrigins:        []string{"http://localhost:3000"},
			requestOrigin:         "http://localhost:3000",
			requestMethod:         "DELETE",
			expectedStatus:        http.StatusNoContent,
			expectCORSHeaders:     true,
			expectedAllowedOrigin: "http://localhost:3000",
		},
	}

	for _, test := range tests {
		suite.Run(test.name, func() {
			// Create a test router with CORS middleware
			router := gin.New()
			router.Use(CORS(test.allowedOrigins))

			// Add a test endpoint (won't be reached for OPTIONS)
			router.POST("/test", func(c *gin.Context) {
				c.Status(http.StatusOK)
			})

			req, err := http.NewRequest(http.MethodOptions, "/test", nil)
			assert.NoError(suite.T(), err)
			req.Header.Set("Origin", test.requestOrigin)
			req.Header.Set("Access-Control-Request-Method", test.requestMethod)

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			// Should return 204 No Content for OPTIONS
			assert.Equal(suite.T(), test.expectedStatus, w.Code)

			// Check CORS headers
			if test.expectCORSHeaders {
				assert.Equal(suite.T(), test.expectedAllowedOrigin, w.Header().Get("Access-Control-Allow-Origin"))
				assert.NotEmpty(suite.T(), w.Header().Get("Access-Control-Allow-Methods"))
				assert.NotEmpty(suite.T(), w.Header().Get("Access-Control-Allow-Headers"))
				assert.Equal(suite.T(), "true", w.Header().Get("Access-Control-Allow-Credentials"))
			} else {
				assert.Empty(suite.T(), w.Header().Get("Access-Control-Allow-Origin"))
			}
		})
	}
}

// TestCORSAllowedMethods tests that correct methods are in Access-Control-Allow-Methods header
func (suite *MiddlewareTestSuite) TestCORSAllowedMethods() {
	router := gin.New()
	router.Use(CORS([]string{"*"}))
	router.GET("/test", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	req, err := http.NewRequest(http.MethodGet, "/test", nil)
	assert.NoError(suite.T(), err)
	req.Header.Set("Origin", "http://example.com")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	allowedMethods := w.Header().Get("Access-Control-Allow-Methods")
	assert.Contains(suite.T(), allowedMethods, "POST")
	assert.Contains(suite.T(), allowedMethods, "GET")
	assert.Contains(suite.T(), allowedMethods, "PUT")
	assert.Contains(suite.T(), allowedMethods, "DELETE")
	assert.Contains(suite.T(), allowedMethods, "OPTIONS")
}

// TestCORSAllowedHeaders tests that correct headers are in Access-Control-Allow-Headers
func (suite *MiddlewareTestSuite) TestCORSAllowedHeaders() {
	router := gin.New()
	router.Use(CORS([]string{"*"}))
	router.GET("/test", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	req, err := http.NewRequest(http.MethodGet, "/test", nil)
	assert.NoError(suite.T(), err)
	req.Header.Set("Origin", "http://example.com")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	allowedHeaders := w.Header().Get("Access-Control-Allow-Headers")
	assert.Contains(suite.T(), allowedHeaders, "Content-Type")
	assert.Contains(suite.T(), allowedHeaders, "Authorization")
	assert.Contains(suite.T(), allowedHeaders, "X-CSRF-Token")
}

// Run the test suite
func TestMiddlewareTestSuite(t *testing.T) {
	suite.Run(t, new(MiddlewareTestSuite))
}

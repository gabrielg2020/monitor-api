// nolint
package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// ConfigTestSuite is the test suite for configuration tests
type ConfigTestSuite struct {
	suite.Suite
	originalEnv map[string]string // Save original env vars
}

// SetupTest runs before each test in the suite
func (suite *ConfigTestSuite) SetupTest() {
	// Save current environment
	suite.originalEnv = make(map[string]string)
	for _, env := range []string{"PORT", "DB_PATH", "GIN_MODE", "ALLOWED_ORIGINS"} {
		suite.originalEnv[env] = os.Getenv(env)
	}
}

// TearDownTest runs after each test in the suite
func (suite *ConfigTestSuite) TearDownTest() {
	// Restore original environment
	for key, value := range suite.originalEnv {
		if value == "" {
			os.Unsetenv(key)
		} else {
			os.Setenv(key, value)
		}
	}
}

// TestLoadDefaultValues tests that default values are applied when optional environment variables are not set
func (suite *ConfigTestSuite) TestLoadDefaultValues() {
	tests := []struct {
		name            string
		envVars         map[string]string
		expectedPort    string
		expectedMode    string
		expectedOrigins []string
	}{
		{
			name: "all_defaults_with_required_db_path",
			envVars: map[string]string{
				"DB_PATH": "/tmp/test.db",
			},
			expectedPort:    "8191",
			expectedMode:    "debug",
			expectedOrigins: []string{"http://localhost"},
		},
		{
			name: "custom_port_other_defaults",
			envVars: map[string]string{
				"DB_PATH": "/tmp/test.db",
				"PORT":    "3000",
			},
			expectedPort:    "3000",
			expectedMode:    "debug",
			expectedOrigins: []string{"http://localhost"},
		},
		{
			name: "custom_mode_other_defaults",
			envVars: map[string]string{
				"DB_PATH":  "/tmp/test.db",
				"GIN_MODE": "release",
			},
			expectedPort:    "8191",
			expectedMode:    "release",
			expectedOrigins: []string{"http://localhost"},
		},
	}

	for _, test := range tests {
		suite.Run(test.name, func() {
			// Set environment variables
			for key, value := range test.envVars {
				os.Setenv(key, value)
			}

			config, err := Load()

			assert.NoError(suite.T(), err)
			assert.NotNil(suite.T(), config)
			assert.Equal(suite.T(), test.expectedPort, config.Server.Port)
			assert.Equal(suite.T(), test.expectedMode, config.Server.Mode)
			assert.Equal(suite.T(), test.expectedOrigins, config.CORS.AllowedOrigins)
		})

		// Reset for next test
		suite.TearDownTest()
		suite.SetupTest()
	}
}

// TestLoadMissingRequiredVariables tests that Load returns an error when required variables are missing
func (suite *ConfigTestSuite) TestLoadMissingRequiredVariables() {
	tests := []struct {
		name         string
		envVars      map[string]string
		expectError  bool
		errorMessage string
	}{
		{
			name:         "missing_db_path",
			envVars:      map[string]string{},
			expectError:  true,
			errorMessage: "DB_PATH",
		},
		{
			name: "empty_db_path",
			envVars: map[string]string{
				"DB_PATH": "",
			},
			expectError:  true,
			errorMessage: "DB_PATH",
		},
		{
			name: "db_path_with_whitespace_only",
			envVars: map[string]string{
				"DB_PATH": "   ",
			},
			expectError:  false, // Empty string check happens before trim
			errorMessage: "",
		},
	}

	for _, test := range tests {
		suite.Run(test.name, func() {
			// Set environment variables
			for key, value := range test.envVars {
				os.Setenv(key, value)
			}

			config, err := Load()

			if test.expectError {
				assert.Error(suite.T(), err)
				assert.Nil(suite.T(), config)
				assert.Contains(suite.T(), err.Error(), test.errorMessage)
			} else {
				assert.NoError(suite.T(), err)
				assert.NotNil(suite.T(), config)
			}
		})

		// Reset for next test
		suite.TearDownTest()
		suite.SetupTest()
	}
}

// TestParseAllowedOrigins tests that CORS origins are correctly parsed from comma-separated strings
func (suite *ConfigTestSuite) TestParseAllowedOrigins() {
	tests := []struct {
		name     string
		input    string
		expected []string
	}{
		{
			name:     "empty_string_returns_default",
			input:    "",
			expected: []string{"http://localhost"},
		},
		{
			name:     "single_origin",
			input:    "http://example.com",
			expected: []string{"http://example.com"},
		},
		{
			name:     "multiple_origins",
			input:    "http://localhost:3000,http://localhost:4000,http://example.com",
			expected: []string{"http://localhost:3000", "http://localhost:4000", "http://example.com"},
		},
		{
			name:     "origins_with_whitespace",
			input:    "http://localhost:3000, http://localhost:4000 , http://example.com",
			expected: []string{"http://localhost:3000", "http://localhost:4000", "http://example.com"},
		},
		{
			name:     "wildcard_origin",
			input:    "*",
			expected: []string{"*"},
		},
		{
			name:     "origins_with_empty_entries",
			input:    "http://localhost:3000,,http://example.com",
			expected: []string{"http://localhost:3000", "http://example.com"},
		},
		{
			name:     "origins_with_trailing_comma",
			input:    "http://localhost:3000,http://example.com,",
			expected: []string{"http://localhost:3000", "http://example.com"},
		},
	}

	for _, test := range tests {
		suite.Run(test.name, func() {
			result := parseAllowedOrigins(test.input)
			assert.Equal(suite.T(), test.expected, result)
		})
	}
}

// TestParseAllowedOriginsWithWhitespace tests that whitespace around origins is properly trimmed
func (suite *ConfigTestSuite) TestParseAllowedOriginsWithWhitespace() {
	tests := []struct {
		name     string
		input    string
		expected []string
	}{
		{
			name:     "leading_whitespace",
			input:    "  http://example.com",
			expected: []string{"http://example.com"},
		},
		{
			name:     "trailing_whitespace",
			input:    "http://example.com  ",
			expected: []string{"http://example.com"},
		},
		{
			name:     "leading_and_trailing_whitespace",
			input:    "  http://example.com  ",
			expected: []string{"http://example.com"},
		},
		{
			name:     "whitespace_between_commas",
			input:    "http://localhost:3000 , http://localhost:4000",
			expected: []string{"http://localhost:3000", "http://localhost:4000"},
		},
		{
			name:     "tabs_and_spaces",
			input:    "\thttp://example.com\t,\thttp://test.com\t",
			expected: []string{"http://example.com", "http://test.com"},
		},
		{
			name:     "whitespace_only_entry",
			input:    "http://example.com,   ,http://test.com",
			expected: []string{"http://example.com", "http://test.com"},
		},
	}

	for _, test := range tests {
		suite.Run(test.name, func() {
			result := parseAllowedOrigins(test.input)
			assert.Equal(suite.T(), test.expected, result)
		})
	}
}

// TestParseAllowedOriginsEmpty tests that default origin is returned when ALLOWED_ORIGINS is empty
func (suite *ConfigTestSuite) TestParseAllowedOriginsEmpty() {
	tests := []struct {
		name     string
		input    string
		expected []string
	}{
		{
			name:     "empty_string",
			input:    "",
			expected: []string{"http://localhost"},
		},
		{
			name:     "whitespace_only",
			input:    "   ",
			expected: []string{"http://localhost"},
		},
		{
			name:     "tabs_only",
			input:    "\t\t\t",
			expected: []string{"http://localhost"},
		},
		{
			name:     "newlines_only",
			input:    "\n\n",
			expected: []string{"http://localhost"},
		},
		{
			name:     "mixed_whitespace",
			input:    " \t \n ",
			expected: []string{"http://localhost"},
		},
	}

	for _, test := range tests {
		suite.Run(test.name, func() {
			result := parseAllowedOrigins(test.input)
			assert.Equal(suite.T(), test.expected, result)
		})
	}
}

// TestParseAllowedOriginsMultiple tests that multiple origins are correctly parsed
func (suite *ConfigTestSuite) TestParseAllowedOriginsMultiple() {
	tests := []struct {
		name          string
		input         string
		expected      []string
		expectedCount int
	}{
		{
			name:          "two_origins",
			input:         "http://localhost:3000,http://localhost:4000",
			expected:      []string{"http://localhost:3000", "http://localhost:4000"},
			expectedCount: 2,
		},
		{
			name:          "three_origins",
			input:         "http://localhost:3000,http://localhost:4000,http://example.com",
			expected:      []string{"http://localhost:3000", "http://localhost:4000", "http://example.com"},
			expectedCount: 3,
		},
		{
			name:          "mixed_protocols",
			input:         "http://example.com,https://secure.com,http://localhost:3000",
			expected:      []string{"http://example.com", "https://secure.com", "http://localhost:3000"},
			expectedCount: 3,
		},
		{
			name:          "different_ports",
			input:         "http://localhost:3000,http://localhost:4000,http://localhost:5000,http://localhost:8080",
			expected:      []string{"http://localhost:3000", "http://localhost:4000", "http://localhost:5000", "http://localhost:8080"},
			expectedCount: 4,
		},
		{
			name:          "with_wildcards",
			input:         "*,http://localhost:3000",
			expected:      []string{"*", "http://localhost:3000"},
			expectedCount: 2,
		},
	}

	for _, test := range tests {
		suite.Run(test.name, func() {
			result := parseAllowedOrigins(test.input)
			assert.Equal(suite.T(), test.expected, result)
			assert.Equal(suite.T(), test.expectedCount, len(result))
		})
	}
}

// TestGetEnv tests that GetEnv returns environment variable value when set, or fallback when not set
func (suite *ConfigTestSuite) TestGetEnv() {
	tests := []struct {
		name     string
		key      string
		envValue string
		fallback string
		expected string
		setEnv   bool
	}{
		{
			name:     "env_var_set_returns_value",
			key:      "TEST_VAR",
			envValue: "test_value",
			fallback: "default_value",
			expected: "test_value",
			setEnv:   true,
		},
		{
			name:     "env_var_not_set_returns_fallback",
			key:      "UNSET_VAR",
			envValue: "",
			fallback: "default_value",
			expected: "default_value",
			setEnv:   false,
		},
		{
			name:     "env_var_empty_string_returns_fallback",
			key:      "EMPTY_VAR",
			envValue: "",
			fallback: "default_value",
			expected: "default_value",
			setEnv:   true,
		},
		{
			name:     "env_var_with_spaces",
			key:      "SPACE_VAR",
			envValue: "  value with spaces  ",
			fallback: "default",
			expected: "  value with spaces  ",
			setEnv:   true,
		},
		{
			name:     "empty_fallback",
			key:      "TEST_VAR",
			envValue: "test_value",
			fallback: "",
			expected: "test_value",
			setEnv:   true,
		},
		{
			name:     "both_empty",
			key:      "EMPTY_VAR",
			envValue: "",
			fallback: "",
			expected: "",
			setEnv:   false,
		},
	}

	for _, test := range tests {
		suite.Run(test.name, func() {
			if test.setEnv {
				os.Setenv(test.key, test.envValue)
			} else {
				os.Unsetenv(test.key)
			}

			result := GetEnv(test.key, test.fallback)
			assert.Equal(suite.T(), test.expected, result)

			// Cleanup
			os.Unsetenv(test.key)
		})
	}
}

// Run the test suite
func TestConfigTestSuite(t *testing.T) {
	suite.Run(t, new(ConfigTestSuite))
}

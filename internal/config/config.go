package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	CORS     CORSConfig
}

type ServerConfig struct {
	Port string
	Mode string
}

type DatabaseConfig struct {
	Path string
}

type CORSConfig struct {
	AllowedOrigins []string
}

// Load loads configuration from environment variables
func Load() (*Config, error) {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8191"
	}

	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		return nil, fmt.Errorf("DB_PATH environment variable is required")
	}

	mode := os.Getenv("GIN_MODE")
	if mode == "" {
		mode = "debug"
	}

	allowedOrigins := parseAllowedOrigins(os.Getenv("ALLOWED_ORIGINS"))

	return &Config{
		Server: ServerConfig{
			Port: port,
			Mode: mode,
		},
		Database: DatabaseConfig{
			Path: dbPath,
		},
		CORS: CORSConfig{
			AllowedOrigins: allowedOrigins,
		},
	}, nil
}

func parseAllowedOrigins(originsStr string) []string {
	trimmed := strings.TrimSpace(originsStr)
	if trimmed == "" {
		return []string{"http://localhost"}
	}

	var origins []string
	for _, origin := range splitAndTrim(trimmed, ",") {
		if origin != "" {
			origins = append(origins, origin)
		}
	}

	return origins
}

func splitAndTrim(s, sep string) []string {
	parts := strings.Split(s, sep)
	result := make([]string, 0, len(parts))
	for _, part := range parts {
		if trimmed := strings.TrimSpace(part); trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}

// GetEnv gets an environment variable with a fallback
func GetEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

// GetEnvAsInt gets an environment variable as integer with a fallback
func GetEnvAsInt(key string, fallback int) int {
	if value := os.Getenv(key); value != "" {
		if intVal, err := strconv.Atoi(value); err == nil {
			return intVal
		}
	}
	return fallback
}

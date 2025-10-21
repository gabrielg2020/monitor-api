package services

// HealthServiceInterface defines methods for health checks
type HealthServiceInterface interface {
	CheckHealth() error
	GetDetailedHealth() (map[string]interface{}, error)
}

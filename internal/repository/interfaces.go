package repository

// HealthRepositoryInterface defines methods for health checks
type HealthRepositoryInterface interface {
	CheckDatabaseConnection() error
	GetDatabaseStats() (map[string]interface{}, error)
	GetTableCounts() (map[string]int, error)
}

var _ HealthRepositoryInterface = (*HealthRepository)(nil)

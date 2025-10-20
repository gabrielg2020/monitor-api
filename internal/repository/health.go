package repository

import (
	"database/sql"
)

type HealthRepository struct {
	db *sql.DB
}

func NewHealthRepository(db *sql.DB) *HealthRepository {
	return &HealthRepository{db: db}
}

// CheckDatabaseConnection verifies the database is accessible
func (repo *HealthRepository) CheckDatabaseConnection() error {
	return repo.db.Ping()
}

// GetDatabaseStats returns database statistics
func (repo *HealthRepository) GetDatabaseStats() (map[string]interface{}, error) {
	stats := repo.db.Stats()

	return map[string]interface{}{
		"open_connections": stats.OpenConnections,
		"in_use":           stats.InUse,
		"idle":             stats.Idle,
		"max_open":         stats.MaxOpenConnections,
	}, nil
}

// GetTableCounts returns record counts for monitoring tables
func (repo *HealthRepository) GetTableCounts() (map[string]int, error) {
	counts := make(map[string]int)

	// Get host count
	var hostCount int
	err := repo.db.QueryRow("SELECT COUNT(*) FROM hosts").Scan(&hostCount)
	if err != nil {
		return nil, err
	}
	counts["hosts"] = hostCount

	// Get metric count
	var metricCount int
	err = repo.db.QueryRow("SELECT COUNT(*) FROM system_metrics").Scan(&metricCount)
	if err != nil {
		return nil, err
	}
	counts["metrics"] = metricCount

	return counts, nil
}

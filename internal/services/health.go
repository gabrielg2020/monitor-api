package services

import (
	"github.com/gabrielg2020/monitor-api/internal/repository"
)

type HealthService struct {
	repo repository.HealthRepositoryInterface
}

func NewHealthService(repo repository.HealthRepositoryInterface) *HealthService {
	return &HealthService{repo: repo}
}

// CheckHealth performs basic health checks
func (service *HealthService) CheckHealth() error {
	// Check database connectivity
	return service.repo.CheckDatabaseConnection()
}

// GetDetailedHealth returns detailed health information
func (service *HealthService) GetDetailedHealth() (map[string]interface{}, error) {
	health := make(map[string]interface{})

	// Check database
	if err := service.repo.CheckDatabaseConnection(); err != nil {
		health["database"] = map[string]interface{}{
			"status": "unhealthy",
			"error":  err.Error(),
		}
		return health, err
	}

	health["database"] = map[string]interface{}{
		"status": "healthy",
	}

	// Get database stats
	stats, err := service.repo.GetDatabaseStats()
	if err == nil {
		health["database_stats"] = stats
	}

	// Get table counts
	counts, err := service.repo.GetTableCounts()
	if err == nil {
		health["table_counts"] = counts
	}

	return health, nil
}

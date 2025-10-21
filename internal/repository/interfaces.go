package repository

import "github.com/gabrielg2020/monitor-api/internal/entities"

// HealthRepositoryInterface defines methods for health checks
type HealthRepositoryInterface interface {
	CheckDatabaseConnection() error
	GetDatabaseStats() (map[string]interface{}, error)
	GetTableCounts() (map[string]int, error)
}

// HostRepositoryInterface defines the contract for host repository operations
type HostRepositoryInterface interface {
	FindAll(limit int) ([]entities.Host, error)
	FindByFilters(params *entities.HostQueryParams) ([]entities.Host, error)
	Create(host *entities.Host) (int64, error)
	Update(id int64, host *entities.Host) error
	UpdateLastSeen(id int64) error
	Delete(id int64) error
}

var _ HealthRepositoryInterface = (*HealthRepository)(nil)
var _ HostRepositoryInterface = (*HostRepository)(nil)

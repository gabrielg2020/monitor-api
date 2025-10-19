package services

import (
	"database/sql"
)

type HealthService struct {
	db *sql.DB
}

func NewHealthService(con *sql.DB) *HealthService {
	return &HealthService{
		db: con,
	}
}

func (service *HealthService) CheckHealth() error {
	// Check database connectivity
	if err := service.db.Ping(); err != nil {
		return err
	}

	// Additional health checks can be added here

	return nil
}

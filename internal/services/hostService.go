package services

import (
	"github.com/gabrielg2020/monitor-api/internal/entities"
	"github.com/gabrielg2020/monitor-api/internal/repository"
)

type HostService struct {
	repo *repository.HostRepository
}

func NewHostService(repo *repository.HostRepository) *HostService {
	return &HostService{repo: repo}
}

// GetHosts retrieves hosts based on query parameters
func (s *HostService) GetHosts(params *entities.HostQueryParams) ([]entities.Host, error) {
	return s.repo.FindByFilters(params)
}

// CreateOrUpdateHost creates a new host or updates if it already exists
func (s *HostService) CreateOrUpdateHost(host *entities.Host) (int64, error) {
	// Business logic: Check if host already exists
	existingHost, err := s.repo.FindByHostnameOrIP(host.Hostname, host.IPAddress)
	if err != nil {
		return 0, err
	}

	// If host doesn't exist, create it
	if existingHost == nil {
		return s.repo.Create(host)
	}

	// If host exists, update it
	err = s.repo.Update(existingHost.ID, host)
	if err != nil {
		return 0, err
	}

	return existingHost.ID, nil
}

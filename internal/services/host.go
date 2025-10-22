package services

import (
	"github.com/gabrielg2020/monitor-api/internal/entities"
	"github.com/gabrielg2020/monitor-api/internal/repository"
)

type HostService struct {
	repo repository.HostRepositoryInterface
}

func NewHostService(repo repository.HostRepositoryInterface) *HostService {
	return &HostService{repo: repo}
}

// CreateHost creates a new host
func (service *HostService) CreateHost(host *entities.Host) (int64, error) {
	return service.repo.Create(host)
}

// GetHosts retrieves hosts based on query parameters
func (service *HostService) GetHosts(params *entities.HostQueryParams) ([]entities.Host, error) {
	return service.repo.FindByFilters(params)
}

// UpdateHost updates an existing host
func (service *HostService) UpdateHost(id int64, host *entities.Host) error {
	return service.repo.Update(id, host)
}

// DeleteHost deletes a host by ID
func (service *HostService) DeleteHost(id int64) error {
	return service.repo.Delete(id)
}

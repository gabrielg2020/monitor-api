package services

import (
	"database/sql"

	"github.com/gabrielg2020/monitor-page/entities"
)

type HostGetService struct {
	db           *sql.DB
	requestQuery *entities.HostQueryParams
}

func NewHostGetService(con *sql.DB) *HostGetService {
	return &HostGetService{
		db: con,
	}
}

func (service *HostGetService) SetQueryParams(params *entities.HostQueryParams) {
	service.requestQuery = params
}

func (service *HostGetService) GetHost() ([]entities.Host, error) {
	querySQL := `
		SELECT id, hostname, ip_address, role
		FROM hosts
		WHERE 1=1` // Dummy WHERE clause for easier appending

	var args []interface{}

	// Add WHERE clauses conditionally
	if service.requestQuery.ID != 0 {
		querySQL += " AND id = ?"
		args = append(args, service.requestQuery.ID)
	}

	if service.requestQuery.Hostname != "" {
		querySQL += " AND hostname = ?"
		args = append(args, service.requestQuery.Hostname)
	}

	if service.requestQuery.IPAddress != "" {
		querySQL += " AND ip_address = ?"
		args = append(args, service.requestQuery.IPAddress)
	}

	rows, err := service.db.Query(querySQL, args...)
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			panic(err)
		}
	}(rows)

	var hosts []entities.Host
	for rows.Next() {
		var host entities.Host
		if err := rows.Scan(
			&host.ID,
			&host.Hostname,
			&host.IPAddress,
			&host.Role,
		); err != nil {
			return nil, err
		}

		hosts = append(hosts, host)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return hosts, nil
}

package services

import (
	"database/sql"

	"github.com/gabrielg2020/monitor-page/entities"
)

type HostPushService struct {
	db           *sql.DB
	requestQuery *entities.HostQueryParams
}

func NewHostPushService(con *sql.DB) *HostPushService {
	return &HostPushService{
		db: con,
	}
}

func (service *HostPushService) SetQueryParams(params *entities.HostQueryParams) {
	service.requestQuery = params
}

func (service *HostPushService) GetHost() (*entities.Host, error) {
	querySQL := `
		SELECT id, hostname, ip_address
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

	var host entities.Host
	if rows.Next() {
		err := rows.Scan(
			&host.ID,
			&host.Hostname,
			&host.IPAddress,
		)
		if err != nil {
			return nil, err
		}
		// Return the first matching host
		return &host, nil
	}

	// No matching host found
	return nil, sql.ErrNoRows
}

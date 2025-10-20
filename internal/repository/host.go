package repository

import (
	"database/sql"
	"errors"
	"time"

	"github.com/gabrielg2020/monitor-api/internal/entities"
)

type HostRepository struct {
	db *sql.DB
}

func NewHostRepository(db *sql.DB) *HostRepository {
	return &HostRepository{db: db}
}

// FindAll retrieves all hosts
func (repo *HostRepository) FindAll(limit int) ([]entities.Host, error) {
	query := `SELECT id, hostname, ip_address, role FROM hosts`

	if limit > 0 {
		query += ` LIMIT ?`
		rows, err := repo.db.Query(query, limit)
		if err != nil {
			return nil, err
		}
		defer rows.Close()
		return repo.scanHosts(rows)
	}

	rows, err := repo.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return repo.scanHosts(rows)
}

// FindByID retrieves a host by ID
func (repo *HostRepository) FindByID(id int64) (*entities.Host, error) {
	querySQL := `
		SELECT id, hostname, ip_address, role
		FROM hosts
		WHERE id = ?`

	var host entities.Host
	err := repo.db.QueryRow(querySQL, id).Scan(
		&host.ID,
		&host.Hostname,
		&host.IPAddress,
		&host.Role,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil // Not found
		}
		return nil, err
	}

	return &host, nil
}

// FindByHostname retrieves a host by hostname
func (repo *HostRepository) FindByHostname(hostname string) (*entities.Host, error) {
	querySQL := `
		SELECT id, hostname, ip_address, role
		FROM hosts
		WHERE hostname = ?`

	var host entities.Host
	err := repo.db.QueryRow(querySQL, hostname).Scan(
		&host.ID,
		&host.Hostname,
		&host.IPAddress,
		&host.Role,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil // Not found
		}
		return nil, err
	}

	return &host, nil
}

// FindByIPAddress retrieves a host by IP address
func (repo *HostRepository) FindByIPAddress(ipAddress string) (*entities.Host, error) {
	querySQL := `
		SELECT id, hostname, ip_address, role
		FROM hosts
		WHERE ip_address = ?`

	var host entities.Host
	err := repo.db.QueryRow(querySQL, ipAddress).Scan(
		&host.ID,
		&host.Hostname,
		&host.IPAddress,
		&host.Role,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil // Not found
		}
		return nil, err
	}

	return &host, nil
}

// FindByHostnameOrIP finds a host by hostname or IP address
func (repo *HostRepository) FindByHostnameOrIP(hostname, ipAddress string) (*entities.Host, error) {
	querySQL := `
		SELECT id, hostname, ip_address, role
		FROM hosts
		WHERE hostname = ? OR ip_address = ?`

	var host entities.Host
	err := repo.db.QueryRow(querySQL, hostname, ipAddress).Scan(
		&host.ID,
		&host.Hostname,
		&host.IPAddress,
		&host.Role,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil // Not found
		}
		return nil, err
	}

	return &host, nil
}

// FindByFilters retrieves hosts based on query parameters
func (repo *HostRepository) FindByFilters(params *entities.HostQueryParams) ([]entities.Host, error) {
	querySQL := `
		SELECT id, hostname, ip_address, role
		FROM hosts
		WHERE 1=1`

	var args []interface{}

	if params.ID != 0 {
		querySQL += " AND id = ?"
		args = append(args, params.ID)
	}

	if params.Hostname != "" {
		querySQL += " AND hostname = ?"
		args = append(args, params.Hostname)
	}

	if params.IPAddress != "" {
		querySQL += " AND ip_address = ?"
		args = append(args, params.IPAddress)
	}

	rows, err := repo.db.Query(querySQL, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return repo.scanHosts(rows)
}

// Create inserts a new host
func (repo *HostRepository) Create(host *entities.Host) (int64, error) {
	timestamp := time.Now().Unix()
	insertSQL := `
		INSERT INTO hosts (hostname, ip_address, role, created_at, last_seen)
		VALUES (?, ?, ?, ?, ?)`

	result, err := repo.db.Exec(insertSQL,
		host.Hostname,
		host.IPAddress,
		host.Role,
		timestamp,
		timestamp,
	)

	if err != nil {
		return 0, err
	}

	return result.LastInsertId()
}

// Update updates an existing host
func (repo *HostRepository) Update(id int64, host *entities.Host) error {
	timestamp := time.Now().Unix()
	updateSQL := `
		UPDATE hosts
		SET role = ?, last_seen = ?
		WHERE id = ?`

	_, err := repo.db.Exec(updateSQL, host.Role, timestamp, id)
	return err
}

// UpdateLastSeen updates the last_seen timestamp
func (repo *HostRepository) UpdateLastSeen(id int64) error {
	timestamp := time.Now().Unix()
	updateSQL := `UPDATE hosts SET last_seen = ? WHERE id = ?`
	_, err := repo.db.Exec(updateSQL, timestamp, id)
	return err
}

// scanHosts is a helper to scan multiple rows into Host slice
func (repo *HostRepository) scanHosts(rows *sql.Rows) ([]entities.Host, error) {
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

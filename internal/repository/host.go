package repository

import (
	"database/sql"
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

// Delete removes a host from the database
func (repo *HostRepository) Delete(id int64) error {
	deleteSQL := `DELETE FROM hosts WHERE id = ?`
	result, err := repo.db.Exec(deleteSQL, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
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

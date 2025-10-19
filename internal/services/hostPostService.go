package services

import (
	"database/sql"
	"errors"
	"time"

	"github.com/gabrielg2020/monitor-api/internal/entities"
)

type HostPostService struct {
	db *sql.DB
}

func NewHostPostService(con *sql.DB) *HostPostService {
	return &HostPostService{
		db: con,
	}
}

func (service *HostPostService) PostHost(host *entities.Host) (int64, error) {
	// First try getting the host by hostname or IP to avoid duplicates
	var existingHostID int64
	timestamp := time.Now().Unix()
	querySQL := `
		SELECT id
		FROM hosts
		WHERE hostname = ? OR ip_address = ?`

	err := service.db.QueryRow(querySQL, host.Hostname, host.IPAddress).Scan(&existingHostID)

	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return 0, err
	}

	// If no existing host, insert new
	if errors.Is(err, sql.ErrNoRows) {
		insertSQL := `
			INSERT INTO hosts (hostname, ip_address, role, created_at, last_seen)
			VALUES (?, ?, ?, ?, ?)`

		result, err := service.db.Exec(insertSQL,
			host.Hostname,
			host.IPAddress,
			host.Role,
			timestamp,
			timestamp,
		)

		if err != nil {
			return 0, err
		}

		lastInsertID, err := result.LastInsertId()
		if err != nil {
			return 0, err
		}

		return lastInsertID, nil
	}

	// If existing host found, update its role and last_seen
	updateSQL := `
		UPDATE hosts
		SET role = ?, last_seen = ?
		WHERE id = ?`

	_, err = service.db.Exec(updateSQL,
		host.Role,
		timestamp,
		existingHostID,
	)

	if err != nil {
		return 0, err
	}

	return existingHostID, nil
}

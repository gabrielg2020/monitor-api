package services

import (
	"database/sql"
	"time"

	"github.com/gabrielg2020/monitor-page/entities"
)

type HostPushService struct {
	db *sql.DB
}

func NewHostPushService(con *sql.DB) *HostPushService {
	return &HostPushService{
		db: con,
	}
}

func (service *HostPushService) PushHost(host *entities.Host) error {
	// Insert into database
	insertSQL := `
	INSERT INTO hosts (
	    hostname, ip_address, role, created_at, last_seen
	) VALUES (?, ?, ?, ?, ?)
	ON CONFLICT(hostname) DO UPDATE SET
		ip_address=excluded.ip_address,
		role=excluded.role,
		last_seen=excluded.last_seen;`

	// Use current time for CreatedAt and LastSeen
	timestamp := time.Now().Unix()

	_, err := service.db.Exec(insertSQL,
		host.Hostname,
		host.IPAddress,
		host.Role,
		timestamp,
		timestamp,
	)

	if err != nil {
		return err
	}

	return nil
}

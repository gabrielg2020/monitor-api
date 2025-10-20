package repository

import (
	"database/sql"
	"log"
)

// closeRows safely closes sql.Rows and logs any errors
func closeRows(rows *sql.Rows) {
	if err := rows.Close(); err != nil {
		log.Printf("failed to close database rows: %v", err)
	}
}

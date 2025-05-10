// db/db.go
// Database initialization and setup

package db

import (
	"database/sql"
	"fmt"

	"ride-sharing/config"
	_ "github.com/go-sql-driver/mysql"
)

// Initialize sets up the MySQL database connection
func Initialize(cfg *config.Config) (*sql.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?parseTime=true",
		cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBName)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %v", err)
	}
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %v", err)
	}
	return db, nil
}
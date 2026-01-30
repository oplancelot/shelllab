package database

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

// MySQLConnection wraps the MySQL database connection
type MySQLConnection struct {
	db *sql.DB
}

// NewMySQLConnection creates a new MySQL database connection
func NewMySQLConnection(dsn string) (*MySQLConnection, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open mysql connection: %w", err)
	}

	// Verify connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping mysql: %w", err)
	}

	// Settings
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(time.Minute * 5)

	return &MySQLConnection{db: db}, nil
}

// Close closes the database connection
func (m *MySQLConnection) Close() error {
	return m.db.Close()
}

// DB returns the underlying sql.DB
func (m *MySQLConnection) DB() *sql.DB {
	return m.db
}

package postgres

import (
	"citary-backend/pkg/constants"
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

// Connection manages the PostgreSQL database connection
type Connection struct {
	DB *sql.DB
}

// NewConnection creates and initializes a new database connection
func NewConnection(databaseURL string) (*Connection, error) {
	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Check connection
	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	// Configure connection pool
	db.SetMaxOpenConns(constants.DatabaseConfig.MaxOpenConnections)
	db.SetMaxIdleConns(constants.DatabaseConfig.MaxIdleConnections)

	log.Println("Connection to PostgreSQL established successfully")

	return &Connection{DB: db}, nil
}

// Close closes the database connection
func (c *Connection) Close() error {
	if c.DB != nil {
		return c.DB.Close()
	}
	return nil
}

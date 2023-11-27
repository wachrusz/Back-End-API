//go:build !exclude_swagger
// +build !exclude_swagger

// Package mydatabase provides database operations functionality.
package mydatabase

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

type Database struct {
	*sql.DB
}

var GlobalDB *Database

func SetDB(db *Database) {
	GlobalDB = db
}

// @Summary Initialize database connection
// @Description Initializes a connection to the database.
// @Tags Database
// @Produce json
// @Param databaseURL query string true "Database URL"
// @Success 200 {string} string "Database connection initialized successfully"
// @Failure 500 {string} string "Error initializing database connection"
// @Router /mydatabase/init [post]
func Init(databaseURL string) (*Database, error) {
	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &Database{db}, nil
}

// @Summary Close database connection
// @Description Closes the connection to the database.
// @Tags Database
// @Produce json
// @Success 200 {string} string "Database connection closed successfully"
// @Failure 500 {string} string "Error closing database connection"
// @Router /mydatabase/close [post]
func (d *Database) Close() {
	if d != nil && d.DB != nil {
		d.DB.Close()
	}
}

// @Summary Execute database query
// @Description Executes a query on the database.
// @Tags Database
// @Produce json
// @Param query body string true "Database query"
// @Success 200 {string} string "Query executed successfully"
// @Failure 400 {string} string "Invalid request payload"
// @Failure 500 {string} string "Error executing database query"
// @Router /mydatabase/exec [post]
func (d *Database) Exec(query string, args ...interface{}) (sql.Result, error) {
	result, err := d.DB.Exec(query, args...)
	if err != nil {
		log.Println("Error executing query:", err)
	}
	return result, err
}

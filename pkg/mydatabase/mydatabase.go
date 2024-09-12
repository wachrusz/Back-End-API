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

func (d *Database) Close() {
	if d != nil && d.DB != nil {
		d.DB.Close()
	}
}

func (d *Database) Exec(query string, args ...interface{}) (sql.Result, error) {
	result, err := d.DB.Exec(query, args...)
	if err != nil {
		log.Println("Error executing query:", err)
	}
	return result, err
}

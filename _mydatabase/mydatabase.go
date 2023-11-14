package mydatabase

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

// Database структура представляет соединение с базой данных.
type Database struct {
	*sql.DB
}

// Init инициализирует соединение с базой данных.
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

// Close закрывает соединение с базой данных.
func (d *Database) Close() {
	if d != nil && d.DB != nil {
		d.DB.Close()
	}
}

// процедура для exec в бд
func (d *Database) Exec(query string, args ...interface{}) (sql.Result, error) {
	result, err := d.DB.Exec(query, args...)
	if err != nil {
		log.Println("Error executing query:", err)
	}
	return result, err
}

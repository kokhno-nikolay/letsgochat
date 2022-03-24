package postgres

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

func NewClient(uri string) (*sql.DB, error) {
	db, err := sql.Open("postgres", uri)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	log.Println("postgres: connected to database has been successfully")
	return db, nil
}

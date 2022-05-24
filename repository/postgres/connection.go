package postgres

import (
	"log"

	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewClient(uri string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(uri), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	log.Println("postgres: connected to database has been successfully")
	return db, nil
}

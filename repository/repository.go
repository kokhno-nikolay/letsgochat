package repository

import (
	"database/sql"

	"github.com/kokhno-nikolay/letsgochat/models"
	"github.com/kokhno-nikolay/letsgochat/repository/postgres"
)

type Users interface {
	FindByUsername(username string) (models.User, error)
	Create(user models.User) error
	GetAllActiveUsers() ([]models.User, error)
	UserExists(username string) (int, error)
}

type Repositories struct {
	Users Users
}

func NewRepositories(db *sql.DB) *Repositories {
	return &Repositories{
		Users: postgres.NewUsersRepo(db),
	}
}

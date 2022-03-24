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
}

type Token interface {
	CheckToken(token string) (bool, error)
	CreateToken(token models.Token) (models.Token, error)
	DeleteToken(token string) error
}

type Repositories struct {
	Users Users
	Token Token
}

func NewRepositories(db *sql.DB) *Repositories {
	return &Repositories{
		Users: postgres.NewUsersRepo(db),
		Token: postgres.NewTokenRepo(db),
	}
}

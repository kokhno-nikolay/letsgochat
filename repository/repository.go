package repository

import (
	"database/sql"

	"github.com/kokhno-nikolay/letsgochat/models"
	"github.com/kokhno-nikolay/letsgochat/repository/postgres"
)

type Users interface {
	FindByID(int) (models.User, error)
	Create(user models.User) error
	GetAllActive() ([]models.User, error)
	SwitchToActive(int) error
	SwitchToInactive(int) error
}

type Messages interface {
	GetAll() ([]models.Message, error)
	Create(message models.Message) error
}

type Repository struct {
	Users    Users
	Messages Messages
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		Users:    postgres.NewUsersRepo(db),
		Messages: postgres.NewMessageRepo(db),
	}
}

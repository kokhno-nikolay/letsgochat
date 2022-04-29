package repository

import (
	"database/sql"
	"github.com/kokhno-nikolay/letsgochat/models"

	"github.com/kokhno-nikolay/letsgochat/repository/postgres"
)

type Users interface {
	Drop() error
	Up() error
	Close()
	FindById(id int) (models.User, error)
	FindByUsername(username string) (models.User, error)
	Create(user models.User) error
	UserExists(username string) (int, error)
	GetAllActive() ([]models.User, error)
	SwitchToActive(userID int) error
	SwitchToInactive(userID int) error
}

type Messages interface {
	GetAll() ([]models.ChatMessage, error)
	Create(message models.Message) error
}

type Repositories struct {
	Users    Users
	Messages Messages
}

func NewRepositories(db *sql.DB) *Repositories {
	return &Repositories{
		Users:    postgres.NewUsersRepo(db),
		Messages: postgres.NewMessagesRepo(db),
	}
}

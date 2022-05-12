package services

import (
	"github.com/kokhno-nikolay/letsgochat/models"
	"github.com/kokhno-nikolay/letsgochat/repository"
)

type Users interface {
	GetActiveUsers() ([]models.User, error)
	FindByUsername(string) (models.User, error)
	FindById(int) (models.User, error)
	Create(models.User) error
	UserExists(string) (bool, error)
	SwitchToActive(userID int) error
	SwitchToInactive(userID int) error
}

type Messages interface {
	GetAll() ([]models.ChatMessage, error)
	Create(message models.Message) error
}

type Services struct {
	Users    Users
	Messages Messages
}

type Deps struct {
	Repos *repository.Repositories
}

func NewServices(deps Deps) *Services {
	usersService := NewUsersService(deps.Repos.Users)
	messagesService := NewMessagesService(deps.Repos.Messages)

	return &Services{
		Users:    usersService,
		Messages: messagesService,
	}
}

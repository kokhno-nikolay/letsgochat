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

type Services struct {
	Users Users
}

type Deps struct {
	Repos *repository.Repositories
}

func NewServices(deps Deps) *Services {
	usersService := NewUsersService(deps.Repos.Users)

	return &Services{
		Users: usersService,
	}
}

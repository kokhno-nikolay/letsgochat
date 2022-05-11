package services

import (
	"github.com/kokhno-nikolay/letsgochat/models"
	"github.com/kokhno-nikolay/letsgochat/repository"
)

type UsersService struct {
	repo repository.Users
}

func NewUsersService(repo repository.Users) *UsersService {
	return &UsersService{
		repo: repo,
	}
}

func (u *UsersService) GetActiveUsers() ([]models.User, error) {
	return u.repo.GetAllActive()
}

func (u *UsersService) FindById(id int) (models.User, error) {
	return u.repo.FindById(id)
}

func (u *UsersService) FindByUsername(username string) (models.User, error) {
	return u.repo.FindByUsername(username)
}

func (u *UsersService) Create(inp models.User) error {
	return u.repo.Create(inp)
}

func (u *UsersService) UserExists(username string) (bool, error) {
	res, err := u.repo.UserExists(username)
	if err != nil {
		return false, err
	}

	if res != 0 {
		return true, nil
	}

	return false, nil
}

func (u *UsersService) SwitchToActive(userID int) error {
	return u.SwitchToActive(userID)
}

func (u *UsersService) SwitchToInactive(userID int) error {
	return u.SwitchToActive(userID)
}

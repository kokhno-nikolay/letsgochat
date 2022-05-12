package services

import (
	"github.com/kokhno-nikolay/letsgochat/models"
	"github.com/kokhno-nikolay/letsgochat/repository"
)

type MessagesService struct {
	repo repository.Messages
}

func NewMessagesService(repo repository.Messages) *MessagesService {
	return &MessagesService{
		repo: repo,
	}
}

func (u *MessagesService) GetAll() ([]models.ChatMessage, error) {
	return u.repo.GetAll()
}

func (u *MessagesService) Create(message models.Message) error {
	return u.repo.Create(message)
}

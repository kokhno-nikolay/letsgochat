package postgres

import (
	"github.com/kokhno-nikolay/letsgochat/models"
	"gorm.io/gorm"
)

type MessagesRepo struct {
	db *gorm.DB
}

func NewMessagesRepo(db *gorm.DB) *MessagesRepo {
	return &MessagesRepo{
		db: db,
	}
}

func (r *MessagesRepo) GetAll() ([]models.ChatMessage, error) {
	messages := make([]models.ChatMessage, 0)

	err := r.db.Find(&messages).Error
	return messages, err
}

func (r *MessagesRepo) Create(message models.Message) error {
	return r.db.Create(&message).Error
}

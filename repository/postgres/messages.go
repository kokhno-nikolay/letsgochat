package postgres

import (
	"context"
	"database/sql"
	"github.com/kokhno-nikolay/letsgochat/models"
	"time"
)

type MessagesRepo struct {
	db *sql.DB
}

func NewMessagesRepo(db *sql.DB) *MessagesRepo {
	return &MessagesRepo{
		db: db,
	}
}

func (r *MessagesRepo) GetAll() ([]models.MessageInp, error) {
	messages := make([]models.MessageInp, 0)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	rows, err := r.db.QueryContext(ctx, "SELECT messages.text, users.username FROM messages JOIN users ON messages.user_id = users.id")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		message := models.MessageInp{}

		err = rows.Scan(&message.Text, &message.Username)
		if err != nil {
			return nil, err
		}

		messages = append(messages, message)
	}

	return messages, nil
}

func (r *MessagesRepo) Create(message models.Message) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := "INSERT INTO messages (text, user_id) VALUES ($1, $2)"
	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, message.Text, message.UserId)
	return err
}

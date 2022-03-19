package postgres

import (
	"context"
	"database/sql"
	"time"

	"github.com/kokhno-nikolay/letsgochat/models"
)

type MessageRepo struct {
	db *sql.DB
}

func NewMessageRepo(db *sql.DB) *MessageRepo {
	return &MessageRepo{
		db: db,
	}
}

func (r *MessageRepo) GetAll() ([]models.Message, error) {
	_, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := "SELECT id, text, user_id FROM messages"
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []models.Message
	for rows.Next() {
		message := models.Message{}
		if err := rows.Scan(&message.Id, &message.Text, &message.UserId); err != nil {
			return nil, err
		}

		messages = append(messages, message)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return messages, nil
}

func (r *MessageRepo) Create(message models.Message) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := "INSERT INTO messages (id, text, user_id) VALUES (?, ?, ?)"
	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, message.Id, message.Text, message.UserId)
	return err
}

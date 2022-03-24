package postgres

import (
	"context"
	"database/sql"
	"github.com/kokhno-nikolay/letsgochat/models"
	"time"
)

type TokenRepo struct {
	db *sql.DB
}

func NewTokenRepo(db *sql.DB) *TokenRepo {
	return &TokenRepo{
		db: db,
	}
}

func (r *TokenRepo) CheckToken(token string) (bool, error) {
	t := models.Token{}

	row := r.db.QueryRow("SELECT uuid, user_id FROM tokens WHERE uuid = $1", token)
	err := row.Scan(&t.UUID, &t.UserId)
	if err != nil {
		return false, err
	}

	if token != t.UUID.String() {
		return false, nil
	}

	return true, nil
}

func (r *TokenRepo) CreateToken(token models.Token) (models.Token, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := "INSERT INTO tokens (uuid, user_id) VALUES ($1, $2)"
	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return models.Token{}, err
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, token.UUID, token.UserId)
	return token, err
}

func (r *TokenRepo) DeleteToken(token string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := "DELETE FROM tokens WHERE uuid = $1"
	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, token)
	return err
}

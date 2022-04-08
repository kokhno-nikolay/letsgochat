package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/jackc/pgerrcode"
	"github.com/lib/pq"
	"time"

	"github.com/kokhno-nikolay/letsgochat/models"
)

type UsersRepo struct {
	db *sql.DB
}

func NewUsersRepo(db *sql.DB) *UsersRepo {
	return &UsersRepo{
		db: db,
	}
}

func (r *UsersRepo) FindByUsername(username string) (models.User, error) {
	user := models.User{}

	row := r.db.QueryRow("SELECT id, username, password, active  FROM users WHERE username = $1", username)
	err := row.Scan(&user.ID, &user.Username, &user.Password, &user.Active)

	return user, err
}

func (r *UsersRepo) Create(user models.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := "INSERT INTO users (username, password, active) VALUES ($1, $2, $3)"
	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, user.Username, user.Password, user.Active)
	if err != nil {
		pqErr := err.(*pq.Error)
		if pqErr != nil && pqErr.Code == pgerrcode.UniqueViolation {
			return fmt.Errorf("%s username already exists", user.Username)
		}
	}

	return err
}

func (r *UsersRepo) UserExists(username string) (bool, error) {
	var exists bool

	row := r.db.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE username = $1)", username)
	err := row.Scan(&exists)

	return exists, err
}

func (r *UsersRepo) GetAllActiveUsers() ([]models.User, error) {
	users := []models.User{}
	rows, err := r.db.Query("SELECT username FROM users JOIN tokens ON tokens.user_id = users.id;")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		user := models.User{}
		if err := rows.Scan(&user.Username); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

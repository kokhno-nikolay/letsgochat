package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/kokhno-nikolay/letsgochat/models"
	"time"

	"github.com/jackc/pgerrcode"
	"github.com/lib/pq"
)

type UsersRepo struct {
	db *sql.DB
}

func NewUsersRepo(db *sql.DB) *UsersRepo {
	return &UsersRepo{
		db: db,
	}
}

func (r *UsersRepo) Drop() error {
	ctx := context.Background()

	query := "DROP TABLE IF EXISTS users"
	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (r *UsersRepo) Up() error {
	ctx := context.Background()

	query :=
		"CREATE TABLE IF NOT EXISTS users (" +
			"id SERIAL PRIMARY KEY," +
			"username VARCHAR(255)," +
			"password VARCHAR(255)," +
			"active BOOLEAN" +
			")"
	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (r *UsersRepo) Close() {
	r.db.Close()
}

func (r *UsersRepo) FindById(id int) (models.User, error) {
	user := models.User{}

	row := r.db.QueryRow("SELECT id, username, password, active  FROM users WHERE id = $1", id)
	err := row.Scan(&user.ID, &user.Username, &user.Password, &user.Active)

	return user, err
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

func (r *UsersRepo) UserExists(username string) (int, error) {
	var count int

	row := r.db.QueryRow("SELECT count(*) from users WHERE username = $1", username)
	err := row.Scan(&count)

	return count, err
}

func (r *UsersRepo) GetAllActive() ([]models.User, error) {
	users := make([]models.User, 0)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	rows, err := r.db.QueryContext(ctx, "SELECT id, username FROM users WHERE active = true")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		user := models.User{}

		err = rows.Scan(&user.ID, &user.Username)
		if err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	return users, nil
}

func (r *UsersRepo) SwitchToActive(userID int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := "UPDATE users SET active = true WHERE id = $1"
	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, userID)
	return err
}

func (r *UsersRepo) SwitchToInactive(userID int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := "UPDATE users SET active = false WHERE id = $1"
	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, userID)
	return err
}

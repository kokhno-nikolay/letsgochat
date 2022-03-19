package postgres

import (
	"context"
	"database/sql"
	"github.com/kokhno-nikolay/letsgochat/models"
	"time"
)

type UsersRepo struct {
	db *sql.DB
}

func NewUsersRepo(db *sql.DB) *UsersRepo {
	return &UsersRepo{
		db: db,
	}
}

func (r *UsersRepo) FindByID(id int) (models.User, error) {
	user := models.User{}

	row := r.db.QueryRow("SELECT username, active  FROM users WHERE id = $1", id)
	err := row.Scan(&user.Username, &user.ID)
	return user, err
}

func (r *UsersRepo) Create(user models.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := "INSERT INTO users (id, username, password, active) VALUES ($1, $2, $3, $4)"
	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, user.ID, user.Username, user.Password, user.Active)
	return err
}

func (r *UsersRepo) GetAllActive() ([]models.User, error) {
	activeUsers := []models.User{}

	rows, err := r.db.Query("SELECT uuid, username, password, active FROM users WHERE active = true")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		activeUser := models.User{}
		if err := rows.Scan(&activeUser.ID, &activeUser.Username, &activeUser.Password, &activeUser.Active); err != nil {
			return nil, err
		}
		activeUsers = append(activeUsers, activeUser)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return activeUsers, err
}

func (r *UsersRepo) SwitchToActive(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := "UPDATE users SET active = true WHERE id = $1"
	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, id)
	return err
}

func (r *UsersRepo) SwitchToInactive(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := "UPDATE users SET active = false WHERE id = $1"
	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, id)
	return err
}

package postgres_test

import (
	"database/sql"
	"log"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"

	"github.com/kokhno-nikolay/letsgochat/models"
	"github.com/kokhno-nikolay/letsgochat/repository"
)

var u = models.User{
	ID:       1,
	Username: "test",
	Password: "213123",
	Active:   false,
}

func NewMock() (*sql.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	return db, mock
}

func TestFindByID(t *testing.T) {
	db, mock := NewMock()
	repo := repository.NewRepository(db)

	query := "SELECT id, username, password, active FROM users WHERE id = \\$1"
	rows := sqlmock.NewRows([]string{"id", "username", "password", "active"}).
		AddRow(u.ID, u.Username, u.Password, u.Active)

	mock.ExpectQuery(query).WithArgs(u.ID).WillReturnRows(rows)

	user, err := repo.Users.FindByID(u.ID)
	assert.NotNil(t, user)
	assert.NoError(t, err)
}

func TestCreate(t *testing.T) {
	db, mock := NewMock()
	repo := repository.NewRepository(db)

	query := "INSERT INTO users \\(id, username, password, active\\) VALUES \\(\\$1, \\$2, \\$3, \\$4\\)"
	prep := mock.ExpectPrepare(query)
	prep.ExpectExec().WithArgs(u.ID, u.Username, u.Password, u.Active).WillReturnResult(sqlmock.NewResult(0, 1))

	err := repo.Users.Create(u)
	assert.NoError(t, err)
}

func TestGetAllActive(t *testing.T) {
	db, mock := NewMock()
	repo := repository.NewRepository(db)

	query := "SELECT id, username, password, active FROM users WHERE active = true"
	rows := sqlmock.NewRows([]string{"id", "username", "password", "active"}).
		AddRow(u.ID, u.Username, u.Password, u.Active)

	mock.ExpectQuery(query).WithArgs(u.ID).WillReturnRows(rows)

	user, err := repo.Users.FindByID(u.ID)
	assert.NotNil(t, user)
	assert.NoError(t, err)
}

func TestSwitchToActive(t *testing.T) {
	db, mock := NewMock()
	repo := repository.NewRepository(db)

	query := "UPDATE users SET active = true WHERE id = \\$1"
	prep := mock.ExpectPrepare(query)
	prep.ExpectExec().WithArgs(u.ID).WillReturnResult(sqlmock.NewResult(0, 1))

	err := repo.Users.SwitchToActive(u.ID)
	assert.NoError(t, err)
}

func TestSwitchToInactive(t *testing.T) {
	db, mock := NewMock()
	repo := repository.NewRepository(db)

	query := "UPDATE users SET active = false WHERE id = \\$1"
	prep := mock.ExpectPrepare(query)
	prep.ExpectExec().WithArgs(u.ID).WillReturnResult(sqlmock.NewResult(0, 1))

	err := repo.Users.SwitchToActive(u.ID)
	assert.NoError(t, err)
}

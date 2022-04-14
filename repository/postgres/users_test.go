package postgres_test

import (
	"database/sql"
	"github.com/stretchr/testify/require"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"

	"github.com/kokhno-nikolay/letsgochat/models"
	"github.com/kokhno-nikolay/letsgochat/repository"
)

var u = models.User{
	ID:       666,
	Username: "test-username",
	Password: "test-password",
	Active:   false,
}

func TestNewUsersRepo(t *testing.T) {
	repos := repository.NewRepositories(&sql.DB{})
	userRepo := repos.Users
	require.IsType(t, userRepo, repos.Users)
}

func TestUsersRepo_TestFindByUsername(t *testing.T) {
	db, mock := NewMock()
	repo := repository.NewRepositories(db)

	query := "SELECT id, username, password, active  FROM users WHERE username = \\\\?"

	rows := sqlmock.NewRows([]string{"id", "username", "password", "active"}).
		AddRow(u.ID, u.Username, u.Password, u.Active)

	mock.ExpectQuery(query).WithArgs(u.Username).WillReturnRows(rows)

	user, err := repo.Users.FindByUsername(u.Username)
	assert.NotNil(t, user)
	assert.NoError(t, err)
}

func TestUsersRepo_TestCreate(t *testing.T) {
	db, mock := NewMock()
	repo := repository.NewRepositories(db)

	query := "INSERT INTO users \\(username, password, active\\) VALUES \\(\\$1, \\$2, \\$3\\)"
	prep := mock.ExpectPrepare(query)
	prep.ExpectExec().WithArgs(u.Username, u.Password, u.Active).WillReturnResult(sqlmock.NewResult(0, 1))

	err := repo.Users.Create(u)
	assert.NoError(t, err)
}

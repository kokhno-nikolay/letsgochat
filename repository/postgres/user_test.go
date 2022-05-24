package postgres_test

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

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
	repos := repository.NewRepositories(&gorm.DB{})
	userRepo := repos.Users
	require.IsType(t, userRepo, repos.Users)
}

func TestUsersRepo_TestFindById(t *testing.T) {
	postgresDB, mock, err := sqlmock.New()
	if err != nil {
		panic(err)
	}
	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: postgresDB,
	}), &gorm.Config{})
	repo := repository.NewRepositories(gormDB)

	query := "SELECT id, username, password, active  FROM users WHERE id = $1 ORDER BY users.id LIMIT 1"

	rows := sqlmock.NewRows([]string{"id", "username", "password", "active"}).
		AddRow(u.ID, u.Username, u.Password, u.Active)

	mock.ExpectQuery(query).WithArgs(u.ID).WillReturnRows(rows)

	user, err := repo.Users.FindById(u.ID)
	assert.NotNil(t, user)
	assert.NoError(t, err)
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

func TestUsersRepo_UserExists(t *testing.T) {
	db, mock := NewMock()
	repo := repository.NewRepositories(db)

	query := "SELECT count(*) from users WHERE username = $1"
	prep := mock.ExpectPrepare(query)
	prep.ExpectExec().WithArgs(u.Username)

	repo.Users.UserExists(u.Username)
}

func TestUsersRepo_GetAllActive(t *testing.T) {
	db, mock := NewMock()
	repo := repository.NewRepositories(db)

	query := "SELECT id, username FROM users WHERE active = true"
	rows := sqlmock.NewRows([]string{"id", "username"}).
		AddRow(u.ID, u.Username)

	mock.ExpectQuery(query).WillReturnRows(rows)

	users, err := repo.Users.GetAllActive()
	assert.NotEmpty(t, users)
	assert.NoError(t, err)
	assert.Len(t, users, 1)
}

func TestUsersRepo_SwitchToActive(t *testing.T) {
	db, mock := NewMock()
	repo := repository.NewRepositories(db)

	query := "UPDATE users SET active = true WHERE id = $1"

	prep := mock.ExpectPrepare(query)
	prep.ExpectExec().WithArgs(u.ID).WillReturnResult(sqlmock.NewResult(0, 1))

	err := repo.Users.SwitchToActive(u.ID)
	assert.Error(t, err)
}

func TestUsersRepo_SwitchToInactive(t *testing.T) {
	db, mock := NewMock()
	repo := repository.NewRepositories(db)

	query := "UPDATE users SET active = false WHERE id = $1"

	prep := mock.ExpectPrepare(query)
	prep.ExpectExec().WithArgs(u.ID).WillReturnResult(sqlmock.NewResult(0, 1))

	err := repo.Users.SwitchToInactive(u.ID)
	assert.Error(t, err)
}

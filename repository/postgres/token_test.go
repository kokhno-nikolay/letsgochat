package postgres_test

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"

	"github.com/kokhno-nikolay/letsgochat/models"
	"github.com/kokhno-nikolay/letsgochat/repository"
)

var token = models.Token{
	UUID:   uuid.New(),
	UserId: 0,
}

func TestCheckToken(t *testing.T) {
	db, mock := NewMock()
	repo := repository.NewRepositories(db)

	query := "SELECT uuid, user_id FROM tokens WHERE uuid = \\$1"
	rows := sqlmock.NewRows([]string{"uuid", "user_id"}).
		AddRow(token.UUID, token.UserId)

	mock.ExpectQuery(query).WithArgs(token.UUID).WillReturnRows(rows)
	ok, err := repo.Token.CheckToken(token.UUID.String())
	if !ok {
		t.Errorf("Received %v (type %v), expected %v (type %v)", true, true, ok, reflect.TypeOf(ok))
	}
	assert.NoError(t, err)
}

func TestCreateToken(t *testing.T) {
	db, mock := NewMock()
	repo := repository.NewRepositories(db)

	query := "INSERT INTO tokens \\(uuid, user_id\\) VALUES \\(\\$1, \\$2\\)"
	prep := mock.ExpectPrepare(query)
	prep.ExpectExec().WithArgs(token.UUID, token.UserId).WillReturnResult(sqlmock.NewResult(0, 1))

	_, err := repo.Token.CreateToken(token)
	assert.NoError(t, err)
}

func TestDeleteToken(t *testing.T) {
	db, mock := NewMock()
	repo := repository.NewRepositories(db)

	query := "DELETE FROM tokens WHERE uuid = \\$1"
	prep := mock.ExpectPrepare(query)
	prep.ExpectExec().WithArgs(token.UUID).WillReturnResult(sqlmock.NewResult(0, 1))

	err := repo.Token.DeleteToken(token.UUID.String())
	assert.NoError(t, err)
}

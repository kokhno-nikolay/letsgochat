package postgres_test

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"

	"github.com/kokhno-nikolay/letsgochat/models"
	"github.com/kokhno-nikolay/letsgochat/repository"
)

var msg = models.Message{
	Id:     666,
	Text:   "test-message",
	UserId: 13903,
}

var msgInp = models.ChatMessage{
	Text:     "Random text message",
	Username: "test",
}

func TestMessagesRepo_GetAll(t *testing.T) {
	postgresDB, mock, err := sqlmock.New()
	if err != nil {
		panic(err)
	}
	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: postgresDB,
	}), &gorm.Config{})
	repo := repository.NewRepositories(gormDB)

	query := "SELECT * FROM messages"
	rows := sqlmock.NewRows([]string{"text", "username"}).
		AddRow(msgInp.Text, msgInp.Username)

	mock.ExpectQuery(query).WillReturnRows(rows)

	msgs, err := repo.Messages.GetAll()
	assert.NoError(t, err)
	assert.Len(t, msgs, 1)
}

func TestMessagesRepo_Create(t *testing.T) {
	db, mock := NewMock()
	repo := repository.NewRepositories(db)

	query := "INSERT INTO messages \\(text, user_id\\) VALUES \\(\\$1, \\$2\\)"
	prep := mock.ExpectPrepare(query)
	prep.ExpectExec().WithArgs(msg.Text, msg.UserId).WillReturnResult(sqlmock.NewResult(0, 1))

	err := repo.Messages.Create(msg)
	assert.NoError(t, err)
}

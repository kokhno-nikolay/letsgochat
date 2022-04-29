package repository_test

import (
	"database/sql"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/kokhno-nikolay/letsgochat/repository"
)

func TestNewRepositories(t *testing.T) {
	repos := repository.NewRepositories(&sql.DB{})
	require.IsType(t, &repository.Repositories{}, repos)
}

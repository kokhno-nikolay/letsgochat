package api_test

import (
	"database/sql"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/kokhno-nikolay/letsgochat/api"
	"github.com/kokhno-nikolay/letsgochat/repository"
)

func TestNewHandler(t *testing.T) {
	repos := repository.NewRepositories(&sql.DB{})
	h := api.NewHandler(api.Deps{repos})
	require.IsType(t, &api.Handler{}, h)
}

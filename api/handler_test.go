package api_test

import (
	"database/sql"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"net/http/httptest"
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

func TestHandler_Init(t *testing.T) {
	db, _, err := sqlmock.New()
	if err != nil {
		assert.Nil(t, err)
	}
	repo := repository.NewRepositories(db)
	handler := api.NewHandler(api.Deps{Repos: repo})

	req := httptest.NewRequest("GET", "/ping", nil)

	w := httptest.NewRecorder()
	handler.Init().ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
}

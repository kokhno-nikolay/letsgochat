package api_test

import (
	"database/sql"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kokhno-nikolay/letsgochat/api"
	"github.com/kokhno-nikolay/letsgochat/repository"
	"github.com/kokhno-nikolay/letsgochat/services"
)

func TestNewHandler(t *testing.T) {
	repos := repository.NewRepositories(&sql.DB{})
	services := services.NewServices(services.Deps{Repos: repos})
	h := api.NewHandler(services)
	require.IsType(t, &api.Handler{}, h)
}

func TestHandler_Init(t *testing.T) {
	db, _, err := sqlmock.New()
	if err != nil {
		assert.Nil(t, err)
	}
	repo := repository.NewRepositories(db)
	services := services.NewServices(services.Deps{Repos: repo})
	handler := api.NewHandler(services)

	req := httptest.NewRequest("GET", "/ping", nil)

	w := httptest.NewRecorder()
	handler.Init().ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
}

package api_test

import (
	"database/sql"
	"github.com/kokhno-nikolay/letsgochat/api"
	"github.com/kokhno-nikolay/letsgochat/repository"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewHandler(t *testing.T) {
	repos := repository.NewRepositories(&sql.DB{})
	h := api.NewHandler(api.Deps{repos})
	require.IsType(t, &api.Handler{repos.Users, repos.Token, 0, ""}, h)
}

func TestHandler_Init(t *testing.T) {
	h := api.NewHandler(api.Deps{})

	router := h.Init()
	ts := httptest.NewServer(router)
	defer ts.Close()

	res, err := http.Get(ts.URL + "/ping")
	if err != nil {
		t.Error(err)
	}

	require.Equal(t, http.StatusOK, res.StatusCode)
}

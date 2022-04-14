package api_test

import (
	"database/sql"
	"net/http"
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
	repos := repository.NewRepositories(&sql.DB{})
	h := api.NewHandler(api.Deps{repos})

	router := h.Init()
	ts := httptest.NewServer(router)
	defer ts.Close()

	res, err := http.Get(ts.URL + "/ping")
	if err != nil {
		t.Error(err)
	}

	require.Equal(t, http.StatusOK, res.StatusCode)
}

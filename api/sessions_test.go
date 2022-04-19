package api_test

import (
	"database/sql"
	"fmt"
	"testing"

	"github.com/kokhno-nikolay/letsgochat/api"
	"github.com/kokhno-nikolay/letsgochat/repository"
)

func TestHandler_CheckUserSession(t *testing.T) {
	repos := repository.NewRepositories(&sql.DB{})
	handler := api.NewHandler(api.Deps{repos})

	for i, tt := range []struct {
		token  string
		userID int
		res    bool
	}{
		{
			"1c7d1048-6263-4baa-aeb2-92f88bbf6486",
			666,
			true,
		},
		{
			"0d7d18e3-2734-414b-a346-89f81b9d4173",
			333,
			false,
		},
		{
			"56863653-3a50-4ffd-b189-9dda062445cb",
			1,
			true,
		},
	} {
		t.Run(fmt.Sprintf("%v", i), func(t *testing.T) {
			if tt.res {
				handler.Sessions[tt.token] = tt.userID
			}

			ok, _ := handler.CheckUserSession(tt.userID)
			if ok != tt.res {
				t.Errorf("want %v; got %v", tt.res, ok)
			}
		})
	}
}

func TestHandler_DeleteSession(t *testing.T) {
	repos := repository.NewRepositories(&sql.DB{})
	handler := api.NewHandler(api.Deps{repos})

	for i, tt := range []struct {
		token  string
		userID int
	}{
		{
			"1c7d1048-6263-4baa-aeb2-92f88bbf6486",
			666,
		},
		{
			"0d7d18e3-2734-414b-a346-89f81b9d4173",
			333,
		},
		{
			"56863653-3a50-4ffd-b189-9dda062445cb",
			1,
		},
	} {
		t.Run(fmt.Sprintf("%v", i), func(t *testing.T) {
			handler.Sessions[tt.token] = tt.userID
			handler.DeleteSession(tt.token)
			_, ok := handler.Sessions[tt.token]

			if ok {
				t.Errorf("error, session was not deleted")
			}
		})
	}
}

func TestHandler_CheckUserToken(t *testing.T) {
	repos := repository.NewRepositories(&sql.DB{})
	handler := api.NewHandler(api.Deps{repos})

	for i, tt := range []struct {
		token  string
		userID int
		res    bool
	}{
		{
			"1c7d1048-6263-4baa-aeb2-92f88bbf6486",
			666,
			true,
		},
		{
			"0d7d18e3-2734-414b-a346-89f81b9d4173",
			333,
			false,
		},
		{
			"56863653-3a50-4ffd-b189-9dda062445cb",
			1,
			true,
		},
	} {
		t.Run(fmt.Sprintf("%v", i), func(t *testing.T) {
			if tt.res {
				handler.Sessions[tt.token] = tt.userID
			}

			ok := handler.CheckUserToken(tt.token)
			if ok != tt.res {
				t.Errorf("want %v; got %v", tt.res, ok)
			}
		})
	}
}

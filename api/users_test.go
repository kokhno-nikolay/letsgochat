package api_test

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"github.com/stretchr/testify/assert"

	"github.com/kokhno-nikolay/letsgochat/api"
	"github.com/kokhno-nikolay/letsgochat/models"
	"github.com/kokhno-nikolay/letsgochat/repository"
	"github.com/kokhno-nikolay/letsgochat/repository/postgres"
)

var (
	repos *repository.Repositories
)

var u = &models.User{
	ID:       1,
	Username: "test_username",
	Password: "test_password",
	Active:   false,
}

var (
	user     = "postgres"
	password = "secret"
	db       = "postgres"
	port     = "5433"
	dialect  = "postgres"
	dsn      = "postgres://%s:%s@localhost:%s/%s?sslmode=disable"
)

func TestMain(m *testing.M) {
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	opts := dockertest.RunOptions{
		Repository: "postgres",
		Tag:        "12.3",
		Env: []string{
			"POSTGRES_USER=" + user,
			"POSTGRES_PASSWORD=" + password,
			"POSTGRES_DB=" + db,
		},
		ExposedPorts: []string{"5432"},
		PortBindings: map[docker.Port][]docker.PortBinding{
			"5432": {
				{HostIP: "0.0.0.0", HostPort: port},
			},
		},
	}

	resource, err := pool.RunWithOptions(&opts)
	if err != nil {
		log.Fatalf("Could not start resource: %s", err.Error())
	}

	dsn = fmt.Sprintf(dsn, user, password, port, db)
	if err = pool.Retry(func() error {
		db, err := postgres.NewClient(dsn)
		repos = repository.NewRepositories(db)
		return err
	}); err != nil {
		log.Fatalf("Could not connect to docker: %s", err.Error())
	}

	defer func() {
		repos.Users.Close()
	}()

	err = repos.Users.Drop()
	if err != nil {
		panic(err)
	}

	err = repos.Users.Up()
	if err != nil {
		panic(err)
	}

	code := m.Run()
	if err := pool.Purge(resource); err != nil {
		log.Fatalf("Could not purge resource: %s", err)
	}

	os.Exit(code)
}

func TestSignUp(t *testing.T) {
	var method = "POST"
	var url = "/user"

	handler := api.NewHandler(api.Deps{Repos: repos})
	router := handler.Init()

	for i, tt := range []struct {
		name       string
		req        string
		statusCode int
	}{
		{
			"success request test",
			fmt.Sprintf("{\"username\": \"%s\", \"password\": \"%s\"}", u.Username, u.Password),
			http.StatusOK,
		},
		{
			"bad request (invalid input request)",
			"",
			http.StatusBadRequest,
		},
		{
			"bad request (username already exists)",
			fmt.Sprintf("{\"username\": \"%s\", \"password\": \"%s\"}", u.Username, u.Password),
			http.StatusBadRequest,
		},
		{
			"bad request (username is too short)",
			fmt.Sprintf("{\"username\": \"%s\", \"password\": \"%s\"}", "", u.Password),
			http.StatusBadRequest,
		},
		{
			"bad request (password is too short)",
			fmt.Sprintf("{\"username\": \"%s\", \"password\": \"%s\"}", u.Username, ""),
			http.StatusBadRequest,
		},
	} {
		t.Run(fmt.Sprintf("%v", i), func(t *testing.T) {
			log.Printf("Test name is \"%s\"", tt.name)

			req, err := http.NewRequest(method, url, strings.NewReader(tt.req))
			assert.NoError(t, err)

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.statusCode, w.Code)
		})
	}
}

func TestSignIn(t *testing.T) {
	var method = "POST"
	var url = "/user/login"

	handler := api.NewHandler(api.Deps{Repos: repos})
	router := handler.Init()

	for i, tt := range []struct {
		name       string
		req        string
		statusCode int
	}{
		{
			"success request test",
			fmt.Sprintf("{\"username\": \"%s\", \"password\": \"%s\"}", u.Username, u.Password),
			http.StatusOK,
		},
		{
			"bad request (invalid input request)",
			"",
			http.StatusBadRequest,
		},
		{
			"bad request (user does not exist)",
			fmt.Sprintf("{\"username\": \"%s\", \"password\": \"%s\"}", "random_username", u.Password),
			http.StatusBadRequest,
		},
		{
			"bad request (username is too short)",
			fmt.Sprintf("{\"username\": \"%s\", \"password\": \"%s\"}", "", u.Password),
			http.StatusBadRequest,
		},
		{
			"bad request (password is too short)",
			fmt.Sprintf("{\"username\": \"%s\", \"password\": \"%s\"}", u.Username, ""),
			http.StatusBadRequest,
		},
		{
			"bad request (invalid password)",
			fmt.Sprintf("{\"username\": \"%s\", \"password\": \"%s\"}", u.Username, u.Password+"1"),
			http.StatusBadRequest,
		},
	} {
		t.Run(fmt.Sprintf("%v", i), func(t *testing.T) {
			log.Println(tt.name)

			if err := repos.Users.Create(models.User{Username: u.Username, Password: u.Password}); err != nil {
				assert.NoError(t, err)
			}

			req, err := http.NewRequest(method, url, strings.NewReader(tt.req))
			assert.NoError(t, err)

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.statusCode, w.Code)
		})
	}
}

func TestGetActiveUsers(t *testing.T) {
	var method = "GET"
	var url = "/user/active"

	handler := api.NewHandler(api.Deps{Repos: repos})
	router := handler.Init()

	if err := repos.Users.Create(models.User{Username: u.Username, Password: u.Password}); err != nil {
		assert.NoError(t, err)
	}

	req, err := http.NewRequest(method, url, nil)
	assert.NoError(t, err)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

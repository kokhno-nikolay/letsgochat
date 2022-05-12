package api_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/kokhno-nikolay/letsgochat/api"
	"github.com/kokhno-nikolay/letsgochat/models"
	"github.com/kokhno-nikolay/letsgochat/services"
	"github.com/kokhno-nikolay/letsgochat/services/mocks"
)

var u = models.User{
	ID:       1,
	Username: "test_username",
	Password: "test_password",
	Active:   false,
}

func TestHandler_SignUp(t *testing.T) {
	var method = "POST"
	var url = "/user"

	c := gomock.NewController(t)
	defer c.Finish()

	userRepo := mocks.NewMockUsers(c)
	userRepo.EXPECT().UserExists(u.Username).Return(false, nil)
	userRepo.EXPECT().Create(u).Return(nil)

	servs := &services.Services{Users: userRepo}
	handler := api.NewHandler(servs)
	router := handler.Init()

	req, err := http.NewRequest(method, url, nil)
	assert.NoError(t, err)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestHandler_SignIn(t *testing.T) {
	var method = "POST"
	var url = "/user/login"

	c := gomock.NewController(t)
	defer c.Finish()

	userRepo := mocks.NewMockUsers(c)
	userRepo.EXPECT().UserExists(u.Username).Return(true, nil)
	userRepo.EXPECT().FindByUsername(u.Username).Return(u, nil)
	userRepo.EXPECT().SwitchToActive(u.ID).Return(nil)

	servs := &services.Services{Users: userRepo}
	handler := api.NewHandler(servs)
	router := handler.Init()

	req, err := http.NewRequest(method, url, nil)
	assert.NoError(t, err)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestHandler_GetActiveUsers(t *testing.T) {
	var users []models.User
	var method = "GET"
	var url = "/user/active"

	c := gomock.NewController(t)
	defer c.Finish()

	userRepo := mocks.NewMockUsers(c)
	users = append(users, u)
	userRepo.EXPECT().GetActiveUsers().Return(users, nil)

	servs := &services.Services{Users: userRepo}
	handler := api.NewHandler(servs)
	router := handler.Init()

	req, err := http.NewRequest(method, url, nil)
	assert.NoError(t, err)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

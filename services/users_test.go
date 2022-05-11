package services_test

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"

	"github.com/kokhno-nikolay/letsgochat/models"
	"github.com/kokhno-nikolay/letsgochat/services/mocks"
)

var u = models.User{
	ID:       666,
	Username: "test-username",
	Password: "test-password",
	Active:   false,
}

func TestUsersService_GetActiveUsers(t *testing.T) {
	var users []models.User

	t.Helper()

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	userRepo := mocks.NewMockUsers(mockCtrl)
	users = append(users, u)
	userRepo.EXPECT().GetActiveUsers().Return(users, nil)

	res, err := userRepo.GetActiveUsers()
	require.Equal(t, u, res[0])
	require.NoError(t, err)
}

func TestUsersService_FindById(t *testing.T) {
	t.Helper()

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	userRepo := mocks.NewMockUsers(mockCtrl)

	userRepo.EXPECT().FindById(u.ID).Return(u, nil)

	user, err := userRepo.FindById(u.ID)
	require.Equal(t, u, user)
	require.NoError(t, err)
}

func TestUsersService_FindByUsername(t *testing.T) {
	t.Helper()

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	userRepo := mocks.NewMockUsers(mockCtrl)

	userRepo.EXPECT().FindByUsername(u.Username).Return(u, nil)

	user, err := userRepo.FindByUsername(u.Username)
	require.Equal(t, u, user)
	require.NoError(t, err)
}

func TestUsersService_Create(t *testing.T) {
	t.Helper()

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	userRepo := mocks.NewMockUsers(mockCtrl)
	userRepo.EXPECT().Create(u).Return(nil)

	err := userRepo.Create(u)
	require.NoError(t, err)
}

func TestUsersService_UserExists(t *testing.T) {
	t.Helper()

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	userRepo := mocks.NewMockUsers(mockCtrl)
	userRepo.EXPECT().UserExists(u.Username).Return(true, nil)

	exists, err := userRepo.UserExists(u.Username)
	require.Equal(t, true, exists)
	require.NoError(t, err)
}

func TestUsersService_SwitchToActive(t *testing.T) {
	t.Helper()

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	userRepo := mocks.NewMockUsers(mockCtrl)
	userRepo.EXPECT().SwitchToActive(u.ID).Return(nil)

	err := userRepo.SwitchToActive(u.ID)
	require.NoError(t, err)
}

func TestUsersService_SwitchToInactive(t *testing.T) {
	t.Helper()

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	userRepo := mocks.NewMockUsers(mockCtrl)
	userRepo.EXPECT().SwitchToInactive(u.ID).Return(nil)

	err := userRepo.SwitchToInactive(u.ID)
	require.NoError(t, err)
}

package services_test

import (
	"github.com/kokhno-nikolay/letsgochat/models"
	"github.com/kokhno-nikolay/letsgochat/services/mocks"
	"github.com/stretchr/testify/require"
	"testing"

	"github.com/golang/mock/gomock"
)

var u = models.User{
	ID:       666,
	Username: "test-username",
	Password: "test-password",
	Active:   false,
}

func TestSignUp(t *testing.T) {
	t.Helper()

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	userRepo := mocks.NewMockUsers(mockCtrl)

	userRepo.EXPECT().SignUp(gomock.Any())

	err := userRepo.SignUp(u)
	require.NoError(t, err)
}

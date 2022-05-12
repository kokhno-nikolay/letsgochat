package services

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"

	"github.com/kokhno-nikolay/letsgochat/models"
	"github.com/kokhno-nikolay/letsgochat/services/mocks"
)

func TestMessagesService_GetAll(t *testing.T) {
	var messages []models.ChatMessage

	t.Helper()

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	messageRepo := mocks.NewMockMessages(mockCtrl)

	messageRepo.EXPECT().GetAll().Return(messages, nil)

	res, err := messageRepo.GetAll()
	require.Equal(t, messages, res)
	require.NoError(t, err)
}

func TestMessagesService_Create(t *testing.T) {
	var msg models.Message = models.Message{UserId: 1, Text: "test"}

	t.Helper()

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	messageRepo := mocks.NewMockMessages(mockCtrl)

	messageRepo.EXPECT().Create(msg).Return(nil)

	err := messageRepo.Create(msg)
	require.NoError(t, err)
}

package wallet

import (
	"avito-shop-service/internal/models"
	"avito-shop-service/internal/usecases/wallet/mock"
	"context"
	"errors"
	"testing"

	"log/slog"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestSendCoin(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSendCoiner := mock.NewMockSendCoiner(ctrl)
	logger := slog.Default()
	usecase := New(logger, mockSendCoiner)
	ctx := context.Background()

	const (
		userID = "user123"
		toUser = "receiver123"
		recvID = "receiver456"
		amount = 50
	)

	t.Run("receiver not found", func(t *testing.T) {
		mockSendCoiner.EXPECT().GetUserByUsername(toUser).Return(nil, nil)
		err := usecase.SendCoin(ctx, toUser, userID, amount)
		assert.ErrorIs(t, err, ErrReceiverNotExist)
	})

	t.Run("error fetching receiver", func(t *testing.T) {
		mockSendCoiner.EXPECT().GetUserByUsername(toUser).Return(nil, errors.New("db error"))
		err := usecase.SendCoin(ctx, toUser, userID, amount)
		assert.Error(t, err)
	})

	t.Run("insufficient funds", func(t *testing.T) {
		mockSendCoiner.EXPECT().GetUserByUsername(toUser).Return(&models.User{ID: recvID}, nil)
		mockSendCoiner.EXPECT().GetUserBalance(ctx, userID).Return(30, nil)
		err := usecase.SendCoin(ctx, toUser, userID, amount)
		assert.ErrorIs(t, err, ErrInsufficientFunds)
	})

	t.Run("error fetching balance", func(t *testing.T) {
		mockSendCoiner.EXPECT().GetUserByUsername(toUser).Return(&models.User{ID: recvID}, nil)
		mockSendCoiner.EXPECT().GetUserBalance(ctx, userID).Return(0, errors.New("db error"))
		err := usecase.SendCoin(ctx, toUser, userID, amount)
		assert.Error(t, err)
	})

	t.Run("transaction error", func(t *testing.T) {
		mockSendCoiner.EXPECT().GetUserByUsername(toUser).Return(&models.User{ID: recvID}, nil)
		mockSendCoiner.EXPECT().GetUserBalance(ctx, userID).Return(100, nil)
		mockSendCoiner.EXPECT().SendCoinToUser(ctx, recvID, userID, amount).Return(errors.New("transfer failed"))
		err := usecase.SendCoin(ctx, toUser, userID, amount)
		assert.Error(t, err)
	})

	t.Run("successful transaction", func(t *testing.T) {
		mockSendCoiner.EXPECT().GetUserByUsername(toUser).Return(&models.User{ID: recvID}, nil)
		mockSendCoiner.EXPECT().GetUserBalance(ctx, userID).Return(100, nil)
		mockSendCoiner.EXPECT().SendCoinToUser(ctx, recvID, userID, amount).Return(nil)
		err := usecase.SendCoin(ctx, toUser, userID, amount)
		assert.NoError(t, err)
	})
}

package shop

import (
	"avito-shop-service/internal/models"
	"avito-shop-service/internal/usecases/shop/mock"
	"context"
	"errors"
	"log/slog"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestBuyItem(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockItemBuyer := mock.NewMockItemBuyer(ctrl)
	logger := slog.Default()
	usecase := New(logger, mockItemBuyer)

	ctx := context.Background()
	userID := "user123"
	itemName := "Sword"
	item := &models.Item{ID: "item123", Name: itemName, Price: 100}

	t.Run("Success - Buy Item", func(t *testing.T) {
		mockItemBuyer.EXPECT().GetItemByName(ctx, itemName).Return(item, nil)
		mockItemBuyer.EXPECT().GetUserBalance(ctx, userID).Return(200, nil)
		mockItemBuyer.EXPECT().AddItemToInventory(ctx, userID, item).Return(nil)

		err := usecase.BuyItem(ctx, userID, itemName)
		assert.NoError(t, err)
	})

	t.Run("Failure - Item Not Found", func(t *testing.T) {
		mockItemBuyer.EXPECT().GetItemByName(ctx, itemName).Return(nil, errors.New("item not found"))

		err := usecase.BuyItem(ctx, userID, itemName)
		assert.Error(t, err)
	})

	t.Run("Failure - Insufficient Funds", func(t *testing.T) {
		mockItemBuyer.EXPECT().GetItemByName(ctx, itemName).Return(item, nil)
		mockItemBuyer.EXPECT().GetUserBalance(ctx, userID).Return(50, nil)

		err := usecase.BuyItem(ctx, userID, itemName)
		assert.ErrorIs(t, err, ErrInsufficientFunds)
	})

	t.Run("Failure - Error Adding to Inventory", func(t *testing.T) {
		mockItemBuyer.EXPECT().GetItemByName(ctx, itemName).Return(item, nil)
		mockItemBuyer.EXPECT().GetUserBalance(ctx, userID).Return(200, nil)
		mockItemBuyer.EXPECT().AddItemToInventory(ctx, userID, item).Return(errors.New("db error"))

		err := usecase.BuyItem(ctx, userID, itemName)
		assert.Error(t, err)
	})

	t.Run("Failure - Error Fetching Balance", func(t *testing.T) {
		mockItemBuyer.EXPECT().GetItemByName(ctx, itemName).Return(item, nil)
		mockItemBuyer.EXPECT().GetUserBalance(ctx, userID).Return(0, errors.New("balance fetch error"))

		err := usecase.BuyItem(ctx, userID, itemName)
		assert.Error(t, err)
	})
}

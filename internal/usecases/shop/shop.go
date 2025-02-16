package shop

import (
	"avito-shop-service/internal/models"
	"context"
	"errors"
	"fmt"
	"log/slog"
)

const (
	errFetchItem    = "failed to fetch item"
	errFetchBalance = "failed to get user balance"
	errTransaction  = "failed to complete transaction"
)

var (
	ErrInsufficientFunds = errors.New("insufficient balance to buy item")
)

//go:generate mockgen -source=shop.go -destination=mock/shop_mock.go -package=mock ItemBuyer
type ItemBuyer interface {
	AddItemToInventory(context.Context, string, *models.Item) error
	GetItemByName(context.Context, string) (*models.Item, error)
	GetUserBalance(context.Context, string) (int, error)
}

type Useacase struct {
	itemBuyer ItemBuyer
	log       *slog.Logger
}

func New(log *slog.Logger, itemBuyer ItemBuyer) *Useacase {
	return &Useacase{
		itemBuyer: itemBuyer,
		log:       log,
	}
}

func (u *Useacase) BuyItem(ctx context.Context, userID, itemName string) error {
	const op = "usecases.shop.BuyItem"

	item, err := u.itemBuyer.GetItemByName(ctx, itemName)
	if err != nil {
		return fmt.Errorf("%s: %s: %w", op, errFetchItem, err)
	}

	balance, err := u.itemBuyer.GetUserBalance(ctx, userID)
	if err != nil {
		return fmt.Errorf("%s: %s: %w", op, errFetchBalance, err)
	}
	if balance < item.Price {
		return fmt.Errorf("%s: %w", op, ErrInsufficientFunds)
	}

	if err := u.itemBuyer.AddItemToInventory(ctx, userID, item); err != nil {
		return fmt.Errorf("%s: %s: %w", op, errTransaction, err)
	}

	return nil
}

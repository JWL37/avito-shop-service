package user

import (
	"avito-shop-service/internal/models"
	"context"
	"fmt"
	"log/slog"
)

const (
	errGetInventory = "failed to get inventory"
	errGetBalance   = "failed to get user balance"
	errGetHistory   = "failed to get coin history"
)

//go:generate mockgen -source=info.go -destination=mock/info_mock.go -package=mock UserInfoGetter
type UserInfoGetter interface {
	GetUserBalance(context.Context, string) (int, error)
	GetUserInventory(context.Context, string) ([]models.InventoryItem, error)
	GetUserTransactions(context.Context, string) (models.CoinHistory, error)
}

type infoUseacase struct {
	userInfoGetter UserInfoGetter
	log            *slog.Logger
}

func NewInfo(log *slog.Logger, userInfoGetter UserInfoGetter) *infoUseacase {
	return &infoUseacase{
		userInfoGetter: userInfoGetter,
		log:            log,
	}
}

func (u *infoUseacase) GetUserInfo(ctx context.Context, userID string) (models.UserInfo, error) {
	const op = "usecases.user.GetUserInfo"

	balance, err := u.userInfoGetter.GetUserBalance(ctx, userID)
	if err != nil {

		return models.UserInfo{}, fmt.Errorf("%s: %s: %w", op, errGetBalance, err)
	}

	inventory, err := u.userInfoGetter.GetUserInventory(ctx, userID)
	if err != nil {

		return models.UserInfo{}, fmt.Errorf("%s: %s: %w", op, errGetInventory, err)
	}

	history, err := u.userInfoGetter.GetUserTransactions(ctx, userID)
	if err != nil {

		return models.UserInfo{}, fmt.Errorf("%s: %s: %w", op, errGetHistory, err)
	}

	return models.UserInfo{
		Coins:       balance,
		Inventory:   inventory,
		CoinHistory: history,
	}, nil
}

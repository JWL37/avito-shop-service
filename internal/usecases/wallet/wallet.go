package wallet

import (
	"avito-shop-service/internal/models"
	"context"
	"errors"
	"fmt"
	"log/slog"
)

const (
	errGetUser     = "failed to get user"
	errGetBalance  = "failed to get user balance"
	errTransaction = "failed to complete transaction"
)

var (
	ErrReceiverNotExist  = errors.New("not exist receiver")
	ErrInsufficientFunds = errors.New("insufficient balance to send coin")
)

//go:generate mockgen -source=wallet.go -destination=mock/wallet_mock.go -package=mock SendCoiner
type SendCoiner interface {
	SendCoinToUser(context.Context, string, string, int) error
	GetUserBalance(context.Context, string) (int, error)
	GetUserByUsername(string) (*models.User, error)
}

type Useacase struct {
	sendCoiner SendCoiner
	log        *slog.Logger
}

func New(log *slog.Logger, sendCoiner SendCoiner) *Useacase {
	return &Useacase{
		sendCoiner: sendCoiner,
		log:        log,
	}
}

func (u *Useacase) SendCoin(ctx context.Context, toUser string, userID string, amount int) error {
	const op = "usecases.wallet.SendCoin"

	receiver, err := u.sendCoiner.GetUserByUsername(toUser)
	if err != nil {

		return fmt.Errorf("%s: %s: %w", op, errGetUser, err)
	}
	if receiver == nil {
		return fmt.Errorf("%s: %w", op, ErrReceiverNotExist)
	}

	balance, err := u.sendCoiner.GetUserBalance(ctx, userID)
	if err != nil {
		return fmt.Errorf("%s: %s: %w", op, errGetBalance, err)
	}

	if balance < amount {
		return fmt.Errorf("%s: %w", op, ErrInsufficientFunds)
	}

	if err := u.sendCoiner.SendCoinToUser(ctx, receiver.ID, userID, amount); err != nil {
		return fmt.Errorf("%s: %s: %w", op, errTransaction, err)

	}
	return nil
}

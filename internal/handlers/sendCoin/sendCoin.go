package sendcoin

import (
	"avito-shop-service/internal/lib/handlers/response"
	"avito-shop-service/internal/middleware"
	"avito-shop-service/internal/usecases/wallet"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
)

const (
	warnGetUserID        = "failed to get user ID from context"
	messageUnauthorized  = "Unauthorized"
	errSendCoin          = "failed to send coin"
	ErrInsufficientFunds = "insufficient balance to send coin"
	successSendCoins     = "send coins successfully"
)

type Request struct {
	ToUser string `json:"toUser"`
	Amount int    `json:"amount"`
}

//go:generate mockgen -source=sendCoin.go -destination=mock/sendCoin_mock.go -package=sendcoin CoinSender
type CoinSender interface {
	SendCoin(context.Context, string, string, int) error
}

func New(log *slog.Logger, coinSender CoinSender) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.sendCoin.New"

		userID, ok := middleware.GetUserID(r)
		if !ok {
			log.Warn(warnGetUserID)

			response.RespondWithError(w, log, http.StatusUnauthorized, messageUnauthorized)
		}

		log := log.With(
			slog.String("op", op),
			slog.String("user_id", userID),
		)

		var req Request
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			log.Error(fmt.Errorf("failed to decode request body: Error: %w", err).Error())

			response.RespondWithError(w, log, http.StatusBadRequest, "invalid request")
			return
		}
		if err := coinSender.SendCoin(r.Context(), req.ToUser, userID, req.Amount); err != nil {
			log.Error(fmt.Errorf("%s Error: %w", errSendCoin, err).Error())

			if errors.Is(err, wallet.ErrInsufficientFunds) {
				response.RespondWithError(w, log, http.StatusBadRequest, ErrInsufficientFunds)
				return
			}

			response.RespondWithError(w, log, http.StatusInternalServerError, errSendCoin)
			return
		}

		w.WriteHeader(http.StatusOK)

		log.Info(successSendCoins)
	}
}

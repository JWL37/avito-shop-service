package buyitem

import (
	"avito-shop-service/internal/lib/handlers/response"
	"avito-shop-service/internal/middleware"
	"avito-shop-service/internal/usecases/shop"
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
)

const (
	warnGetUserID        = "failed to get user ID from context"
	messageUnauthorized  = "Unauthorized"
	errBuyItem           = "failed to buy item"
	ErrInsufficientFunds = "insufficient balance to buy item"
	successBuy           = "item purchased successfully"
	paramItem            = "item"
)

//go:generate mockgen -source=buyItem.go -destination=mock/buyItem_mock.go -package=mock ItemBuyer
type ItemBuyer interface {
	BuyItem(context.Context, string, string) error
}

func New(log *slog.Logger, itemBuyer ItemBuyer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.buyItem.New"

		userID, ok := middleware.GetUserID(r)
		if !ok {
			log.Warn("%s: %s", op, warnGetUserID)

			response.RespondWithError(w, log, http.StatusUnauthorized, messageUnauthorized)
		}

		log := log.With(
			slog.String("op", op),
			slog.String("user_id", userID),
		)

		itemName := chi.URLParam(r, paramItem)
		log.Info("Processing item purchase", slog.String(paramItem, itemName))

		if err := itemBuyer.BuyItem(r.Context(), userID, itemName); err != nil {
			log.Error(fmt.Errorf("%s Error: %w", errBuyItem, err).Error())

			if errors.Is(err, shop.ErrInsufficientFunds) {
				response.RespondWithError(w, log, http.StatusBadRequest, ErrInsufficientFunds)
				return
			}

			response.RespondWithError(w, log, http.StatusInternalServerError, errBuyItem)
			return
		}

		w.WriteHeader(http.StatusOK)

		log.Info(successBuy)
	}
}

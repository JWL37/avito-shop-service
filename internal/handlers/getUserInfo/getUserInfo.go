package getuserinfo

import (
	"avito-shop-service/internal/lib/handlers/response"
	"avito-shop-service/internal/middleware"
	"avito-shop-service/internal/models"
	"context"
	"log/slog"
	"net/http"
)

const (
	warnGetUserID       = "failed to get user ID from context"
	messageUnauthorized = "Unauthorized"
	successGetUserInfo  = "get user info successfully"
)

//go:generate mockgen -source=getUserInfo.go -destination=mock/getUserInfo_mock.go -package=mock UserInfoGetter
type UserInfoGetter interface {
	GetUserInfo(context.Context, string) (models.UserInfo, error)
}

func New(log *slog.Logger, userInfoGetter UserInfoGetter) http.HandlerFunc {
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
		info, err := userInfoGetter.GetUserInfo(r.Context(), userID)
		if err != nil {
			// log.Error(fmt.Errorf("%s Error: %w", errSendCoin, err).Error())

			// response.RespondWithError(w, log, http.StatusInternalServerError, errSendCoin)
			return
		}
		response.RespondWithJSON(w, log, http.StatusOK, info)

		log.Info(successGetUserInfo)
	}
}

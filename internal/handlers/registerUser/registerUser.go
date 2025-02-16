package registeruser

import (
	"avito-shop-service/internal/lib/handlers/response"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
)

type Request struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type ResponseOK struct {
	Token string `json:"token"`
}

//go:generate mockgen -source=registerUser.go -destination=mock/registerUser_mock.go -package=registeruser UserRegistrater
type UserRegistrater interface {
	Register(string, string) (string, error)
}

func New(log *slog.Logger, userRegistrater UserRegistrater) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.registerUser.New"

		log := log.With(
			slog.String("op", op),
		)

		var req Request
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			log.Error(fmt.Errorf("failed to decode request body: Error: %w", err).Error())

			response.RespondWithError(w, log, http.StatusBadRequest, "invalid request")
			return
		}

		token, err := userRegistrater.Register(req.Username, req.Password)
		if err != nil {
			log.Error(fmt.Errorf("failed registration: Error: %w", err).Error())

			response.RespondWithError(w, log, http.StatusBadRequest, "invalid username or password")
			return
		}
		response.RespondWithJSON(w, log, http.StatusOK, ResponseOK{Token: token})

	}
}

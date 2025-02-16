package app

import (
	"avito-shop-service/internal/config"
	buyitem "avito-shop-service/internal/handlers/buyItem"
	getuserinfo "avito-shop-service/internal/handlers/getUserInfo"
	registeruser "avito-shop-service/internal/handlers/registerUser"
	sendcoin "avito-shop-service/internal/handlers/sendCoin"
	custommiddleware "avito-shop-service/internal/middleware"
	"avito-shop-service/internal/storage/postgresql"
	"avito-shop-service/internal/usecases/shop"
	"avito-shop-service/internal/usecases/user"
	"avito-shop-service/internal/usecases/wallet"
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type App struct {
	httpServer *http.Server
	log        *slog.Logger
	storage    *postgresql.Storage
}

// NewApp создает и инициализирует приложение
func NewApp(log *slog.Logger, cfg *config.Config) *App {
	storage, err := postgresql.ConnectAndNew(log, &cfg.Database)
	if err != nil {
		log.Error("Failed to create DB:")
		os.Exit(1)
	}

	authUsecase := user.New(log, storage, cfg.Secret)
	shopUsecase := shop.New(log, storage)
	walletUsecase := wallet.New(log, storage)
	infoUsecase := user.NewInfo(log, storage)

	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	router.Route("/api", func(r chi.Router) {
		r.Use(custommiddleware.Auth(log, cfg.Secret))
		r.Get("/buy/{item}", buyitem.New(log, shopUsecase))
		r.Post("/sendCoin", sendcoin.New(log, walletUsecase))
		r.Get("/info", getuserinfo.New(log, infoUsecase))
	})
	router.Post("/api/auth", registeruser.New(log, authUsecase))

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%s", cfg.ServerPort),
		Handler:      router,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	return &App{
		httpServer: srv,
		log:        log,
		storage:    storage,
	}
}

// Run запускает HTTP сервер
func (a *App) Run() error {
	a.log.Info("Starting server ", slog.String("port", a.httpServer.Addr))
	return a.httpServer.ListenAndServe()
}

// Shutdown останавливает сервер
func (a *App) Shutdown(ctx context.Context) error {
	a.log.Info("Shutting down server...")
	err := a.httpServer.Shutdown(ctx)
	if err != nil {
		return err
	}

	a.storage.Stop()
	a.log.Info("Database connection closed.")

	return nil
}

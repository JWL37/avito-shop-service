package postgresql

import (
	"avito-shop-service/internal/config"
	"database/sql"
	"fmt"
	"log/slog"

	_ "github.com/jackc/pgx/v5/stdlib"
)

const UnableToConnectDatabase = "Unable to connect to database: "

type storage struct {
	DB  *sql.DB
	log *slog.Logger
}

func ConnectAndNew(log *slog.Logger, cfg *config.DatabaseConfig) (*storage, error) {
	const op = "storage.postgresql.New"

	log = log.With(
		slog.String("op", op),
	)

	dsn := getDSN(cfg)

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		log.Error(UnableToConnectDatabase, "error", err)
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		log.Error(UnableToConnectDatabase, "error", err)
		return nil, err
	}

	db.SetMaxOpenConns(40)
	db.SetMaxIdleConns(40)

	storage := &storage{
		DB:  db,
		log: log,
	}

	return storage, nil
}

func NewRep(db *sql.DB, log *slog.Logger) *storage {
	return &storage{
		DB:  db,
		log: log,
	}
}

func getDSN(cfg *config.DatabaseConfig) string {
	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s target_session_attrs=read-write", cfg.Host, cfg.User, cfg.Password, cfg.Name, cfg.Port)
}

func (s *storage) Stop() {
	s.DB.Close()
}

package postgresql

import (
	"avito-shop-service/internal/models"
	"database/sql"
	"errors"
	"fmt"
)

const (
	errFailedToSetBalance = " set default balance"
	errCreateUser         = "failed to create user"

	defaultBalance = 1000

	querryCreateUser       = `INSERT INTO users (username, password_hash) VALUES ($1, $2) RETURNING id`
	querryAuthorizeUser    = `SELECT id, username, password_hash FROM users WHERE username = $1`
	querySetDefaultBalance = `INSERT INTO coins (user_id, balance) VALUES ($1, $2) ON CONFLICT (user_id) DO NOTHING`
)

func (s *storage) Create(username, hashedPassword string) (*models.User, error) {
	const op = "storage.postgresql.Create"

	var userID string

	err := s.DB.QueryRow(querryCreateUser, username, hashedPassword).Scan(&userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("no rows affected")
		}
		return nil, fmt.Errorf("%s: %s: %w", op, errCreateUser, err)
	}

	if err := s.GiveDefaultBalanceToUser(userID); err != nil {
		return nil, fmt.Errorf("%s: failed to give default balance: %w", op, err)
	}

	return &models.User{
		ID:           userID,
		Username:     username,
		PasswordHash: hashedPassword,
	}, nil
}

func (s *storage) GetUserByUsername(username string) (*models.User, error) {
	const op = "storage.postgresql.GetUserByUsername"

	user := models.User{}

	row := s.DB.QueryRow(querryAuthorizeUser, username)

	err := row.Scan(&user.ID, &user.Username, &user.PasswordHash)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("%s: request execution error: %w", op, err)
	}

	return &user, nil
}

func (s *storage) GiveDefaultBalanceToUser(userID string) error {
	const op = "storage.postgresql.GiveDefaultBalanceToUser"

	_, err := s.DB.Exec(querySetDefaultBalance, userID, defaultBalance)
	if err != nil {
		return fmt.Errorf("%s: %s: %w", op, errFailedToSetBalance, err)
	}

	return nil
}

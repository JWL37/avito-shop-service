package postgresql

import (
	"avito-shop-service/internal/models"
	"context"
	"database/sql"
	"errors"
	"fmt"
)

const (
	queryUpdateBalance   = "UPDATE coins SET balance = balance - $1 WHERE user_id = $2"
	queryAddToInventory  = "INSERT INTO inventory (user_id, item_id, quantity) VALUES ($1, $2, 1) ON CONFLICT (user_id, item_id) DO UPDATE SET quantity = inventory.quantity + 1"
	queryGetItemByName   = `SELECT id, name, price FROM items WHERE name = $1`
	querryGetUserBalance = `SELECT balance FROM coins WHERE user_id = $1`

	errFetchBalance    = "failed to fetch user balance"
	errFetchItemCost   = "failed to fetch item cost"
	errUpdateBalance   = "failed to update user balance"
	errUpdateInventory = "failed to update inventory"
	errQuery           = "database query error"
)

func (s *Storage) AddItemToInventory(ctx context.Context, userID string, item *models.Item) error {
	const op = "storage.postgresql.AddItemToInventory"

	tx, err := s.DB.BeginTx(ctx, nil)
	if err != nil {

		return fmt.Errorf("%s: start transaction: %w", op, err)
	}

	defer func() {
		var e error
		if err == nil {
			e = tx.Commit()
		} else {
			e = tx.Rollback()
		}

		if err == nil && e != nil {
			err = fmt.Errorf("finishing transaction: %w", e)
		}
	}()

	if _, err = tx.ExecContext(ctx, queryUpdateBalance, item.Price, userID); err != nil {

		return fmt.Errorf("%s: %s: %w", op, errUpdateBalance, err)
	}

	if _, err = tx.ExecContext(ctx, queryAddToInventory, userID, item.ID); err != nil {

		return fmt.Errorf("%s: %s: %w", op, errUpdateInventory, err)
	}

	return nil
}

func (s *Storage) GetItemByName(ctx context.Context, itemName string) (*models.Item, error) {
	const op = "storage.postgresql.GetItemByName"

	item := &models.Item{}

	if err := s.DB.QueryRowContext(ctx, queryGetItemByName, itemName).Scan(&item.ID, &item.Name, &item.Price); err != nil {

		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("%s: item does not exist :%w", op, err)
		}

		return nil, fmt.Errorf("%s: %s: %w", op, errFetchItemCost, err)
	}
	return item, nil
}

func (s *Storage) GetUserBalance(ctx context.Context, userID string) (int, error) {
	const op = "storage.postgresql.GetUserBalance"

	var balance int

	if err := s.DB.QueryRowContext(ctx, querryGetUserBalance, userID).Scan(&balance); err != nil {
		if errors.Is(err, sql.ErrNoRows) {

			return 0, fmt.Errorf("%s: user balance does not exist: %w", op, err)
		}

		return 0, fmt.Errorf("%s: %s %w", op, errFetchBalance, err)
	}
	return balance, nil
}

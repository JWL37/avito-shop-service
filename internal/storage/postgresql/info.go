package postgresql

import (
	"avito-shop-service/internal/models"
	"context"
	"log/slog"
)

func (s *storage) GetUserInventory(ctx context.Context, userID string) ([]models.InventoryItem, error) {
	query := `
		SELECT i.name, inv.quantity
		FROM inventory inv
		JOIN items i ON inv.item_id = i.id
		WHERE inv.user_id = $1
	`

	rows, err := s.DB.QueryContext(ctx, query, userID)

	if err != nil {
		s.log.Error("failed to query user inventory", slog.Any("error", err))
		return nil, err
	}

	defer rows.Close()

	var inventory []models.InventoryItem

	for rows.Next() {

		var item models.InventoryItem

		if err := rows.Scan(&item.Type, &item.Quantity); err != nil {
			s.log.Error("failed to scan inventory row", slog.Any("error", err))
			return nil, err
		}

		inventory = append(inventory, item)
	}

	if err = rows.Err(); err != nil {
		s.log.Error("error iterating over rows", slog.Any("error", err))
		return nil, err
	}

	return inventory, nil
}

func (s *storage) GetUserTransactions(ctx context.Context, userID string) (models.CoinHistory, error) {
	queryReceived := `
		SELECT u.username, t.amount
		FROM transactions t
		JOIN users u ON t.sender_id = u.id
		WHERE t.receiver_id = $1
		ORDER BY t.transaction_time DESC
	`

	querySent := `
		SELECT u.username, t.amount 
		FROM transactions t
		JOIN users u ON t.receiver_id = u.id
		WHERE t.sender_id = $1
		ORDER BY t.transaction_time DESC
	`

	var history models.CoinHistory

	receivedRows, err := s.DB.QueryContext(ctx, queryReceived, userID)
	if err != nil {
		s.log.Error("failed to query received transactions", slog.Any("error", err))
		return models.CoinHistory{}, err
	}
	defer receivedRows.Close()

	for receivedRows.Next() {
		var transaction models.CoinTransaction
		if err := receivedRows.Scan(&transaction.FromUser, &transaction.Amount); err != nil {
			s.log.Error("failed to scan received transaction", slog.Any("error", err))
			return models.CoinHistory{}, err
		}
		history.Received = append(history.Received, transaction)
	}

	sentRows, err := s.DB.QueryContext(ctx, querySent, userID)
	if err != nil {
		s.log.Error("failed to query sent transactions", slog.Any("error", err))
		return models.CoinHistory{}, err
	}
	defer sentRows.Close()

	for sentRows.Next() {
		var transaction models.CoinTransaction
		if err := sentRows.Scan(&transaction.ToUser, &transaction.Amount); err != nil {
			s.log.Error("failed to scan sent transaction", slog.Any("error", err))
			return models.CoinHistory{}, err
		}
		history.Sent = append(history.Sent, transaction)
	}

	return history, nil
}

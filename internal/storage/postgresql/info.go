package postgresql

import (
	"avito-shop-service/internal/models"
	"context"
	"fmt"
	"log/slog"
)

func (s *Storage) GetUserInventory(ctx context.Context, userID string) ([]models.InventoryItem, error) {
	const op = "storage.postgresql.GetUserInventory"

	query := `
		SELECT i.name, inv.quantity
		FROM inventory inv
		JOIN items i ON inv.item_id = i.id
		WHERE inv.user_id = $1
	`

	rows, err := s.DB.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to query user inventory: %w", op, err)
	}

	defer rows.Close()

	var inventory []models.InventoryItem

	for rows.Next() {

		var item models.InventoryItem

		if err := rows.Scan(&item.Type, &item.Quantity); err != nil {
			return nil, fmt.Errorf("%s: failed to scan inventory row: %w", op, err)
		}

		inventory = append(inventory, item)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: error iterating over rows: %w", op, err)
	}

	return inventory, nil
}

func (s *Storage) getReceivedTransactions(ctx context.Context, userID string) ([]models.CoinTransaction, error) {
	const op = "storage.postgresql.getReceivedTransactions"

	query := `
		SELECT u.username, t.amount
		FROM transactions t
		JOIN users u ON t.sender_id = u.id
		WHERE t.receiver_id = $1
		ORDER BY t.transaction_time DESC
	`

	receivedRows, err := s.DB.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("%s failed to query received transactions: %w", op, err)
	}
	defer receivedRows.Close()

	var received []models.CoinTransaction

	for receivedRows.Next() {

		var transaction models.CoinTransaction

		if err := receivedRows.Scan(&transaction.FromUser, &transaction.Amount); err != nil {
			s.log.Error("failed to scan received transaction", slog.Any("error", err))
			return nil, err
		}
		received = append(received, transaction)
	}

	return received, nil
}

func (s *Storage) getSentTransactions(ctx context.Context, userID string) ([]models.CoinTransaction, error) {
	const op = "storage.postgresql.getSentTransactions"

	query := `
		SELECT u.username, t.amount
		FROM transactions t
		JOIN users u ON t.receiver_id = u.id
		WHERE t.sender_id = $1
		ORDER BY t.transaction_time DESC
	`

	sentRows, err := s.DB.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to query sent transactions:%w", op, err)
	}
	defer sentRows.Close()

	var sent []models.CoinTransaction

	for sentRows.Next() {

		var transaction models.CoinTransaction

		if err := sentRows.Scan(&transaction.ToUser, &transaction.Amount); err != nil {
			s.log.Error("failed to scan sent transaction", slog.Any("error", err))
			return nil, err
		}
		sent = append(sent, transaction)
	}

	return sent, nil
}

func (s *Storage) GetUserTransactions(ctx context.Context, userID string) (models.CoinHistory, error) {
	const op = "storage.postgresql.GetUserTransactions"

	received, err := s.getReceivedTransactions(ctx, userID)
	if err != nil {
		return models.CoinHistory{}, fmt.Errorf("%s failed to get received transactions: %w", op, err)
	}

	sent, err := s.getSentTransactions(ctx, userID)
	if err != nil {
		return models.CoinHistory{}, fmt.Errorf("%s failed to get sent transactions: %w", op, err)
	}

	return models.CoinHistory{Received: received, Sent: sent}, nil
}

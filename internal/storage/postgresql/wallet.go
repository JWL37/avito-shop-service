package postgresql

import (
	"context"
	"fmt"
)

func (s *storage) SendCoinToUser(ctx context.Context, receiverID string, senderID string, amount int) error {
	const op = "storage.postgresql.SendCoinToUser"

	const queryUpdateSender = `UPDATE coins SET balance = balance - $1 WHERE user_id = $2 AND balance >= $1`
	const queryUpdateReceiver = `UPDATE coins SET balance = balance + $1 WHERE user_id = $2 `
	const queryInsertTransaction = `INSERT INTO transactions (sender_id, receiver_id, amount) VALUES ($1, $2, $3)`

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

	if _, err = tx.ExecContext(ctx, queryUpdateSender, amount, senderID); err != nil {
		return fmt.Errorf("%s: update sender balance: %w", op, err)
	}

	if _, err = tx.ExecContext(ctx, queryUpdateReceiver, amount, receiverID); err != nil {
		return fmt.Errorf("%s: update receiver balance: %w", op, err)
	}

	if _, err = tx.ExecContext(ctx, queryInsertTransaction, senderID, receiverID, amount); err != nil {
		return fmt.Errorf("%s: insert transaction: %w", op, err)
	}

	return nil
}

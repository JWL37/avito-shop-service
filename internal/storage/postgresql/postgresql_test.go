package postgresql

import (
	"avito-shop-service/internal/models"
	"context"
	"database/sql"
	"errors"
	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
)

func TestCreateUser(t *testing.T) {
	const username = "testuser"

	tests := []struct {
		name          string
		setupMock     func(mock sqlmock.Sqlmock)
		expectError   bool
		errorContains string
	}{
		{
			name: "Successful user creation",
			setupMock: func(mock sqlmock.Sqlmock) {
				hashedPassword := "hashedpassword"
				userID := "123e4567-e89b-12d3-a456-426614174000"

				mock.ExpectQuery("INSERT INTO users").
					WithArgs(username, hashedPassword).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(userID))

				mock.ExpectExec("INSERT INTO coins").
					WithArgs(userID, defaultBalance).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			expectError: false,
		},
		{
			name: "Failed user insertion",
			setupMock: func(mock sqlmock.Sqlmock) {
				hashedPassword := "hashedpassword"

				mock.ExpectQuery("INSERT INTO users ").
					WithArgs(username, hashedPassword).
					WillReturnError(errors.New("failed to insert"))
			},
			expectError:   true,
			errorContains: "failed to create user",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			assert.NoError(t, err)
			defer db.Close()

			store := NewRep(db, slog.Default())
			tt.setupMock(mock)

			user, err := store.Create("testuser", "hashedpassword")

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, user)
				assert.Contains(t, err.Error(), tt.errorContains)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, user)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestGetUserByUsername(t *testing.T) {
	const username = "testuser"

	tests := []struct {
		name          string
		setupMock     func(mock sqlmock.Sqlmock)
		expectUser    *models.User
		expectError   bool
		errorContains string
	}{
		{
			name: "User found",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT id, username, password_hash FROM users WHERE").
					WithArgs(username).
					WillReturnRows(sqlmock.NewRows([]string{"id", "username", "password_hash"}).
						AddRow("123e4567-e89b-12d3-a456-426614174000", username, "hashedpassword"))
			},
			expectUser: &models.User{
				ID:           "123e4567-e89b-12d3-a456-426614174000",
				Username:     "testuser",
				PasswordHash: "hashedpassword",
			},
			expectError: false,
		},
		{
			name: "User not found",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT id, username, password_hash FROM users WHERE ").
					WithArgs(username).
					WillReturnError(sql.ErrNoRows)
			},
			expectUser:  nil,
			expectError: false,
		},
		{
			name: "Query execution error",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT id, username, password_hash FROM users WHERE ").
					WithArgs(username).
					WillReturnError(errors.New("query error"))
			},
			expectUser:    nil,
			expectError:   true,
			errorContains: "request execution error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			assert.NoError(t, err)
			defer db.Close()

			store := NewRep(db, slog.Default())
			tt.setupMock(mock)

			user, err := store.GetUserByUsername("testuser")

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, user)
				assert.Contains(t, err.Error(), tt.errorContains)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectUser, user)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestSendCoinToUser(t *testing.T) {
	tests := []struct {
		name          string
		setupMock     func(mock sqlmock.Sqlmock)
		expectError   bool
		errorContains string
	}{
		{
			name: "Successful transaction",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE coins SET balance = balance - ").
					WithArgs(100, "sender_id").
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectExec("UPDATE coins SET balance = balance + ").
					WithArgs(100, "receiver_id").
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectExec("INSERT INTO transactions ").
					WithArgs("sender_id", "receiver_id", 100).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			expectError: false,
		},
		{
			name: "Insufficient balance",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE coins SET balance = balance -").
					WithArgs(100, "sender_id").
					WillReturnError(errors.New("db error"))
				mock.ExpectRollback()
			},
			expectError:   true,
			errorContains: "update sender balance",
		},

		{
			name: "Database error on sender update",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE coins SET balance = balance -").
					WithArgs(100, "sender_id").
					WillReturnError(errors.New("db error"))
				mock.ExpectRollback()
			},
			expectError:   true,
			errorContains: "update sender balance",
		},
		{
			name: "Database error on receiver update",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE coins SET balance = balance - ").
					WithArgs(100, "sender_id").
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectExec("UPDATE coins SET balance = balance").
					WithArgs(100, "receiver_id").
					WillReturnError(errors.New("db error"))
				mock.ExpectRollback()
			},
			expectError:   true,
			errorContains: "update receiver balance",
		},
		{
			name: "Database error on transaction insert",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE coins SET balance = balance - ").
					WithArgs(100, "sender_id").
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectExec("UPDATE coins SET balance = balance ").
					WithArgs(100, "receiver_id").
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectExec("INSERT INTO transactions ").
					WithArgs("sender_id", "receiver_id", 100).
					WillReturnError(errors.New("db error"))
				mock.ExpectRollback()
			},
			expectError:   true,
			errorContains: "insert transaction",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			assert.NoError(t, err)
			defer db.Close()

			store := NewRep(db, slog.Default())
			tt.setupMock(mock)

			err = store.SendCoinToUser(context.Background(), "receiver_id", "sender_id", 100)

			if tt.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorContains)
			} else {
				assert.NoError(t, err)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestGetItemByName(t *testing.T) {
	tests := []struct {
		name          string
		itemName      string
		setupMock     func(mock sqlmock.Sqlmock)
		expectError   bool
		errorContains string
		expectedItem  *models.Item
	}{
		{
			name:     "Successful item retrieval",
			itemName: "TestItem",
			setupMock: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "name", "price"}).
					AddRow("item123", "TestItem", 100)
				mock.ExpectQuery("SELECT id, name, price FROM items WHERE").
					WithArgs("TestItem").
					WillReturnRows(rows)
			},
			expectError:  false,
			expectedItem: &models.Item{ID: "item123", Name: "TestItem", Price: 100},
		},
		{
			name:     "Item does not exist",
			itemName: "NonExistentItem",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT id, name, price FROM items WHERE name").
					WithArgs("NonExistentItem").
					WillReturnError(sql.ErrNoRows)
			},
			expectError:   true,
			errorContains: "item does not exist",
		},
		{
			name:     "Database error",
			itemName: "TestItem",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT id, name, price FROM items WHERE name =").
					WithArgs("TestItem").
					WillReturnError(errors.New("db error"))
			},
			expectError:   true,
			errorContains: "fetch item cost",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			assert.NoError(t, err)
			defer db.Close()

			store := NewRep(db, slog.Default())
			tt.setupMock(mock)

			item, err := store.GetItemByName(context.Background(), tt.itemName)

			if tt.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorContains)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedItem, item)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestGetUserBalance(t *testing.T) {
	tests := []struct {
		name            string
		setupMock       func(mock sqlmock.Sqlmock)
		expectError     bool
		errorContains   string
		expectedBalance int
	}{
		{
			name: "Successful balance retrieval",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT balance FROM coins WHERE").
					WithArgs("user_id").
					WillReturnRows(sqlmock.NewRows([]string{"balance"}).AddRow(500))
			},
			expectError:     false,
			expectedBalance: 500,
		},
		{
			name: "User balance does not exist",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT balance FROM coins WHERE").
					WithArgs("user_id").
					WillReturnError(sql.ErrNoRows)
			},
			expectError:     true,
			errorContains:   "user balance does not exist",
			expectedBalance: 0,
		},
		{
			name: "Database error",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT balance FROM coins WHERE").
					WithArgs("user_id").
					WillReturnError(errors.New("db error"))
			},
			expectError:     true,
			errorContains:   "failed to fetch user balance",
			expectedBalance: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			assert.NoError(t, err)
			defer db.Close()

			store := NewRep(db, slog.Default())
			tt.setupMock(mock)

			balance, err := store.GetUserBalance(context.Background(), "user_id")

			if tt.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorContains)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedBalance, balance)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestAddItemToInventory(t *testing.T) {
	tests := []struct {
		name          string
		setupMock     func(mock sqlmock.Sqlmock)
		expectError   bool
		errorContains string
	}{
		{
			name: "Successful item addition",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE coins SET balance = balance -").
					WithArgs(100, "user_id").
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectExec("INSERT INTO inventory").
					WithArgs("user_id", "item_id").
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			expectError: false,
		},
		{
			name: "Error updating balance",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE coins SET balance = balance -").
					WithArgs(100, "user_id").
					WillReturnError(errors.New("balance update error"))
				mock.ExpectRollback()
			},
			expectError:   true,
			errorContains: "failed to update user balance",
		},
		{
			name: "Error inserting into inventory",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE coins SET balance = balance -").
					WithArgs(100, "user_id").
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectExec("INSERT INTO inventory").
					WithArgs("user_id", "item_id").
					WillReturnError(errors.New("inventory update error"))
				mock.ExpectRollback()
			},
			expectError:   true,
			errorContains: "update inventory",
		},
		{
			name: "Transaction start error",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin().WillReturnError(errors.New("transaction error"))
			},
			expectError:   true,
			errorContains: "start transaction",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			assert.NoError(t, err)
			defer db.Close()

			store := NewRep(db, slog.Default())
			tt.setupMock(mock)

			item := &models.Item{ID: "item_id", Name: "Test Item", Price: 100}
			err = store.AddItemToInventory(context.Background(), "user_id", item)

			if tt.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorContains)
			} else {
				assert.NoError(t, err)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

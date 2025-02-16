package user

import (
	"avito-shop-service/internal/models"
	"avito-shop-service/internal/usecases/user/mock"
	"context"
	"fmt"
	"log/slog"

	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

func TestRegister(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAuth := mock.NewMockUserAuthenticater(ctrl)
	logger := slog.Default()
	secret := "testsecret"
	u := New(logger, mockAuth, secret)

	t.Run("Success - User Already Exists", func(t *testing.T) {
		password := "securepassword"
		//nolint:errcheck
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
		existingUser := &models.User{Username: "testuser", PasswordHash: string(hashedPassword)}

		mockAuth.EXPECT().GetUserByUsername("testuser").Return(existingUser, nil)

		token, err := u.Register("testuser", password)
		assert.NoError(t, err)
		assert.NotEmpty(t, token)
	})

	t.Run("Failure - Invalid Credentials", func(t *testing.T) {
		mockAuth.EXPECT().GetUserByUsername("testuser").Return(&models.User{Username: "testuser", PasswordHash: "wronghash"}, nil)

		token, err := u.Register("testuser", "wrongpassword")
		assert.Error(t, err)
		assert.Empty(t, token)
	})

	t.Run("Success - New User Registration", func(t *testing.T) {
		mockAuth.EXPECT().GetUserByUsername("newuser").Return(nil, nil)
		password := "securepassword"
		//nolint:errcheck
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
		mockAuth.EXPECT().Create("newuser", gomock.Any()).Return(&models.User{Username: "newuser", PasswordHash: string(hashedPassword)}, nil)

		token, err := u.Register("newuser", password)
		assert.NoError(t, err)
		assert.NotEmpty(t, token)
	})

	t.Run("Failure - Error Creating User", func(t *testing.T) {
		mockAuth.EXPECT().GetUserByUsername("newuser").Return(nil, nil)
		mockAuth.EXPECT().Create("newuser", gomock.Any()).Return(nil, errors.New("create error"))

		token, err := u.Register("newuser", "newpassword")
		assert.Error(t, err)
		assert.Empty(t, token)
	})

	t.Run("Failure - Error get user by username", func(t *testing.T) {
		mockAuth.EXPECT().GetUserByUsername("newuser").Return(nil, errors.New("create error"))

		token, err := u.Register("newuser", "newpassword")
		assert.Error(t, err)
		assert.Empty(t, token)
	})
}

func TestGetUserInfo(t *testing.T) {
	tests := []struct {
		name             string
		mockSetup        func(mockUserInfoGetter *mock.MockUserInfoGetter)
		expectedUserInfo models.UserInfo
		expectedError    string
	}{
		{
			name: "Success",
			mockSetup: func(mockUserInfoGetter *mock.MockUserInfoGetter) {
				mockUserInfoGetter.EXPECT().GetUserBalance(gomock.Any(), "user123").Return(100, nil)
				mockUserInfoGetter.EXPECT().GetUserInventory(gomock.Any(), "user123").Return([]models.InventoryItem{
					{Type: "item1", Quantity: 2},
					{Type: "item2", Quantity: 5},
				}, nil)
				mockUserInfoGetter.EXPECT().GetUserTransactions(gomock.Any(), "user123").Return(models.CoinHistory{
					Received: []models.CoinTransaction{{FromUser: "user456", ToUser: "user123", Amount: 10}},
					Sent:     []models.CoinTransaction{{FromUser: "user123", ToUser: "user789", Amount: 20}},
				}, nil)
			},
			expectedUserInfo: models.UserInfo{
				Coins: 100,
				Inventory: []models.InventoryItem{
					{Type: "item1", Quantity: 2},
					{Type: "item2", Quantity: 5},
				},
				CoinHistory: models.CoinHistory{
					Received: []models.CoinTransaction{{FromUser: "user456", ToUser: "user123", Amount: 10}},
					Sent:     []models.CoinTransaction{{FromUser: "user123", ToUser: "user789", Amount: 20}},
				},
			},
			expectedError: "",
		},
		{
			name: "Fail to get balance",
			mockSetup: func(mockUserInfoGetter *mock.MockUserInfoGetter) {
				mockUserInfoGetter.EXPECT().GetUserBalance(gomock.Any(), "user123").Return(0, fmt.Errorf("some error"))
			},
			expectedUserInfo: models.UserInfo{},
			expectedError:    "usecases.user.GetUserInfo: failed to get user balance: some error",
		},
		{
			name: "Fail to get inventory",
			mockSetup: func(mockUserInfoGetter *mock.MockUserInfoGetter) {
				mockUserInfoGetter.EXPECT().GetUserBalance(gomock.Any(), "user123").Return(100, nil)
				mockUserInfoGetter.EXPECT().GetUserInventory(gomock.Any(), "user123").Return(nil, fmt.Errorf("some error"))
			},
			expectedUserInfo: models.UserInfo{},
			expectedError:    "usecases.user.GetUserInfo: failed to get inventory: some error",
		},
		{
			name: "Fail to get coin history",
			mockSetup: func(mockUserInfoGetter *mock.MockUserInfoGetter) {
				mockUserInfoGetter.EXPECT().GetUserBalance(gomock.Any(), "user123").Return(100, nil)
				mockUserInfoGetter.EXPECT().GetUserInventory(gomock.Any(), "user123").Return([]models.InventoryItem{{Type: "item1", Quantity: 2}}, nil)
				mockUserInfoGetter.EXPECT().GetUserTransactions(gomock.Any(), "user123").Return(models.CoinHistory{}, fmt.Errorf("some error"))
			},
			expectedUserInfo: models.UserInfo{},
			expectedError:    "usecases.user.GetUserInfo: failed to get coin history: some error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockUserInfoGetter := mock.NewMockUserInfoGetter(ctrl)

			tt.mockSetup(mockUserInfoGetter)

			infoUsecase := NewInfo(slog.Default(), mockUserInfoGetter)

			userInfo, err := infoUsecase.GetUserInfo(context.Background(), "user123")

			if tt.expectedError != "" {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError, err.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedUserInfo, userInfo)
			}
		})
	}
}

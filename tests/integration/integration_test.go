package integration

import (
	"encoding/json"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	baseURL = "http://localhost:8080" // url test_service
)

type AuthResponse struct {
	Token string `json:"token"`
}

type ErrorResponse struct {
	Errors string `json:"errors"`
}

type SendCoinRequest struct {
	ToUser string `json:"toUser"`
	Amount int    `json:"amount"`
}

func TestBuyMerch(t *testing.T) {
	// Регистрация пользователя
	authBody := `{"username": "ttessttuser", "password": "testpassword"}`
	req, err := http.NewRequest(http.MethodPost, baseURL+"/api/auth", strings.NewReader(authBody))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var authResp AuthResponse
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&authResp))
	require.NotEmpty(t, authResp.Token)

	// Покупка товара
	buyURL := baseURL + "/api/buy/socks"
	req, err = http.NewRequest(http.MethodGet, buyURL, nil)
	require.NoError(t, err)
	req.Header.Set("Authorization", "Bearer "+authResp.Token)

	resp, err = http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	// Проверяем, что ответ успешный
	assert.Equal(t, http.StatusOK, resp.StatusCode, "Покупка товара не удалась")
}

func TestSendCoins(t *testing.T) {
	// Регистрация первого пользователя
	authBody1 := `{"username": "ttessttuser1", "password": "password1"}`
	req, err := http.NewRequest(http.MethodPost, baseURL+"/api/auth", strings.NewReader(authBody1))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var authResp1 AuthResponse
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&authResp1))
	require.NotEmpty(t, authResp1.Token)

	// Регистрация второго пользователя
	authBody2 := `{"username": "ttessttuser2", "password": "password2"}`
	req, err = http.NewRequest(http.MethodPost, baseURL+"/api/auth", strings.NewReader(authBody2))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	resp, err = http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var authResp2 AuthResponse
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&authResp2))
	require.NotEmpty(t, authResp2.Token)

	// Отправка монет
	sendCoinBody := `{"toUser": "ttessttuser2", "amount": 10}`
	req, err = http.NewRequest(http.MethodPost, baseURL+"/api/sendCoin", strings.NewReader(sendCoinBody))
	require.NoError(t, err)
	req.Header.Set("Authorization", "Bearer "+authResp1.Token)
	req.Header.Set("Content-Type", "application/json")

	resp, err = http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	// Проверяем, что ответ успешный
	assert.Equal(t, http.StatusOK, resp.StatusCode, "Передача монет не удалась")
}

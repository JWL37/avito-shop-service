package user

import (
	"avito-shop-service/internal/models"
	"errors"
	"fmt"

	jwt "github.com/dgrijalva/jwt-go"

	"log/slog"
	"time"

	"golang.org/x/crypto/bcrypt"
)

var (
	errInvalidCredentials = errors.New("invalid credentials")
)

//go:generate mockgen -source=auth.go -destination=mock/auth_mock.go -package=mock UserAuthenticater
type UserAuthenticater interface {
	Create(string, string) (*models.User, error)
	GetUserByUsername(string) (*models.User, error)
}

type Useacase struct {
	userAuthenticater UserAuthenticater
	log               *slog.Logger
	secret            string
}

func New(log *slog.Logger, userAuthenticater UserAuthenticater, secret string) *Useacase {
	return &Useacase{
		userAuthenticater: userAuthenticater,
		log:               log,
		secret:            secret,
	}
}

func (u *Useacase) Register(username, password string) (string, error) {
	const op = "usecases.user.Register"

	user, err := u.userAuthenticater.GetUserByUsername(username)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	if user == nil {

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
		if err != nil {
			return "", fmt.Errorf("%s: %w", op, err)
		}

		user, err = u.userAuthenticater.Create(username, string(hashedPassword))
		if err != nil {
			return "", fmt.Errorf("%s: %w", op, err)
		}
	}

	if bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)) != nil {
		return "", fmt.Errorf("%s: %w", op, errInvalidCredentials)
	}

	token, err := u.GenerateJWTtoken(user)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return token, nil
}

func (u *Useacase) GenerateJWTtoken(user *models.User) (string, error) {
	const op = "usecases.user.GenerateJWTtoken"

	now := time.Now()
	exp := now.AddDate(0, 0, 7)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user": map[string]string{
			"username": user.Username,
			"id":       user.ID,
		},
		"iat": now.Unix(),
		"exp": exp.Unix(),
	})

	tokenString, err := token.SignedString([]byte(u.secret))
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return tokenString, nil
}

package service

import (
	"crypto/sha1"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/child6yo/y-lms-discalc/orchestrator"
)

const (
	salt       = "w4rio4893jkwlq3"
	signingKey = "qrkjk#4#35FSFJlja#4353KSFjH"
	tokenTTL   = 3 * time.Hour
)

type tokenClaims struct {
	UserID int `json:"user_id"`
	jwt.RegisteredClaims
}

// CreateUser создает нового пользователя.
// На вход принимает модель пользователя.
func (s *MainService) CreateUser(user orchestrator.User) (int, error) {
	user.Password = generatePasswordHash(user.Password)

	return s.repo.CreateUser(user)
}

// GenerateToken генерирует и возвращает новый JWT.
// На вход принимает логин и пароль пользователя.
func (s *MainService) GenerateToken(username, password string) (string, error) {
	user, err := s.repo.GetUser(username, generatePasswordHash(password))
	if err != nil {
		return "", err
	}

	claims := tokenClaims{
		user.ID,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(tokenTTL)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(signingKey))
}

// ParseToken обрабатывает JWT и возвращает из него айди пользователя.
// На вход принимает JWT.
func (s *MainService) ParseToken(accessToken string) (int, error) {
	token, err := jwt.ParseWithClaims(accessToken, &tokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}

		return []byte(signingKey), nil
	})
	if err != nil {
		return 0, err
	}

	claims, ok := token.Claims.(*tokenClaims)
	if !ok {
		return 0, errors.New("token claims are not of type *tokenClaims")
	}

	return claims.UserID, nil
}

func generatePasswordHash(password string) string {
	hash := sha1.New()
	hash.Write([]byte(password))

	return fmt.Sprintf("%x", hash.Sum([]byte(salt)))
}

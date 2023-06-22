package service

import (
	"crypto/sha256"
	"fmt"
	"time"

	"github.com/gogapopp/gophermart/internal/app/models"
	"github.com/gogapopp/gophermart/internal/app/storage"
	"github.com/golang-jwt/jwt/v4"
)

const (
	secretKey  = "secretKey"
	signingKey = "signingKey"
)

type tokenClaims struct {
	jwt.RegisteredClaims
	UserID int `json:"id"`
}

type AuthService struct {
	storage storage.Auth
}

// NewAuthService возвращает структуру AuthService
func NewAuthService(storage storage.Auth) *AuthService {
	return &AuthService{storage: storage}
}

// CreateUser хеширует пароль который передал юзер и передаёт дальше на слой storage
func (s *AuthService) CreateUser(user models.User) (int, error) {
	user.Password = generatePasswordHash(user.Password)
	return s.storage.CreateUser(user)
}

// GenerateToken создаёт юзера и токен авторизации
func (s *AuthService) GenerateToken(login, password string) (string, error) {
	user, err := s.storage.GetUser(login, generatePasswordHash(password))
	if err != nil {
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
		jwt.RegisteredClaims{
			IssuedAt: jwt.NewNumericDate(time.Now()),
		},
		user.ID,
	})

	return token.SignedString([]byte(signingKey))
}

// ParseToken разбирает jwt токен на claims
func (s *AuthService) ParseToken(accessToken string) (int, error) {
	token, err := jwt.ParseWithClaims(accessToken, &tokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}

		return []byte(signingKey), nil
	})
	if err != nil {
		return 0, err
	}

	claims, ok := token.Claims.(*tokenClaims)
	if !ok {
		return 0, jwt.ErrTokenInvalidClaims
	}

	return claims.UserID, nil
}

// generatePasswordHas хеширует пароль
func generatePasswordHash(password string) string {
	hash := sha256.New()
	hash.Write([]byte(password))
	return fmt.Sprintf("%x", hash.Sum([]byte(secretKey)))
}

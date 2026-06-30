package service

import (
	"cinefinder/internal/model"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtKey = func() []byte {
	key := os.Getenv("JWT_SECRET")
	if key == "" {
		key = "chave_secreta_dev_fallback"
	}
	return []byte(key)
}()

type AuthService struct{}

func (s *AuthService) GenerateToken(user model.User) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
		"exp":     expirationTime.Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}

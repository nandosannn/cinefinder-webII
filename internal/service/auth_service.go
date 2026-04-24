package service

import (
	"time"
	"github.com/golang-jwt/jwt/v5"
	"cinefinder/internal/model"
)

var jwtKey = []byte("chave_secreta")

func (s *AuthService) GenerateToken(user model.User) (string, error){
	expirationTime := time.Now().add(24 * time.Hour)
	claims := &jwt.MapClaims{
		"user_id": user.ID,
		"email": user.Email,
		"exp": expirationTime.Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigninMethodHS256, claims)
	return token.SignedString(jwtKey)
}
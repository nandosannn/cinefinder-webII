package service

import (
	"cinefinder/internal/model"
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var jwtKey = func() []byte {
	key := os.Getenv("JWT_SECRET")
	if key == "" {
		key = "chave_secreta_dev_fallback"
	}
	return []byte(key)
}()

type RefreshServiceInterface interface {
	GenerateToken(user model.User) (string, error)
	GenerateRefreshToken(ctx context.Context, userID int) (string, error)
	ValidateRefreshToken(ctx context.Context, token string) (*model.RefreshToken, error)
	RevokeRefreshToken(ctx context.Context, token string) error
}

type AuthService struct {
	DB *pgxpool.Pool
}

func NewAuthService(db *pgxpool.Pool) *AuthService {
	return &AuthService{DB: db}
}

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

func (s *AuthService) GenerateRefreshToken(ctx context.Context, userID int) (string, error) {
	if s.DB == nil {
		return "", errors.New("banco não configurado")
	}
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("erro ao gerar token: %w", err)
	}
	token := hex.EncodeToString(b)
	expiresAt := time.Now().Add(7 * 24 * time.Hour)
	_, err := s.DB.Exec(ctx,
		`INSERT INTO refresh_tokens (user_id, token, expires_at) VALUES ($1, $2, $3)`,
		userID, token, expiresAt,
	)
	if err != nil {
		return "", fmt.Errorf("erro ao salvar refresh token: %w", err)
	}
	return token, nil
}

func (s *AuthService) ValidateRefreshToken(ctx context.Context, token string) (*model.RefreshToken, error) {
	if s.DB == nil {
		return nil, errors.New("banco não configurado")
	}
	var rt model.RefreshToken
	err := s.DB.QueryRow(ctx,
		`SELECT id, user_id, token, expires_at, revoked FROM refresh_tokens WHERE token = $1`,
		token,
	).Scan(&rt.ID, &rt.UserID, &rt.Token, &rt.ExpiresAt, &rt.Revoked)
	if err != nil {
		return nil, errors.New("refresh token inválido")
	}
	if rt.Revoked {
		return nil, errors.New("refresh token foi revogado")
	}
	if time.Now().After(rt.ExpiresAt) {
		return nil, errors.New("refresh token expirado")
	}
	return &rt, nil
}

func (s *AuthService) RevokeRefreshToken(ctx context.Context, token string) error {
	if s.DB == nil {
		return errors.New("banco não configurado")
	}
	_, err := s.DB.Exec(ctx,
		`UPDATE refresh_tokens SET revoked = true WHERE token = $1`, token,
	)
	return err
}

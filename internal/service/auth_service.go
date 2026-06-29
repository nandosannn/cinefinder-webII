package service

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"cinefinder/internal/config"
	"cinefinder/internal/model"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AuthServiceInterface interface {
	GenerateToken(user model.User) (string, error)
	GenerateRefreshToken(ctx context.Context, userID int) (string, error)
	ValidateRefreshToken(ctx context.Context, token string) (*model.RefreshToken, error)
	RevokeRefreshToken(ctx context.Context, token string) error
}

type AuthService struct {
	db *pgxpool.Pool
}

func NewAuthService(db *pgxpool.Pool) *AuthService {
	return &AuthService{db: db}
}

func (s *AuthService) GenerateToken(user model.User) (string, error) {
	claims := &jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
		"exp":     time.Now().Add(15 * time.Minute).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(config.GetJWTKey())
}

func (s *AuthService) GenerateRefreshToken(ctx context.Context, userID int) (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("erro ao gerar bytes aleatórios: %w", err)
	}
	token := hex.EncodeToString(b)
	expiresAt := time.Now().Add(7 * 24 * time.Hour)

	_, err := s.db.Exec(ctx,
		`INSERT INTO refresh_tokens (user_id, token, expires_at) VALUES ($1, $2, $3)`,
		userID, token, expiresAt,
	)
	if err != nil {
		return "", fmt.Errorf("erro ao salvar refresh token: %w", err)
	}
	return token, nil
}

func (s *AuthService) ValidateRefreshToken(ctx context.Context, token string) (*model.RefreshToken, error) {
	var rt model.RefreshToken
	err := s.db.QueryRow(ctx,
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
	_, err := s.db.Exec(ctx,
		`UPDATE refresh_tokens SET revoked = true WHERE token = $1`,
		token,
	)
	return err
}

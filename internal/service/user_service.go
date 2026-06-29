package service

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"strings"

	"cinefinder/internal/model"

	"github.com/jackc/pgx/v5/pgxpool"
)

type UserServiceInterface interface {
	Create(user model.User) (*model.User, error)
	List() ([]model.User, error)
	GetByID(id int) (*model.User, error)
	ValidateUser(email, password string) (*model.User, error)
}

type UserService struct {
	db *pgxpool.Pool
}

func NewUserService(db *pgxpool.Pool) *UserService {
	return &UserService{db: db}
}

func hashPassword(password string) (string, error) {
	salt := make([]byte, 16)
	if _, err := rand.Read(salt); err != nil {
		return "", err
	}
	saltHex := hex.EncodeToString(salt)
	hash := sha256.Sum256([]byte(saltHex + password))
	return saltHex + ":" + hex.EncodeToString(hash[:]), nil
}

func checkPassword(password, stored string) bool {
	parts := strings.Split(stored, ":")
	if len(parts) != 2 {
		return false
	}
	hash := sha256.Sum256([]byte(parts[0] + password))
	return hex.EncodeToString(hash[:]) == parts[1]
}

func (s *UserService) Create(user model.User) (*model.User, error) {
	if user.Name == "" || user.Email == "" || user.Password == "" {
		return nil, errors.New("nome, email e senha são obrigatórios")
	}

	var count int
	err := s.db.QueryRow(context.Background(),
		`SELECT COUNT(*) FROM users WHERE email = $1`, user.Email,
	).Scan(&count)
	if err != nil {
		return nil, err
	}
	if count > 0 {
		return nil, errors.New("Usuário já cadastrado")
	}

	hashed, err := hashPassword(user.Password)
	if err != nil {
		return nil, err
	}

	err = s.db.QueryRow(context.Background(),
		`INSERT INTO users (name, email, password, created_at)
		 VALUES ($1, $2, $3, NOW())
		 RETURNING id, name, email, password, created_at`,
		user.Name, user.Email, hashed,
	).Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.CreatedAt)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (s *UserService) List() ([]model.User, error) {
	rows, err := s.db.Query(context.Background(),
		`SELECT id, name, email, password, created_at FROM users`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []model.User
	for rows.Next() {
		var u model.User
		if err := rows.Scan(&u.ID, &u.Name, &u.Email, &u.Password, &u.CreatedAt); err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, nil
}

func (s *UserService) GetByID(id int) (*model.User, error) {
	var u model.User
	err := s.db.QueryRow(context.Background(),
		`SELECT id, name, email, password, created_at FROM users WHERE id = $1`, id,
	).Scan(&u.ID, &u.Name, &u.Email, &u.Password, &u.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (s *UserService) ValidateUser(email, password string) (*model.User, error) {
	var u model.User
	err := s.db.QueryRow(context.Background(),
		`SELECT id, name, email, password, created_at FROM users WHERE email = $1`, email,
	).Scan(&u.ID, &u.Name, &u.Email, &u.Password, &u.CreatedAt)
	if err != nil {
		return nil, errors.New("credenciais inválidas")
	}
	if !checkPassword(password, u.Password) {
		return nil, errors.New("credenciais inválidas")
	}
	return &u, nil
}

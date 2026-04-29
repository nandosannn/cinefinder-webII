package service

import (
	"context"
	"errors"

	"cinefinder/internal/model"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Interface (IMPORTANTE para testes)
type UserServiceInterface interface {
	Create(user model.User) (*model.User, error)
	List() ([]model.User, error)
	GetByID(id int) (*model.User, error)
	ValidateUser(email, password string) (*model.User, error)
}

// Implementação real
type UserService struct {
	db *pgxpool.Pool
}

func NewUserService(db *pgxpool.Pool) *UserService {
	return &UserService{db: db}
}

func (s *UserService) Create(user model.User) (*model.User, error) {
	var userCount int

	checkQuery := "SELECT COUNT(*) FROM users WHERE id = $1 OR email = $2"

	err := s.db.QueryRow(context.Background(), checkQuery, user.ID, user.Email).Scan(&userCount)
	if err != nil {
		return nil, err
	}

	if userCount > 0 {
		return nil, errors.New("Usuário já cadastrado")
	}

	query := `
	INSERT INTO users (name, email, password, created_at)
	VALUES ($1, $2, $3, NOW())
	RETURNING id, name, email, password, created_at;
	`

	err = s.db.QueryRow(context.Background(), query,
		user.Name,
		user.Email,
		user.Password,
	).Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.CreatedAt)

	if err != nil {
		return nil, err
	}

	return &user, nil
}
func (s *UserService) List() ([]model.User, error) {
	rows, err := s.db.Query(context.Background(),
		"SELECT id, name, email, password, created_at FROM users",
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := []model.User{}

	for rows.Next() {
		var u model.User
		err := rows.Scan(&u.ID, &u.Name, &u.Email, &u.Password, &u.CreatedAt)
		if err != nil {
			return nil, err
		}
		users = append(users, u)
	}

	return users, nil
}

func (s *UserService) GetByID(id int) (*model.User, error) {
	query := `
	SELECT id, name, email, password, created_at
	FROM users
	WHERE id = $1;
	`

	var u model.User

	err := s.db.QueryRow(context.Background(), query, id).
		Scan(&u.ID, &u.Name, &u.Email, &u.Password, &u.CreatedAt)

	if err != nil {
		return nil, err
	}

	return &u, nil
}

// Necessário pro auth
func (s *UserService) ValidateUser(email, password string) (*model.User, error) {
	query := "SELECT id, name, email, password, created_at FROM users WHERE email = $1 AND password = $2"
	var u model.User
	err := s.db.QueryRow(context.Background(), query, email, password).
		Scan(&u.ID, &u.Name, &u.Email, &u.Password, &u.CreatedAt)

	if err != nil {
		return nil, errors.New("credenciais inválidas")
	}
	return &u, nil
}

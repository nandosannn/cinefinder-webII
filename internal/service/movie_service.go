package service

import (
	"context"
	"cinefinder/internal/model"
	"github.com/jackc/pgx/v5/pgxpool"
)

type MovieService struct {
	db *pgxpool.Pool
}

func NewMovieService(db *pgxpool.Pool) *MovieService {
	return &MovieService{db: db}
}

func (s *MovieService) Create(movie model.Movie) (*model.Movie, error) {
	query := `
	INSERT INTO movies (title, director, year, genre)
	VALUES ($1, $2, $3, $4)
	RETURNING id;
	`

	err := s.db.QueryRow(context.Background(), query,
		movie.Title,
		movie.Director,
		movie.Year,
		movie.Genre,
	).Scan(&movie.ID)

	if err != nil {
		return nil, err
	}

	return &movie, nil
}
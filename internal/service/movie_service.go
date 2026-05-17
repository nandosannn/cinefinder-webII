package service

import (
	"context"
	"errors"

	"cinefinder/internal/model"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Interface (IMPORTANTE para testes)
type MovieServiceInterface interface {
	Create(movie model.Movie) (*model.Movie, error)
	List() ([]model.Movie, error)
	GetByID(id int) (*model.Movie, error)
}

// Implementação real
type MovieService struct {
	db *pgxpool.Pool
}

func NewMovieService(db *pgxpool.Pool) *MovieService {
	return &MovieService{db: db}
}

func (s *MovieService) Create(movie model.Movie) (*model.Movie, error) {
	var movieCount int
	
	checkQuery := "SELECT COUNT(*) FROM movies WHERE title = $1 AND director = $2 AND year = $3 AND genre = $4"
	
	err := s.db.QueryRow(context.Background(), checkQuery, movie.Title, movie.Director, movie.Year, movie.Genre).Scan(&movieCount)
	if err != nil {
		return nil, err
	}
	
	if movieCount > 0 {
		return nil, errors.New("Filme já cadastrado")
	}
	
	query := `
	INSERT INTO movies (title, director, year, genre, available)
	VALUES ($1, $2, $3, $4, true)
	RETURNING id;
	`

	err = s.db.QueryRow(context.Background(), query,
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

func (s *MovieService) List() ([]model.Movie, error) {
	rows, err := s.db.Query(context.Background(),
		"SELECT id, title, director, year, genre FROM movies",
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	movies := []model.Movie{}

	for rows.Next() {
		var m model.Movie
		err := rows.Scan(&m.ID, &m.Title, &m.Director, &m.Year, &m.Genre)
		if err != nil {
			return nil, err
		}
		movies = append(movies, m)
	}

	return movies, nil
}

func (s *MovieService) GetByID(id int) (*model.Movie, error) {
	query := `
	SELECT id, title, director, year, genre
	FROM movies
	WHERE id = $1;
	`

	var m model.Movie

	err := s.db.QueryRow(context.Background(), query, id).
		Scan(&m.ID, &m.Title, &m.Director, &m.Year, &m.Genre)

	if err != nil {
		return nil, err
	}

	return &m, nil
}

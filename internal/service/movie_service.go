package service

import (
	"context"
	"errors"

	"cinefinder/internal/database"
	"cinefinder/internal/model"
)

type MovieServiceInterface interface {
	Create(movie model.Movie) (*model.Movie, error)
	List() ([]model.Movie, error)
	GetByID(id int) (*model.Movie, error)
	Update(id int, movie model.Movie) (*model.Movie, error)
	Delete(id int) error
}

type MovieService struct {
	queries *database.Queries
}

func NewMovieService(queries *database.Queries) *MovieService {
	return &MovieService{queries: queries}
}

func (s *MovieService) Create(movie model.Movie) (*model.Movie, error) {
	movies, err := s.queries.ListMovies(context.Background())
	if err != nil {
		return nil, err
	}

	for _, m := range movies {
		if m.Title == movie.Title &&
			m.Director == movie.Director &&
			int(m.ReleaseYear) == movie.Year &&
			m.Genre == movie.Genre {
			return nil, errors.New("Filme já cadastrado")
		}
	}

	created, err := s.queries.CreateMovie(context.Background(), database.CreateMovieParams{
		Title:           movie.Title,
		Director:        movie.Director,
		ReleaseYear:     int32(movie.Year),
		Genre:           movie.Genre,
		AvailableCopies: int32(movie.Available),
	})

	if err != nil {
		return nil, err
	}

	return &model.Movie{
		ID:       int(created.ID),
		Title:    created.Title,
		Director: created.Director,
		Year:     int(created.ReleaseYear),
		Genre:    created.Genre,
		Available: int(created.AvailableCopies),
	}, nil
}

func (s *MovieService) List() ([]model.Movie, error) {
	dbMovies, err := s.queries.ListMovies(context.Background())
	if err != nil {
		return nil, err
	}

	movies := []model.Movie{}

	for _, m := range dbMovies {
		movies = append(movies, model.Movie{
			ID:       int(m.ID),
			Title:    m.Title,
			Director: m.Director,
			Year:     int(m.ReleaseYear),
			Genre:    m.Genre,
			Available: int(m.AvailableCopies),
		})
	}

	return movies, nil
}

func (s *MovieService) GetByID(id int) (*model.Movie, error) {
	m, err := s.queries.GetMovieByID(context.Background(), int32(id))
	if err != nil {
		return nil, err
	}

	return &model.Movie{
		ID:       int(m.ID),
		Title:    m.Title,
		Director: m.Director,
		Year:     int(m.ReleaseYear),
		Genre:    m.Genre,
		Available: int(m.AvailableCopies),
	}, nil
}

func (s *MovieService) Update(id int, movie model.Movie) (*model.Movie, error) {
	updated, err := s.queries.UpdateMovie(context.Background(), database.UpdateMovieParams{
		ID:              int32(id),
		Title:           movie.Title,
		Director:        movie.Director,
		ReleaseYear:     int32(movie.Year),
		Genre:           movie.Genre,
		AvailableCopies: int32(movie.Available),
	})

	if err != nil {
		return nil, err
	}

	return &model.Movie{
		ID:       int(updated.ID),
		Title:    updated.Title,
		Director: updated.Director,
		Year:     int(updated.ReleaseYear),
		Genre:    updated.Genre,
		Available: int(updated.AvailableCopies),
	}, nil
}

func (s *MovieService) Delete(id int) error {
	return s.queries.DeleteMovie(context.Background(), int32(id))
}
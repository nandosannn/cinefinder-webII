package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"cinefinder/internal/model"
)

// Mock do service
type mockMovieService struct{}

func (m *mockMovieService) Create(movie model.Movie) (*model.Movie, error) {
	movie.ID = 1
	return &movie, nil
}

func TestCreateMovie_Success(t *testing.T) {
	mockService := &mockMovieService{}
	handler := NewMovieHandler(mockService)

	body := model.Movie{
		Title:    "Matrix",
		Director: "Wachowski",
		Year:     1999,
		Genre:    "Sci-Fi",
	}

	jsonBody, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPost, "/movies", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()

	handler.Create(w, req)

	resp := w.Result()

	if resp.StatusCode != http.StatusCreated {
		t.Errorf("esperado 201, veio %d", resp.StatusCode)
	}

	var response model.Movie
	json.NewDecoder(resp.Body).Decode(&response)

	if response.ID != 1 {
		t.Errorf("esperado ID 1, veio %d", response.ID)
	}
}

func TestCreateMovie_InvalidJSON(t *testing.T) {
	mockService := &mockMovieService{}
	handler := NewMovieHandler(mockService)

	req := httptest.NewRequest(http.MethodPost, "/movies", bytes.NewBuffer([]byte("json inválido")))
	w := httptest.NewRecorder()

	handler.Create(w, req)

	if w.Result().StatusCode != http.StatusBadRequest {
		t.Errorf("esperado 400, veio %d", w.Result().StatusCode)
	}
}
func (m *mockMovieService) List() ([]model.Movie, error) {
	return []model.Movie{
		{
			ID:       1,
			Title:    "Matrix",
			Director: "Wachowski",
			Year:     1999,
			Genre:    "Sci-Fi",
		},
	}, nil
}
func (m *mockMovieService) GetByID(id int) (*model.Movie, error) {
	return &model.Movie{
		ID:       id,
		Title:    "Matrix",
		Director: "Wachowski",
		Year:     1999,
		Genre:    "Sci-Fi",
	}, nil
}
func (m *mockMovieService) Update(id int, movie model.Movie) (*model.Movie, error) {
	movie.ID = id
	return &movie, nil
}

func (m *mockMovieService) Delete(id int) error {
	return nil
}

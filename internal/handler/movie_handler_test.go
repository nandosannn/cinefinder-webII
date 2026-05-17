package handler

import (
	"bytes"
	"encoding/json"
	"errors"
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

func (m *mockMovieService) Search(title, genre string) ([]model.Movie, error) {
	all := []model.Movie{
		{ID: 1, Title: "Matrix", Director: "Wachowski", Year: 1999, Genre: "Sci-Fi"},
		{ID: 2, Title: "Matrix Reloaded", Director: "Wachowski", Year: 2003, Genre: "Sci-Fi"},
		{ID: 3, Title: "Inception", Director: "Nolan", Year: 2010, Genre: "Sci-Fi"},
	}

	var results []model.Movie
	for _, m := range all {
		matchTitle := title == "" || contains(m.Title, title)
		matchGenre := genre == "" || eqIgnoreCase(m.Genre, genre)
		if matchTitle && matchGenre {
			results = append(results, m)
		}
	}
	return results, nil
}

func contains(s, sub string) bool {
	return len(s) >= len(sub) && (s == sub || len(sub) == 0 ||
		func() bool {
			for i := 0; i <= len(s)-len(sub); i++ {
				if s[i:i+len(sub)] == sub {
					return true
				}
			}
			return false
		}())
}

func eqIgnoreCase(a, b string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		ca, cb := a[i], b[i]
		if ca >= 'A' && ca <= 'Z' {
			ca += 32
		}
		if cb >= 'A' && cb <= 'Z' {
			cb += 32
		}
		if ca != cb {
			return false
		}
	}
	return true
}

// --- Testes existentes ---

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

// --- Testes de busca ---

func TestListMovies_SemFiltro(t *testing.T) {
	handler := NewMovieHandler(&mockMovieService{})

	req := httptest.NewRequest(http.MethodGet, "/movies", nil)
	w := httptest.NewRecorder()

	handler.List(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("esperado 200, veio %d", w.Code)
	}

	var movies []model.Movie
	json.NewDecoder(w.Body).Decode(&movies)

	if len(movies) == 0 {
		t.Error("esperado pelo menos 1 filme na listagem")
	}
}

func TestListMovies_FiltroPorTitulo(t *testing.T) {
	handler := NewMovieHandler(&mockMovieService{})

	req := httptest.NewRequest(http.MethodGet, "/movies?title=Matrix", nil)
	w := httptest.NewRecorder()

	handler.List(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("esperado 200, veio %d", w.Code)
	}

	var movies []model.Movie
	json.NewDecoder(w.Body).Decode(&movies)

	if len(movies) != 2 {
		t.Errorf("esperado 2 filmes com 'Matrix' no título, veio %d", len(movies))
	}
}

func TestListMovies_FiltroPorGenero(t *testing.T) {
	handler := NewMovieHandler(&mockMovieService{})

	req := httptest.NewRequest(http.MethodGet, "/movies?genre=sci-fi", nil)
	w := httptest.NewRecorder()

	handler.List(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("esperado 200, veio %d", w.Code)
	}

	var movies []model.Movie
	json.NewDecoder(w.Body).Decode(&movies)

	if len(movies) == 0 {
		t.Error("esperado filmes do gênero Sci-Fi")
	}
}

func TestListMovies_FiltroSemResultado(t *testing.T) {
	handler := NewMovieHandler(&mockMovieService{})

	req := httptest.NewRequest(http.MethodGet, "/movies?title=Naoexiste", nil)
	w := httptest.NewRecorder()

	handler.List(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("esperado 200, veio %d", w.Code)
	}

	var movies []model.Movie
	json.NewDecoder(w.Body).Decode(&movies)

	if movies == nil {
		movies = []model.Movie{}
	}
	if len(movies) != 0 {
		t.Errorf("esperado 0 filmes, veio %d", len(movies))
	}
}

// mockMovieServiceError para simular falhas
type mockMovieServiceError struct{ mockMovieService }

func (m *mockMovieServiceError) Search(title, genre string) ([]model.Movie, error) {
	return nil, errors.New("erro simulado")
}
func (m *mockMovieServiceError) List() ([]model.Movie, error) {
	return nil, errors.New("erro simulado")
}

func TestListMovies_ErroInterno(t *testing.T) {
	handler := NewMovieHandler(&mockMovieServiceError{})

	req := httptest.NewRequest(http.MethodGet, "/movies", nil)
	w := httptest.NewRecorder()

	handler.List(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("esperado 500, veio %d", w.Code)
	}
}

package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"cinefinder/internal/model"
	"cinefinder/internal/service"

	"github.com/go-chi/chi/v5"
)

type MovieHandler struct {
	service service.MovieServiceInterface
}

func NewMovieHandler(s service.MovieServiceInterface) *MovieHandler {
	return &MovieHandler{service: s}
}

func (h *MovieHandler) Create(w http.ResponseWriter, r *http.Request) {
	var movie model.Movie

	if err := json.NewDecoder(r.Body).Decode(&movie); err != nil {
		http.Error(w, "JSON inválido", http.StatusBadRequest)
		return
	}

	createdMovie, err := h.service.Create(movie)
	if err != nil {
		if err.Error() == "Filme já cadastrado" {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		http.Error(w, "Erro ao salvar filme: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdMovie)
}

// List retorna todos os filmes. Se os query params ?title= ou ?genre= forem
// fornecidos, delega para Search filtrando os resultados.
func (h *MovieHandler) List(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Query().Get("title")
	genre := r.URL.Query().Get("genre")

	var (
		movies []model.Movie
		err    error
	)

	if title != "" || genre != "" {
		movies, err = h.service.Search(title, genre)
	} else {
		movies, err = h.service.List()
	}

	if err != nil {
		http.Error(w, "Erro ao buscar filmes", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(movies)
}

func (h *MovieHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")

	id, err := strconv.Atoi(idParam)
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	movie, err := h.service.GetByID(id)
	if err != nil {
		http.Error(w, "Filme não encontrado", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(movie)
}

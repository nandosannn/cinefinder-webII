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
		http.Error(w, "Erro ao salvar filme: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdMovie)
}
func (h *MovieHandler) List(w http.ResponseWriter, r *http.Request) {
	movies, err := h.service.List()
	if err != nil {
		http.Error(w, "Erro ao buscar filmes", http.StatusInternalServerError)
		return
	}

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

	json.NewEncoder(w).Encode(movie)
}

func (h *MovieHandler) Update(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")

	id, err := strconv.Atoi(idParam)
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	var movie model.Movie

	err = json.NewDecoder(r.Body).Decode(&movie)
	if err != nil {
		http.Error(w, "JSON inválido", http.StatusBadRequest)
		return
	}

	updatedMovie, err := h.service.Update(id, movie)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedMovie)
}

func (h *MovieHandler) Delete(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")

	id, err := strconv.Atoi(idParam)
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	err = h.service.Delete(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

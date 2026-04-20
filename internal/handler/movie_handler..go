package handler

import (
	"encoding/json"
	"net/http"

	"cinefinder/internal/model"
	"cinefinder/internal/service"
)

type MovieHandler struct {
	service *service.MovieService
}

func NewMovieHandler(s *service.MovieService) *MovieHandler {
	return &MovieHandler{service: s}
}

func (h *MovieHandler) Create(w http.ResponseWriter, r *http.Request) {
	var movie model.Movie

	// ler JSON do body
	if err := json.NewDecoder(r.Body).Decode(&movie); err != nil {
		http.Error(w, "JSON inválido", http.StatusBadRequest)
		return
	}

	// salvar no banco
	createdMovie, err := h.service.Create(movie)
	if err != nil {
		http.Error(w, "Erro ao salvar filme", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdMovie)
}
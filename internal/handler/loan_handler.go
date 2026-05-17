package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"cinefinder/internal/model"
	"cinefinder/internal/service"

	"github.com/go-chi/chi/v5"
)

type LoanHandler struct {
	service service.LoanServiceInterface
}

func NewLoanHandler(s service.LoanServiceInterface) *LoanHandler {
	return &LoanHandler{service: s}
}

func (h *LoanHandler) Create(w http.ResponseWriter, r *http.Request) {
	var loan model.Loan

	if err := json.NewDecoder(r.Body).Decode(&loan); err != nil {
		http.Error(w, "JSON inválido", http.StatusBadRequest)
		return
	}

	createdLoan, err := h.service.Create(loan)
	if err != nil {

		if err.Error() == "Usuário possui empréstimo em aberto" ||
			err.Error() == "Filme indisponível" {

			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		http.Error(w, "Erro ao criar empréstimo: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(createdLoan)
}

func (h *LoanHandler) List(w http.ResponseWriter, r *http.Request) {
	loans, err := h.service.List()
	if err != nil {
		http.Error(w, "Erro ao buscar empréstimos", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(loans)
}

func (h *LoanHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")

	id, err := strconv.Atoi(idParam)
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	loan, err := h.service.GetByID(id)
	if err != nil {
		http.Error(w, "Empréstimo não encontrado", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(loan)
}
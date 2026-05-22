package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

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

func (h *LoanHandler) ReturnMovie(w http.ResponseWriter, r *http.Request) {

	idParam := chi.URLParam(r, "id")

	id, err := strconv.Atoi(idParam)
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	err = h.service.ReturnMovie(id)
	if err != nil {

		if err.Error() == "Empréstimo já devolvido" ||
			err.Error() == "Empréstimo não encontrado" {

			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		http.Error(
			w,
			"Erro ao devolver filme: "+err.Error(),
			http.StatusInternalServerError,
		)

		return
	}

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(map[string]string{
		"message": "Filme devolvido com sucesso",
	})
}

func (h *LoanHandler) History(w http.ResponseWriter, r *http.Request) {
	var filter model.LoanFilter

	// Filtro por movie_id
	movieIDParam := r.URL.Query().Get("movie_id")
	if movieIDParam != "" {
		movieID, err := strconv.Atoi(movieIDParam)
		if err != nil {
			http.Error(w, "movie_id inválido", http.StatusBadRequest)
			return
		}
		filter.MovieID = movieID
	}

	// Filtro por data inicial (loan_date >= start_date)
	startDateParam := r.URL.Query().Get("start_date")
	if startDateParam != "" {
		startDate, err := time.Parse("2006-01-02", startDateParam)
		if err != nil {
			http.Error(w, "start_date inválido (formato esperado: YYYY-MM-DD)", http.StatusBadRequest)
			return
		}
		filter.StartDate = &startDate
	}

	// Filtro por data final (loan_date <= end_date)
	endDateParam := r.URL.Query().Get("end_date")
	if endDateParam != "" {
		endDate, err := time.Parse("2006-01-02", endDateParam)
		if err != nil {
			http.Error(w, "end_date inválido (formato esperado: YYYY-MM-DD)", http.StatusBadRequest)
			return
		}
		// Ajustar para final do dia (23:59:59)
		endDate = endDate.Add(23*time.Hour + 59*time.Minute + 59*time.Second)
		filter.EndDate = &endDate
	}

	history, err := h.service.ListHistory(filter)
	if err != nil {
		http.Error(w, "Erro ao buscar histórico de empréstimos", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(history)
}
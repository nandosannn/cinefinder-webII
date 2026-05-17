package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"cinefinder/internal/model"

	"github.com/go-chi/chi/v5"
)

// Mock do service
type mockLoanService struct{}

func (m *mockLoanService) Create(loan model.Loan) (*model.Loan, error) {
	loan.ID = 1
	return &loan, nil
}

func (m *mockLoanService) List() ([]model.Loan, error) {
	return []model.Loan{
		{
			ID:         1,
			LoanDate:   time.Now(),
			ReturnDate: time.Now().Add(24 * time.Hour),
			Price:      20,
			Returned:   false,
			UserID:     1,
			MovieID:    1,
		},
	}, nil
}

func (m *mockLoanService) GetByID(id int) (*model.Loan, error) {
	return &model.Loan{
		ID:         id,
		LoanDate:   time.Now(),
		ReturnDate: time.Now().Add(24 * time.Hour),
		Price:      20,
		Returned:   false,
		UserID:     1,
		MovieID:    1,
	}, nil
}

func (m *mockLoanService) ReturnMovie(id int) error {
	return nil
}

// =====================================
// TEST CREATE SUCCESS
// =====================================

func TestCreateLoan_Success(t *testing.T) {

	mockService := &mockLoanService{}
	handler := NewLoanHandler(mockService)

	body := model.Loan{
		ReturnDate: time.Now().Add(24 * time.Hour),
		Price:      15.5,
		Returned:   false,
		UserID:     1,
		MovieID:    1,
	}

	jsonBody, _ := json.Marshal(body)

	req := httptest.NewRequest(
		http.MethodPost,
		"/loans",
		bytes.NewBuffer(jsonBody),
	)

	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()

	handler.Create(w, req)

	resp := w.Result()

	if resp.StatusCode != http.StatusCreated {
		t.Errorf("esperado status 201, veio %d", resp.StatusCode)
	}

	var response model.Loan

	json.NewDecoder(resp.Body).Decode(&response)

	if response.ID != 1 {
		t.Errorf("esperado ID 1, veio %d", response.ID)
	}
}

// =====================================
// TEST CREATE INVALID JSON
// =====================================

func TestCreateLoan_InvalidJSON(t *testing.T) {

	mockService := &mockLoanService{}
	handler := NewLoanHandler(mockService)

	req := httptest.NewRequest(
		http.MethodPost,
		"/loans",
		bytes.NewBuffer([]byte("json inválido")),
	)

	w := httptest.NewRecorder()

	handler.Create(w, req)

	if w.Result().StatusCode != http.StatusBadRequest {
		t.Errorf("esperado status 400, veio %d", w.Result().StatusCode)
	}
}

// =====================================
// TEST LIST SUCCESS
// =====================================

func TestListLoans_Success(t *testing.T) {

	mockService := &mockLoanService{}
	handler := NewLoanHandler(mockService)

	req := httptest.NewRequest(
		http.MethodGet,
		"/loans",
		nil,
	)

	w := httptest.NewRecorder()

	handler.List(w, req)

	resp := w.Result()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("esperado status 200, veio %d", resp.StatusCode)
	}

	var response []model.Loan

	json.NewDecoder(resp.Body).Decode(&response)

	if len(response) == 0 {
		t.Errorf("esperado lista com empréstimos")
	}
}

// =====================================
// TEST GET BY ID SUCCESS
// =====================================

func TestGetLoanByID_Success(t *testing.T) {

	mockService := &mockLoanService{}
	handler := NewLoanHandler(mockService)

	req := httptest.NewRequest(
		http.MethodGet,
		"/loans/1",
		nil,
	)

	w := httptest.NewRecorder()

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "1")

	req = req.WithContext(
		context.WithValue(
			req.Context(),
			chi.RouteCtxKey,
			rctx,
		),
	)

	handler.GetByID(w, req)

	resp := w.Result()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("esperado status 200, veio %d", resp.StatusCode)
	}

	var response model.Loan

	json.NewDecoder(resp.Body).Decode(&response)

	if response.ID != 1 {
		t.Errorf("esperado ID 1, veio %d", response.ID)
	}
}

// =====================================
// TEST GET BY ID INVALID
// =====================================

func TestGetLoanByID_InvalidID(t *testing.T) {

	mockService := &mockLoanService{}
	handler := NewLoanHandler(mockService)

	req := httptest.NewRequest(
		http.MethodGet,
		"/loans/abc",
		nil,
	)

	w := httptest.NewRecorder()

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "abc")

	req = req.WithContext(
		context.WithValue(
			req.Context(),
			chi.RouteCtxKey,
			rctx,
		),
	)

	handler.GetByID(w, req)

	if w.Result().StatusCode != http.StatusBadRequest {
		t.Errorf("esperado status 400, veio %d", w.Result().StatusCode)
	}
}

// =====================================
// TEST RETURN MOVIE SUCCESS
// =====================================

func TestReturnMovie_Success(t *testing.T) {

	mockService := &mockLoanService{}
	handler := NewLoanHandler(mockService)

	req := httptest.NewRequest(
		http.MethodPut,
		"/loans/1/return",
		nil,
	)

	w := httptest.NewRecorder()

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "1")

	req = req.WithContext(
		context.WithValue(
			req.Context(),
			chi.RouteCtxKey,
			rctx,
		),
	)

	handler.ReturnMovie(w, req)

	resp := w.Result()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("esperado status 200, veio %d", resp.StatusCode)
	}

	var response map[string]string

	json.NewDecoder(resp.Body).Decode(&response)

	if response["message"] != "Filme devolvido com sucesso" {
		t.Errorf("mensagem inesperada")
	}
}

// =====================================
// MOCK COM ERRO
// =====================================

type mockLoanServiceError struct{}

func (m *mockLoanServiceError) Create(loan model.Loan) (*model.Loan, error) {
	return nil, errors.New("Usuário possui empréstimo em aberto")
}

func (m *mockLoanServiceError) List() ([]model.Loan, error) {
	return nil, errors.New("erro ao listar")
}

func (m *mockLoanServiceError) GetByID(id int) (*model.Loan, error) {
	return nil, errors.New("não encontrado")
}

func (m *mockLoanServiceError) ReturnMovie(id int) error {
	return errors.New("Empréstimo já devolvido")
}

// =====================================
// TEST RETURN MOVIE ERROR
// =====================================

func TestReturnMovie_AlreadyReturned(t *testing.T) {

	mockService := &mockLoanServiceError{}
	handler := NewLoanHandler(mockService)

	req := httptest.NewRequest(
		http.MethodPut,
		"/loans/1/return",
		nil,
	)

	w := httptest.NewRecorder()

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "1")

	req = req.WithContext(
		context.WithValue(
			req.Context(),
			chi.RouteCtxKey,
			rctx,
		),
	)

	handler.ReturnMovie(w, req)

	if w.Result().StatusCode != http.StatusBadRequest {
		t.Errorf("esperado status 400, veio %d", w.Result().StatusCode)
	}
}

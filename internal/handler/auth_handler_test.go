package handler

import (
	"bytes"
	"cinefinder/internal/model"
	"cinefinder/internal/service"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

type mockAuthUserService struct {
	mockUserService
}

func (m *mockAuthUserService) ValidateUser(email, password string) (*model.User, error) {
	if email == "admin@cinefinder.com" && password == "123456" {
		return &model.User{ID: 1, Email: email, Name: "Admin"}, nil
	}
	return nil, errors.New("credenciais inválidas")
}

func TestLoginHandler_Success(t *testing.T) { // Corrigido Sucess
	userService := &mockAuthUserService{}
	authService := &service.AuthService{}

	h := LoginHandler(authService, userService)

	payload := map[string]string{
		"email":    "admin@cinefinder.com",
		"password": "123456",
	}
	body, _ := json.Marshal(payload)

	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	h.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Esperado status 200, recebeu %d", w.Code)
	}

	var response map[string]string
	json.NewDecoder(w.Body).Decode(&response)
	if _, ok := response["token"]; !ok {
		t.Error("Resposta não contém o campo 'token'")
	}
}

func TestLoginHandler_InvalidCredentials(t *testing.T) {
	userService := &mockAuthUserService{}
	authService := &service.AuthService{}

	h := LoginHandler(authService, userService)

	payload := map[string]string{
		"email":    "errado@cinefinder.com",
		"password": "senha_errada",
	}
	body, _ := json.Marshal(payload)

	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(body)) // Corrigido Mehtod
	w := httptest.NewRecorder()

	h.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Esperado status 401, recebeu %d", w.Code)
	}
}

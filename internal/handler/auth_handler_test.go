package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"cinefinder/internal/model"
	"cinefinder/internal/service"
)

type mockAuthUserService struct {
	mockUserService
}

func (m *mockAuthUserService) ValidateUser(email, password string) (*model.User, error) {
	if email == "admin@cinefinder.com" && password == "123456"{
		return &model.User{ID: 1, Email: email, Name: "Admin"}, nil
	}
	return nil, service.Err.InvalidCredentials
}

func TestLoginHandler_Sucess(t *testing.T){
	userService := &mockAuthUserService{}
	authService := &service.AuthService{}

	h := LoginHandler(authService, userService.UserServiceInterface)

	payload := map[string]string{
		"email": "admin@cinefinder.com",
		"password": "123456",
	}
	body, _ := json.Marshal(payload)

	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	h.ServeHTTP(w, req)

	if w.Code != http.StatusOK{
		t.Error("Esperado status 200, recebeu %d", w.Code)

		var response map[string]string
		json.NewDecoder(w.Body).Decode(&response)

		if _, ok := response["token"]; !ok {
			t.Error("Resposta não contém o campo 'token'")
		}
	}
}

func TestLoginHandler_InvalidCredentials(t *testing.T){
	userService := &mockAuthUserService{}
	authService := &service.AuthService{}

	h := LoginHandler(authService, userService.UserServiceInterface)

	payload := map[string]string{
		"email": "errado@cinefinder.com",
		"password": "senha_errada",
	}
	body, _ := json.Marshal(payload)

	req := httptest.NewRequest(http.MehtodPost, "/login", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	h.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized{
		t.Errorf("Esperado status 401, recebeu %d", w.Code)
	}
}


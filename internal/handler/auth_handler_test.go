package handler

import (
	"bytes"
	"cinefinder/internal/model"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

type mockAuthService struct{}

func (m *mockAuthService) GenerateToken(user model.User) (string, error) {
	return "mock.jwt.token", nil
}
func (m *mockAuthService) GenerateRefreshToken(_ context.Context, _ int) (string, error) {
	return "mock_refresh_token", nil
}
func (m *mockAuthService) ValidateRefreshToken(_ context.Context, token string) (*model.RefreshToken, error) {
	if token == "valid_refresh_token" {
		return &model.RefreshToken{ID: 1, UserID: 1, Token: token}, nil
	}
	return nil, errors.New("refresh token inválido")
}
func (m *mockAuthService) RevokeRefreshToken(_ context.Context, _ string) error {
	return nil
}

type mockAuthUserService struct {
	mockUserService
}

func (m *mockAuthUserService) ValidateUser(email, password string) (*model.User, error) {
	if email == "admin@cinefinder.com" && password == "123456" {
		return &model.User{ID: 1, Email: email, Name: "Admin"}, nil
	}
	return nil, errors.New("credenciais inválidas")
}

func TestLoginHandler_Success(t *testing.T) {
	h := LoginHandler(&mockAuthService{}, &mockAuthUserService{})

	body, _ := json.Marshal(map[string]string{
		"email": "admin@cinefinder.com", "password": "123456",
	})
	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(body))
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("esperado 200, recebeu %d", w.Code)
	}
	var resp map[string]string
	json.NewDecoder(w.Body).Decode(&resp)
	if _, ok := resp["token"]; !ok {
		t.Error("resposta não contém 'token'")
	}
	if _, ok := resp["refresh_token"]; !ok {
		t.Error("resposta não contém 'refresh_token'")
	}
}

func TestLoginHandler_CredenciaisInvalidas(t *testing.T) {
	h := LoginHandler(&mockAuthService{}, &mockAuthUserService{})

	body, _ := json.Marshal(map[string]string{
		"email": "errado@cinefinder.com", "password": "errada",
	})
	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(body))
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("esperado 401, recebeu %d", w.Code)
	}
}

func TestLoginHandler_CamposObrigatorios(t *testing.T) {
	h := LoginHandler(&mockAuthService{}, &mockAuthUserService{})

	body, _ := json.Marshal(map[string]string{"email": "test@test.com"}) // sem senha
	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(body))
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("esperado 400, recebeu %d", w.Code)
	}
}

func TestRefreshHandler_Success(t *testing.T) {
	h := RefreshHandler(&mockAuthService{}, &mockAuthUserService{})

	body, _ := json.Marshal(map[string]string{"refresh_token": "valid_refresh_token"})
	req := httptest.NewRequest(http.MethodPost, "/refresh", bytes.NewBuffer(body))
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("esperado 200, recebeu %d", w.Code)
	}
	var resp map[string]string
	json.NewDecoder(w.Body).Decode(&resp)
	if _, ok := resp["token"]; !ok {
		t.Error("resposta não contém 'token'")
	}
	if _, ok := resp["refresh_token"]; !ok {
		t.Error("resposta não contém 'refresh_token'")
	}
}

func TestRefreshHandler_TokenInvalido(t *testing.T) {
	h := RefreshHandler(&mockAuthService{}, &mockAuthUserService{})

	body, _ := json.Marshal(map[string]string{"refresh_token": "token_invalido"})
	req := httptest.NewRequest(http.MethodPost, "/refresh", bytes.NewBuffer(body))
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("esperado 401, recebeu %d", w.Code)
	}
}

func TestRefreshHandler_TokenAusente(t *testing.T) {
	h := RefreshHandler(&mockAuthService{}, &mockAuthUserService{})

	body, _ := json.Marshal(map[string]string{})
	req := httptest.NewRequest(http.MethodPost, "/refresh", bytes.NewBuffer(body))
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("esperado 400, recebeu %d", w.Code)
	}
}

func TestLogoutHandler_Success(t *testing.T) {
	h := LogoutHandler(&mockAuthService{})

	body, _ := json.Marshal(map[string]string{"refresh_token": "qualquer_token"})
	req := httptest.NewRequest(http.MethodPost, "/logout", bytes.NewBuffer(body))
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("esperado 200, recebeu %d", w.Code)
	}
}

func TestLogoutHandler_TokenAusente(t *testing.T) {
	h := LogoutHandler(&mockAuthService{})

	body, _ := json.Marshal(map[string]string{})
	req := httptest.NewRequest(http.MethodPost, "/logout", bytes.NewBuffer(body))
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("esperado 400, recebeu %d", w.Code)
	}
}

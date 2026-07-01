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

type mockRefreshService struct{}

func (m *mockRefreshService) GenerateToken(user model.User) (string, error) {
	return "new.jwt.token", nil
}
func (m *mockRefreshService) GenerateRefreshToken(_ context.Context, _ int) (string, error) {
	return "new_refresh_token", nil
}
func (m *mockRefreshService) ValidateRefreshToken(_ context.Context, token string) (*model.RefreshToken, error) {
	if token == "valid_refresh_token" {
		return &model.RefreshToken{ID: 1, UserID: 1, Token: token}, nil
	}
	return nil, errors.New("refresh token inválido")
}
func (m *mockRefreshService) RevokeRefreshToken(_ context.Context, _ string) error {
	return nil
}

func TestRefreshHandler_Success(t *testing.T) {
	h := RefreshHandler(&mockRefreshService{}, &mockUserService{})

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
	h := RefreshHandler(&mockRefreshService{}, &mockUserService{})

	body, _ := json.Marshal(map[string]string{"refresh_token": "token_invalido"})
	req := httptest.NewRequest(http.MethodPost, "/refresh", bytes.NewBuffer(body))
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("esperado 401, recebeu %d", w.Code)
	}
}

func TestRefreshHandler_TokenAusente(t *testing.T) {
	h := RefreshHandler(&mockRefreshService{}, &mockUserService{})

	body, _ := json.Marshal(map[string]string{})
	req := httptest.NewRequest(http.MethodPost, "/refresh", bytes.NewBuffer(body))
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("esperado 400, recebeu %d", w.Code)
	}
}

func TestLogoutHandler_Success(t *testing.T) {
	h := LogoutHandler(&mockRefreshService{})

	body, _ := json.Marshal(map[string]string{"refresh_token": "qualquer_token"})
	req := httptest.NewRequest(http.MethodPost, "/logout", bytes.NewBuffer(body))
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("esperado 200, recebeu %d", w.Code)
	}
}

func TestLogoutHandler_TokenAusente(t *testing.T) {
	h := LogoutHandler(&mockRefreshService{})

	body, _ := json.Marshal(map[string]string{})
	req := httptest.NewRequest(http.MethodPost, "/logout", bytes.NewBuffer(body))
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("esperado 400, recebeu %d", w.Code)
	}
}

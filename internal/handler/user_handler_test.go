package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"cinefinder/internal/model"
)

// Mock do service
type mockUserService struct {
	users []model.User
}

func (m *mockUserService) Create(user model.User) (*model.User, error) {
	for _, u := range m.users {
		if u.Email == user.Email {
			return nil, errors.New("Usuário já cadastrado")
		}
	}
	user.ID = len(m.users) + 1
	m.users = append(m.users, user)
	return &user, nil
}

// Para implementar ValidateUser
func (m *mockUserService) ValidateUser(email, password string) (*model.User, error) {
	for _, u := range m.users {
		if u.Email == email && u.Password == password {
			return &u, nil
		}
	}
	return nil, errors.New("credenciais inválidas")
}

func TestCreateUser_Success(t *testing.T) {
	mockService := &mockUserService{}
	handler := NewUserHandler(mockService)

	body := model.User{
		Name:     "Artemis",
		Email:    "artemis@nasa.com",
		Password: "password",
	}

	jsonBody, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()

	handler.Create(w, req)

	resp := w.Result()

	if resp.StatusCode != http.StatusCreated {
		t.Errorf("esperado 201, veio %d", resp.StatusCode)
	}

	var response model.User
	json.NewDecoder(resp.Body).Decode(&response)

	if response.ID != 1 {
		t.Errorf("esperado ID 1, veio %d", response.ID)
	}
}

func TestCreateUser_AlreadyExists(t *testing.T) {
	mockService := &mockUserService{}
	handler := NewUserHandler(mockService)

	body := model.User{
		Name:     "Artemis",
		Email:    "artemis@nasa.com",
		Password: "password",
	}

	jsonBody, _ := json.Marshal(body)

	// Primeira requisição: deve criar com sucesso (201)
	req1 := httptest.NewRequest(http.MethodPost, "/users", bytes.NewBuffer(jsonBody))
	req1.Header.Set("Content-Type", "application/json")
	w1 := httptest.NewRecorder()
	handler.Create(w1, req1)

	if w1.Result().StatusCode != http.StatusCreated {
		t.Errorf("primeira requisição: esperado 201, veio %d", w1.Result().StatusCode)
	}

	// Segunda requisição (mesmo e-mail): deve falhar (400)
	req2 := httptest.NewRequest(http.MethodPost, "/users", bytes.NewBuffer(jsonBody))
	req2.Header.Set("Content-Type", "application/json")
	w2 := httptest.NewRecorder()
	handler.Create(w2, req2)

	resp2 := w2.Result()

	if resp2.StatusCode != http.StatusBadRequest {
		t.Errorf("segunda requisição: esperado 400, veio %d", resp2.StatusCode)
	}

	buf := new(bytes.Buffer)
	buf.ReadFrom(resp2.Body)
	responseStr := buf.String()

	if !strings.Contains(responseStr, "Usuário já cadastrado") {
		t.Errorf("esperado mensagem de erro 'Usuário já cadastrado', veio: %s", responseStr)
	}
}

func TestCreateUser_InvalidJSON(t *testing.T) {
	mockService := &mockUserService{}
	handler := NewUserHandler(mockService)

	req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewBuffer([]byte("json inválido")))
	w := httptest.NewRecorder()

	handler.Create(w, req)

	if w.Result().StatusCode != http.StatusBadRequest {
		t.Errorf("esperado 400, veio %d", w.Result().StatusCode)
	}
}

func (m *mockUserService) List() ([]model.User, error) {
	return []model.User{
		{
			ID:       1,
			Name:     "Artemis",
			Email:    "artemis@nasa.com",
			Password: "password",
		},
	}, nil
}

func (m *mockUserService) GetByID(id int) (*model.User, error) {
	return &model.User{
		ID:       id,
		Name:     "Artemis",
		Email:    "artemis@nasa.com",
		Password: "password",
	}, nil
}

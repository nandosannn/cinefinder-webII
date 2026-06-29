package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"cinefinder/internal/config"

	"github.com/golang-jwt/jwt/v5"
)

func gerarTokenTeste(userID int) string {
	claims := jwt.MapClaims{
		"user_id": float64(userID),
		"email":   "test@cinefinder.com",
		"exp":     time.Now().Add(time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, _ := token.SignedString(config.GetJWTKey())
	return signed
}

func TestAuthMiddleware_SemToken(t *testing.T) {
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	h := AuthMiddleware(next)

	req := httptest.NewRequest(http.MethodGet, "/movies", nil)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("esperado 401, veio %d", w.Code)
	}
}

func TestAuthMiddleware_FormatoInvalido(t *testing.T) {
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	h := AuthMiddleware(next)

	req := httptest.NewRequest(http.MethodGet, "/movies", nil)
	req.Header.Set("Authorization", "TokenSemBearer")
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("esperado 401, veio %d", w.Code)
	}
}

func TestAuthMiddleware_TokenValido(t *testing.T) {
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Context().Value(UserIdKey) == nil {
			t.Error("user_id não encontrado no contexto")
		}
		w.WriteHeader(http.StatusOK)
	})
	h := AuthMiddleware(next)

	token := gerarTokenTeste(42)
	req := httptest.NewRequest(http.MethodGet, "/movies", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("esperado 200, veio %d", w.Code)
	}
}

func TestAuthMiddleware_TokenExpirado(t *testing.T) {
	claims := jwt.MapClaims{
		"user_id": float64(1),
		"email":   "test@cinefinder.com",
		"exp":     time.Now().Add(-time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, _ := token.SignedString(config.GetJWTKey())

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	h := AuthMiddleware(next)

	req := httptest.NewRequest(http.MethodGet, "/movies", nil)
	req.Header.Set("Authorization", "Bearer "+tokenStr)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("esperado 401, veio %d", w.Code)
	}
}

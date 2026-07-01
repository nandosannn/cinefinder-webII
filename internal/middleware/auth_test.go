package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func gerarTokenTeste(userID int, expirado bool) string {
	exp := time.Now().Add(time.Hour)
	if expirado {
		exp = time.Now().Add(-time.Hour)
	}
	claims := jwt.MapClaims{
		"user_id": float64(userID),
		"email":   "test@cinefinder.com",
		"exp":     exp.Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, _ := token.SignedString(jwtKey)
	return signed
}

func TestAuthMiddleware_SemToken(t *testing.T) {
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusOK) })
	h := AuthMiddleware(next)

	req := httptest.NewRequest(http.MethodGet, "/movies", nil)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("esperado 401, veio %d", w.Code)
	}
}

func TestAuthMiddleware_FormatoInvalido(t *testing.T) {
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusOK) })
	h := AuthMiddleware(next)

	req := httptest.NewRequest(http.MethodGet, "/movies", nil)
	req.Header.Set("Authorization", "SemPrefixoBearer")
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

	req := httptest.NewRequest(http.MethodGet, "/movies", nil)
	req.Header.Set("Authorization", "Bearer "+gerarTokenTeste(42, false))
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("esperado 200, veio %d", w.Code)
	}
}

func TestAuthMiddleware_TokenExpirado(t *testing.T) {
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusOK) })
	h := AuthMiddleware(next)

	req := httptest.NewRequest(http.MethodGet, "/movies", nil)
	req.Header.Set("Authorization", "Bearer "+gerarTokenTeste(1, true))
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("esperado 401, veio %d", w.Code)
	}
}

package handler

import (
	"cinefinder/internal/service"
	"encoding/json"
	"net/http"
)

func RefreshHandler(authService service.RefreshServiceInterface, userService service.UserServiceInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req RefreshRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Requisição inválida", http.StatusBadRequest)
			return
		}
		if req.RefreshToken == "" {
			http.Error(w, "Refresh token é obrigatório", http.StatusBadRequest)
			return
		}

		rt, err := authService.ValidateRefreshToken(r.Context(), req.RefreshToken)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		if err := authService.RevokeRefreshToken(r.Context(), req.RefreshToken); err != nil {
			http.Error(w, "Erro interno", http.StatusInternalServerError)
			return
		}

		user, err := userService.GetByID(rt.UserID)
		if err != nil {
			http.Error(w, "Usuário não encontrado", http.StatusUnauthorized)
			return
		}

		newToken, err := authService.GenerateToken(*user)
		if err != nil {
			http.Error(w, "Erro ao gerar token", http.StatusInternalServerError)
			return
		}

		newRefreshToken, err := authService.GenerateRefreshToken(r.Context(), user.ID)
		if err != nil {
			http.Error(w, "Erro ao gerar refresh token", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"token":         newToken,
			"refresh_token": newRefreshToken,
		})
	}
}

func LogoutHandler(authService service.RefreshServiceInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req LogoutRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Requisição inválida", http.StatusBadRequest)
			return
		}
		if req.RefreshToken == "" {
			http.Error(w, "Refresh token é obrigatório", http.StatusBadRequest)
			return
		}

		if err := authService.RevokeRefreshToken(r.Context(), req.RefreshToken); err != nil {
			http.Error(w, "Erro ao revogar token", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"message": "Logout realizado com sucesso"})
	}
}

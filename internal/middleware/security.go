package middleware

import "net/http"

func SecurityHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Previne MIME sniffing
		w.Header().Set("X-Content-Type-Options", "nosniff")
		// Protege contra clickjacking
		w.Header().Set("X-Frame-Options", "DENY")
		// Proteção XSS em browsers antigos
		w.Header().Set("X-XSS-Protection", "1; mode=block")
		// Força HTTPS por 1 ano
		w.Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains; preload")
		// Restringe origens de conteúdo
		w.Header().Set("Content-Security-Policy", "default-src 'self'; frame-ancestors 'none'")
		// Controle de referrer
		w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")
		// Desabilita recursos sensíveis do browser
		w.Header().Set("Permissions-Policy", "geolocation=(), microphone=(), camera=()")

		next.ServeHTTP(w, r)
	})
}

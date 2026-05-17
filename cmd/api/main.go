package main

import (
	"cinefinder/internal/db"
	"cinefinder/internal/handler"
	"cinefinder/internal/middleware"
	"cinefinder/internal/service"
	"net/http"

	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
)

func main() {

	// carregar .env
	if err := godotenv.Load(); err != nil {
		println("Aviso: .env não carregado")
	}

	// conectar banco
	dbPool := db.NewDB()
	defer dbPool.Close()

	// criar tabelas
	db.RunMigrations(dbPool)

	// services
	movieService := service.NewMovieService(dbPool)
	loanService := service.NewLoanService(dbPool)
	userService := service.NewUserService(dbPool)
	authService := &service.AuthService{}

	// handlers
	movieHandler := handler.NewMovieHandler(movieService)
	loanHandler := handler.NewLoanHandler(loanService)
	userHandler := handler.NewUserHandler(userService)

	// router
	r := chi.NewRouter()
	r.Use(chiMiddleware.Logger)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status": "ok", "message": "Cinefinder API is running 🚀"}`))
	})

	// Autenticação 
	r.Post("/login", handler.LoginHandler(authService, userService))

	// Rotas públicas de filmes
	r.Get("/movies", movieHandler.List)
	r.Get("/movies/{id}", movieHandler.GetByID)

	// Rotas protegidas — exigem token JWT no header Authorization: Bearer <token>
	r.Group(func(r chi.Router) {
		r.Use(middleware.AuthMiddleware)
		r.Post("/movies", movieHandler.Create)
	})

	// Usuários
	r.Post("/users", userHandler.Create)
	r.Get("/users", userHandler.List)
	r.Get("/users/{id}", userHandler.GetByID)

	r.Post("/loans", loanHandler.Create)
	r.Get("/loans", loanHandler.List)
	r.Get("/loans/{id}", loanHandler.GetByID)

	// subir servidor
	println("Servidor rodando em http://localhost:3000 🚀")
	http.ListenAndServe(":3000", r)
}

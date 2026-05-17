package main

import (
	"net/http"

	"cinefinder/internal/database"
	"cinefinder/internal/db"
	"cinefinder/internal/handler"
	"cinefinder/internal/middleware"
	"cinefinder/internal/service"

	"github.com/go-chi/chi/v5"
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

	// criar queries do sqlc
	queries := database.New(dbPool)

	// criar tabela
	db.RunMigrations(dbPool)

	// service + handler
	movieService := service.NewMovieService(queries)
	movieHandler := handler.NewMovieHandler(movieService)

	loanService := service.NewLoanService(dbPool)
	loanHandler := handler.NewLoanHandler(loanService)

	userService := service.NewUserService(dbPool)
	userHandler := handler.NewUserHandler(userService)

	authService := &service.AuthService{}
	loginHandler := handler.LoginHandler(authService, userService)

	// router
	r := chi.NewRouter()

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status": "ok", "message": "Cinefinder API is running 🚀"}`))
	})

	r.Post("/users", userHandler.Create)
	r.Post("/login", loginHandler)

	r.Group(func(r chi.Router) {
		r.Use(middleware.AuthMiddleware)

		r.Post("/movies", movieHandler.Create)
		r.Get("/movies", movieHandler.List)
		r.Get("/movies/{id}", movieHandler.GetByID)

		r.Post("/loans", loanHandler.Create)
		r.Get("/loans", loanHandler.List)
		r.Get("/loans/{id}", loanHandler.GetByID)
		r.Put("/loans/{id}/return", loanHandler.ReturnMovie)

		r.Get("/users", userHandler.List)
		r.Get("/users/{id}", userHandler.GetByID)
		r.Put("/movies/{id}", movieHandler.Update)
		r.Delete("/movies/{id}", movieHandler.Delete)
	})

	// subir servidor
	println("Servidor rodando em http://localhost:3000")
	if err := http.ListenAndServe(":3000", r); err != nil {
		println("Erro ao iniciar servidor:", err.Error())
	}
}

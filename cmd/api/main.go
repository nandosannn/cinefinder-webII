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

	authService := service.NewAuthService(dbPool)

	// router
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.SecurityHeaders)
	r.Use(middleware.RateLimit)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"status": "ok", "message": "Cinefinder API is running 🚀"}`))
	})

	r.Post("/users", userHandler.Create)
	r.Post("/login", handler.LoginHandler(authService, userService))
	r.Post("/refresh", handler.RefreshHandler(authService, userService))
	r.Post("/logout", handler.LogoutHandler(authService))

	r.Group(func(r chi.Router) {
		r.Use(middleware.AuthMiddleware)

		r.Post("/movies", movieHandler.Create)
		r.Get("/movies", movieHandler.List)
		r.Get("/movies/{id}", movieHandler.GetByID)
		r.Put("/movies/{id}", movieHandler.Update)
		r.Delete("/movies/{id}", movieHandler.Delete)

		r.Post("/loans", loanHandler.Create)
		r.Get("/loans", loanHandler.List)
		r.Get("/loans/history", loanHandler.History)
		r.Get("/loans/{id}", loanHandler.GetByID)
		r.Put("/loans/{id}/return", loanHandler.ReturnMovie)

		r.Get("/users", userHandler.List)
		r.Get("/users/{id}", userHandler.GetByID)
	})

	// subir servidor
	println("Servidor rodando em http://localhost:3000")
	if err := http.ListenAndServe(":3000", r); err != nil {
		println("Erro ao iniciar servidor:", err.Error())
	}
}

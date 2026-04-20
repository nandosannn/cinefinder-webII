package main

import (
	"cinefinder/internal/db"
	"cinefinder/internal/handler"
	"cinefinder/internal/service"
	"net/http"
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

	// criar tabela
	db.RunMigrations(dbPool)

	// service + handler
	movieService := service.NewMovieService(dbPool)
	movieHandler := handler.NewMovieHandler(movieService)

	// router
	r := chi.NewRouter()
	r.Post("/movies", movieHandler.Create)

	// subir servidor
	println("Servidor rodando em http://localhost:3000 🚀")
	http.ListenAndServe(":3000", r)
}

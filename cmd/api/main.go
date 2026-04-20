package main

import (
	"cinefinder/internal/db"

	"github.com/joho/godotenv"
)

func main() {

	if err := godotenv.Load(); err != nil {
		println("Aviso: .env não carregado")
	}

	dbPool := db.NewDB()
	defer dbPool.Close()

	db.RunMigrations(dbPool)
}

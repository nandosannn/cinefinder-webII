package db

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

func RunMigrations(pool *pgxpool.Pool) {

	// tabela movies
	query := `
	CREATE TABLE IF NOT EXISTS movies (
		id SERIAL PRIMARY KEY,
		title TEXT NOT NULL,
		director TEXT NOT NULL,
		release_year INT NOT NULL,
		genre TEXT NOT NULL,
		available_copies INT NOT NULL DEFAULT 1
	);
	`

	_, err := pool.Exec(context.Background(), query)
	if err != nil {
		log.Fatalf("Erro ao criar tabela movies: %v", err)
	}

	log.Println("Tabela movies pronta ")

	// tabela users
	query = `
	CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		name TEXT NOT NULL,
		email TEXT NOT NULL,
		password TEXT NOT NULL,
		created_at TIMESTAMP DEFAULT NOW()
	);
	`

	_, err = pool.Exec(context.Background(), query)
	if err != nil {
		log.Fatalf("Erro ao criar tabela users: %v", err)
	}

	log.Println("Tabela users pronta ")

	// tabela loans
	query = `
	CREATE TABLE IF NOT EXISTS loans (
		id SERIAL PRIMARY KEY,
		user_id INT NOT NULL REFERENCES users(id),
		movie_id INT NOT NULL REFERENCES movies(id),
		loan_date TIMESTAMP DEFAULT NOW(),
		return_date TIMESTAMP DEFAULT NULL,
		price DECIMAL(10, 2) NOT NULL,
		returned BOOLEAN DEFAULT false
	);
	`

	_, err = pool.Exec(context.Background(), query)
	if err != nil {
		log.Fatalf("Erro ao criar tabela loans: %v", err)
	}

	log.Println("Tabela loans pronta ")
}

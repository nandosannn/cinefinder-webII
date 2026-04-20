package db

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"os"
	"time"
)

func NewDB() *pgxpool.Pool {
	databaseURL := os.Getenv("DATABASE_URL")

	if databaseURL == "" {
		log.Fatal("DATABASE_URL não definida")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pool, err := pgxpool.New(ctx, databaseURL)
	if err != nil {
		log.Fatalf("Erro ao conectar no banco: %v", err)
	}

	if err := pool.Ping(ctx); err != nil {
		log.Fatalf("Erro ao pingar banco: %v", err)
	}

	log.Println("Conectado ao PostgreSQL 🚀")

	return pool
}

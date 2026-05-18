package db

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

var schema string

func RunMigrations(pool *pgxpool.Pool) {
	_, err := pool.Exec(context.Background(), schema)
	if err != nil {
		log.Fatalf("Erro ao aplicar schema: %v", err)
	}
	log.Println("Schema aplicado com sucesso")
}

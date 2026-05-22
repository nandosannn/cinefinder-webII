package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

type Movie struct {
	Title           string
	Director        string
	ReleaseYear     int
	Genre           string
	AvailableCopies int
}

type User struct {
	Name     string
	Email    string
	Password string
}

func main() {
	// Carregar .env
	if err := godotenv.Load(); err != nil {
		log.Println("Aviso: .env não carregado, usando variáveis de ambiente do sistema")
	}

	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		log.Fatal("Erro: DATABASE_URL não definida no arquivo .env ou nas variáveis de ambiente")
	}

	ctx := context.Background()

	// Conectar ao banco
	pool, err := pgxpool.New(ctx, databaseURL)
	if err != nil {
		log.Fatalf("Erro ao conectar no banco: %v", err)
	}
	defer pool.Close()

	if err := pool.Ping(ctx); err != nil {
		log.Fatalf("Erro ao pingar o banco: %v", err)
	}

	fmt.Println("🌱 Iniciando a geração de dados mockados no PostgreSQL...")

	// 1. Limpar dados antigos
	fmt.Println("🧹 Limpando dados antigos (TRUNCATE)...")
	_, err = pool.Exec(ctx, "TRUNCATE TABLE loans, users, movies RESTART IDENTITY CASCADE;")
	if err != nil {
		log.Fatalf("Erro ao limpar tabelas: %v", err)
	}

	// 2. Inserir Filmes
	movies := []Movie{
		{"Interestelar", "Christopher Nolan", 2014, "Ficção Científica", 3},
		{"O Poderoso Chefão", "Francis Ford Coppola", 1972, "Drama", 2},
		{"Matrix", "Lana Wachowski", 1999, "Ficção Científica", 4},
		{"A Origem", "Christopher Nolan", 2010, "Ação", 3},
		{"Pulp Fiction", "Quentin Tarantino", 1994, "Crime", 2},
		{"Cidade de Deus", "Fernando Meirelles", 2002, "Drama", 3},
		{"Clube da Luta", "David Fincher", 1999, "Drama", 2},
		{"Vingadores: Ultimato", "Anthony Russo", 2019, "Aventura", 5},
		{"O Rei Leão", "Roger Allers", 1994, "Animação", 3},
		{"Parasita", "Bong Joon Ho", 2019, "Suspense", 2},
	}

	fmt.Println("🎬 Cadastrando filmes...")
	movieIDs := []int{}
	for _, m := range movies {
		var id int
		err := pool.QueryRow(ctx, `
			INSERT INTO movies (title, director, release_year, genre, available_copies)
			VALUES ($1, $2, $3, $4, $5) RETURNING id;`,
			m.Title, m.Director, m.ReleaseYear, m.Genre, m.AvailableCopies,
		).Scan(&id)
		if err != nil {
			log.Fatalf("Erro ao inserir filme %s: %v", m.Title, err)
		}
		movieIDs = append(movieIDs, id)
	}
	fmt.Printf("✅ %d filmes cadastrados com sucesso!\n", len(movies))

	// 3. Inserir Usuários
	users := []User{
		{"Administrador", "admin@email.com", "123456"},
		{"João Silva", "joao@email.com", "123456"},
		{"Maria Souza", "maria@email.com", "123456"},
		{"Carlos Santos", "carlos@email.com", "123456"},
		{"Ana Oliveira", "ana@email.com", "123456"},
	}

	fmt.Println("👥 Cadastrando usuários...")
	userIDs := []int{}
	for _, u := range users {
		var id int
		err := pool.QueryRow(ctx, `
			INSERT INTO users (name, email, password, created_at)
			VALUES ($1, $2, $3, NOW() - INTERVAL '10 days') RETURNING id;`,
			u.Name, u.Email, u.Password,
		).Scan(&id)
		if err != nil {
			log.Fatalf("Erro ao inserir usuário %s: %v", u.Name, err)
		}
		userIDs = append(userIDs, id)
	}
	fmt.Printf("✅ %d usuários cadastrados com sucesso!\n", len(users))

	// 4. Inserir Empréstimos (Loans)
	// Vamos criar um histórico rico cobrindo vários meses, devoluções feitas e algumas em aberto.
	fmt.Println("📅 Gerando histórico de empréstimos...")
	now := time.Now()

	// Lista de empréstimos fictícios estruturados para teste de filtragem de datas e filmes
	type loanData struct {
		userID    int
		movieID   int
		loanDate  time.Time
		retDate   *time.Time
		price     float64
		returned  bool
	}

	// Helper para criar ponteiro de time.Time
	tPtr := func(t time.Time) *time.Time { return &t }

	loansToInsert := []loanData{
		// Janeiro 2026
		{userIDs[1], movieIDs[0], now.AddDate(0, -4, -15), tPtr(now.AddDate(0, -4, -8)), 12.50, true},
		{userIDs[2], movieIDs[1], now.AddDate(0, -4, -10), tPtr(now.AddDate(0, -4, -3)), 15.00, true},
		
		// Fevereiro 2026
		{userIDs[3], movieIDs[2], now.AddDate(0, -3, -20), tPtr(now.AddDate(0, -3, -15)), 10.00, true},
		{userIDs[4], movieIDs[0], now.AddDate(0, -3, -5), tPtr(now.AddDate(0, -3, -1)), 12.50, true},
		
		// Março 2026
		{userIDs[1], movieIDs[3], now.AddDate(0, -2, -12), tPtr(now.AddDate(0, -2, -5)), 15.00, true},
		{userIDs[2], movieIDs[4], now.AddDate(0, -2, -2), tPtr(now.AddDate(0, -2, 5)), 12.50, true},
		
		// Abril 2026
		{userIDs[3], movieIDs[5], now.AddDate(0, -1, -25), tPtr(now.AddDate(0, -1, -18)), 10.00, true},
		{userIDs[4], movieIDs[1], now.AddDate(0, -1, -10), tPtr(now.AddDate(0, -1, -3)), 15.00, true},
		{userIDs[1], movieIDs[6], now.AddDate(0, -1, -5), tPtr(now.AddDate(0, -1, 0)), 12.50, true},

		// Maio 2026 (Atuais)
		{userIDs[2], movieIDs[7], now.AddDate(0, 0, -8), tPtr(now.AddDate(0, 0, -2)), 20.00, true},
		
		// Empréstimos Ativos (Não devolvidos)
		{userIDs[3], movieIDs[8], now.AddDate(0, 0, -3), nil, 15.00, false},
		{userIDs[4], movieIDs[9], now.AddDate(0, 0, -1), nil, 15.00, false},
	}

	// Adicionar mais alguns empréstimos randômicos históricos para volumetria
	r := rand.New(rand.NewSource(42))
	for i := 0; i < 15; i++ {
		uIdx := r.Intn(len(userIDs))
		mIdx := r.Intn(len(movieIDs))
		
		// Dias atrás (entre 10 e 120 dias atrás)
		daysAgo := r.Intn(110) + 10
		loanD := now.AddDate(0, 0, -daysAgo)
		
		// Devolvido após 3 a 10 dias
		duration := r.Intn(7) + 3
		retD := loanD.AddDate(0, 0, duration)
		
		price := 10.0 + float64(r.Intn(11)) // 10.00 a 20.00
		
		loansToInsert = append(loansToInsert, loanData{
			userID:   userIDs[uIdx],
			movieID:  movieIDs[mIdx],
			loanDate: loanD,
			retDate:  &retD,
			price:    price,
			returned: true,
		})
	}

	count := 0
	for _, l := range loansToInsert {
		_, err := pool.Exec(ctx, `
			INSERT INTO loans (user_id, movie_id, loan_date, return_date, price, returned)
			VALUES ($1, $2, $3, $4, $5, $6);`,
			l.userID, l.movieID, l.loanDate, l.retDate, l.price, l.returned,
		)
		if err != nil {
			log.Printf("Aviso: Falha ao inserir empréstimo do usuário %d com filme %d: %v", l.userID, l.movieID, err)
			continue
		}
		count++
	}

	fmt.Printf("✅ %d registros de empréstimos inseridos com sucesso!\n", count)
	fmt.Println("\n🎉 Geração de dados de teste finalizada com sucesso!")
	fmt.Println("🚀 Agora você pode filtrar o histórico por filme ou intervalo de datas perfeitamente!")
}

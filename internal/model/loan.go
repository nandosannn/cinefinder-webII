package model

import "time"

type Loan struct {
	ID         int       `json:"id"`
	LoanDate   time.Time `json:"loan_date"`
	ReturnDate time.Time `json:"return_date"`
	Price      float64   `json:"price"`
	Returned   bool      `json:"returned"`

	UserID  int `json:"user_id"`
	MovieID int `json:"movie_id"`

	User  User  `json:"user,omitempty"`
	Movie Movie `json:"movie,omitempty"`
}

// LoanFilter contém os parâmetros de filtragem para o histórico de empréstimos
type LoanFilter struct {
	MovieID   int        // filtrar por filme específico
	StartDate *time.Time // data inicial (loan_date >= start_date)
	EndDate   *time.Time // data final (loan_date <= end_date)
}

// LoanHistory representa um empréstimo com dados enriquecidos para o histórico
type LoanHistory struct {
	ID         int        `json:"id"`
	LoanDate   time.Time  `json:"loan_date"`
	ReturnDate *time.Time `json:"return_date"`
	Price      float64    `json:"price"`
	Returned   bool       `json:"returned"`

	UserID    int    `json:"user_id"`
	UserName  string `json:"user_name"`
	UserEmail string `json:"user_email"`

	MovieID    int    `json:"movie_id"`
	MovieTitle string `json:"movie_title"`
}
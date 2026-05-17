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
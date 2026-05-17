package model

import "time"

type Reservation struct {
	ID              int       `json:"id"`
	UserID          int       `json:"user_id"`
	MovieID         int       `json:"movie_id"`
	ReservationDate time.Time `json:"reservation_date"`
	ReturnDate      time.Time `json:"return_date"`
	ReservationFee  float64   `json:"reservation_fee"`

	// Relacionamentos
	User  User  `json:"user"`
	Movie Movie `json:"movie"`
}

package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"cinefinder/internal/model"

	"github.com/jackc/pgx/v5/pgxpool"
)

type LoanServiceInterface interface {
	Create(loan model.Loan) (*model.Loan, error)
	List() ([]model.Loan, error)
	GetByID(id int) (*model.Loan, error)
	ReturnMovie(id int) error
	ListHistory(filter model.LoanFilter) ([]model.LoanHistory, error)
}

type LoanService struct {
	db *pgxpool.Pool
}

func NewLoanService(db *pgxpool.Pool) *LoanService {
	return &LoanService{db: db}
}

func (s *LoanService) Create(loan model.Loan) (*model.Loan, error) {

	// verificar empréstimo em aberto
	var unreturnedCount int

	checkLoanQuery := `
	SELECT COUNT(*)
	FROM loans
	WHERE user_id = $1
	AND returned = false
	`

	err := s.db.QueryRow(
		context.Background(),
		checkLoanQuery,
		loan.UserID,
	).Scan(&unreturnedCount)

	if err != nil {
		return nil, err
	}

	if unreturnedCount > 0 {
		return nil, errors.New("Usuário possui empréstimo em aberto")
	}

	// buscar quantidade total de cópias do filme
	var totalCopies int

	checkMovieQuery := `
	SELECT available_copies
	FROM movies
	WHERE id = $1
	`

	err = s.db.QueryRow(
		context.Background(),
		checkMovieQuery,
		loan.MovieID,
	).Scan(&totalCopies)

	if err != nil {
		return nil, err
	}

	// verificar quantidade de empréstimos ativos
	var activeLoans int

	activeLoansQuery := `
	SELECT COUNT(*)
	FROM loans
	WHERE movie_id = $1
	AND returned = false
	`

	err = s.db.QueryRow(
		context.Background(),
		activeLoansQuery,
		loan.MovieID,
	).Scan(&activeLoans)

	if err != nil {
		return nil, err
	}

	// verificar disponibilidade real
	if activeLoans >= totalCopies {
		return nil, errors.New(
			"Filme indisponível",
		)
	}

	// validar datas
	if !loan.ReturnDate.After(time.Now()) {
		return nil, errors.New("A data de devolução deve ser posterior à data atual")
	}

	// criar empréstimo
	insertQuery := `
	INSERT INTO loans (
		user_id,
		movie_id,
		loan_date,
		return_date,
		price,
		returned
	)
	VALUES ($1, $2, NOW(), $3, $4, $5)
	RETURNING id, loan_date;
	`

	err = s.db.QueryRow(
		context.Background(),
		insertQuery,
		loan.UserID,
		loan.MovieID,
		loan.ReturnDate,
		loan.Price,
		loan.Returned,
	).Scan(
		&loan.ID,
		&loan.LoanDate,
	)

	if err != nil {
		return nil, err
	}

	return &loan, nil
}

func (s *LoanService) List() ([]model.Loan, error) {

	query := `
	SELECT
		id,
		user_id,
		movie_id,
		loan_date,
		return_date,
		price,
		returned
	FROM loans
	`

	rows, err := s.db.Query(context.Background(), query)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	loans := []model.Loan{}

	for rows.Next() {

		var l model.Loan

		err := rows.Scan(
			&l.ID,
			&l.UserID,
			&l.MovieID,
			&l.LoanDate,
			&l.ReturnDate,
			&l.Price,
			&l.Returned,
		)

		if err != nil {
			return nil, err
		}

		loans = append(loans, l)
	}

	return loans, nil
}

func (s *LoanService) GetByID(id int) (*model.Loan, error) {

	query := `
	SELECT
		id,
		user_id,
		movie_id,
		loan_date,
		return_date,
		price,
		returned
	FROM loans
	WHERE id = $1
	`

	var l model.Loan

	err := s.db.QueryRow(
		context.Background(),
		query,
		id,
	).Scan(
		&l.ID,
		&l.UserID,
		&l.MovieID,
		&l.LoanDate,
		&l.ReturnDate,
		&l.Price,
		&l.Returned,
	)

	if err != nil {
		return nil, err
	}

	return &l, nil
}

func (s *LoanService) ReturnMovie(id int) error {

	var loan model.Loan

	checkQuery := `
	SELECT id, movie_id, returned
	FROM loans
	WHERE id = $1
	`

	err := s.db.QueryRow(
		context.Background(),
		checkQuery,
		id,
	).Scan(
		&loan.ID,
		&loan.MovieID,
		&loan.Returned,
	)

	if err != nil {
		return errors.New("Empréstimo não encontrado")
	}

	if loan.Returned {
		return errors.New("Empréstimo já devolvido")
	}

	updateLoanQuery := `
	UPDATE loans
	SET returned = true
	WHERE id = $1
	`

	_, err = s.db.Exec(
		context.Background(),
		updateLoanQuery,
		id,
	)

	if err != nil {
		return err
	}

	return nil
}

func (s *LoanService) ListHistory(filter model.LoanFilter) ([]model.LoanHistory, error) {
	query := `
	SELECT
		l.id,
		l.loan_date,
		l.return_date,
		l.price,
		l.returned,
		l.user_id,
		u.name AS user_name,
		u.email AS user_email,
		l.movie_id,
		m.title AS movie_title
	FROM loans l
	JOIN users u ON u.id = l.user_id
	JOIN movies m ON m.id = l.movie_id
	WHERE 1=1
	`

	args := []interface{}{}
	argIndex := 1

	if filter.MovieID > 0 {
		query += fmt.Sprintf(" AND l.movie_id = $%d", argIndex)
		args = append(args, filter.MovieID)
		argIndex++
	}

	if filter.StartDate != nil {
		query += fmt.Sprintf(" AND l.loan_date >= $%d", argIndex)
		args = append(args, *filter.StartDate)
		argIndex++
	}

	if filter.EndDate != nil {
		query += fmt.Sprintf(" AND l.loan_date <= $%d", argIndex)
		args = append(args, *filter.EndDate)
		argIndex++
	}

	query += " ORDER BY l.loan_date DESC"

	rows, err := s.db.Query(context.Background(), query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	history := []model.LoanHistory{}

	for rows.Next() {
		var h model.LoanHistory
		var returnDate *time.Time

		err := rows.Scan(
			&h.ID,
			&h.LoanDate,
			&returnDate,
			&h.Price,
			&h.Returned,
			&h.UserID,
			&h.UserName,
			&h.UserEmail,
			&h.MovieID,
			&h.MovieTitle,
		)

		if err != nil {
			return nil, err
		}

		h.ReturnDate = returnDate
		history = append(history, h)
	}

	return history, nil
}

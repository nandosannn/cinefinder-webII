-- name: CreateLoan :one
INSERT INTO loans (movie_id, user_id, return_date)
VALUES ($1, $2, $3)
RETURNING *;

-- name: ListLoans :many
SELECT * FROM loans
ORDER BY id;

-- name: GetLoanByID :one
SELECT * FROM loans
WHERE id = $1;

-- name: ListLoansByMovie :many
SELECT * FROM loans
WHERE movie_id = $1
ORDER BY id;

-- name: GetMovieWithLoans :many
SELECT
    m.id AS movie_id,
    m.title,
    m.director,
    m.release_year,
    m.genre,
    m.available_copies,
    l.id AS loan_id,
    l.user_id,
    l.loan_date,
    l.return_date
FROM movies m
LEFT JOIN loans l ON l.movie_id = m.id
WHERE m.id = $1
ORDER BY l.id;
-- name: CreateMovie :one
INSERT INTO movies (title, director, release_year, genre, available_copies)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: ListMovies :many
SELECT * FROM movies
ORDER BY id;

-- name: GetMovieByID :one
SELECT * FROM movies
WHERE id = $1;

-- name: UpdateMovie :one
UPDATE movies
SET title = $2,
    director = $3,
    release_year = $4,
    genre = $5,
    available_copies = $6
WHERE id = $1
RETURNING *;

-- name: DeleteMovie :exec
DELETE FROM movies
WHERE id = $1;
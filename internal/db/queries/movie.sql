-- name: CreateMovie :one
INSERT INTO movies (title, director, year, genre)
VALUES ($1, $2, $3, $4)
RETURNING id, title, director, year, genre;

-- name: ListMovies :many
SELECT id, title, director, year, genre
FROM movies
ORDER BY title;

-- name: GetMovieByID :one
SELECT id, title, director, year, genre
FROM movies
WHERE id = $1;

-- name: SearchMovies :many
-- Busca parcial e case-insensitive por título e/ou gênero (P3 - Sprint 2)
SELECT id, title, director, year, genre
FROM movies
WHERE ($1 = '' OR title ILIKE '%' || $1 || '%')
  AND ($2 = '' OR genre ILIKE $2)
ORDER BY title;

-- name: CheckMovieDuplicate :one
SELECT COUNT(*)
FROM movies
WHERE title = $1 AND director = $2 AND year = $3 AND genre = $4;

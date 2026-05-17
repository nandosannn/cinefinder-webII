-- name: CreateUser :one
INSERT INTO users (name, email, password, created_at)
VALUES ($1, $2, $3, NOW())
RETURNING id, name, email, password, created_at;

-- name: ListUsers :many
SELECT id, name, email, created_at
FROM users
ORDER BY name;

-- name: GetUserByID :one
SELECT id, name, email, created_at
FROM users
WHERE id = $1;

-- name: GetUserByEmailAndPassword :one
-- Usado na autenticação (P2 - Sprint 2)
SELECT id, name, email, password, created_at
FROM users
WHERE email = $1 AND password = $2;

-- name: CheckUserDuplicate :one
SELECT COUNT(*)
FROM users
WHERE email = $1;

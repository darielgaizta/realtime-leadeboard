-- name: CreateGame :one
INSERT INTO games (name)
VALUES ($1)
RETURNING *;

-- name: GetGames :many
SELECT * FROM games;

-- name: GetGameByID :one
SELECT * FROM games WHERE id = $1;
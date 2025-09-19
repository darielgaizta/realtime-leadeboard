-- name: CreateUserScore :one
INSERT INTO user_scores (user_id, game_id, value)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetUserScoreByGame :one
SELECT * FROM user_scores WHERE user_id = $1 AND game_id = $2;
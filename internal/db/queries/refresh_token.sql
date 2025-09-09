-- name: CreateRefreshToken :one
INSERT INTO refresh_tokens (token_id, user_id, token_hash, expires_at, device_info, ip_address)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: GetRefreshTokenByTokenID :one
SELECT * FROM refresh_tokens WHERE token_id = $1 AND expires_at > NOW() AND is_revoked = FALSE;

-- name: GetRefreshTokensByUser :many
SELECT * FROM refresh_tokens WHERE user_id = $1 AND expires_at > NOW() AND is_revoked = FALSE;

-- name: GetUserByTokenID :one
SELECT u.* FROM users u
JOIN refresh_tokens r ON u.id = r.user_id
WHERE r.token_id = $1 AND r.expires_at > NOW() AND r.is_revoked = FALSE;

-- name: RevokeRefreshTokenByTokenID :exec
UPDATE refresh_tokens
SET is_revoked = TRUE,
    updated_at = NOW()
WHERE token_id = $1;

-- name: RevokeRefreshTokensByUser :exec
UPDATE refresh_tokens
SET is_revoked = TRUE,
    updated_at = NOW()
WHERE user_id = $1;

-- name: DeleteExpiredRefreshTokens :exec
DELETE FROM refresh_tokens WHERE expires_at < NOW();
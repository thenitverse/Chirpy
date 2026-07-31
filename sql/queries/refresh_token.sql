-- name: CreateRefreshToken :one
INSERT INTO refresh_tokens (token,created_at,updated_at,expires_at,
user_id,revoked_at)
VALUES($1,
$2,
$3,
$4,
$5,
$6
)
RETURNING *;

-- name: GetUserFromRefreshToken :one
SELECT users.id FROM refresh_tokens JOIN users ON refresh_tokens.user_id = users.id WHERE
 token = $1 AND refresh_tokens.revoked_at IS NULL
  AND refresh_tokens.expires_at > CURRENT_TIMESTAMP;

-- name: RevokeRefreshToken :exec
UPDATE refresh_tokens
SET 
revoked_at = CURRENT_TIMESTAMP,
updated_at = CURRENT_TIMESTAMP
WHERE token = $1;

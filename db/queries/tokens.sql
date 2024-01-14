-- name: CreateToken :one
INSERT INTO tokens (
    token, user_id, expires_at
) VALUES (
             $1, $2, $3
         ) RETURNING *;

-- name: GetToken :one
SELECT * FROM tokens
WHERE token = $1 LIMIT 1;

-- name: DeleteToken :exec
DELETE FROM tokens
WHERE id = $1;

-- name: DeleteUsersTokens :exec
DELETE FROM tokens
WHERE user_id = $1;
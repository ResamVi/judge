-- name: GetUser :one
SELECT * FROM users
WHERE username = $1 LIMIT 1;

-- name: CreateUser :one
INSERT INTO users (
  username, password
) VALUES (
  $1, $2
)
RETURNING *;

-- name: UpsertUser :exec
INSERT INTO users (
  username, password, approved
) VALUES (
  $1, $2, $3
)
ON CONFLICT DO NOTHING;


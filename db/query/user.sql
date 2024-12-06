-- name: CreateUser :one
INSERT INTO users (
  username,
  password
) VALUES (
  $1, $2
) RETURNING *;

-- name: GetAccount :one
SELECT * FROM users
WHERE username = $1 LIMIT 1;
-- name: AddUser :one
INSERT INTO users (
  id, email, name
) VALUES (
  $1, $2, $3
)
RETURNING *;

-- name: GetUserByEmail :one
SELECT * FROM users WHERE email=$1;

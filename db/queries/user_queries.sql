-- name: CreateUser :one
INSERT INTO users (
    username, password_hash
) VALUES (
    $1, $2
)
RETURNING id, username;

-- name: GetUser :one
SELECT id, username, password_hash FROM users
WHERE username = $1
LIMIT 1;


-- name: UpdateUsername :one
UPDATE users
SET username = $2
WHERE id = $1
RETURNING id, username;

-- name: UpdatePassword :one
UPDATE users
SET password_hash = $2
WHERE id = $1
RETURNING id, username;

-- name: GetById :one
SELECT id, username, password_hash FROM users
WHERE id = $1
LIMIT 1;

-- name: DeleteUser :one
DELETE FROM users
WHERE id = $1
RETURNING id;



-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, email, hashed_password, is_chirpy_red)
VALUES (
    gen_random_uuid(),
    NOW(),
    NOW(),
    $1, -- email parameter provided by the application
    COALESCE($2, 'unset'), -- password_hash parameter provided by the application
    FALSE
)
RETURNING *;

-- name: GetUsers :many
SELECT * FROM users;


-- name: GetUserByEmail :one
SELECT * FROM users WHERE email = $1 LIMIT 1;
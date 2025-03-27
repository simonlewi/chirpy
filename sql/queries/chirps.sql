-- name: CreateChirp :one
INSERT INTO chirps (id, created_at, updated_at, body, user_id)
VALUES (
    gen_random_uuid(),
    NOW(),
    NOW(),
    $1, -- body parameter provided by the application
    $2 -- user_id parameter provided by the application
)
RETURNING *;

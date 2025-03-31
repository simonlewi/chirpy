
-- name: UpgradeUserToChirpyRed :one
UPDATE users
SET
    is_chirpy_red = true,
    updated_at = NOW()
WHERE id = $1
RETURNING is_chirpy_red, updated_at, id;
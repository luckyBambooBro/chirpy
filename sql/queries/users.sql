-- name: CreateUser :one
INSERT INTO users(id, created_at, updated_at, email, hashed_password)
VALUES(
    gen_random_uuid(),
    NOW(),
    NOW(),
    $1,
    $2
) RETURNING *;

-- name: DeleteAllUsers :exec
DELETE FROM users;

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = $1;

-- name: UpdateEmailAndPassword :one
UPDATE users
SET 
hashed_password = $1, 
email = $2,
updated_at = NOW()
WHERE id = $3
RETURNING id, created_at, updated_at, email, hashed_password;

-- name: UpdateChirpyRed :one
UPDATE users
SET is_chirpy_red = true,
updated_at = NOW()
WHERE id = $1
RETURNING id, created_at, updated_at, email, is_chirpy_red;
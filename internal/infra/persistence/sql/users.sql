-- name: CreateUser :exec
INSERT INTO users (
    id,
    username,
    email,
    password_hash,
    created_at,
    updated_at
) VALUES (
    @id,
    @username,
    @email,
    @password_hash,
    @created_at,
    @updated_at
);

-- name: GetUserByID :one
SELECT
    id,
    username,
    email,
    password_hash,
    created_at,
    updated_at
FROM users
WHERE id = @id;

-- name: GetUserByEmail :one
SELECT
    id,
    username,
    email,
    password_hash,
    created_at,
    updated_at
FROM users
WHERE email = @email;

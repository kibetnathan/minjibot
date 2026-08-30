-- name: CreateUser :one
INSERT INTO users (user_id, email, passwordhash)
VALUES (@user_id, @email, @passwordhash)
RETURNING *;

-- name: GetUser :one
SELECT * FROM users WHERE id = @id;

-- name: GetUserByDiscordID :one
SELECT * FROM users WHERE user_id = @user_id;

-- name: GetUserByEmail :one
SELECT * FROM users WHERE email = @email;

-- name: ListUsers :many
SELECT * FROM users ORDER BY created_at DESC;

-- name: UpdateUser :one
UPDATE users
SET email = COALESCE(@email, email),
    passwordhash = COALESCE(@passwordhash, passwordhash),
    is_active = COALESCE(@is_active, is_active),
    updated_at = CURRENT_TIMESTAMP
WHERE id = @id
RETURNING *;

-- name: SetUserPassword :one
UPDATE users
SET passwordhash = @passwordhash,
    updated_at = CURRENT_TIMESTAMP
WHERE id = @id
RETURNING *;

-- name: DeactivateUser :exec
UPDATE users
SET is_active = FALSE,
    updated_at = CURRENT_TIMESTAMP
WHERE id = @id;

-- name: ReactivateUser :exec
UPDATE users
SET is_active = TRUE,
    updated_at = CURRENT_TIMESTAMP
WHERE id = @id;

-- name: DeleteUser :exec
DELETE FROM users WHERE id = @id;

-- name: CountUsers :one
SELECT COUNT(*) FROM users;

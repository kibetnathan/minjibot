-- name: CreateGuild :one
INSERT INTO guilds (id, name, premium_tier)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetGuild :one
SELECT * FROM guilds WHERE id = $1;

-- name: ListGuilds :many
SELECT * FROM guilds ORDER BY created_at DESC;

-- name: UpdateGuild :one
UPDATE guilds
SET name = @name,
    premium_tier = @premium_tier
WHERE id = @id
RETURNING *;

-- name: DeleteGuild :exec
DELETE FROM guilds WHERE id = $1;

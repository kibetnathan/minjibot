-- name: GetGuildBirthdaySettings :one
SELECT * FROM guild_birthday_settings WHERE guild_id = $1;

-- name: UpsertGuildBirthdayChannel :one
INSERT INTO guild_birthday_settings (guild_id, channel_id)
VALUES (@guild_id, @channel_id)
ON CONFLICT (guild_id) DO UPDATE
SET channel_id = EXCLUDED.channel_id,
    updated_at = NOW()
RETURNING *;

-- name: UpsertGuildBirthdayRole :one
INSERT INTO guild_birthday_settings (guild_id, role_id)
VALUES (@guild_id, @role_id)
ON CONFLICT (guild_id) DO UPDATE
SET role_id = EXCLUDED.role_id,
    updated_at = NOW()
RETURNING *;

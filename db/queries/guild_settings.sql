-- name: GetGuildSettings :one
SELECT * FROM guild_settings WHERE guild_id = $1;

-- name: UpsertGuildSettings :one
INSERT INTO guild_settings (guild_id, prefix, language, auto_moderation_enabled, logging_channel_id)
VALUES (@guild_id, @prefix, @language, @auto_moderation_enabled, @logging_channel_id)
ON CONFLICT (guild_id) DO UPDATE
SET prefix = EXCLUDED.prefix,
    language = EXCLUDED.language,
    auto_moderation_enabled = EXCLUDED.auto_moderation_enabled,
    logging_channel_id = EXCLUDED.logging_channel_id,
    updated_at = NOW()
RETURNING *;

-- name: UpdateGuildSettings :one
UPDATE guild_settings
SET prefix = @prefix,
    language = @language,
    auto_moderation_enabled = @auto_moderation_enabled,
    logging_channel_id = @logging_channel_id,
    updated_at = NOW()
WHERE guild_id = @guild_id
RETURNING *;

-- name: DeleteGuildSettings :exec
DELETE FROM guild_settings WHERE guild_id = $1;

-- +goose Up
CREATE TABLE guild_settings (
    guild_id VARCHAR(20) PRIMARY KEY REFERENCES guilds(id) ON DELETE CASCADE,
    prefix VARCHAR(10) NOT NULL DEFAULT '!',
    language VARCHAR(5) NOT NULL DEFAULT 'en',
    auto_moderation_enabled BOOLEAN NOT NULL DEFAULT TRUE,
    logging_channel_id VARCHAR(20),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- +goose Down
DROP TABLE IF EXISTS guild_settings;

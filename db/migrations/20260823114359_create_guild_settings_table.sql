-- +goose Up
CREATE TABLE user_permissions (
    id BIGSERIAL PRIMARY KEY,
    user_id VARCHAR(20) NOT NULL,
    guild_id VARCHAR(20) NOT NULL REFERENCES guilds(id) ON DELETE CASCADE,
    role VARCHAR(50) NOT NULL,
    permissions_json JSONB NOT NULL DEFAULT '{}',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT unique_user_guild_role UNIQUE (user_id, guild_id, role)
);
-- +goose Down
DROP TABLE IF EXISTS guild_settings;

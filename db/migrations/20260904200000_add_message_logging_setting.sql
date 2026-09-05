-- +goose Up
-- Per-guild opt-in for logging every message's content. Defaults to FALSE so
-- message content is not stored unless a guild explicitly enables it.
ALTER TABLE guild_settings
    ADD COLUMN message_logging_enabled BOOLEAN NOT NULL DEFAULT FALSE;

-- +goose Down
ALTER TABLE guild_settings
    DROP COLUMN message_logging_enabled;

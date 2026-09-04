-- +goose Up
CREATE TABLE deleted_messages (
    id              BIGSERIAL    PRIMARY KEY,
    guild_id        VARCHAR(20)  NOT NULL REFERENCES guilds(id) ON DELETE CASCADE,
    channel_id      VARCHAR(20)  NOT NULL,
    message_id      VARCHAR(20)  NOT NULL,
    author_id       VARCHAR(20)  NOT NULL DEFAULT 'unknown',
    author_name     TEXT         NOT NULL DEFAULT '',
    content         TEXT         NOT NULL DEFAULT '',
    attachments     JSONB,
    deleted_by      VARCHAR(20)  NOT NULL DEFAULT 'unknown',
    deleted_by_name TEXT         NOT NULL DEFAULT '',
    created_at      TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_deleted_messages_guild ON deleted_messages(guild_id, created_at DESC);

-- +goose Down
DROP TABLE IF EXISTS deleted_messages;

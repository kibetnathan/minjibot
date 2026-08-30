-- +goose Up
-- Private per-user journal entries. Each row is an entry owned by a single
-- user and shown only to them (`diary add`, `diary view`, `diary delete`).
CREATE TABLE diary_entries (
    id         BIGSERIAL PRIMARY KEY,
    user_id    VARCHAR(20) NOT NULL,
    content    TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Used to show a user's entries newest-first.
CREATE INDEX idx_diary_entries_user ON diary_entries (user_id, created_at DESC);

-- +goose Down
DROP TABLE IF EXISTS diary_entries;

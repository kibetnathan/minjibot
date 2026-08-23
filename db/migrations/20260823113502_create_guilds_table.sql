-- +goose Up
CREATE TABLE guilds (
    id VARCHAR(20) PRIMARY KEY,
    name TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    premium_tier INTEGER NOT NULL DEFAULT 0
);

-- +goose Down
DROP TABLE IF EXISTS guilds;

-- +goose Up
ALTER TABLE audit_logs
    ADD COLUMN actor_name TEXT NOT NULL DEFAULT '',
    ADD COLUMN target_name TEXT NOT NULL DEFAULT '';

-- +goose Down
ALTER TABLE audit_logs
    DROP COLUMN IF EXISTS actor_name,
    DROP COLUMN IF EXISTS target_name;

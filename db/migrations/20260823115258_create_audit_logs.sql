-- +goose Up
CREATE TABLE audit_logs (
    id BIGSERIAL PRIMARY KEY,
    guild_id VARCHAR(20) NOT NULL REFERENCES guilds(id) ON DELETE CASCADE,
    action VARCHAR(50) NOT NULL,
    actor_id VARCHAR(20) NOT NULL,
    target_id VARCHAR(20),
    metadata JSONB,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_audit_logs_guild ON audit_logs(guild_id, created_at DESC);

-- +goose Down
DROP TABLE IF EXISTS audit_logs;

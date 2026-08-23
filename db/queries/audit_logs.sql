-- name: CreateAuditLog :one
INSERT INTO audit_logs (guild_id, action, actor_id, target_id, metadata)
VALUES (@guild_id, @action, @actor_id, @target_id, @metadata)
RETURNING *;

-- name: ListAuditLogsForGuild :many
SELECT * FROM audit_logs
WHERE guild_id = @guild_id
ORDER BY created_at DESC
LIMIT @page_size OFFSET @page_offset;

-- name: ListAuditLogsByActor :many
SELECT * FROM audit_logs
WHERE guild_id = @guild_id AND actor_id = @actor_id
ORDER BY created_at DESC
LIMIT @page_size OFFSET @page_offset;

-- name: CountAuditLogsForGuild :one
SELECT COUNT(*) FROM audit_logs WHERE guild_id = @guild_id;

-- name: DeleteAuditLogsBefore :exec
DELETE FROM audit_logs WHERE created_at < @cutoff;

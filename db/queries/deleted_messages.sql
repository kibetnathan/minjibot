-- name: InsertDeletedMessage :one
INSERT INTO deleted_messages (guild_id, channel_id, message_id, author_id, author_name, content, attachments, deleted_by, deleted_by_name)
VALUES (@guild_id, @channel_id, @message_id, @author_id, @author_name, @content, @attachments, @deleted_by, @deleted_by_name)
RETURNING *;

-- name: ListDeletedMessagesForGuild :many
SELECT * FROM deleted_messages
WHERE guild_id = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: ListDeletedMessagesForChannel :many
SELECT * FROM deleted_messages
WHERE guild_id = $1 AND channel_id = $2
ORDER BY created_at DESC
LIMIT $3 OFFSET $4;

-- name: CountDeletedMessagesForGuild :one
SELECT COUNT(*) FROM deleted_messages
WHERE guild_id = $1;

-- name: DeleteDeletedMessagesBefore :exec
DELETE FROM deleted_messages WHERE created_at < $1;

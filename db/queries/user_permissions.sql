-- name: UpsertUserPermission :one
INSERT INTO user_permissions (user_id, guild_id, role, permissions_json)
VALUES (@user_id, @guild_id, @role, @permissions_json)
ON CONFLICT (user_id, guild_id, role) DO UPDATE
SET permissions_json = EXCLUDED.permissions_json
RETURNING *;

-- name: GetUserPermission :one
SELECT * FROM user_permissions
WHERE user_id = @user_id AND guild_id = @guild_id AND role = @role;

-- name: ListUserPermissionsForUser :many
SELECT * FROM user_permissions
WHERE user_id = @user_id AND guild_id = @guild_id
ORDER BY role;

-- name: ListUserPermissionsForGuild :many
SELECT * FROM user_permissions
WHERE guild_id = @guild_id
ORDER BY user_id, role;

-- name: DeleteUserPermission :exec
DELETE FROM user_permissions
WHERE user_id = @user_id AND guild_id = @guild_id AND role = @role;

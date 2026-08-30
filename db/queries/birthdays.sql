-- name: UpsertBirthday :one
INSERT INTO birthdays (guild_id, user_id, birthday)
VALUES (@guild_id, @user_id, @birthday)
ON CONFLICT (guild_id, user_id) DO UPDATE
SET birthday = EXCLUDED.birthday,
    updated_at = NOW()
RETURNING *;

-- name: GetBirthday :one
SELECT * FROM birthdays WHERE guild_id = $1 AND user_id = $2;

-- name: ListBirthdaysByGuild :many
SELECT * FROM birthdays
WHERE guild_id = $1
ORDER BY birthday ASC;

-- name: ListBirthdaysTodayByGuild :many
SELECT * FROM birthdays
WHERE guild_id = $1
  AND EXTRACT(MONTH FROM birthday) = @month::int
  AND EXTRACT(DAY FROM birthday) = @day::int;

-- name: DeleteBirthday :exec
DELETE FROM birthdays WHERE guild_id = $1 AND user_id = $2;

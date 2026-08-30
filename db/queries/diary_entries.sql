-- name: CreateDiaryEntry :one
INSERT INTO diary_entries (user_id, content)
VALUES (@user_id, @content)
RETURNING *;

-- name: ListDiaryEntriesByUser :many
SELECT * FROM diary_entries
WHERE user_id = $1
ORDER BY created_at DESC;

-- name: GetDiaryEntry :one
SELECT * FROM diary_entries WHERE id = $1 AND user_id = $2;

-- name: DeleteDiaryEntry :exec
DELETE FROM diary_entries WHERE id = $1 AND user_id = $2;

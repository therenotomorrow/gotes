-- name: SelectNotes :many
SELECT id, title, content, created_at, updated_at
FROM notes
ORDER BY id;

-- name: SelectNote :one
SELECT id, title, content, created_at, updated_at
FROM notes
WHERE id = @id;

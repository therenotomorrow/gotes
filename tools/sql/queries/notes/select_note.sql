-- name: SelectNote :one
SELECT *
FROM notes
WHERE id = @id;

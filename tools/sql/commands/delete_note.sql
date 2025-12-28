-- name: DeleteNote :exec
DELETE
FROM notes
WHERE id = @id;

-- name: InsertNote :one
INSERT INTO notes (title, content, created_at, updated_at)
VALUES (@title, @content, @created_at, @updated_at)
RETURNING id;

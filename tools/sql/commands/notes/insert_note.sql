-- name: InsertNote :one
INSERT INTO notes (title, content, user_id, created_at, updated_at)
VALUES (@title, @content, @user_id, @created_at, @updated_at)
RETURNING id;

-- name: InsertUser :one
INSERT INTO users(name, email, password, token, created_at, updated_at)
VALUES (@name, @email, @password, @token, @created_at, @updated_at)
RETURNING id;

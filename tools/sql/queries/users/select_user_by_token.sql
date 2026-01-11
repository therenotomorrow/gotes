-- name: SelectUserByToken :one
SELECT *
FROM users
WHERE token = @token;

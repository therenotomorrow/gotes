-- name: SelectUserByEmail :one
SELECT *
FROM users
WHERE email = @email;

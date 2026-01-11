-- name: UpdateUser :exec
UPDATE users
SET name       = @name,
    email      = @email,
    password   = @password,
    token      = @token,
    updated_at = @updated_at
WHERE id = @id;

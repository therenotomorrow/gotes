-- name: SelectNotesByUser :many
SELECT *
FROM notes
WHERE user_id = @user_id
ORDER BY notes.id;

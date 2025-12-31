package domain

type Error string

func (e Error) Error() string {
	return string(e)
}

const (
	ErrNoteNotFound     Error = "note not found"
	ErrUserNotFound     Error = "user not found"
	ErrPermissionDenied Error = "permission denied"
)

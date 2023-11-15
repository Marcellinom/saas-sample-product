package errors

type ForbiddenError struct {
	msg     string
	details string
}

func NewForbiddenError(msg string, details string) ForbiddenError {
	return ForbiddenError{msg, details}
}

func (e ForbiddenError) Error() string {
	return e.msg
}

func (e ForbiddenError) Details() string {
	return e.details
}

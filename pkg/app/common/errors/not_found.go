package errors

const NotFoundErrorDefaultMessage = "not_found"

type NotFoundError struct {
	code int
	msg  string
}

func NewNotFoundErrorWithCode(code int, msg string) NotFoundError {
	if msg == "" {
		msg = NotFoundErrorDefaultMessage
	}
	return NotFoundError{code, msg}
}

func NewNotFoundError(msg string) NotFoundError {
	return NewNotFoundErrorWithCode(404, msg)
}

func (e NotFoundError) Code() int {
	return e.code
}

func (e NotFoundError) Error() string {
	return e.msg
}

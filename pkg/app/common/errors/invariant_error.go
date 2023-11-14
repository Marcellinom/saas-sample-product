package errors

type InvariantError struct {
	code int
	msg  string
}

func NewInvariantErrorWithCode(code int, msg string) InvariantError {
	return InvariantError{code, msg}
}

// DEPRECATED: Jangan dipakai lagi, pakai NewInvariantErrorWithCode() saja
func NewInvariantError(msg string) InvariantError {
	return InvariantError{code: 400, msg: msg}
}

func (e InvariantError) Code() int {
	return e.code
}

func (e InvariantError) Error() string {
	return e.msg
}

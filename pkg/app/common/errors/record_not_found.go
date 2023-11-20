package errors

const RecordNotFoundErrorDefaultMessage = "record_not_found"

type RecordNotFoundError struct {
	message string
}

func NewRecordNotFoundError(msg string) RecordNotFoundError {
	if msg == "" {
		msg = RecordNotFoundErrorDefaultMessage
	}
	return RecordNotFoundError{msg}
}

func (e RecordNotFoundError) Error() string {
	return e.message
}

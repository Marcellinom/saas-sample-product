package errors

const (
	AggregateVersionMismatchErrorMessage = "aggregate_version_mismatch"
)

type AggregateVersionMismatchError struct {
	code int
	msg  string
}

func NewAggregateVersionMismatchErrorrWithCode(code int, msg string) AggregateVersionMismatchError {
	return AggregateVersionMismatchError{code, msg}
}

func NewAggregateVersionMismatchError() AggregateVersionMismatchError {
	return AggregateVersionMismatchError{409, AggregateVersionMismatchErrorMessage}
}

func (e AggregateVersionMismatchError) Code() int {
	return e.code
}

func (e AggregateVersionMismatchError) Error() string {
	return e.msg
}

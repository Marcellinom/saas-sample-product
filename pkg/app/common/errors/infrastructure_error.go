package errors

type InfrastructureError struct {
	msg    string
	detail string
}

func NewInfrastructureError(msg string, detail string) *InfrastructureError {
	return &InfrastructureError{msg: msg, detail: detail}
}

func (e *InfrastructureError) Error() string {
	return e.msg
}

func (e *InfrastructureError) Detail() string {
	return e.detail
}

package errors

type BadRequestError struct {
	code    int
	message string
	data    map[string]interface{}
}

func NewBadRequestError(code int, message string, data map[string]interface{}) BadRequestError {
	return BadRequestError{
		code:    code,
		message: message,
		data:    data,
	}
}

func (e BadRequestError) Code() int {
	return e.code
}

func (e BadRequestError) Message() string {
	return e.message
}

func (e BadRequestError) Data() map[string]interface{} {
	return e.data
}

func (e BadRequestError) Error() string {
	return e.message
}

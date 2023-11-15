package web

const (
	validationError     string = "validation_error"
	internalServerError string = "internal_server_error"
	forbiddenError      string = "forbidden"
)

var statusCode = map[string]int{
	validationError:     9001,
	internalServerError: 9002,
	forbiddenError:      9003,
}

package apperrors

import (
	"fmt"
	"strings"
)

type BusinessError struct {
	Message   string `json:"message"`
	ErrorCode string `json:"error_code"`
}

type NotFoundError struct {
	BusinessError
}

func Business(message, code string) BusinessError {
	return BusinessError{Message: message, ErrorCode: code}
}

func (e BusinessError) Msgf(params ...interface{}) BusinessError {
	e.Message = fmt.Sprintf(e.Message, params...)

	return e
}

// NotFound error.
func NotFound(resource, code string) NotFoundError {
	return NotFoundError{
		BusinessError: Business("Not found: "+resource, code),
	}
}

// Error string.
func (e BusinessError) Error() string {
	var builder strings.Builder

	if len(e.ErrorCode) > 0 {
		builder.WriteString(e.ErrorCode)
		builder.WriteString(": ")
	}

	builder.WriteString(e.Message)

	return builder.String()
}

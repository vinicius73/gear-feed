package apperrors

type ValidationError struct {
	BusinessError
	Details []error `json:"details,omitempty"`
}

// Validation error build.
func Validation(message, code string, details ...error) ValidationError {
	return ValidationError{
		BusinessError: Business(message, code),
		Details:       details,
	}
}

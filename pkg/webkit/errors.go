package webkit

import "fmt"

// Error represents an HTTP error.
type Error struct {
	Code    int
	Message string
}

// NewError creates a new Error type by HTTP code and a
// user defined message.
func NewError(code int, message string) *Error {
	return &Error{
		Code:    code,
		Message: message,
	}
}

// Error implements the Error type and returns the
// formatted error message.
func (e *Error) Error() string {
	return fmt.Sprintf("code=%d, message=%v", e.Code, e.Message)
}

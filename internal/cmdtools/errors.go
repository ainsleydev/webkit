package cmdtools

// ExitError is an error that causes the program to exit
// with a non-zero code without printing an error message
type ExitError struct {
	Code int
}

func (e ExitError) Error() string {
	return "" // Empty string means no error output
}

// ExitWithCode returns a silent exit error
func ExitWithCode(code int) *ExitError {
	return &ExitError{Code: code}
}

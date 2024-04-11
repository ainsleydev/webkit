package httputil

// Is2xx returns true if the status code is in the 2xx range.
func Is2xx(statusCode int) bool {
	return statusCode >= 200 && statusCode < 300
}

// Is3xx returns true if the status code is in the 3xx range.
func Is3xx(statusCode int) bool {
	return statusCode >= 300 && statusCode < 400
}

// Is4xx returns true if the status code is in the 4xx range.
func Is4xx(statusCode int) bool {
	return statusCode >= 400 && statusCode < 500
}

// Is5xx returns true if the status code is in the 5xx range.
func Is5xx(statusCode int) bool {
	return statusCode >= 500 && statusCode < 600
}

// IsError returns true if the status code is in the 4xx or 5xx range.
func IsError(statusCode int) bool {
	return Is4xx(statusCode) || Is5xx(statusCode)
}

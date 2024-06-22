package stringutil

// IsNotEmpty checks if a string pointer is empty.
func IsNotEmpty(in *string) bool {
	return in != nil && *in != ""
}

// IsEmpty checks if a string pointer is empty.
func IsEmpty(in *string) bool {
	return in == nil || *in == ""
}

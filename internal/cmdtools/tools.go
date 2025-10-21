// Package cmdtools provides utilities for command-line interface operations,
// including error handling and input/output management.
package cmdtools

import (
	"fmt"
	"os"
)

// Exit prints an error message if one exists and exits the program.
func Exit(err error) {
	if err != nil {
		fmt.Println(err.Error()) //nolint
	}
	os.Exit(0)
}

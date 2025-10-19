package enforce

import (
	"fmt"
	"os"
	"reflect"

	"github.com/ainsleydev/webkit/internal/printer"
)

var console = printer.New(os.Stderr)
var exit = os.Exit

// Equal asserts that two values are equal. If they are not equal,
// it prints an error message and exits the program with status 1.
//
// Example:
//
//	enforce.Equal(got, want, "database connection failed")
func Equal(got, want any, msg string) {
	if !reflect.DeepEqual(got, want) {
		console.Error(fmt.Sprintf("Enforcement failed: %s", msg))
		console.Printf("   Got:  %+v", got)
		console.Printf("   Want: %+v", want)
		exit(1)
	}
}

// NotEqual asserts that two values are not equal. If they are equal,
// it prints an error message and exits the program with status 1.
//
// Example:
//
//	enforce.NotEqual(password, "", "password cannot be empty")
func NotEqual(got, notWant any, msg string) {
	if reflect.DeepEqual(got, notWant) {
		console.Error(fmt.Sprintf("Enforcement failed: %s", msg))
		console.Printf("   Got:  %+v (must not equal %+v)", got, notWant)
		exit(1)
	}
}

// True asserts that a condition is true. If it is false,
// it prints an error message and exits the program with status 1.
//
// Example:
//
//	enforce.True(len(items) > 0, "items list cannot be empty")
func True(condition bool, msg string) {
	if !condition {
		console.Error(fmt.Sprintf("Enforcement failed: %s", msg))
		exit(1)
	}
}

// NoError asserts that an error is nil. If the error is not nil,
// it prints the error and exits the program with status 1.
//
// Example:
//
//	enforce.NoError(db.Connect(), "failed to connect to database")
func NoError(err error, msg string) {
	if err != nil {
		console.Error(fmt.Sprintf("Enforcement failed: %s", msg))
		console.Printf("   Error: %v", err)
		exit(1)
	}
}

// NotNil asserts that a value is not nil. If it is nil,
// it prints an error message and exits the program with status 1.
//
// Example:
//
//	enforce.NotNil(db, "database connection cannot be nil")
func NotNil(value any, msg string) {
	if value == nil || reflect.ValueOf(value).IsNil() {
		console.Error(fmt.Sprintf("Enforcement failed: %s", msg))
		exit(1)
	}
}

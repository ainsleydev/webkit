package jsonformat

import (
	"reflect"

	"github.com/ainsleydev/webkit/internal/appdef"
)

func init() {
	// Register the Definition type to scan for inline tags.
	RegisterType(reflect.TypeOf(appdef.Definition{}))
}

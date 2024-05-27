package forms

import (
	"bytes"
	"context"
	"fmt"
	"testing"

	"github.com/ainsleydev/webkit/pkg/markup/base"
)

func TestTemp(t *testing.T) {
	buf := &bytes.Buffer{}
	Render(InputFieldProps{
		Name:  "Test",
		Value: "Fuck",
		ElementProps: base.ElementProps{
			Classes: []string{"Hey"},
		},
	}).Render(context.TODO(), buf)
	fmt.Println(buf.String())
}

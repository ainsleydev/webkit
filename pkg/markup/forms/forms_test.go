package forms

import (
	"bytes"
	"context"
	"fmt"
	"testing"

	"github.com/a-h/templ"
	"github.com/stretchr/testify/require"

	"github.com/ainsleydev/webkit/pkg/util/ptr"
)

func TestForms(t *testing.T) {
	buf := &bytes.Buffer{}

	props := FormGroupProps{
		Label:       "Label",
		Description: "Description",
		ID:          "ID",
		Error:       "Error",
		Width:       ptr.IntPtr(12),
	}

	ctx := templ.WithChildren(context.Background(), InputField(InputFieldProps{
		ID: "ID",
	}))

	err := FormGroup(props).Render(ctx, buf)
	require.NoError(t, err)

	fmt.Println(buf.String())
}

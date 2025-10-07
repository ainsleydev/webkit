package appdef

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAppType_String(t *testing.T) {
	t.Parallel()

	got := AppTypeGoLang.String()
	assert.Equal(t, "golang", got)
	assert.IsType(t, "", got)
}

func TestApp_Language(t *testing.T) {
	t.Parallel()

	tt := []struct {
		input AppType
		want  string
	}{
		{input: AppTypeGoLang, want: "go"},
		{input: AppTypePayload, want: "js"},
		{input: AppTypeSvelteKit, want: "js"},
	}

	for _, test := range tt {
		t.Run(test.input.String(), func(t *testing.T) {
			t.Parallel()
			a := App{Type: test.input}
			got := a.Language()
			assert.Equal(t, test.want, got)
		})
	}
}

func TestDomainType_String(t *testing.T) {
	t.Parallel()

	got := DomainTypePrimary.String()
	assert.Equal(t, "primary", got)
	assert.IsType(t, "", got)
}

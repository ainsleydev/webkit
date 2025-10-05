package appdef

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGithubLabels(t *testing.T) {
	t.Parallel()

	tt := map[string]struct {
		input Definition
		want  []string
	}{
		"No Apps": {
			input: Definition{Apps: nil},
			want:  []string{"webkit"},
		},
		"Single App": {
			input: Definition{
				Apps: []App{
					{Type: AppTypeSvelteKit},
				},
			},
			want: []string{"webkit", AppTypeSvelteKit.String()},
		},
		"Multiple Apps": {
			input: Definition{
				Apps: []App{
					{Type: AppTypeGoLang},
					{Type: AppTypePayload},
				},
			},
			want: []string{
				"webkit",
				AppTypeGoLang.String(),
				AppTypePayload.String(),
			},
		},
	}

	for name, test := range tt {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			got := test.input.GithubLabels()
			assert.Equal(t, test.want, got)
		})
	}
}

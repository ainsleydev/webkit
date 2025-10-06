package appdef

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestResourceTypeString(t *testing.T) {
	t.Parallel()

	got := ResourceTypePostgres.String()
	assert.Equal(t, "postgres", got)
	assert.IsType(t, "", got)

}

func TestResourceProviderString(t *testing.T) {
	t.Parallel()

	got := ResourceProviderDigitalOcean.String()
	assert.Equal(t, "digital-ocean", got)
	assert.IsType(t, "", got)
}

func TestResourceApplyDefaults(t *testing.T) {
	t.Parallel()

	tt := map[string]struct {
		input Resource
		want  Resource
	}{
		"Nil Config And Outputs": {
			input: Resource{},
			want: Resource{
				Config:  make(map[string]any),
				Outputs: []string{},
				Backup: ResourceBackupConfig{
					Enabled: true,
				},
			},
		},
		"Existing Config And Outputs": {
			input: Resource{
				Config:  map[string]any{"size": "small"},
				Outputs: []string{"url"},
			},
			want: Resource{
				Config:  map[string]any{"size": "small"},
				Outputs: []string{"url"},
				Backup: ResourceBackupConfig{
					Enabled: true,
				},
			},
		},
	}

	for name, test := range tt {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			err := test.input.applyDefaults()
			assert.NoError(t, err)
			assert.Equal(t, test.want, test.input)
		})
	}
}

package appdef

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMonitorType_String(t *testing.T) {
	t.Parallel()

	tt := map[string]struct {
		input MonitorType
		want  string
	}{
		"HTTP":     {input: MonitorTypeHTTP, want: "http"},
		"Postgres": {input: MonitorTypePostgres, want: "postgres"},
		"Push":     {input: MonitorTypePush, want: "push"},
	}

	for name, test := range tt {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			got := test.input.String()
			assert.Equal(t, test.want, got)
		})
	}
}

func TestApp_GenerateMonitors(t *testing.T) {
	t.Parallel()

	t.Run("Monitoring Disabled", func(t *testing.T) {
		t.Parallel()

		app := &App{
			Name: "web",
			Domains: []Domain{
				{Name: "example.com", Type: DomainTypePrimary},
			},
			Monitoring: Monitoring{Enabled: false},
		}

		monitors := app.GenerateMonitors()
		assert.Empty(t, monitors)
	})

	t.Run("No Domains", func(t *testing.T) {
		t.Parallel()

		app := &App{
			Name:       "web",
			Domains:    []Domain{},
			Monitoring: Monitoring{Enabled: true},
		}

		monitors := app.GenerateMonitors()
		assert.Empty(t, monitors)
	})

	t.Run("Single Primary Domain", func(t *testing.T) {
		t.Parallel()

		app := &App{
			Name: "web",
			Domains: []Domain{
				{Name: "example.com", Type: DomainTypePrimary},
			},
			Monitoring: Monitoring{Enabled: true},
		}

		monitors := app.GenerateMonitors()
		require.Len(t, monitors, 1)

		m := monitors[0]
		assert.Equal(t, "web-example-com", m.Name)
		assert.Equal(t, MonitorTypeHTTP, m.Type)
		assert.Equal(t, "https://example.com", m.URL)
		assert.Equal(t, "GET", m.Method)
	})

	t.Run("Multiple Domains Primary And Alias", func(t *testing.T) {
		t.Parallel()

		app := &App{
			Name: "api",
			Domains: []Domain{
				{Name: "api.example.com", Type: DomainTypePrimary},
				{Name: "www.api.example.com", Type: DomainTypeAlias},
			},
			Infra:      Infra{Config: nil}, // Default health check path.
			Monitoring: Monitoring{Enabled: true},
		}

		monitors := app.GenerateMonitors()
		require.Len(t, monitors, 2)

		assert.Equal(t, "api-api-example-com", monitors[0].Name)
		assert.Equal(t, "https://api.example.com", monitors[0].URL)

		assert.Equal(t, "api-www-api-example-com", monitors[1].Name)
		assert.Equal(t, "https://www.api.example.com", monitors[1].URL)
	})

	t.Run("Unmanaged Domains Skipped", func(t *testing.T) {
		t.Parallel()

		app := &App{
			Name: "web",
			Domains: []Domain{
				{Name: "example.com", Type: DomainTypePrimary},
				{Name: "unmanaged.com", Type: DomainTypeUnmanaged},
				{Name: "www.example.com", Type: DomainTypeAlias},
			},
			Infra:      Infra{},
			Monitoring: Monitoring{Enabled: true},
		}

		monitors := app.GenerateMonitors()
		require.Len(t, monitors, 2) // Only primary and alias.

		assert.Equal(t, "web-example-com", monitors[0].Name)
		assert.Equal(t, "web-www-example-com", monitors[1].Name)
	})
}

func TestSanitiseMonitorName(t *testing.T) {
	t.Parallel()

	tt := map[string]struct {
		domain string
		want   string
	}{
		"Simple Domain":       {domain: "example.com", want: "example-com"},
		"Subdomain":           {domain: "api.example.com", want: "api-example-com"},
		"Deep Subdomain":      {domain: "v1.api.example.com", want: "v1-api-example-com"},
		"Multiple Subdomains": {domain: "auth.api.v2.example.com", want: "auth-api-v2-example-com"},
		"No Dots":             {domain: "localhost", want: "localhost"},
	}

	for name, test := range tt {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := sanitiseMonitorName(test.domain)
			assert.Equal(t, test.want, got)
		})
	}
}

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
		"DNS":      {input: MonitorTypeDNS, want: "dns"},
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
			Name:  "web",
			Title: "Web",
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
			Title:      "Web",
			Domains:    []Domain{},
			Monitoring: Monitoring{Enabled: true},
		}

		monitors := app.GenerateMonitors()
		assert.Empty(t, monitors)
	})

	t.Run("Single Primary Domain", func(t *testing.T) {
		t.Parallel()

		app := &App{
			Name:  "web",
			Title: "Web",
			Domains: []Domain{
				{Name: "example.com", Type: DomainTypePrimary},
			},
			Monitoring: Monitoring{Enabled: true},
		}

		monitors := app.GenerateMonitors()
		require.Len(t, monitors, 2) // HTTP + DNS

		// HTTP monitor.
		assert.Equal(t, "Web - example.com", monitors[0].Name)
		assert.Equal(t, MonitorTypeHTTP, monitors[0].Type)
		assert.Equal(t, "https://example.com", monitors[0].URL)
		assert.Equal(t, "GET", monitors[0].Method)

		// DNS monitor.
		assert.Equal(t, "Web DNS - example.com", monitors[1].Name)
		assert.Equal(t, MonitorTypeDNS, monitors[1].Type)
		assert.Equal(t, "example.com", monitors[1].Domain)
	})

	t.Run("Multiple Domains Primary And Alias", func(t *testing.T) {
		t.Parallel()

		app := &App{
			Name:  "api",
			Title: "API",
			Domains: []Domain{
				{Name: "api.example.com", Type: DomainTypePrimary},
				{Name: "www.api.example.com", Type: DomainTypeAlias},
			},
			Infra:      Infra{Config: nil}, // Default health check path.
			Monitoring: Monitoring{Enabled: true},
		}

		monitors := app.GenerateMonitors()
		require.Len(t, monitors, 4) // 2 domains × 2 types (HTTP + DNS)

		// First domain - HTTP.
		assert.Equal(t, "API - api.example.com", monitors[0].Name)
		assert.Equal(t, MonitorTypeHTTP, monitors[0].Type)
		assert.Equal(t, "https://api.example.com", monitors[0].URL)

		// First domain - DNS.
		assert.Equal(t, "API DNS - api.example.com", monitors[1].Name)
		assert.Equal(t, MonitorTypeDNS, monitors[1].Type)
		assert.Equal(t, "api.example.com", monitors[1].Domain)

		// Second domain - HTTP.
		assert.Equal(t, "API - www.api.example.com", monitors[2].Name)
		assert.Equal(t, MonitorTypeHTTP, monitors[2].Type)
		assert.Equal(t, "https://www.api.example.com", monitors[2].URL)

		// Second domain - DNS.
		assert.Equal(t, "API DNS - www.api.example.com", monitors[3].Name)
		assert.Equal(t, MonitorTypeDNS, monitors[3].Type)
		assert.Equal(t, "www.api.example.com", monitors[3].Domain)
	})

	t.Run("Unmanaged Domains Skipped", func(t *testing.T) {
		t.Parallel()

		app := &App{
			Name:  "web",
			Title: "Web",
			Domains: []Domain{
				{Name: "example.com", Type: DomainTypePrimary},
				{Name: "unmanaged.com", Type: DomainTypeUnmanaged},
				{Name: "www.example.com", Type: DomainTypeAlias},
			},
			Infra:      Infra{},
			Monitoring: Monitoring{Enabled: true},
		}

		monitors := app.GenerateMonitors()
		require.Len(t, monitors, 4) // 2 managed domains × 2 types (HTTP + DNS)

		// First managed domain monitors.
		assert.Equal(t, "Web - example.com", monitors[0].Name)
		assert.Equal(t, MonitorTypeHTTP, monitors[0].Type)
		assert.Equal(t, "Web DNS - example.com", monitors[1].Name)
		assert.Equal(t, MonitorTypeDNS, monitors[1].Type)

		// Second managed domain monitors.
		assert.Equal(t, "Web - www.example.com", monitors[2].Name)
		assert.Equal(t, MonitorTypeHTTP, monitors[2].Type)
		assert.Equal(t, "Web DNS - www.example.com", monitors[3].Name)
		assert.Equal(t, MonitorTypeDNS, monitors[3].Type)
	})
}

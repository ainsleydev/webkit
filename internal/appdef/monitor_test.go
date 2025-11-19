package appdef

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ainsleydev/webkit/pkg/env"
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
			Infra:      Infra{Config: map[string]any{"health_check_path": "/health"}},
			Monitoring: Monitoring{Enabled: true},
		}

		monitors := app.GenerateMonitors()
		require.Len(t, monitors, 1)

		m := monitors[0]
		assert.Equal(t, "web-example-com", m.Name)
		assert.Equal(t, MonitorTypeHTTP, m.Type)
		assert.Equal(t, "https://example.com/health", m.URL)
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
		assert.Equal(t, "https://api.example.com/", monitors[0].URL)

		assert.Equal(t, "api-www-api-example-com", monitors[1].Name)
		assert.Equal(t, "https://www.api.example.com/", monitors[1].URL)
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

func TestApp_healthCheckPath(t *testing.T) {
	t.Parallel()

	tt := map[string]struct {
		config map[string]any
		want   string
	}{
		"Default Slash":       {config: nil, want: "/"},
		"Empty Config":        {config: map[string]any{}, want: "/"},
		"Custom Path":         {config: map[string]any{"health_check_path": "/health"}, want: "/health"},
		"Custom API Path":     {config: map[string]any{"health_check_path": "/api/health"}, want: "/api/health"},
		"Empty String Path":   {config: map[string]any{"health_check_path": ""}, want: "/"},
		"Non String Value":    {config: map[string]any{"health_check_path": 123}, want: "/"},
		"Other Config Fields": {config: map[string]any{"port": 8080}, want: "/"},
	}

	for name, test := range tt {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			app := &App{
				Infra: Infra{Config: test.config},
			}

			got := app.healthCheckPath()
			assert.Equal(t, test.want, got)
		})
	}
}

func TestResource_GenerateMonitors(t *testing.T) {
	t.Parallel()

	mockURLGen := func(r *Resource, enviro env.Environment, output string) string {
		return "${module.resources." + r.Name + "_" + enviro.String() + "_" + output + "}"
	}

	t.Run("Monitoring Disabled", func(t *testing.T) {
		t.Parallel()

		resource := &Resource{
			Name:       "db",
			Type:       ResourceTypePostgres,
			Monitoring: Monitoring{Enabled: false},
		}

		monitors := resource.GenerateMonitors(env.Production, mockURLGen)
		assert.Empty(t, monitors)
	})

	t.Run("Non Postgres Resource", func(t *testing.T) {
		t.Parallel()

		resource := &Resource{
			Name:       "bucket",
			Type:       ResourceTypeS3,
			Monitoring: Monitoring{Enabled: true},
		}

		monitors := resource.GenerateMonitors(env.Production, mockURLGen)
		assert.Empty(t, monitors)
	})

	t.Run("Postgres Resource Production", func(t *testing.T) {
		t.Parallel()

		resource := &Resource{
			Name:       "db",
			Type:       ResourceTypePostgres,
			Monitoring: Monitoring{Enabled: true},
		}

		monitors := resource.GenerateMonitors(env.Production, mockURLGen)
		require.Len(t, monitors, 1)

		m := monitors[0]
		assert.Equal(t, "db-production", m.Name)
		assert.Equal(t, MonitorTypePostgres, m.Type)
		assert.Equal(t, "${module.resources.db_production_connection_url}", m.URL)
		assert.Equal(t, "", m.Method)
	})

	t.Run("Postgres Resource Staging", func(t *testing.T) {
		t.Parallel()

		resource := &Resource{
			Name:       "analytics-db",
			Type:       ResourceTypePostgres,
			Monitoring: Monitoring{Enabled: true},
		}

		monitors := resource.GenerateMonitors(env.Staging, mockURLGen)
		require.Len(t, monitors, 1)

		m := monitors[0]
		assert.Equal(t, "analytics-db-staging", m.Name)
		assert.Equal(t, "${module.resources.analytics-db_staging_connection_url}", m.URL)
	})
}

func TestResource_GenerateHeartbeatMonitor(t *testing.T) {
	t.Parallel()

	t.Run("Backup Disabled", func(t *testing.T) {
		t.Parallel()

		resource := &Resource{
			Name:   "db",
			Backup: ResourceBackupConfig{Enabled: false},
		}

		monitor := resource.GenerateHeartbeatMonitor("0 2 * * *")
		assert.Empty(t, monitor.Name)
		assert.Empty(t, monitor.Type)
	})

	t.Run("Backup Enabled", func(t *testing.T) {
		t.Parallel()

		resource := &Resource{
			Name:   "db",
			Backup: ResourceBackupConfig{Enabled: true},
		}

		monitor := resource.GenerateHeartbeatMonitor("0 2 * * *")

		assert.Equal(t, "backup-db", monitor.Name)
		assert.Equal(t, MonitorTypePush, monitor.Type)
		assert.Empty(t, monitor.URL)
		assert.Empty(t, monitor.Method)
	})

	t.Run("Multiple Resources", func(t *testing.T) {
		t.Parallel()

		r1 := &Resource{Name: "primary-db", Backup: ResourceBackupConfig{Enabled: true}}
		r2 := &Resource{Name: "analytics-db", Backup: ResourceBackupConfig{Enabled: true}}

		m1 := r1.GenerateHeartbeatMonitor("0 2 * * *")
		m2 := r2.GenerateHeartbeatMonitor("0 2 * * *")

		assert.Equal(t, "backup-primary-db", m1.Name)
		assert.Equal(t, "backup-analytics-db", m2.Name)
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

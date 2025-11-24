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

func TestDefinition_GenerateMonitors(t *testing.T) {
	t.Parallel()

	t.Run("Monitoring Disabled", func(t *testing.T) {
		t.Parallel()

		def := &Definition{
			Project: Project{Title: "My Project"},
			Apps: []App{
				{
					Name:  "web",
					Title: "Web",
					Domains: []Domain{
						{Name: "example.com", Type: DomainTypePrimary},
					},
					Monitoring: false,
				},
			},
		}

		monitors := def.GenerateMonitors()
		// Codebase backup monitor is always generated.
		require.Len(t, monitors, 1)
		assert.Equal(t, "Backup - Codebase", monitors[0].Name)
		assert.Equal(t, MonitorTypePush, monitors[0].Type)
		assert.Equal(t, MonitorIntervalBackup, monitors[0].Interval)
	})

	t.Run("No Domains", func(t *testing.T) {
		t.Parallel()

		def := &Definition{
			Project: Project{Title: "My Project"},
			Apps: []App{
				{
					Name:       "web",
					Title:      "Web",
					Domains:    []Domain{},
					Monitoring: true,
				},
			},
		}

		monitors := def.GenerateMonitors()
		// Codebase backup monitor is always generated.
		require.Len(t, monitors, 1)
		assert.Equal(t, "Backup - Codebase", monitors[0].Name)
		assert.Equal(t, MonitorTypePush, monitors[0].Type)
		assert.Equal(t, MonitorIntervalBackup, monitors[0].Interval)
	})

	t.Run("Single Primary Domain", func(t *testing.T) {
		t.Parallel()

		def := &Definition{
			Project: Project{Title: "My Project"},
			Apps: []App{
				{
					Name:  "web",
					Title: "Web",
					Domains: []Domain{
						{Name: "example.com", Type: DomainTypePrimary},
					},
					Monitoring: true,
				},
			},
		}

		monitors := def.GenerateMonitors()
		require.Len(t, monitors, 3) // HTTP + DNS + Codebase Backup

		// HTTP monitor.
		assert.Equal(t, "HTTP - example.com", monitors[0].Name)
		assert.Equal(t, MonitorTypeHTTP, monitors[0].Type)
		assert.Equal(t, "https://example.com", monitors[0].URL)
		assert.Equal(t, "GET", monitors[0].Method)
		assert.Equal(t, MonitorIntervalHTTP, monitors[0].Interval)

		// DNS monitor.
		assert.Equal(t, "DNS - example.com", monitors[1].Name)
		assert.Equal(t, MonitorTypeDNS, monitors[1].Type)
		assert.Equal(t, "example.com", monitors[1].Domain)
		assert.Equal(t, MonitorIntervalDNS, monitors[1].Interval)

		// Codebase backup monitor.
		assert.Equal(t, "Backup - Codebase", monitors[2].Name)
		assert.Equal(t, MonitorTypePush, monitors[2].Type)
		assert.Equal(t, MonitorIntervalBackup, monitors[2].Interval)
	})

	t.Run("Multiple Domains Primary And Alias", func(t *testing.T) {
		t.Parallel()

		def := &Definition{
			Project: Project{Title: "My Project"},
			Apps: []App{
				{
					Name:  "api",
					Title: "API",
					Domains: []Domain{
						{Name: "api.example.com", Type: DomainTypePrimary},
						{Name: "www.api.example.com", Type: DomainTypeAlias},
					},
					Infra:      Infra{Config: nil}, // Default health check path.
					Monitoring: true,
				},
			},
		}

		monitors := def.GenerateMonitors()
		require.Len(t, monitors, 5) // 2 domains × 2 types (HTTP + DNS) + Codebase Backup

		// First domain (primary) - HTTP.
		assert.Equal(t, "HTTP - api.example.com", monitors[0].Name)
		assert.Equal(t, MonitorTypeHTTP, monitors[0].Type)
		assert.Equal(t, "https://api.example.com", monitors[0].URL)
		assert.Equal(t, MonitorIntervalHTTP, monitors[0].Interval)

		// First domain - DNS.
		assert.Equal(t, "DNS - api.example.com", monitors[1].Name)
		assert.Equal(t, MonitorTypeDNS, monitors[1].Type)
		assert.Equal(t, "api.example.com", monitors[1].Domain)
		assert.Equal(t, MonitorIntervalDNS, monitors[1].Interval)

		// Second domain (alias) - HTTP.
		assert.Equal(t, "HTTP - www.api.example.com", monitors[2].Name)
		assert.Equal(t, MonitorTypeHTTP, monitors[2].Type)
		assert.Equal(t, "https://www.api.example.com", monitors[2].URL)
		assert.Equal(t, MonitorIntervalHTTP, monitors[2].Interval)

		// Second domain - DNS.
		assert.Equal(t, "DNS - www.api.example.com", monitors[3].Name)
		assert.Equal(t, MonitorTypeDNS, monitors[3].Type)
		assert.Equal(t, "www.api.example.com", monitors[3].Domain)
		assert.Equal(t, MonitorIntervalDNS, monitors[3].Interval)

		// Codebase backup monitor.
		assert.Equal(t, "Backup - Codebase", monitors[4].Name)
		assert.Equal(t, MonitorTypePush, monitors[4].Type)
		assert.Equal(t, MonitorIntervalBackup, monitors[4].Interval)
	})

	t.Run("Unmanaged Domains Skipped", func(t *testing.T) {
		t.Parallel()

		def := &Definition{
			Project: Project{Title: "My Project"},
			Apps: []App{
				{
					Name:  "web",
					Title: "Web",
					Domains: []Domain{
						{Name: "example.com", Type: DomainTypePrimary},
						{Name: "unmanaged.com", Type: DomainTypeUnmanaged},
						{Name: "www.example.com", Type: DomainTypeAlias},
					},
					Infra:      Infra{},
					Monitoring: true,
				},
			},
		}

		monitors := def.GenerateMonitors()
		require.Len(t, monitors, 5) // 2 managed domains × 2 types (HTTP + DNS) + Codebase Backup

		// First managed domain (primary).
		assert.Equal(t, "HTTP - example.com", monitors[0].Name)
		assert.Equal(t, MonitorTypeHTTP, monitors[0].Type)
		assert.Equal(t, MonitorIntervalHTTP, monitors[0].Interval)
		assert.Equal(t, "DNS - example.com", monitors[1].Name)
		assert.Equal(t, MonitorTypeDNS, monitors[1].Type)
		assert.Equal(t, MonitorIntervalDNS, monitors[1].Interval)

		// Second managed domain (alias).
		assert.Equal(t, "HTTP - www.example.com", monitors[2].Name)
		assert.Equal(t, MonitorTypeHTTP, monitors[2].Type)
		assert.Equal(t, MonitorIntervalHTTP, monitors[2].Interval)
		assert.Equal(t, "DNS - www.example.com", monitors[3].Name)
		assert.Equal(t, MonitorTypeDNS, monitors[3].Type)
		assert.Equal(t, MonitorIntervalDNS, monitors[3].Interval)

		// Codebase backup monitor.
		assert.Equal(t, "Backup - Codebase", monitors[4].Name)
		assert.Equal(t, MonitorTypePush, monitors[4].Type)
		assert.Equal(t, MonitorIntervalBackup, monitors[4].Interval)
	})

	t.Run("Multiple Apps", func(t *testing.T) {
		t.Parallel()

		def := &Definition{
			Project: Project{Title: "My Project"},
			Apps: []App{
				{
					Name:  "web",
					Title: "Web",
					Domains: []Domain{
						{Name: "example.com", Type: DomainTypePrimary},
					},
					Monitoring: true,
				},
				{
					Name:  "api",
					Title: "API",
					Domains: []Domain{
						{Name: "api.example.com", Type: DomainTypePrimary},
					},
					Monitoring: true,
				},
			},
		}

		monitors := def.GenerateMonitors()
		require.Len(t, monitors, 5) // 2 apps × 2 types (HTTP + DNS) + Codebase Backup

		// First app - HTTP.
		assert.Equal(t, "HTTP - example.com", monitors[0].Name)
		assert.Equal(t, MonitorTypeHTTP, monitors[0].Type)
		assert.Equal(t, MonitorIntervalHTTP, monitors[0].Interval)

		// First app - DNS.
		assert.Equal(t, "DNS - example.com", monitors[1].Name)
		assert.Equal(t, MonitorTypeDNS, monitors[1].Type)
		assert.Equal(t, MonitorIntervalDNS, monitors[1].Interval)

		// Second app - HTTP.
		assert.Equal(t, "HTTP - api.example.com", monitors[2].Name)
		assert.Equal(t, MonitorTypeHTTP, monitors[2].Type)
		assert.Equal(t, MonitorIntervalHTTP, monitors[2].Interval)

		// Second app - DNS.
		assert.Equal(t, "DNS - api.example.com", monitors[3].Name)
		assert.Equal(t, MonitorTypeDNS, monitors[3].Type)
		assert.Equal(t, MonitorIntervalDNS, monitors[3].Interval)

		// Codebase backup monitor.
		assert.Equal(t, "Backup - Codebase", monitors[4].Name)
		assert.Equal(t, MonitorTypePush, monitors[4].Type)
		assert.Equal(t, MonitorIntervalBackup, monitors[4].Interval)
	})

	t.Run("Globally disabled", func(t *testing.T) {
		t.Parallel()

		def := &Definition{
			Project: Project{Title: "My Project"},
			Monitoring: Monitoring{
				Enabled: boolPtr(false),
			},
			Apps: []App{
				{
					Name:  "web",
					Title: "Web",
					Domains: []Domain{
						{Name: "example.com", Type: DomainTypePrimary},
					},
					Monitoring: true,
				},
			},
			Resources: []Resource{
				{
					Name:  "db",
					Title: "Database",
					Backup: Backup{
						Enabled: boolPtr(true),
					},
					Monitoring: true,
				},
			},
		}

		monitors := def.GenerateMonitors()
		assert.Empty(t, monitors, "Expected no monitors when globally disabled")
	})

	t.Run("Globally disabled with custom monitors", func(t *testing.T) {
		t.Parallel()

		def := &Definition{
			Project: Project{Title: "My Project"},
			Monitoring: Monitoring{
				Enabled: boolPtr(false),
				Custom: []Monitor{
					{
						Name:     "Custom HTTP Monitor",
						Type:     MonitorTypeHTTP,
						URL:      "https://example.com/health",
						Method:   "GET",
						Interval: 60,
					},
				},
			},
			Apps: []App{
				{
					Name:  "web",
					Title: "Web",
					Domains: []Domain{
						{Name: "example.com", Type: DomainTypePrimary},
					},
					Monitoring: true,
				},
			},
		}

		monitors := def.GenerateMonitors()
		assert.Empty(t, monitors, "Expected no monitors when globally disabled, even with custom monitors defined")
	})

	t.Run("Globally enabled", func(t *testing.T) {
		t.Parallel()

		def := &Definition{
			Project: Project{Title: "My Project"},
			Monitoring: Monitoring{
				Enabled: boolPtr(true),
			},
			Apps: []App{
				{
					Name:  "web",
					Title: "Web",
					Domains: []Domain{
						{Name: "example.com", Type: DomainTypePrimary},
					},
					Monitoring: true,
				},
			},
		}

		monitors := def.GenerateMonitors()
		// Should have HTTP + DNS + Codebase Backup monitors.
		require.Len(t, monitors, 3)
		assert.Equal(t, "HTTP - example.com", monitors[0].Name)
		assert.Equal(t, "DNS - example.com", monitors[1].Name)
		assert.Equal(t, "Backup - Codebase", monitors[2].Name)
	})

	t.Run("Default enabled when nil", func(t *testing.T) {
		t.Parallel()

		def := &Definition{
			Project:    Project{Title: "My Project"},
			Monitoring: Monitoring{Enabled: nil},
			Apps: []App{
				{
					Name:  "web",
					Title: "Web",
					Domains: []Domain{
						{Name: "example.com", Type: DomainTypePrimary},
					},
					Monitoring: true,
				},
			},
		}

		monitors := def.GenerateMonitors()
		// Should have HTTP + DNS + Codebase Backup monitors.
		require.Len(t, monitors, 3)
		assert.Equal(t, "HTTP - example.com", monitors[0].Name)
		assert.Equal(t, "DNS - example.com", monitors[1].Name)
		assert.Equal(t, "Backup - Codebase", monitors[2].Name)
	})
}

func TestMonitor_VariableName(t *testing.T) {
	t.Parallel()

	tt := map[string]struct {
		monitor  Monitor
		envShort string
		want     string
	}{
		"Backup monitor with db identifier": {
			monitor: Monitor{
				Name:       "Backup - Database",
				Identifier: "db",
			},
			envShort: "prod",
			want:     "PROD_DB_BACKUP_PING_URL",
		},
		"Backup monitor with codebase identifier": {
			monitor: Monitor{
				Name:       "Backup - Codebase",
				Identifier: "codebase",
			},
			envShort: "prod",
			want:     "PROD_CODEBASE_BACKUP_PING_URL",
		},
		"Maintenance monitor": {
			monitor: Monitor{
				Name:       "Maintenance - Web",
				Identifier: "web",
			},
			envShort: "prod",
			want:     "PROD_WEB_MAINTENANCE_PING_URL",
		},
		"Staging environment": {
			monitor: Monitor{
				Name:       "Backup - Database",
				Identifier: "db",
			},
			envShort: "stag",
			want:     "STAG_DB_BACKUP_PING_URL",
		},
		"Empty identifier returns empty string": {
			monitor: Monitor{
				Name:       "Backup - Database",
				Identifier: "",
			},
			envShort: "prod",
			want:     "",
		},
		"Identifier with spaces": {
			monitor: Monitor{
				Name:       "Backup - User Data",
				Identifier: "user data",
			},
			envShort: "prod",
			want:     "PROD_USER_DATA_BACKUP_PING_URL",
		},
	}

	for name, test := range tt {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			got := test.monitor.VariableName(test.envShort)
			assert.Equal(t, test.want, got)
		})
	}
}

func TestMonitoring_IsEnabled(t *testing.T) {
	t.Parallel()

	tt := map[string]struct {
		monitoring Monitoring
		want       bool
	}{
		"Nil defaults to true": {
			monitoring: Monitoring{Enabled: nil},
			want:       true,
		},
		"Explicit true": {
			monitoring: Monitoring{Enabled: boolPtr(true)},
			want:       true,
		},
		"Explicit false": {
			monitoring: Monitoring{Enabled: boolPtr(false)},
			want:       false,
		},
	}

	for name, test := range tt {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			got := test.monitoring.IsEnabled()
			assert.Equal(t, test.want, got)
		})
	}
}

func TestMonitoring_ApplyDefaults(t *testing.T) {
	t.Parallel()

	tt := map[string]struct {
		monitoring Monitoring
		want       bool
	}{
		"Nil becomes true": {
			monitoring: Monitoring{Enabled: nil},
			want:       true,
		},
		"Explicit true unchanged": {
			monitoring: Monitoring{Enabled: boolPtr(true)},
			want:       true,
		},
		"Explicit false unchanged": {
			monitoring: Monitoring{Enabled: boolPtr(false)},
			want:       false,
		},
	}

	for name, test := range tt {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			test.monitoring.applyDefaults()
			require.NotNil(t, test.monitoring.Enabled)
			assert.Equal(t, test.want, *test.monitoring.Enabled)
		})
	}
}

// boolPtr is a helper function to create a pointer to a boolean value.
func boolPtr(b bool) *bool {
	return &b
}

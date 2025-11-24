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
		url, _ := monitors[0].GetConfigString("url")
		assert.Equal(t, "https://example.com", url)
		method, _ := monitors[0].GetConfigString("method")
		assert.Equal(t, "GET", method)
		assert.Equal(t, MonitorIntervalHTTP, monitors[0].Interval)

		// DNS monitor.
		assert.Equal(t, "DNS - example.com", monitors[1].Name)
		assert.Equal(t, MonitorTypeDNS, monitors[1].Type)
		domain, _ := monitors[1].GetConfigString("domain")
		assert.Equal(t, "example.com", domain)
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
		url0, _ := monitors[0].GetConfigString("url")
		assert.Equal(t, "https://api.example.com", url0)
		assert.Equal(t, MonitorIntervalHTTP, monitors[0].Interval)

		// First domain - DNS.
		assert.Equal(t, "DNS - api.example.com", monitors[1].Name)
		assert.Equal(t, MonitorTypeDNS, monitors[1].Type)
		domain1, _ := monitors[1].GetConfigString("domain")
		assert.Equal(t, "api.example.com", domain1)
		assert.Equal(t, MonitorIntervalDNS, monitors[1].Interval)

		// Second domain (alias) - HTTP.
		assert.Equal(t, "HTTP - www.api.example.com", monitors[2].Name)
		assert.Equal(t, MonitorTypeHTTP, monitors[2].Type)
		url2, _ := monitors[2].GetConfigString("url")
		assert.Equal(t, "https://www.api.example.com", url2)
		assert.Equal(t, MonitorIntervalHTTP, monitors[2].Interval)

		// Second domain - DNS.
		assert.Equal(t, "DNS - www.api.example.com", monitors[3].Name)
		assert.Equal(t, MonitorTypeDNS, monitors[3].Type)
		domain3, _ := monitors[3].GetConfigString("domain")
		assert.Equal(t, "www.api.example.com", domain3)
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

func TestMonitor_ValidateConfig(t *testing.T) {
	t.Parallel()

	tt := map[string]struct {
		monitor Monitor
		wantErr bool
	}{
		"Valid HTTP monitor": {
			monitor: Monitor{
				Name:     "Test HTTP",
				Type:     MonitorTypeHTTP,
				Interval: 60,
				Config: map[string]any{
					"url":    "https://example.com",
					"method": "GET",
				},
			},
			wantErr: false,
		},
		"Valid HTTP keyword monitor": {
			monitor: Monitor{
				Name:     "Test Keyword",
				Type:     MonitorTypeHTTPKeyword,
				Interval: 120,
				Config: map[string]any{
					"url":            "https://example.com",
					"method":         "GET",
					"keyword":        "success",
					"invert_keyword": false,
				},
			},
			wantErr: false,
		},
		"Valid HTTP keyword monitor with inverted": {
			monitor: Monitor{
				Name:     "Test Inverted Keyword",
				Type:     MonitorTypeHTTPKeyword,
				Interval: 120,
				Config: map[string]any{
					"url":            "https://example.com",
					"method":         "GET",
					"keyword":        "error",
					"invert_keyword": true,
				},
			},
			wantErr: false,
		},
		"Valid DNS monitor": {
			monitor: Monitor{
				Name:     "Test DNS",
				Type:     MonitorTypeDNS,
				Interval: 300,
				Config: map[string]any{
					"domain": "example.com",
				},
			},
			wantErr: false,
		},
		"Valid Postgres monitor": {
			monitor: Monitor{
				Name:     "Test Postgres",
				Type:     MonitorTypePostgres,
				Interval: 60,
				Config: map[string]any{
					"connection_string": "postgresql://localhost:5432/db",
				},
			},
			wantErr: false,
		},
		"Valid Push monitor without config": {
			monitor: Monitor{
				Name:       "Test Push",
				Type:       MonitorTypePush,
				Interval:   90000,
				Identifier: "backup",
			},
			wantErr: false,
		},
		"HTTP monitor missing url": {
			monitor: Monitor{
				Name:     "Invalid HTTP",
				Type:     MonitorTypeHTTP,
				Interval: 60,
				Config: map[string]any{
					"method": "GET",
				},
			},
			wantErr: true,
		},
		"HTTP monitor missing method": {
			monitor: Monitor{
				Name:     "Invalid HTTP",
				Type:     MonitorTypeHTTP,
				Interval: 60,
				Config: map[string]any{
					"url": "https://example.com",
				},
			},
			wantErr: true,
		},
		"HTTP keyword monitor missing keyword": {
			monitor: Monitor{
				Name:     "Invalid Keyword",
				Type:     MonitorTypeHTTPKeyword,
				Interval: 120,
				Config: map[string]any{
					"url":    "https://example.com",
					"method": "GET",
				},
			},
			wantErr: true,
		},
		"DNS monitor missing domain": {
			monitor: Monitor{
				Name:     "Invalid DNS",
				Type:     MonitorTypeDNS,
				Interval: 300,
				Config:   map[string]any{},
			},
			wantErr: true,
		},
		"Postgres monitor missing connection_string": {
			monitor: Monitor{
				Name:     "Invalid Postgres",
				Type:     MonitorTypePostgres,
				Interval: 60,
				Config:   map[string]any{},
			},
			wantErr: true,
		},
		"HTTP monitor with nil config": {
			monitor: Monitor{
				Name:     "Invalid HTTP",
				Type:     MonitorTypeHTTP,
				Interval: 60,
				Config:   nil,
			},
			wantErr: true,
		},
	}

	for name, test := range tt {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			err := test.monitor.ValidateConfig()
			assert.Equal(t, test.wantErr, err != nil)
		})
	}
}

func TestMonitor_GetConfigString(t *testing.T) {
	t.Parallel()

	tt := map[string]struct {
		monitor Monitor
		key     string
		want    string
		wantOk  bool
	}{
		"Get existing string value": {
			monitor: Monitor{
				Config: map[string]any{
					"url": "https://example.com",
				},
			},
			key:    "url",
			want:   "https://example.com",
			wantOk: true,
		},
		"Get non-existent key": {
			monitor: Monitor{
				Config: map[string]any{
					"url": "https://example.com",
				},
			},
			key:    "missing",
			want:   "",
			wantOk: false,
		},
		"Get with nil config": {
			monitor: Monitor{
				Config: nil,
			},
			key:    "url",
			want:   "",
			wantOk: false,
		},
		"Get non-string value": {
			monitor: Monitor{
				Config: map[string]any{
					"port": 8080,
				},
			},
			key:    "port",
			want:   "",
			wantOk: false,
		},
	}

	for name, test := range tt {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			got, ok := test.monitor.GetConfigString(test.key)
			assert.Equal(t, test.want, got)
			assert.Equal(t, test.wantOk, ok)
		})
	}
}

func TestMonitor_GetConfigInt(t *testing.T) {
	t.Parallel()

	tt := map[string]struct {
		monitor Monitor
		key     string
		want    int
		wantOk  bool
	}{
		"Get existing int value": {
			monitor: Monitor{
				Config: map[string]any{
					"max_redirects": 5,
				},
			},
			key:    "max_redirects",
			want:   5,
			wantOk: true,
		},
		"Get non-existent key": {
			monitor: Monitor{
				Config: map[string]any{
					"max_redirects": 5,
				},
			},
			key:    "missing",
			want:   0,
			wantOk: false,
		},
		"Get with nil config": {
			monitor: Monitor{
				Config: nil,
			},
			key:    "max_redirects",
			want:   0,
			wantOk: false,
		},
	}

	for name, test := range tt {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			got, ok := test.monitor.GetConfigInt(test.key)
			assert.Equal(t, test.want, got)
			assert.Equal(t, test.wantOk, ok)
		})
	}
}

func TestMonitor_GetConfigBool(t *testing.T) {
	t.Parallel()

	tt := map[string]struct {
		monitor Monitor
		key     string
		want    bool
		wantOk  bool
	}{
		"Get existing bool value true": {
			monitor: Monitor{
				Config: map[string]any{
					"invert_keyword": true,
				},
			},
			key:    "invert_keyword",
			want:   true,
			wantOk: true,
		},
		"Get existing bool value false": {
			monitor: Monitor{
				Config: map[string]any{
					"invert_keyword": false,
				},
			},
			key:    "invert_keyword",
			want:   false,
			wantOk: true,
		},
		"Get non-existent key": {
			monitor: Monitor{
				Config: map[string]any{
					"invert_keyword": true,
				},
			},
			key:    "missing",
			want:   false,
			wantOk: false,
		},
		"Get with nil config": {
			monitor: Monitor{
				Config: nil,
			},
			key:    "invert_keyword",
			want:   false,
			wantOk: false,
		},
	}

	for name, test := range tt {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			got, ok := test.monitor.GetConfigBool(test.key)
			assert.Equal(t, test.want, got)
			assert.Equal(t, test.wantOk, ok)
		})
	}
}

func TestMonitorType_String_HTTPKeyword(t *testing.T) {
	t.Parallel()

	got := MonitorTypeHTTPKeyword.String()
	assert.Equal(t, "http-keyword", got)
}

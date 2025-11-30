package appdef

import (
	"errors"
	"fmt"
	"strings"
)

type (
	// Monitoring is the project-level monitoring configuration.
	// It consolidates status page settings and custom monitors.
	//
	// Monitoring is enabled by default (opt-out pattern).
	Monitoring struct {
		Enabled    *bool      `json:"enabled,omitempty" description:"Enable or disable all monitoring globally (defaults to true)"`
		StatusPage StatusPage `json:"statusPage,omitempty" description:"Public status page configuration"`
		Custom     []Monitor  `json:"custom,omitempty" description:"Custom monitors beyond auto-generated ones"`
	}
	// StatusPage defines the configuration for a project's status page.
	// This information is used for custom domain setup and status page presentation.
	StatusPage struct {
		Domain string `json:"domain,omitempty" validate:"omitempty,fqdn" description:"Custom domain for the status page (e.g., status.example.com). If not set, Terraform will not configure a custom domain."`
		Slug   string `json:"slug,omitempty" validate:"omitempty,lowercase,alphanumdash" description:"URL slug for the status page (e.g., my-project). If not set, defaults to a kebab-case version of the project name."`
		Theme  string `json:"theme,omitempty" validate:"omitempty,oneof=auto light dark" description:"Theme for the status page (auto, light, dark). Defaults to auto."`
	}
	// Monitor contains minimal monitoring configuration.
	//
	// Config field usage by monitor type:
	// - HTTP monitors: {url, method, max_redirects}
	// - HTTP-Keyword monitors: {url, method, keyword, invert_keyword, max_redirects}
	// - DNS monitors: {domain, resolver_type}
	// - Postgres monitors: {connection_string}
	// - Push monitors: No config required
	Monitor struct {
		Name       string      `json:"name" validate:"required" description:"Unique monitor name"`
		Type       MonitorType `json:"type" validate:"required,oneof=http http-keyword dns postgres push" description:"Monitor type (http, http-keyword, dns, postgres, push)"`
		Interval   int         `json:"interval,omitempty" description:"Interval in seconds between checks (defaults based on monitor type if not specified)"`
		Identifier string      `json:"identifier,omitempty" description:"Machine-readable identifier for variable naming (e.g., 'db' for database). Used by VariableName() method."`
		Config     Config      `json:"config,omitempty" description:"Type-specific monitor configuration (e.g., url, method, keyword, domain)"`
	}
	// MonitorType defines the type of monitor.
	MonitorType string
)

// MonitorType constants.
const (
	MonitorTypeHTTP        MonitorType = "http"
	MonitorTypeHTTPKeyword MonitorType = "http-keyword"
	MonitorTypeDNS         MonitorType = "dns"
	MonitorTypePostgres    MonitorType = "postgres"
	MonitorTypePush        MonitorType = "push"
)

// Monitor interval constants (in seconds).
const (
	// MonitorIntervalMin is the minimum interval required by Peekaping provider.
	MonitorIntervalMin = 20
	// MonitorIntervalHTTP is 1 minute for HTTP health checks.
	MonitorIntervalHTTP = 60
	// MonitorIntervalDNS is 5 minutes for DNS resolution checks.
	MonitorIntervalDNS = 300
	// MonitorIntervalBackup is 25 hours (90000s) for daily backup heartbeats with 1 hour buffer.
	MonitorIntervalBackup = 90000
	// MonitorIntervalMaintenance is 8 days (691200s) for weekly maintenance with 1 day buffer.
	MonitorIntervalMaintenance = 691200
)

// Monitor config constants.
const (
	// MonitorMaxRedirectsDefault is the default maximum redirects for HTTP monitors.
	MonitorMaxRedirectsDefault = 3
)

// String implements fmt.Stringer on MonitorType.
func (m MonitorType) String() string {
	return string(m)
}

// IsEnabled returns whether monitoring is globally enabled.
// Defaults to true when the field is nil or explicitly set to true.
func (m *Monitoring) IsEnabled() bool {
	if m.Enabled == nil {
		return true
	}
	return *m.Enabled
}

// applyDefaults sets default values for monitoring configuration.
func (m *Monitoring) applyDefaults() {
	// Default monitoring to enabled (opt-out).
	if m.Enabled == nil {
		enabled := true
		m.Enabled = &enabled
	}
}

// VariableName returns the GitHub Actions variable name for this monitor's ping URL.
// Format: {ENV}_{IDENTIFIER}_{TYPE}_PING_URL (e.g., PROD_DB_BACKUP_PING_URL).
// Only applicable for push monitors with an Identifier set.
func (m *Monitor) VariableName(envShort string) string {
	if m.Identifier == "" {
		return ""
	}

	// Extract type from monitor name (e.g., "Backup" from "Backup - Database").
	monitorType := "UNKNOWN"
	if parts := strings.SplitN(m.Name, " - ", 2); len(parts) > 0 {
		monitorType = strings.ToUpper(strings.ReplaceAll(parts[0], " ", "_"))
	}

	return fmt.Sprintf("%s_%s_%s_PING_URL",
		strings.ToUpper(envShort),
		strings.ToUpper(strings.ReplaceAll(m.Identifier, " ", "_")),
		monitorType,
	)
}

// monitorValidators maps monitor types to their validation functions.
var monitorValidators = map[MonitorType]func(*Monitor) error{
	MonitorTypeHTTP: func(m *Monitor) error {
		if m.Config == nil {
			return errors.New("http monitor requires config")
		}
		if _, ok := m.Config.String("url"); !ok {
			return errors.New("http monitor requires 'url' in config")
		}
		if _, ok := m.Config.String("method"); !ok {
			return errors.New("http monitor requires 'method' in config")
		}
		return nil
	},
	MonitorTypeHTTPKeyword: func(m *Monitor) error {
		if m.Config == nil {
			return errors.New("http-keyword monitor requires config")
		}
		if _, ok := m.Config.String("url"); !ok {
			return errors.New("http-keyword monitor requires 'url' in config")
		}
		if _, ok := m.Config.String("method"); !ok {
			return errors.New("http-keyword monitor requires 'method' in config")
		}
		if _, ok := m.Config.String("keyword"); !ok {
			return errors.New("http-keyword monitor requires 'keyword' in config")
		}
		return nil
	},
	MonitorTypeDNS: func(m *Monitor) error {
		if m.Config == nil {
			return errors.New("dns monitor requires config")
		}
		if _, ok := m.Config.String("domain"); !ok {
			return errors.New("dns monitor requires 'domain' in config")
		}
		return nil
	},
	MonitorTypePostgres: func(m *Monitor) error {
		if m.Config == nil {
			return errors.New("postgres monitor requires config")
		}
		if _, ok := m.Config.String("connection_string"); !ok {
			return errors.New("postgres monitor requires 'connection_string' in config")
		}
		return nil
	},
	MonitorTypePush: func(m *Monitor) error {
		// Push monitors don't require config.
		return nil
	},
}

// ValidateConfig ensures the monitor has the required config fields for its type.
func (m *Monitor) ValidateConfig() error {
	// Validate interval if explicitly set.
	if m.Interval != 0 && m.Interval < MonitorIntervalMin {
		return fmt.Errorf("monitor interval must be at least %d seconds (got %d)", MonitorIntervalMin, m.Interval)
	}

	validator, ok := monitorValidators[m.Type]
	if !ok {
		return fmt.Errorf("unknown monitor type: %s", m.Type)
	}
	return validator(m)
}

// applyDefaults sets default values for the monitor.
// If interval is not set (0), applies sensible defaults based on monitor type.
// Also sets max_redirects for HTTP/HTTP-keyword monitors if not provided.
func (m *Monitor) applyDefaults() {
	// Apply default interval if not explicitly set (0).
	if m.Interval == 0 {
		switch m.Type {
		case MonitorTypeHTTP, MonitorTypeHTTPKeyword, MonitorTypePostgres:
			m.Interval = MonitorIntervalHTTP // 60 seconds
		case MonitorTypeDNS:
			m.Interval = MonitorIntervalDNS // 300 seconds (5 minutes)
		case MonitorTypePush:
			m.Interval = MonitorIntervalBackup // 90000 seconds (25 hours)
		}
	}

	// Apply max_redirects default for HTTP and HTTP-keyword monitors.
	if m.Type == MonitorTypeHTTP || m.Type == MonitorTypeHTTPKeyword {
		if m.Config != nil {
			if _, ok := m.Config.Int("max_redirects"); !ok {
				m.Config["max_redirects"] = MonitorMaxRedirectsDefault
			}
		}
	}
}

// GenerateMonitors creates all monitors for the definition.
// This includes:
// - HTTP and DNS monitors for app domains
// - Backup monitors for resources
// - Codebase backup monitor (always generated)
// - Maintenance monitors for VM apps
// - Custom monitors from project configuration
func (d *Definition) GenerateMonitors() []Monitor {
	// Return empty slice if monitoring is globally disabled.
	if !d.Monitoring.IsEnabled() {
		return []Monitor{}
	}

	monitors := make([]Monitor, 0)

	// Generate HTTP and DNS monitors for all apps.
	monitors = append(monitors, d.generateHTTPDNSMonitors()...)

	// Generate backup monitors for all resources.
	monitors = append(monitors, d.generateResourceBackupMonitors()...)

	// Generate codebase backup monitor (always generated).
	monitors = append(monitors, Monitor{
		Name:       "Backup - Codebase",
		Type:       MonitorTypePush,
		Interval:   MonitorIntervalBackup,
		Identifier: "codebase",
	})

	// Generate maintenance monitors for all apps.
	monitors = append(monitors, d.generateMaintenanceMonitors()...)

	// Append custom monitors from root configuration.
	// Apply defaults to custom monitors before adding them.
	for i := range d.Monitoring.Custom {
		d.Monitoring.Custom[i].applyDefaults()
		monitors = append(monitors, d.Monitoring.Custom[i])
	}

	return monitors
}

// generateHTTPDNSMonitors creates HTTP and DNS monitors for all apps in the definition.
// It generates two monitors per domain (HTTP + DNS) for primary and alias domains,
// excluding unmanaged domains. Monitoring must be explicitly enabled in each app's configuration.
func (d *Definition) generateHTTPDNSMonitors() []Monitor {
	monitors := make([]Monitor, 0)

	for _, app := range d.Apps {
		// Only skip if monitoring is explicitly disabled
		if app.Monitoring != nil && !*app.Monitoring {
			continue
		}

		for _, domain := range app.Domains {
			if domain.Type == DomainTypeUnmanaged {
				continue
			}

			// HTTP monitor - checks the availability of the web application.
			monitors = append(monitors, Monitor{
				Name:     fmt.Sprintf("HTTP - %s", domain.Name),
				Type:     MonitorTypeHTTP,
				Interval: MonitorIntervalHTTP,
				Config: map[string]any{
					"url":           fmt.Sprintf("https://%s", domain.Name),
					"method":        "GET",
					"max_redirects": 3,
				},
			})

			// DNS monitor - checks domain name resolution.
			monitors = append(monitors, Monitor{
				Name:     fmt.Sprintf("DNS - %s", domain.Name),
				Type:     MonitorTypeDNS,
				Interval: MonitorIntervalDNS,
				Config: map[string]any{
					"domain": domain.Name,
				},
			})
		}
	}

	return monitors
}

// generateResourceBackupMonitors creates push monitors for resource backup workflows.
func (d *Definition) generateResourceBackupMonitors() []Monitor {
	monitors := make([]Monitor, 0)

	for _, resource := range d.Resources {
		// Only generate backup monitor if both backup and monitoring are enabled.
		if resource.Backup == nil || !resource.Backup.Enabled {
			continue
		}
		// Only skip if monitoring is explicitly disabled
		if resource.Monitoring != nil && !*resource.Monitoring {
			continue
		}

		monitors = append(monitors, Monitor{
			Name:       fmt.Sprintf("Backup - %s", resource.Title),
			Type:       MonitorTypePush,
			Interval:   MonitorIntervalBackup,
			Identifier: resource.Name,
		})
	}

	return monitors
}

// generateMaintenanceMonitors creates push monitors for VM app maintenance workflows.
func (d *Definition) generateMaintenanceMonitors() []Monitor {
	monitors := make([]Monitor, 0)

	for _, app := range d.Apps {
		// Only generate maintenance monitor for VM apps with monitoring enabled.
		// Only skip if monitoring is explicitly disabled
		if app.Monitoring != nil && !*app.Monitoring {
			continue
		}
		if app.Infra.Type != "vm" {
			continue
		}

		monitors = append(monitors, Monitor{
			Name:       fmt.Sprintf("Maintenance - %s", app.Title),
			Type:       MonitorTypePush,
			Interval:   MonitorIntervalMaintenance,
			Identifier: app.Name,
		})
	}

	return monitors
}

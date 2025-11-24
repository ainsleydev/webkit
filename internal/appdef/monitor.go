package appdef

import (
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
		Name       string         `json:"name" validate:"required" description:"Unique monitor name"`
		Type       MonitorType    `json:"type" validate:"required,oneof=http http-keyword dns postgres push" description:"Monitor type (http, http-keyword, dns, postgres, push)"`
		Interval   int            `json:"interval,omitempty" description:"Interval in seconds between checks (defaults based on monitor type if not specified)"`
		Identifier string         `json:"identifier,omitempty" description:"Machine-readable identifier for variable naming (e.g., 'db' for database). Used by VariableName() method."`
		Config     map[string]any `json:"config,omitempty" description:"Type-specific monitor configuration (e.g., url, method, keyword, domain)"`
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
	// MonitorIntervalHTTP is 1 minute for HTTP health checks.
	MonitorIntervalHTTP = 60
	// MonitorIntervalDNS is 5 minutes for DNS resolution checks.
	MonitorIntervalDNS = 300
	// MonitorIntervalBackup is 25 hours (90000s) for daily backup heartbeats with 1 hour buffer.
	MonitorIntervalBackup = 90000
	// MonitorIntervalMaintenance is 8 days (691200s) for weekly maintenance with 1 day buffer.
	MonitorIntervalMaintenance = 691200
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
func (m Monitor) VariableName(envShort string) string {
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

// GetConfigString safely retrieves a string value from the monitor config.
// Returns the value and true if found, empty string and false otherwise.
func (m Monitor) GetConfigString(key string) (string, bool) {
	if m.Config == nil {
		return "", false
	}
	val, ok := m.Config[key].(string)
	return val, ok
}

// GetConfigInt safely retrieves an int value from the monitor config.
// Returns the value and true if found, 0 and false otherwise.
func (m Monitor) GetConfigInt(key string) (int, bool) {
	if m.Config == nil {
		return 0, false
	}
	val, ok := m.Config[key].(int)
	return val, ok
}

// GetConfigBool safely retrieves a bool value from the monitor config.
// Returns the value and true if found, false and false otherwise.
func (m Monitor) GetConfigBool(key string) (bool, bool) {
	if m.Config == nil {
		return false, false
	}
	val, ok := m.Config[key].(bool)
	return val, ok
}

// ValidateConfig ensures the monitor has the required config fields for its type.
func (m Monitor) ValidateConfig() error {
	// Push monitors don't require config.
	if m.Type == MonitorTypePush {
		return nil
	}

	if m.Config == nil {
		return fmt.Errorf("%s monitor requires config", m.Type)
	}

	switch m.Type {
	case MonitorTypeHTTP:
		if _, ok := m.GetConfigString("url"); !ok {
			return fmt.Errorf("http monitor requires 'url' in config")
		}
		if _, ok := m.GetConfigString("method"); !ok {
			return fmt.Errorf("http monitor requires 'method' in config")
		}

	case MonitorTypeHTTPKeyword:
		if _, ok := m.GetConfigString("url"); !ok {
			return fmt.Errorf("http-keyword monitor requires 'url' in config")
		}
		if _, ok := m.GetConfigString("method"); !ok {
			return fmt.Errorf("http-keyword monitor requires 'method' in config")
		}
		if _, ok := m.GetConfigString("keyword"); !ok {
			return fmt.Errorf("http-keyword monitor requires 'keyword' in config")
		}

	case MonitorTypeDNS:
		if _, ok := m.GetConfigString("domain"); !ok {
			return fmt.Errorf("dns monitor requires 'domain' in config")
		}

	case MonitorTypePostgres:
		if _, ok := m.GetConfigString("connection_string"); !ok {
			return fmt.Errorf("postgres monitor requires 'connection_string' in config")
		}

	default:
		return fmt.Errorf("unknown monitor type: %s", m.Type)
	}

	return nil
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
	monitors = append(monitors, d.Monitoring.Custom...)

	return monitors
}

// generateHTTPDNSMonitors creates HTTP and DNS monitors for all apps in the definition.
// It generates two monitors per domain (HTTP + DNS) for primary and alias domains,
// excluding unmanaged domains. Monitoring must be explicitly enabled in each app's configuration.
func (d *Definition) generateHTTPDNSMonitors() []Monitor {
	monitors := make([]Monitor, 0)

	for _, app := range d.Apps {
		if !app.Monitoring {
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
		if !resource.Backup.Enabled || !resource.Monitoring {
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
		if !app.Monitoring || app.Infra.Type != "vm" {
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

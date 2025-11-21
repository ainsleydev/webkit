package appdef

import (
	"fmt"
)

type (
	// Monitoring is the user-facing config in app.json.
	// It's intentionally simple - just an enabled flag.
	Monitoring struct {
		Enabled bool `json:"enabled" description:"Whether to enable uptime monitoring for this app or resource (defaults to true)"`
	}
	// Monitor contains minimal monitoring configuration.
	// Defaults are applied by the Terraform layer based on monitor type.
	//
	// Field usage by monitor type:
	// - HTTP monitors: URL contains the full URL (including path), Method contains HTTP method
	// - DNS monitors: Domain contains the domain name to check
	// - Database monitors: URL contains database connection string or Terraform reference, Method is empty
	// - Push monitors: URL and Method are empty
	Monitor struct {
		Name   string      // Unique monitor name.
		Type   MonitorType // Monitor type (http, dns, postgres, push).
		URL    string      // URL for HTTP monitors, database connection string for postgres monitors.
		Method string      // HTTP method for HTTP monitors (e.g., "GET"), empty for other types.
		Domain string      // Domain name for DNS monitors.
	}
	// MonitorType defines the type of monitor.
	MonitorType string
)

// MonitorType constants.
const (
	MonitorTypeHTTP     MonitorType = "http"
	MonitorTypeDNS      MonitorType = "dns"
	MonitorTypePostgres MonitorType = "postgres"
	MonitorTypePush     MonitorType = "push"
)

// String implements fmt.Stringer on MonitorType.
func (m MonitorType) String() string {
	return string(m)
}

// GenerateMonitors creates all monitors for the definition.
// This includes:
// - HTTP and DNS monitors for app domains
// - Backup monitors for resources
// - Codebase backup monitor (always generated)
// - Maintenance monitors for VM apps
func (d *Definition) GenerateMonitors() []Monitor {
	monitors := make([]Monitor, 0)

	// Generate HTTP and DNS monitors for all apps.
	monitors = append(monitors, d.generateHTTPDNSMonitors()...)

	// Generate backup monitors for all resources.
	monitors = append(monitors, d.generateResourceBackupMonitors()...)

	// Generate codebase backup monitor (always generated).
	monitors = append(monitors, Monitor{
		Name: fmt.Sprintf("%s - Codebase Backup", d.Project.Title),
		Type: MonitorTypePush,
	})

	// Generate maintenance monitors for all apps.
	for _, app := range d.Apps {
		if monitor := app.GenerateMaintenanceMonitor(d.Project.Title); monitor != nil {
			monitors = append(monitors, *monitor)
		}
	}

	return monitors
}

// generateHTTPDNSMonitors creates HTTP and DNS monitors for all apps in the definition.
// It generates two monitors per domain (HTTP + DNS) for primary and alias domains,
// excluding unmanaged domains. Monitoring must be explicitly enabled in each app's configuration.
func (d *Definition) generateHTTPDNSMonitors() []Monitor {
	monitors := make([]Monitor, 0)

	for _, app := range d.Apps {
		if !app.Monitoring.Enabled {
			continue
		}

		for _, domain := range app.Domains {
			if domain.Type == DomainTypeUnmanaged {
				continue
			}

			// HTTP monitor - checks the availability of the web application.
			monitors = append(monitors, Monitor{
				Name:   fmt.Sprintf("%s, %s - %s", d.Project.Title, app.Title, domain.Name),
				Type:   MonitorTypeHTTP,
				URL:    fmt.Sprintf("https://%s", domain.Name),
				Method: "GET",
			})

			// DNS monitor - checks domain name resolution.
			monitors = append(monitors, Monitor{
				Name:   fmt.Sprintf("%s, %s DNS - %s", d.Project.Title, app.Title, domain.Name),
				Type:   MonitorTypeDNS,
				Domain: domain.Name,
			})
		}
	}

	return monitors
}

// generateResourceBackupMonitors creates push monitors for resource backup workflows.
func (d *Definition) generateResourceBackupMonitors() []Monitor {
	monitors := make([]Monitor, 0)

	for _, resource := range d.Resources {
		if monitor := resource.GenerateBackupMonitor(d.Project.Title); monitor != nil {
			monitors = append(monitors, *monitor)
		}
	}

	return monitors
}

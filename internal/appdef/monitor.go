package appdef

import (
	"fmt"
)

type (
	// Monitoring is the project-level monitoring configuration.
	// It consolidates status page settings and custom monitors.
	Monitoring struct {
		StatusPage StatusPage `json:"statusPage,omitempty" description:"Public status page configuration"`
		Custom     []Monitor  `json:"custom,omitempty" description:"Custom monitors beyond auto-generated ones"`
	}
	// StatusPage defines the configuration for a project's status page.
	// This information is used for custom domain setup and status page presentation.
	StatusPage struct {
		Domain string `json:"domain,omitempty" validate:"omitempty,fqdn" description:"Custom domain for the status page (e.g., status.example.com). If not set, Terraform will not configure a custom domain."`
	}
	// Monitor contains minimal monitoring configuration.
	//
	// Field usage by monitor type:
	// - HTTP monitors: URL contains the full URL (including path), Method contains HTTP method
	// - DNS monitors: Domain contains the domain name to check
	// - Database monitors: URL contains database connection string or Terraform reference, Method is empty
	// - Push monitors: URL and Method are empty
	Monitor struct {
		Name     string      `json:"name" validate:"required" description:"Unique monitor name"`
		Type     MonitorType `json:"type" validate:"required,oneof=http dns postgres push" description:"Monitor type (http, dns, postgres, push)"`
		URL      string      `json:"url,omitempty" description:"URL for HTTP monitors or database connection string for postgres monitors"`
		Method   string      `json:"method,omitempty" description:"HTTP method for HTTP monitors (e.g., GET, POST)"`
		Domain   string      `json:"domain,omitempty" description:"Domain name for DNS monitors"`
		Interval int         `json:"interval,omitempty" description:"Interval in seconds between checks (defaults based on monitor type if not specified)"`
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

// GenerateMonitors creates all monitors for the definition.
// This includes:
// - HTTP and DNS monitors for app domains
// - Backup monitors for resources
// - Codebase backup monitor (always generated)
// - Maintenance monitors for VM apps
// - Custom monitors from project configuration
func (d *Definition) GenerateMonitors() []Monitor {
	monitors := make([]Monitor, 0)

	// Generate HTTP and DNS monitors for all apps.
	monitors = append(monitors, d.generateHTTPDNSMonitors()...)

	// Generate backup monitors for all resources.
	monitors = append(monitors, d.generateResourceBackupMonitors()...)

	// Generate codebase backup monitor (always generated).
	monitors = append(monitors, Monitor{
		Name:     fmt.Sprintf("%s - Codebase Backup", d.Project.Title),
		Type:     MonitorTypePush,
		Interval: MonitorIntervalBackup,
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
				Name:     fmt.Sprintf("%s, %s - %s", d.Project.Title, app.Title, domain.Name),
				Type:     MonitorTypeHTTP,
				URL:      fmt.Sprintf("https://%s", domain.Name),
				Method:   "GET",
				Interval: MonitorIntervalHTTP,
			})

			// DNS monitor - checks domain name resolution.
			monitors = append(monitors, Monitor{
				Name:     fmt.Sprintf("%s, %s DNS - %s", d.Project.Title, app.Title, domain.Name),
				Type:     MonitorTypeDNS,
				Domain:   domain.Name,
				Interval: MonitorIntervalDNS,
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
			Name:     fmt.Sprintf("%s - %s Backup", d.Project.Title, resource.Title),
			Type:     MonitorTypePush,
			Interval: MonitorIntervalBackup,
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
			Name:     fmt.Sprintf("%s - %s Maintenance", d.Project.Title, app.Title),
			Type:     MonitorTypePush,
			Interval: MonitorIntervalMaintenance,
		})
	}

	return monitors
}

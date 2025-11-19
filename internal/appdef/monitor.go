package appdef

import (
	"fmt"
	"strings"

	"github.com/ainsleydev/webkit/pkg/env"
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
	// - Database monitors: URL contains database connection string or Terraform reference, Method is empty
	// - Push monitors: URL and Method are empty
	Monitor struct {
		Name   string      // Unique monitor name.
		Type   MonitorType // Monitor type (http, postgres, push).
		URL    string      // URL for HTTP monitors, database connection string for postgres monitors.
		Method string      // HTTP method for HTTP monitors (e.g., "GET"), empty for other types.
	}
	// MonitorType defines the type of monitor.
	MonitorType string
)

// MonitorType constants.
const (
	MonitorTypeHTTP     MonitorType = "http"
	MonitorTypePostgres MonitorType = "postgres"
	MonitorTypePush     MonitorType = "push"
)

// String implements fmt.Stringer on MonitorType.
func (m MonitorType) String() string {
	return string(m)
}

// GenerateMonitors creates HTTP monitors for all domains in the app.
// It generates one monitor per domain (primary + aliases), excluding unmanaged domains.
// Monitoring must be explicitly enabled in the app configuration.
func (a *App) GenerateMonitors() []Monitor {
	if !a.Monitoring.Enabled {
		return nil
	}

	monitors := make([]Monitor, 0)

	// Create HTTP monitor for each domain (primary + aliases).
	for _, domain := range a.Domains {
		if domain.Type == DomainTypeUnmanaged {
			continue
		}

		monitors = append(monitors, Monitor{
			Name:   fmt.Sprintf("%s-%s", a.Name, sanitiseMonitorName(domain.Name)),
			Type:   MonitorTypeHTTP,
			URL:    fmt.Sprintf("https://%s%s", domain.Name, a.healthCheckPath()),
			Method: "GET",
		})
	}

	return monitors
}

// healthCheckPath extracts the health check path from the app's infra config.
// It defaults to "/" if not specified.
func (a *App) healthCheckPath() string {
	if a.Infra.Config == nil {
		return "/"
	}

	if path, ok := a.Infra.Config["health_check_path"].(string); ok && path != "" {
		return path
	}

	return "/"
}

// GenerateMonitors creates monitors for resources based on their type.
// Currently only Postgres databases are supported for monitoring.
// Monitoring must be explicitly enabled in the resource configuration.
func (r *Resource) GenerateMonitors(enviro env.Environment, dbURLGenerator func(*Resource, env.Environment, string) string) []Monitor {
	if !r.Monitoring.Enabled {
		return nil
	}

	// Only Postgres supported for now (Uptime Kuma limitation).
	if r.Type != ResourceTypePostgres {
		return nil
	}

	return []Monitor{
		{
			Name:   fmt.Sprintf("%s-%s", r.Name, enviro),
			Type:   MonitorTypePostgres,
			URL:    dbURLGenerator(r, enviro, "connection_url"),
			Method: "", // Empty for database monitors.
		},
	}
}

// GenerateHeartbeatMonitor creates a push monitor for backup job heartbeats.
// The monitor expects a heartbeat signal after each successful backup.
// Note: Push monitors don't use URL or Method fields.
func (r *Resource) GenerateHeartbeatMonitor(cronSchedule string) Monitor {
	if !r.Backup.Enabled {
		return Monitor{}
	}

	return Monitor{
		Name:   fmt.Sprintf("backup-%s", r.Name),
		Type:   MonitorTypePush,
		URL:    "", // Empty for push monitors.
		Method: "", // Empty for push monitors.
	}
}

// sanitiseMonitorName converts a domain name to a valid monitor name component.
// It replaces dots with hyphens to create Terraform-safe resource names.
//
// Example:
//
//	sanitiseMonitorName("api.example.com") -> "api-example-com"
func sanitiseMonitorName(domain string) string {
	return strings.ReplaceAll(domain, ".", "-")
}

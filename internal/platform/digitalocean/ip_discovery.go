package digitalocean

import (
	"context"
	"fmt"

	"github.com/ainsleydev/webkit/internal/appdef"
	"github.com/ainsleydev/webkit/internal/infra"
)

// IPDiscovery handles discovery of IP addresses for apps deployed
// on DigitalOcean infrastructure.
type IPDiscovery struct {
	apiToken string
}

// NewIPDiscovery creates a new IP discovery service.
func NewIPDiscovery(apiToken string) *IPDiscovery {
	return &IPDiscovery{
		apiToken: apiToken,
	}
}

// DiscoverAppIPs discovers the IP addresses for an app based on its
// infrastructure type (VM or container).
//
// For VMs: Returns the droplet's public IPv4 address from terraform outputs.
// For Containers: Returns the DigitalOcean App Platform egress IPs for the region.
func (d *IPDiscovery) DiscoverAppIPs(ctx context.Context, app appdef.App, tfOutput infra.OutputResult) ([]string, error) {
	switch app.Infra.Type {
	case "vm":
		return d.getVMIP(app, tfOutput)
	case "container":
		return d.getContainerEgressIPs(ctx, app)
	default:
		return nil, fmt.Errorf("unsupported infra type: %s", app.Infra.Type)
	}
}

// getVMIP retrieves the IP address of a VM (droplet) from terraform outputs.
func (d *IPDiscovery) getVMIP(app appdef.App, tfOutput infra.OutputResult) ([]string, error) {
	// Look for the app in terraform outputs
	appOutput, ok := tfOutput.Apps[app.Name]
	if !ok {
		return nil, fmt.Errorf("app %q not found in terraform outputs", app.Name)
	}

	// Get the IP address
	ip, ok := appOutput["ipv4_address"]
	if !ok {
		return nil, fmt.Errorf("ipv4_address not found in terraform outputs for app %q", app.Name)
	}

	ipStr, ok := ip.(string)
	if !ok {
		return nil, fmt.Errorf("ipv4_address is not a string for app %q", app.Name)
	}

	return []string{ipStr}, nil
}

// getContainerEgressIPs retrieves the egress IP addresses for an App Platform
// container based on its region.
//
// DigitalOcean App Platform apps share egress IPs within their datacenter region.
// Currently, this returns an error as automatic discovery for containers requires
// additional DO API integration. Users should manually add IPs to the resource config.
func (d *IPDiscovery) getContainerEgressIPs(ctx context.Context, app appdef.App) ([]string, error) {
	region, ok := app.Infra.Config["region"].(string)
	if !ok || region == "" {
		return nil, fmt.Errorf("region not found in app %q config", app.Name)
	}

	// TODO: Implement automatic container IP discovery.
	// Options:
	// 1. Query DO API for app egress IPs (requires godo library).
	// 2. Use known datacenter egress IP ranges for the region.
	// 3. Make a test connection from the app and capture the IP.

	return nil, fmt.Errorf("automatic IP discovery for containers (App Platform) not yet implemented for app %q (region: %s) - please manually add IPs to postgres resource config", app.Name, region)
}

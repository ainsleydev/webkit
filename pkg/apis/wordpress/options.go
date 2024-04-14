package wordpress

import "fmt"

// Options represents options for configuring the WordPress client.
type Options struct {
	baseURL      string
	hasBasicAuth bool
	authUser     string
	authPassword string
}

// NewOptions creates a new Options instance.
func NewOptions() *Options {
	return &Options{}
}

// Validate checks if all required fields in Options are set.
func (o Options) Validate() error {
	if o.baseURL == "" {
		return fmt.Errorf("baseURL is required")
	}
	return nil
}

// WithBaseURL sets the base URL for the WordPress client.
func (o Options) WithBaseURL(baseURL string) Options {
	o.baseURL = baseURL
	return o
}

// WithBasicAuth sets the basic authentication credentials for the WordPress client.
func (o Options) WithBasicAuth(consumerKey, consumerSecret string) Options {
	o.authUser = consumerKey
	o.authPassword = consumerSecret
	o.hasBasicAuth = true
	return o
}

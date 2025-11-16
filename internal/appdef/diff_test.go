package appdef

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCompare(t *testing.T) {
	t.Parallel()

	tt := map[string]struct {
		current    *Definition
		previous   *Definition
		wantSkip   bool
		wantReason string
	}{
		"Identical definitions": {
			current:    makeTestDefinition("container", "digitalocean", makeEnv("KEY1", "value1")),
			previous:   makeTestDefinition("container", "digitalocean", makeEnv("KEY1", "value1")),
			wantSkip:   true,
			wantReason: "app.json unchanged",
		},
		"Domain changed": {
			current: &Definition{
				Apps: []App{
					{
						Name:  "web",
						Infra: Infra{Type: "container", Provider: ResourceProviderDigitalOcean},
						Domains: []Domain{
							{Name: "new.example.com", Type: DomainTypePrimary},
						},
					},
				},
			},
			previous: &Definition{
				Apps: []App{
					{
						Name:  "web",
						Infra: Infra{Type: "container", Provider: ResourceProviderDigitalOcean},
						Domains: []Domain{
							{Name: "old.example.com", Type: DomainTypePrimary},
						},
					},
				},
			},
			wantSkip:   false,
			wantReason: "Infrastructure config changed (domains/sizes/regions/resources/etc)",
		},
		"DO container env values unchanged": {
			current:    makeTestDefinition("container", "digitalocean", makeEnv("KEY1", "value1")),
			previous:   makeTestDefinition("container", "digitalocean", makeEnv("KEY1", "value1")),
			wantSkip:   true,
			wantReason: "app.json unchanged",
		},
		"DO container env values changed": {
			current:    makeTestDefinition("container", "digitalocean", makeEnv("KEY1", "value2")),
			previous:   makeTestDefinition("container", "digitalocean", makeEnv("KEY1", "value1")),
			wantSkip:   false,
			wantReason: "DigitalOcean container app env values changed",
		},
		"VM app env changed": {
			current:    makeTestDefinition("vm", "digitalocean", makeEnv("KEY1", "value2")),
			previous:   makeTestDefinition("vm", "digitalocean", makeEnv("KEY1", "value1")),
			wantSkip:   false,
			wantReason: "VM or non-DigitalOcean container app env changes detected",
		},
		"New env var added to DO container": {
			current:    makeTestDefinition("container", "digitalocean", makeEnv("KEY1", "value1", "KEY2", "value2")),
			previous:   makeTestDefinition("container", "digitalocean", makeEnv("KEY1", "value1")),
			wantSkip:   false,
			wantReason: "DigitalOcean container app env values changed",
		},
		"Resource added": {
			current: &Definition{
				Apps: []App{
					{Name: "web", Infra: Infra{Type: "container", Provider: ResourceProviderDigitalOcean}},
				},
				Resources: []Resource{
					{Name: "db", Type: ResourceTypePostgres, Provider: ResourceProviderDigitalOcean},
				},
			},
			previous: &Definition{
				Apps: []App{
					{Name: "web", Infra: Infra{Type: "container", Provider: ResourceProviderDigitalOcean}},
				},
				Resources: []Resource{},
			},
			wantSkip:   false,
			wantReason: "Infrastructure config changed (domains/sizes/regions/resources/etc)",
		},
	}

	for name, test := range tt {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := Compare(test.current, test.previous)

			assert.Equal(t, test.wantSkip, got.Skip)
			assert.Equal(t, test.wantReason, got.Reason)
		})
	}
}

// Helper functions for creating test data.

//nolint:unparam // provider parameter designed for reusability across different test scenarios.
func makeTestDefinition(infraType, provider string, env Environment) *Definition {
	providerType := ResourceProviderDigitalOcean
	if provider != "digitalocean" {
		providerType = ResourceProvider(provider)
	}

	return &Definition{
		Apps: []App{
			{
				Name:  "web",
				Infra: Infra{Type: infraType, Provider: providerType},
				Env:   env,
			},
		},
	}
}

func makeEnv(keyValuePairs ...string) Environment {
	env := Environment{
		Production: make(map[string]EnvValue),
	}

	for i := 0; i < len(keyValuePairs); i += 2 {
		key := keyValuePairs[i]
		value := keyValuePairs[i+1]
		env.Production[key] = EnvValue{
			Source: EnvSourceValue,
			Value:  value,
		}
	}

	return env
}

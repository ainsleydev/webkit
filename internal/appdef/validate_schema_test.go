package appdef

import (
	"encoding/json"
	"testing"

	"github.com/ainsleydev/webkit/pkg/util/ptr"
)

// NOTE: The kaptinlin/jsonschema validator currently only enforces a subset of JSON Schema constraints.
// Specifically, it validates required fields and types, but may not enforce pattern, enum, minimum/maximum,
// or format constraints. These constraints are still included in the schema for IDE autocomplete and validation.
func TestValidateAgainstSchema(t *testing.T) {
	tt := map[string]struct {
		json    string
		wantErr bool
	}{
		"valid minimal definition": {
			json: `{
				"webkit_version": "1.0.0",
				"project": {
					"name": "test-project",
					"title": "Test Project",
					"description": "A test project",
					"repo": {
						"owner": "testowner",
						"name": "testrepo"
					}
				},
				"apps": [{
					"name": "test-app",
					"title": "Test App",
					"type": "golang",
					"path": "./apps/test",
					"infra": {
						"provider": "digitalocean",
						"type": "vm"
					}
				}]
			}`,
			wantErr: false,
		},
		"missing required field webkit_version": {
			json: `{
				"project": {
					"name": "test-project",
					"title": "Test Project",
					"description": "A test project",
					"repo": {
						"owner": "testowner",
						"name": "testrepo"
					}
				},
				"apps": [{
					"name": "test-app",
					"title": "Test App",
					"type": "golang",
					"path": "./apps/test",
					"infra": {
						"provider": "digitalocean",
						"type": "vm"
					}
				}]
			}`,
			wantErr: true,
		},
		"missing required field project": {
			json: `{
				"webkit_version": "1.0.0",
				"apps": [{
					"name": "test-app",
					"title": "Test App",
					"type": "golang",
					"path": "./apps/test",
					"infra": {
						"provider": "digitalocean",
						"type": "vm"
					}
				}]
			}`,
			wantErr: true,
		},
		"missing required field apps": {
			json: `{
				"webkit_version": "1.0.0",
				"project": {
					"name": "test-project",
					"title": "Test Project",
					"description": "A test project",
					"repo": {
						"owner": "testowner",
						"name": "testrepo"
					}
				}
			}`,
			wantErr: true,
		},
		"valid with all fields": {
			json: `{
				"webkit_version": "1.0.0",
				"project": {
					"name": "test-project",
					"title": "Test Project",
					"description": "A comprehensive test project with all fields",
					"repo": {
						"owner": "testowner",
						"name": "testrepo"
					}
				},
				"notifications": {
					"webhook_url": "https://hooks.slack.com/services/XXX/YYY/ZZZ"
				},
				"shared": {
					"env": {
						"default": {
							"API_URL": {
								"source": "value",
								"value": "https://api.example.com"
							}
						}
					}
				},
				"resources": [{
					"name": "test-db",
					"type": "postgres",
					"provider": "digitalocean"
				}],
				"apps": [{
					"name": "test-app",
					"title": "Test App",
					"type": "golang",
					"description": "A test application",
					"path": "./apps/test",
					"build": {
						"dockerfile": "Dockerfile",
						"port": 8080
					},
					"infra": {
						"provider": "digitalocean",
						"type": "vm"
					},
					"domains": [{
						"name": "test.example.com",
						"type": "primary"
					}]
				}]
			}`,
			wantErr: false,
		},
	}

	for name, tc := range tt {
		t.Run(name, func(t *testing.T) {
			errs := ValidateAgainstSchema([]byte(tc.json))

			if tc.wantErr && errs == nil {
				t.Fatal("expected error, got nil")
			}
			if !tc.wantErr && errs != nil {
				t.Fatalf("expected no error, got: %v", errs)
			}
		})
	}
}

func TestValidateAgainstSchema_Integration(t *testing.T) {
	// Test with actual Definition struct to ensure schema generation matches
	def := &Definition{
		WebkitVersion: "1.0.0",
		Project: Project{
			Name:        "test-project",
			Title:       "Test Project",
			Description: "A test project",
			Repo: GitHubRepo{
				Owner: "testowner",
				Name:  "testrepo",
			},
		},
		Apps: []App{
			{
				Name:  "test-app",
				Title: "Test App",
				Type:  AppTypeGoLang,
				Path:  "./apps/test",
				Build: Build{
					Dockerfile: "Dockerfile",
					Port:       8080,
					Release:    ptr.BoolPtr(true),
				},
				Infra: Infra{
					Provider: ResourceProviderDigitalOcean,
					Type:     "vm",
				},
			},
		},
	}

	// Marshal to JSON
	data, err := json.Marshal(def)
	if err != nil {
		t.Fatalf("marshaling definition: %v", err)
	}

	// Validate against schema
	errs := ValidateAgainstSchema(data)
	if errs != nil {
		t.Fatalf("validation failed: %v", errs)
	}
}

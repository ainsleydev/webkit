package appdef

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ainsleydev/webkit/pkg/util/ptr"
)

// NOTE: The kaptinlin/jsonschema validator currently only enforces a subset of JSON Schema constraints.
// Specifically, it validates required fields and types, but may not enforce pattern, enum, minimum/maximum,
// or format constraints. These constraints are still included in the schema for IDE autocomplete and validation.
func TestValidateAgainstSchema2(t *testing.T) {
	t.Parallel()

	tt := map[string]struct {
		json    string
		wantErr bool
	}{
		"Valid minimal definition": {
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
		"Missing required field webkit_version": {
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
		"Missing required field project": {
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
		"Missing required field apps": {
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
		"Valid with all fields": {
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

	for name, test := range tt {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := ValidateAgainstSchema([]byte(test.json))
			fmt.Print(got)
			assert.Equal(t, test.wantErr, len(got) != 0)
		})
	}
}

func TestValidateAgainstSchema_Integration2(t *testing.T) {
	t.Parallel()

	// Test with actual Definition struct to ensure schema generation matches.
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

	data, err := json.Marshal(def)
	require.NoError(t, err)

	errs := ValidateAgainstSchema(data)
	assert.Len(t, errs, 0)
}

func TestValidateAgainstSchema_Constraints(t *testing.T) {
	t.Parallel()

	tt := map[string]struct {
		json    string
		wantErr bool
	}{
		"Valid all constraints met": {
			json: `{
				"webkit_version": "1.0.0",
				"project": {
					"name": "test-project",
					"title": "Test",
					"description": "Valid description",
					"repo": {"owner": "test", "name": "test"}
				},
				"apps": [{
					"name": "test-app",
					"title": "Test",
					"type": "golang",
					"path": "./test",
					"build": {"port": 8080},
					"infra": {"provider": "digitalocean", "type": "vm"}
				}]
			}`,
			wantErr: false,
		},
		"Invalid pattern uppercase in project name": {
			json: `{
				"webkit_version": "1.0.0",
				"project": {
					"name": "TEST-Project",
					"title": "Test",
					"description": "Invalid pattern",
					"repo": {"owner": "test", "name": "test"}
				},
				"apps": [{
					"name": "test-app",
					"title": "Test",
					"type": "golang",
					"path": "./test",
					"infra": {"provider": "digitalocean", "type": "vm"}
				}]
			}`,
			wantErr: true,
		},
		"Invalid pattern uppercase in app name": {
			json: `{
				"webkit_version": "1.0.0",
				"project": {
					"name": "test-project",
					"title": "Test",
					"description": "Valid",
					"repo": {"owner": "test", "name": "test"}
				},
				"apps": [{
					"name": "TEST-App",
					"title": "Test",
					"type": "golang",
					"path": "./test",
					"infra": {"provider": "digitalocean", "type": "vm"}
				}]
			}`,
			wantErr: true,
		},
		"Invalid pattern starts with number": {
			json: `{
				"webkit_version": "1.0.0",
				"project": {
					"name": "123-project",
					"title": "Test",
					"description": "Invalid pattern",
					"repo": {"owner": "test", "name": "test"}
				},
				"apps": [{
					"name": "test-app",
					"title": "Test",
					"type": "golang",
					"path": "./test",
					"infra": {"provider": "digitalocean", "type": "vm"}
				}]
			}`,
			wantErr: true,
		},
		"Invalid enum wrong app type": {
			json: `{
				"webkit_version": "1.0.0",
				"project": {
					"name": "test-project",
					"title": "Test",
					"description": "Invalid enum",
					"repo": {"owner": "test", "name": "test"}
				},
				"apps": [{
					"name": "test-app",
					"title": "Test",
					"type": "react",
					"path": "./test",
					"infra": {"provider": "digitalocean", "type": "vm"}
				}]
			}`,
			wantErr: true,
		},
		"Invalid enum wrong domain type": {
			json: `{
				"webkit_version": "1.0.0",
				"project": {
					"name": "test-project",
					"title": "Test",
					"description": "Valid",
					"repo": {"owner": "test", "name": "test"}
				},
				"apps": [{
					"name": "test-app",
					"title": "Test",
					"type": "golang",
					"path": "./test",
					"infra": {"provider": "digitalocean", "type": "vm"},
					"domains": [{"name": "example.com", "type": "invalid-type"}]
				}]
			}`,
			wantErr: true,
		},
		"Invalid enum wrong env source": {
			json: `{
				"webkit_version": "1.0.0",
				"project": {
					"name": "test-project",
					"title": "Test",
					"description": "Valid",
					"repo": {"owner": "test", "name": "test"}
				},
				"shared": {
					"env": {
						"default": {
							"API_URL": {
								"source": "invalid-source",
								"value": "test"
							}
						}
					}
				},
				"apps": [{
					"name": "test-app",
					"title": "Test",
					"type": "golang",
					"path": "./test",
					"infra": {"provider": "digitalocean", "type": "vm"}
				}]
			}`,
			wantErr: true,
		},
		"Invalid maximum port too high": {
			json: `{
				"webkit_version": "1.0.0",
				"project": {
					"name": "test-project",
					"title": "Test",
					"description": "Invalid max",
					"repo": {"owner": "test", "name": "test"}
				},
				"apps": [{
					"name": "test-app",
					"title": "Test",
					"type": "golang",
					"path": "./test",
					"build": {"port": 99999},
					"infra": {"provider": "digitalocean", "type": "vm"}
				}]
			}`,
			wantErr: true,
		},
		"Invalid minimum port zero": {
			json: `{
				"webkit_version": "1.0.0",
				"project": {
					"name": "test-project",
					"title": "Test",
					"description": "Invalid min",
					"repo": {"owner": "test", "name": "test"}
				},
				"apps": [{
					"name": "test-app",
					"title": "Test",
					"type": "golang",
					"path": "./test",
					"build": {"port": 0},
					"infra": {"provider": "digitalocean", "type": "vm"}
				}]
			}`,
			wantErr: true,
		},
		"Invalid minimum port negative": {
			json: `{
				"webkit_version": "1.0.0",
				"project": {
					"name": "test-project",
					"title": "Test",
					"description": "Invalid min",
					"repo": {"owner": "test", "name": "test"}
				},
				"apps": [{
					"name": "test-app",
					"title": "Test",
					"type": "golang",
					"path": "./test",
					"build": {"port": -1},
					"infra": {"provider": "digitalocean", "type": "vm"}
				}]
			}`,
			wantErr: true,
		},
		"Invalid maxLength description too long": {
			json: `{
				"webkit_version": "1.0.0",
				"project": {
					"name": "test-project",
					"title": "Test",
					"description": "` + generateString(201) + `",
					"repo": {"owner": "test", "name": "test"}
				},
				"apps": [{
					"name": "test-app",
					"title": "Test",
					"type": "golang",
					"path": "./test",
					"infra": {"provider": "digitalocean", "type": "vm"}
				}]
			}`,
			wantErr: true,
		},
		"Invalid maxLength app description too long": {
			json: `{
				"webkit_version": "1.0.0",
				"project": {
					"name": "test-project",
					"title": "Test",
					"description": "Valid",
					"repo": {"owner": "test", "name": "test"}
				},
				"apps": [{
					"name": "test-app",
					"title": "Test",
					"type": "golang",
					"description": "` + generateString(201) + `",
					"path": "./test",
					"infra": {"provider": "digitalocean", "type": "vm"}
				}]
			}`,
			wantErr: true,
		},
		"Invalid format webhook url not uri": {
			json: `{
				"webkit_version": "1.0.0",
				"project": {
					"name": "test-project",
					"title": "Test",
					"description": "Valid",
					"repo": {"owner": "test", "name": "test"}
				},
				"notifications": {
					"webhook_url": "not-a-valid-url"
				},
				"apps": [{
					"name": "test-app",
					"title": "Test",
					"type": "golang",
					"path": "./test",
					"infra": {"provider": "digitalocean", "type": "vm"}
				}]
			}`,
			wantErr: true,
		},
		"Valid format webhook url with https": {
			json: `{
				"webkit_version": "1.0.0",
				"project": {
					"name": "test-project",
					"title": "Test",
					"description": "Valid",
					"repo": {"owner": "test", "name": "test"}
				},
				"notifications": {
					"webhook_url": "https://hooks.slack.com/services/XXX/YYY/ZZZ"
				},
				"apps": [{
					"name": "test-app",
					"title": "Test",
					"type": "golang",
					"path": "./test",
					"infra": {"provider": "digitalocean", "type": "vm"}
				}]
			}`,
			wantErr: false,
		},
		"Valid boundary port minimum": {
			json: `{
				"webkit_version": "1.0.0",
				"project": {
					"name": "test-project",
					"title": "Test",
					"description": "Valid",
					"repo": {"owner": "test", "name": "test"}
				},
				"apps": [{
					"name": "test-app",
					"title": "Test",
					"type": "golang",
					"path": "./test",
					"build": {"port": 1},
					"infra": {"provider": "digitalocean", "type": "vm"}
				}]
			}`,
			wantErr: false,
		},
		"Valid boundary port maximum": {
			json: `{
				"webkit_version": "1.0.0",
				"project": {
					"name": "test-project",
					"title": "Test",
					"description": "Valid",
					"repo": {"owner": "test", "name": "test"}
				},
				"apps": [{
					"name": "test-app",
					"title": "Test",
					"type": "golang",
					"path": "./test",
					"build": {"port": 65535},
					"infra": {"provider": "digitalocean", "type": "vm"}
				}]
			}`,
			wantErr: false,
		},
		"Valid boundary description max length": {
			json: `{
				"webkit_version": "1.0.0",
				"project": {
					"name": "test-project",
					"title": "Test",
					"description": "` + generateString(200) + `",
					"repo": {"owner": "test", "name": "test"}
				},
				"apps": [{
					"name": "test-app",
					"title": "Test",
					"type": "golang",
					"path": "./test",
					"infra": {"provider": "digitalocean", "type": "vm"}
				}]
			}`,
			wantErr: false,
		},
	}

	for name, test := range tt {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := ValidateAgainstSchema([]byte(test.json))
			assert.Equal(t, test.wantErr, len(got) != 0, "Validation errors: %v", got)
		})
	}
}

func generateString(length int) string {
	result := ""
	for i := 0; i < length; i++ {
		result += "a"
	}
	return result
}

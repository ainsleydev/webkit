//go:build !race

package infra

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ainsleydev/webkit/internal/appdef"
	"github.com/ainsleydev/webkit/pkg/env"
)

//nolint:tparallel // Cannot use t.Parallel() due to t.Setenv() usage in setup
func TestTerraform_Resources(t *testing.T) {
	appDef := &appdef.Definition{
		Project: appdef.Project{
			Name:  "project",
			Title: "Project",
			Repo: appdef.GitHubRepo{
				Owner: "ainsley-dev",
				Name:  "project",
			},
		},
		Resources: []appdef.Resource{
			{
				Name:     "db",
				Title:    "Database",
				Type:     appdef.ResourceTypePostgres,
				Provider: appdef.ResourceProviderDigitalOcean,
				Config: map[string]any{
					"pg_version": "18",
					"size":       "db-s-1vcpu-2gb",
					"region":     "ams3",
					"node_count": 2,
				},
				Backup: appdef.ResourceBackupConfig{
					Enabled: true,
				},
			},
			{
				Name:     "cache",
				Title:    "Cache Database",
				Type:     appdef.ResourceTypeSQLite,
				Provider: appdef.ResourceProviderTurso,
				Config: map[string]any{
					"organisation": "test-org",
					"group":        "default",
				},
				Backup: appdef.ResourceBackupConfig{
					Enabled: true,
				},
			},
			{
				Name:     "storage",
				Title:    "Storage",
				Type:     appdef.ResourceTypeS3,
				Provider: appdef.ResourceProviderDigitalOcean,
				Config: map[string]any{
					"acl":    "private",
					"region": "ams3",
				},
			},
			{
				Name:     "backups",
				Title:    "Backups",
				Type:     appdef.ResourceTypeS3,
				Provider: appdef.ResourceProviderBackBlaze,
				Config: map[string]any{
					"acl": "allPrivate",
				},
			},
		},
	}

	tf, teardown := setup(t, appDef)
	t.Cleanup(teardown)

	err := tf.Init(t.Context())
	require.NoError(t, err)

	got, err := tf.Plan(t.Context(), env.Production, false)
	require.NoError(t, err)
	require.NotNil(t, got)
	require.True(t, got.HasChanges, "Plan should have changes")

	t.Run("Digital Ocean Postgres", func(t *testing.T) {
		t.Parallel()

		t.Log("Database Cluster Configuration")
		{
			var dbCluster map[string]any
			for _, rc := range got.Plan.ResourceChanges {
				if rc.Type == "digitalocean_database_cluster" && rc.Name == "this" {
					dbCluster = rc.Change.After.(map[string]any)
					break
				}
			}
			require.NotNil(t, dbCluster, "Database cluster resource should be planned")

			assert.Equal(t, "pg", dbCluster["engine"])
			assert.Equal(t, "project-db", dbCluster["name"])
			assert.Equal(t, float64(2), dbCluster["node_count"])
			assert.Equal(t, "ams3", dbCluster["region"])
			assert.Equal(t, "db-s-1vcpu-2gb", dbCluster["size"])
			assert.Equal(t, "18", dbCluster["version"])

			tags := dbCluster["tags"].([]any)
			assert.Contains(t, tags, "production")
			assert.Contains(t, tags, "project")
			assert.Contains(t, tags, "terraform")
		}

		t.Log("Database and User")
		{
			var dbDatabase map[string]any
			var dbUser map[string]any

			for _, rc := range got.Plan.ResourceChanges {
				if rc.Type == "digitalocean_database_db" && rc.Name == "this" {
					dbDatabase = rc.Change.After.(map[string]any)
				}
				if rc.Type == "digitalocean_database_user" && rc.Name == "this" {
					dbUser = rc.Change.After.(map[string]any)
				}
			}

			assert.NotNil(t, dbDatabase, "Database should be planned")
			assert.Equal(t, "project_db", dbDatabase["name"])

			assert.NotNil(t, dbUser, "Database user should be planned")
			assert.Equal(t, "project_db_admin", dbUser["name"])
		}

		t.Log("Connection Pool")
		{
			var connPool map[string]any
			for _, rc := range got.Plan.ResourceChanges {
				if rc.Type == "digitalocean_database_connection_pool" && rc.Name == "this" {
					connPool = rc.Change.After.(map[string]any)
					break
				}
			}

			assert.NotNil(t, connPool, "Connection pool should be planned")
			assert.Equal(t, "project_db", connPool["db_name"])
			assert.Equal(t, "transaction", connPool["mode"])
			assert.Equal(t, "project_db_pool", connPool["name"])
			assert.Equal(t, float64(20), connPool["size"])
			assert.Equal(t, "project_db_admin", connPool["user"])
		}

		t.Log("Postgres GitHub Secrets")
		{
			var secretNames []string
			for _, rc := range got.Plan.ResourceChanges {
				if rc.Type == "github_actions_secret" {
					after := rc.Change.After.(map[string]any)
					secretName := after["secret_name"].(string)
					secretNames = append(secretNames, secretName)
				}
			}

			assert.Contains(t, secretNames, "TF_PROD_DB_CONNECTION_URL")
			assert.Contains(t, secretNames, "TF_PROD_DB_ID")
			assert.Contains(t, secretNames, "TF_PROD_DB_URN")
		}
	})

	t.Run("Turso SQLite", func(t *testing.T) {
		t.Parallel()

		t.Log("Turso Database Configuration")
		{
			var tursoDb map[string]any
			for _, rc := range got.Plan.ResourceChanges {
				if rc.Type == "turso_database" && rc.Name == "this" {
					tursoDb = rc.Change.After.(map[string]any)
					break
				}
			}
			require.NotNil(t, tursoDb, "Turso database resource should be planned")

			assert.Equal(t, "default", tursoDb["group"])
			assert.Equal(t, "project-cache", tursoDb["name"])
			assert.Equal(t, "test-org", tursoDb["organization_name"])
		}

		t.Log("Turso Database Token")
		{
			var tursoToken map[string]any
			for _, rc := range got.Plan.ResourceChanges {
				if rc.Type == "turso_database_token" && rc.Name == "this" {
					tursoToken = rc.Change.After.(map[string]any)
					break
				}
			}
			assert.NotNil(t, tursoToken, "Turso token should be planned")
		}

		t.Log("Turso GitHub Secrets")
		{
			var secretNames []string
			for _, rc := range got.Plan.ResourceChanges {
				if rc.Type == "github_actions_secret" {
					after := rc.Change.After.(map[string]any)
					secretName := after["secret_name"].(string)
					secretNames = append(secretNames, secretName)
				}
			}

			assert.Contains(t, secretNames, "TF_PROD_CACHE_CONNECTION_URL")
			assert.Contains(t, secretNames, "TF_PROD_CACHE_AUTH_TOKEN")
			assert.Contains(t, secretNames, "TF_PROD_CACHE_HOST")
			assert.Contains(t, secretNames, "TF_PROD_CACHE_DATABASE")
			assert.Contains(t, secretNames, "TF_PROD_CACHE_ID")
		}
	})

	t.Run("Digital Ocean Spaces", func(t *testing.T) {
		t.Parallel()

		t.Log("Spaces Bucket Configuration")
		{
			var bucket map[string]any
			for _, rc := range got.Plan.ResourceChanges {
				if rc.Type == "digitalocean_spaces_bucket" && rc.Name == "this" {
					bucket = rc.Change.After.(map[string]any)
					break
				}
			}
			require.NotNil(t, bucket, "Spaces bucket resource should be planned")

			assert.Equal(t, "private", bucket["acl"])
			assert.Equal(t, "project-storage", bucket["name"])
			assert.Equal(t, "ams3", bucket["region"])
			assert.Equal(t, false, bucket["force_destroy"])
		}

		t.Log("CDN Configuration")
		{
			var cdn map[string]any
			for _, rc := range got.Plan.ResourceChanges {
				if rc.Type == "digitalocean_cdn" && rc.Name == "this" {
					cdn = rc.Change.After.(map[string]any)
					break
				}
			}
			require.NotNil(t, cdn, "CDN resource should be planned")
			assert.Nil(t, cdn["custom_domain"], "No custom domain should be set")
		}

		t.Log("CORS Configuration")
		{
			var cors map[string]any
			for _, rc := range got.Plan.ResourceChanges {
				if rc.Type == "digitalocean_spaces_bucket_cors_configuration" && rc.Name == "this" {
					cors = rc.Change.After.(map[string]any)
					break
				}
			}
			require.NotNil(t, cors, "CORS configuration should be planned")

			assert.Equal(t, "ams3", cors["region"])

			corsRules := cors["cors_rule"].([]any)
			assert.Len(t, corsRules, 1)

			rule := corsRules[0].(map[string]any)
			assert.Contains(t, rule["allowed_methods"].([]any), "GET")
			assert.Contains(t, rule["allowed_origins"].([]any), "*")
			assert.Contains(t, rule["allowed_headers"].([]any), "*")
			assert.Equal(t, float64(31536000), rule["max_age_seconds"])
		}

		t.Log("Spaces GitHub Secrets")
		{
			var secretNames []string
			for _, rc := range got.Plan.ResourceChanges {
				if rc.Type == "github_actions_secret" {
					after := rc.Change.After.(map[string]any)
					secretName := after["secret_name"].(string)
					secretNames = append(secretNames, secretName)
				}
			}

			assert.Contains(t, secretNames, "TF_PROD_STORAGE_BUCKET_NAME")
			assert.Contains(t, secretNames, "TF_PROD_STORAGE_BUCKET_URL")
			assert.Contains(t, secretNames, "TF_PROD_STORAGE_ENDPOINT")
			assert.Contains(t, secretNames, "TF_PROD_STORAGE_ID")
			assert.Contains(t, secretNames, "TF_PROD_STORAGE_REGION")
			assert.Contains(t, secretNames, "TF_PROD_STORAGE_URN")
		}
	})

	t.Run("Backblaze B2 Bucket", func(t *testing.T) {
		t.Parallel()

		t.Log("B2 Bucket Configuration")
		{
			var b2Buckets []map[string]any
			for _, rc := range got.Plan.ResourceChanges {
				if rc.Type == "b2_bucket" && rc.Name == "this" {
					b2Buckets = append(b2Buckets, rc.Change.After.(map[string]any))
				}
			}

			var backupsBucket map[string]any
			for _, bucket := range b2Buckets {
				if bucket["bucket_name"] == "backups" {
					backupsBucket = bucket
					break
				}
			}
			require.NotNil(t, backupsBucket, "B2 backups bucket should be planned")

			assert.Equal(t, "backups", backupsBucket["bucket_name"])
			assert.Equal(t, "allPrivate", backupsBucket["bucket_type"])
		}

		t.Log("B2 GitHub Secrets")
		{
			var secretNames []string
			for _, rc := range got.Plan.ResourceChanges {
				if rc.Type == "github_actions_secret" {
					after := rc.Change.After.(map[string]any)
					secretName := after["secret_name"].(string)
					secretNames = append(secretNames, secretName)
				}
			}

			assert.Contains(t, secretNames, "TF_PROD_BACKUPS_ID")
			assert.Contains(t, secretNames, "TF_PROD_BACKUPS_BUCKET_NAME")
		}
	})

	t.Run("Resource Outputs", func(t *testing.T) {
		t.Parallel()

		resources := got.Plan.PlannedValues.Outputs["resources"]
		require.NotNil(t, resources)
		assert.True(t, resources.Sensitive, "Resources output should be sensitive")

		resourceNames := got.Plan.PlannedValues.Outputs["resource_names"]
		require.NotNil(t, resourceNames)
		names := resourceNames.Value.([]any)
		assert.Len(t, names, 4)
		assert.Contains(t, names, "db")
		assert.Contains(t, names, "cache")
		assert.Contains(t, names, "storage")
		assert.Contains(t, names, "backups")
	})
}

//nolint:tparallel // Cannot use t.Parallel() due to t.Setenv() usage in setup
func TestTerraform_Apps(t *testing.T) {
	appDef := &appdef.Definition{
		Project: appdef.Project{
			Name:  "project",
			Title: "Project",
			Repo: appdef.GitHubRepo{
				Owner: "ainsley-dev",
				Name:  "project",
			},
		},
		Apps: []appdef.App{
			{
				Name: "api",
				Type: appdef.AppTypeSvelteKit,
				Infra: appdef.Infra{
					Provider: appdef.ResourceProviderDigitalOcean,
					Type:     "container",
					Config: map[string]any{
						"region":            "lon",
						"size":              "apps-s-1vcpu-1gb",
						"instance_count":    1,
						"port":              3000,
						"health_check_path": "/health",
					},
				},
			},
			{
				Name: "cms",
				Type: appdef.AppTypePayload,
				Infra: appdef.Infra{
					Provider: appdef.ResourceProviderDigitalOcean,
					Type:     "vm",
					Config: map[string]any{
						"size":   "s-1vcpu-1gb",
						"region": "lon1",
					},
				},
			},
			{
				Name: "worker",
				Type: appdef.AppTypeGoLang,
				Infra: appdef.Infra{
					Provider: appdef.ResourceProviderHetzner,
					Type:     "vm",
					Config: map[string]any{
						"size":   "cx11",
						"region": "nbg1",
					},
				},
			},
		},
	}

	tf, teardown := setup(t, appDef)
	t.Cleanup(teardown)

	err := tf.Init(t.Context())
	require.NoError(t, err)

	got, err := tf.Plan(t.Context(), env.Production, false)
	require.NoError(t, err)
	require.NotNil(t, got)
	require.True(t, got.HasChanges, "Plan should have changes")

	t.Run("Digital Ocean App Platform", func(t *testing.T) {
		t.Parallel()

		t.Log("App Configuration")
		{
			var app map[string]any
			for _, rc := range got.Plan.ResourceChanges {
				if rc.Type == "digitalocean_app" && rc.Name == "this" {
					app = rc.Change.After.(map[string]any)
					break
				}
			}
			require.NotNil(t, app, "App Platform app should be planned")

			spec := app["spec"].([]any)[0].(map[string]any)
			assert.Equal(t, "project-api", spec["name"])
			assert.Equal(t, "lon", spec["region"])

			// Note: Spec-level deployment alerts (DEPLOYMENT_FAILED, DEPLOYMENT_LIVE)
			// are not currently implemented in the terraform module
		}

		t.Log("Service Configuration")
		{
			var app map[string]any
			for _, rc := range got.Plan.ResourceChanges {
				if rc.Type == "digitalocean_app" && rc.Name == "this" {
					app = rc.Change.After.(map[string]any)
					break
				}
			}

			spec := app["spec"].([]any)[0].(map[string]any)
			services := spec["service"].([]any)
			require.Len(t, services, 1)

			service := services[0].(map[string]any)
			assert.Equal(t, "svelte-kit", service["name"]) // Service name is app type, not app name
			assert.Equal(t, float64(3000), service["http_port"])
			assert.Equal(t, float64(1), service["instance_count"])
			assert.Equal(t, "apps-s-1vcpu-1gb", service["instance_size_slug"])

			healthChecks := service["health_check"].([]any)
			require.Len(t, healthChecks, 1)
			healthCheck := healthChecks[0].(map[string]any)
			assert.Equal(t, "/health", healthCheck["http_path"])
			assert.Equal(t, float64(10), healthCheck["failure_threshold"])

			serviceAlerts := service["alert"].([]any)
			assert.Len(t, serviceAlerts, 3)

			serviceAlertRules := make(map[string]map[string]any)
			for _, alert := range serviceAlerts {
				alertMap := alert.(map[string]any)
				rule := alertMap["rule"].(string)
				serviceAlertRules[rule] = alertMap
			}

			assert.Contains(t, serviceAlertRules, "CPU_UTILIZATION")
			assert.Equal(t, float64(80), serviceAlertRules["CPU_UTILIZATION"]["value"])

			assert.Contains(t, serviceAlertRules, "MEM_UTILIZATION")
			assert.Equal(t, float64(80), serviceAlertRules["MEM_UTILIZATION"]["value"])

			assert.Contains(t, serviceAlertRules, "RESTART_COUNT")
			assert.Equal(t, float64(3), serviceAlertRules["RESTART_COUNT"]["value"])
		}
	})

	t.Run("Digital Ocean VM", func(t *testing.T) {
		t.Parallel()

		t.Log("Droplet Configuration")
		{
			var droplet map[string]any
			for _, rc := range got.Plan.ResourceChanges {
				if rc.Type == "digitalocean_droplet" && rc.Name == "this" {
					droplet = rc.Change.After.(map[string]any)
					break
				}
			}
			require.NotNil(t, droplet, "Droplet should be planned")

			assert.Equal(t, "project-cms", droplet["name"])
			assert.Equal(t, "lon1", droplet["region"])
			assert.Equal(t, "s-1vcpu-1gb", droplet["size"])

			tags := droplet["tags"].([]any)
			assert.Contains(t, tags, "production")
			assert.Contains(t, tags, "project")
		}

		t.Log("DO VM GitHub Secrets")
		{
			var secretNames []string
			for _, rc := range got.Plan.ResourceChanges {
				if rc.Type == "github_actions_secret" {
					after := rc.Change.After.(map[string]any)
					secretName := after["secret_name"].(string)
					secretNames = append(secretNames, secretName)
				}
			}

			assert.Contains(t, secretNames, "TF_PROD_CMS_IP_ADDRESS")
			assert.Contains(t, secretNames, "TF_PROD_CMS_SSH_PRIVATE_KEY")
			assert.Contains(t, secretNames, "TF_PROD_CMS_SERVER_USER")
		}
	})

	t.Run("Hetzner VM", func(t *testing.T) {
		t.Parallel()

		t.Log("Server Configuration")
		{
			var server map[string]any
			for _, rc := range got.Plan.ResourceChanges {
				if rc.Type == "hcloud_server" && rc.Name == "this" {
					server = rc.Change.After.(map[string]any)
					break
				}
			}
			require.NotNil(t, server, "Hetzner server should be planned")

			assert.Equal(t, "project-worker", server["name"])
			assert.Equal(t, "nbg1", server["location"])
			assert.Equal(t, "cx11", server["server_type"])

			labels := server["labels"].(map[string]any)
			// Hetzner labels are created from tags where each tag becomes a key with value "true"
			// Default tags are: project name, environment (e.g., "production"), and "terraform"
			assert.Equal(t, "true", labels["production"])
			assert.Equal(t, "true", labels["project"])
			assert.Equal(t, "true", labels["terraform"])
		}

		t.Log("Hetzner VM GitHub Secrets")
		{
			var secretNames []string
			for _, rc := range got.Plan.ResourceChanges {
				if rc.Type == "github_actions_secret" {
					after := rc.Change.After.(map[string]any)
					secretName := after["secret_name"].(string)
					secretNames = append(secretNames, secretName)
				}
			}

			assert.Contains(t, secretNames, "TF_PROD_WORKER_IP_ADDRESS")
			assert.Contains(t, secretNames, "TF_PROD_WORKER_SSH_PRIVATE_KEY")
			assert.Contains(t, secretNames, "TF_PROD_WORKER_SERVER_USER")
		}
	})

	t.Run("App Outputs", func(t *testing.T) {
		t.Parallel()

		appNames := got.Plan.PlannedValues.Outputs["app_names"]
		require.NotNil(t, appNames)
		names := appNames.Value.([]any)
		assert.Len(t, names, 3)
		assert.Contains(t, names, "api")
		assert.Contains(t, names, "cms")
		assert.Contains(t, names, "worker")
	})
}

//nolint:tparallel // Cannot use t.Parallel() due to t.Setenv() usage in setup
func TestTerraform_Monitoring(t *testing.T) {
	enabled := true
	appDef := &appdef.Definition{
		Project: appdef.Project{
			Name:  "project",
			Title: "Project",
			Repo: appdef.GitHubRepo{
				Owner: "ainsley-dev",
				Name:  "project",
			},
		},
		Monitoring: appdef.Monitoring{
			Enabled: &enabled,
			StatusPage: appdef.StatusPage{
				Theme: "auto",
				Slug:  "project-status",
			},
			Custom: []appdef.Monitor{
				{
					Name:     "api-health",
					Type:     appdef.MonitorTypeHTTP,
					Interval: 60,
					Config: map[string]any{
						"url": "https://api.example.com/health",
					},
				},
				{
					Name:     "keyword-check",
					Type:     appdef.MonitorTypeHTTPKeyword,
					Interval: 120,
					Config: map[string]any{
						"url":     "https://example.com",
						"keyword": "Welcome",
					},
				},
			},
		},
		Apps: []appdef.App{
			{
				Name:       "web",
				Type:       appdef.AppTypeSvelteKit,
				Monitoring: true,
				Infra: appdef.Infra{
					Provider: appdef.ResourceProviderDigitalOcean,
					Type:     "container",
					Config: map[string]any{
						"region": "lon",
					},
				},
				Domains: []appdef.Domain{
					{
						Name: "example.com",
						Type: appdef.DomainTypePrimary,
					},
				},
			},
		},
		Resources: []appdef.Resource{
			{
				Name:     "db",
				Type:     appdef.ResourceTypePostgres,
				Provider: appdef.ResourceProviderDigitalOcean,
				Backup: appdef.ResourceBackupConfig{
					Enabled: true,
				},
			},
		},
	}

	tf, teardown := setup(t, appDef)
	t.Cleanup(teardown)

	err := tf.Init(t.Context())
	require.NoError(t, err)

	got, err := tf.Plan(t.Context(), env.Production, false)
	require.NoError(t, err)
	require.NotNil(t, got)
	require.True(t, got.HasChanges, "Plan should have changes")

	t.Run("Peekaping Project Tag", func(t *testing.T) {
		t.Parallel()

		var projectTag map[string]any
		for _, rc := range got.Plan.ResourceChanges {
			if rc.Type == "peekaping_project_tag" && rc.Name == "this" {
				projectTag = rc.Change.After.(map[string]any)
				break
			}
		}
		require.NotNil(t, projectTag, "Peekaping project tag should be planned")
		assert.Equal(t, "project", projectTag["name"])
	})

	t.Run("HTTP Monitors", func(t *testing.T) {
		t.Parallel()

		var httpMonitors []map[string]any
		for _, rc := range got.Plan.ResourceChanges {
			if rc.Type == "peekaping_http_monitor" {
				httpMonitors = append(httpMonitors, rc.Change.After.(map[string]any))
			}
		}

		assert.GreaterOrEqual(t, len(httpMonitors), 2, "Should have at least 2 HTTP monitors (custom + app)")

		var customMonitor map[string]any
		for _, mon := range httpMonitors {
			if mon["name"] == "project-api-health" {
				customMonitor = mon
				break
			}
		}
		require.NotNil(t, customMonitor, "Custom HTTP monitor should be planned")
		assert.Equal(t, "https://api.example.com/health", customMonitor["url"])
		assert.Equal(t, float64(60), customMonitor["interval"])
	})

	t.Run("HTTP Keyword Monitors", func(t *testing.T) {
		t.Parallel()

		var keywordMonitors []map[string]any
		for _, rc := range got.Plan.ResourceChanges {
			if rc.Type == "peekaping_http_keyword_monitor" {
				keywordMonitors = append(keywordMonitors, rc.Change.After.(map[string]any))
			}
		}

		assert.GreaterOrEqual(t, len(keywordMonitors), 1, "Should have at least 1 keyword monitor")

		var customMonitor map[string]any
		for _, mon := range keywordMonitors {
			if mon["name"] == "project-keyword-check" {
				customMonitor = mon
				break
			}
		}
		require.NotNil(t, customMonitor, "Custom keyword monitor should be planned")
		assert.Equal(t, "https://example.com", customMonitor["url"])
		assert.Equal(t, "Welcome", customMonitor["keyword"])
		assert.Equal(t, float64(120), customMonitor["interval"])
	})

	t.Run("DNS Monitors", func(t *testing.T) {
		t.Parallel()

		var dnsMonitors []map[string]any
		for _, rc := range got.Plan.ResourceChanges {
			if rc.Type == "peekaping_dns_monitor" {
				dnsMonitors = append(dnsMonitors, rc.Change.After.(map[string]any))
			}
		}

		assert.GreaterOrEqual(t, len(dnsMonitors), 1, "Should have at least 1 DNS monitor for app domain")
	})

	t.Run("Push Monitors", func(t *testing.T) {
		t.Parallel()

		var pushMonitors []map[string]any
		for _, rc := range got.Plan.ResourceChanges {
			if rc.Type == "peekaping_push_monitor" {
				pushMonitors = append(pushMonitors, rc.Change.After.(map[string]any))
			}
		}

		assert.GreaterOrEqual(t, len(pushMonitors), 2, "Should have push monitors for codebase and db backup")

		var dbBackupMonitor map[string]any
		for _, mon := range pushMonitors {
			if mon["name"] == "project-db-backup" {
				dbBackupMonitor = mon
				break
			}
		}
		require.NotNil(t, dbBackupMonitor, "DB backup push monitor should be planned")
		assert.Equal(t, float64(90000), dbBackupMonitor["interval"])
	})

	t.Run("Status Page", func(t *testing.T) {
		t.Parallel()

		var statusPage map[string]any
		for _, rc := range got.Plan.ResourceChanges {
			if rc.Type == "peekaping_status_page" && rc.Name == "this" {
				statusPage = rc.Change.After.(map[string]any)
				break
			}
		}
		require.NotNil(t, statusPage, "Status page should be planned")
		assert.Equal(t, "project-status", statusPage["slug"])
		assert.Equal(t, "auto", statusPage["theme"])
	})

	t.Run("GitHub Variables for Push Monitors", func(t *testing.T) {
		t.Parallel()

		var variableNames []string
		for _, rc := range got.Plan.ResourceChanges {
			if rc.Type == "github_actions_variable" {
				after := rc.Change.After.(map[string]any)
				variableName := after["variable_name"].(string)
				variableNames = append(variableNames, variableName)
			}
		}

		assert.Contains(t, variableNames, "PROD_CODEBASE_BACKUP_PING_URL")
		assert.Contains(t, variableNames, "PROD_DB_BACKUP_PING_URL")
	})
}

//nolint:tparallel // Cannot use t.Parallel() due to t.Setenv() usage in setup
func TestTerraform_Defaults(t *testing.T) {
	appDef := &appdef.Definition{
		Project: appdef.Project{
			Name:  "project",
			Title: "Project",
			Repo: appdef.GitHubRepo{
				Owner: "ainsley-dev",
				Name:  "project",
			},
		},
		Resources: []appdef.Resource{},
		Apps:      []appdef.App{},
	}

	tf, teardown := setup(t, appDef)
	t.Cleanup(teardown)

	err := tf.Init(t.Context())
	require.NoError(t, err)

	got, err := tf.Plan(t.Context(), env.Production, false)
	require.NoError(t, err)
	require.NotNil(t, got)

	t.Run("Default B2 Bucket", func(t *testing.T) {
		t.Parallel()

		t.Log("B2 Bucket Configuration")
		{
			var b2Bucket map[string]any
			for _, rc := range got.Plan.ResourceChanges {
				if rc.Type == "b2_bucket" && rc.Name == "this" {
					b2Bucket = rc.Change.After.(map[string]any)
					break
				}
			}
			require.NotNil(t, b2Bucket, "B2 bucket resource should be planned")

			assert.Equal(t, "project", b2Bucket["bucket_name"])
			assert.Equal(t, "allPrivate", b2Bucket["bucket_type"])

			lifecycleRules := b2Bucket["lifecycle_rules"].([]any)
			require.Len(t, lifecycleRules, 1, "Should have exactly one lifecycle rule")

			rule := lifecycleRules[0].(map[string]any)
			assert.Equal(t, float64(1), rule["days_from_hiding_to_deleting"])
			assert.Equal(t, float64(0), rule["days_from_uploading_to_hiding"])
			assert.Equal(t, "", rule["file_name_prefix"])
		}

		t.Log("B2 Bucket Output")
		{
			defaultBucket := got.Plan.PlannedValues.Outputs["default_b2_bucket"]
			require.NotNil(t, defaultBucket, "Default B2 bucket output should exist")
			assert.True(t, defaultBucket.Sensitive, "Default B2 bucket output should be sensitive")
		}
	})

	t.Run("Slack Channel", func(t *testing.T) {
		t.Parallel()

		t.Log("Slack Channel Configuration")
		{
			var slackChannel map[string]any
			for _, rc := range got.Plan.ResourceChanges {
				if rc.Type == "slack_conversation" {
					slackChannel = rc.Change.After.(map[string]any)
					break
				}
			}
			require.NotNil(t, slackChannel, "Slack channel should be planned")

			assert.Equal(t, "alerts-project", slackChannel["name"])
			assert.Equal(t, false, slackChannel["is_private"])
			assert.Equal(t, "archive", slackChannel["action_on_destroy"])
		}

		t.Log("Slack Channel GitHub Secret")
		{
			var secretNames []string
			for _, rc := range got.Plan.ResourceChanges {
				if rc.Type == "github_actions_secret" {
					after := rc.Change.After.(map[string]any)
					secretName := after["secret_name"].(string)
					secretNames = append(secretNames, secretName)
				}
			}

			assert.Contains(t, secretNames, "TF_SLACK_CHANNEL_ID")
		}
	})

	t.Run("DigitalOcean Project", func(t *testing.T) {
		t.Parallel()

		t.Log("Project Configuration")
		{
			var doProject map[string]any
			for _, rc := range got.Plan.ResourceChanges {
				if rc.Type == "digitalocean_project" && rc.Name == "this" {
					doProject = rc.Change.After.(map[string]any)
					break
				}
			}
			require.NotNil(t, doProject, "DigitalOcean project should be planned")

			assert.Equal(t, "Project", doProject["name"])
			assert.Equal(t, "Production", doProject["environment"])
			assert.Equal(t, "Web Application", doProject["purpose"])
		}

		t.Log("Project Output")
		{
			projectID := got.Plan.PlannedValues.Outputs["digitalocean_project_id"]
			require.NotNil(t, projectID, "DigitalOcean project ID output should exist")
		}
	})

	t.Run("GitHub Secrets Count", func(t *testing.T) {
		t.Parallel()

		secrets := got.Plan.PlannedValues.Outputs["github_secrets_created"]
		require.NotNil(t, secrets)

		secretsCount := got.Plan.PlannedValues.Outputs["github_secrets_count"]
		require.NotNil(t, secretsCount)

		assert.Equal(t, float64(1), secretsCount.Value.(float64), "Should have 1 secret (Slack channel)")
	})
}

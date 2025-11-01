//go:build !race

package infra

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ainsleydev/webkit/internal/appdef"
	"github.com/ainsleydev/webkit/pkg/env"
)

func TestTerraform_Resources(t *testing.T) {
	t.Run("Digital Ocean - Postgres - Basic", func(t *testing.T) {
		appDef := &appdef.Definition{
			Project: appdef.Project{
				Name: "project",
				Repo: appdef.GitHubRepo{
					Owner: "ainsley-dev",
					Name:  "project",
				},
			},
			Resources: []appdef.Resource{
				{
					Name:     "db",
					Type:     appdef.ResourceTypePostgres,
					Provider: appdef.ResourceProviderDigitalOcean,
					Config: map[string]any{
						"pg_version": "18",
						"size":       "db-s-1vcpu-2gb",
						"region":     "ams3",
						"node_count": 2,
					},
					Backup: appdef.ResourceBackupConfig{},
				},
			},
		}

		tf, teardown := setup(t, appDef)
		defer teardown()

		err := tf.Init(t.Context())
		require.NoError(t, err)

		got, err := tf.Plan(t.Context(), env.Production)
		require.NoError(t, err)
		require.NotNil(t, got)
		require.True(t, got.HasChanges, "Plan should have changes")

		t.Log("Plan Summary")
		{
			require.Equal(t, 8, len(got.Plan.ResourceChanges), "Should plan to create 8 resources")
		}

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

		t.Log("Resource Outputs")
		{
			resources := got.Plan.PlannedValues.Outputs["resources"]
			require.NotNil(t, resources)
			assert.True(t, resources.Sensitive, "Resources output should be sensitive")

			resourceNames := got.Plan.PlannedValues.Outputs["resource_names"]
			require.NotNil(t, resourceNames)
			names := resourceNames.Value.([]any)
			assert.Len(t, names, 1)
			assert.Equal(t, "db", names[0])
		}

		t.Log("Permissions Grant")
		{
			var grantResource map[string]any
			for _, rc := range got.Plan.ResourceChanges {
				if rc.Type == "null_resource" && rc.Name == "grant_permissions" {
					grantResource = rc.Change.After.(map[string]any)
					break
				}
			}
			assert.NotNil(t, grantResource, "Permission grant resource should be planned")
		}

		t.Log("GitHub Secrets")
		{
			secrets := got.Plan.PlannedValues.Outputs["github_secrets_created"]
			require.NotNil(t, secrets)
			require.Equal(t, 3, int(got.Plan.PlannedValues.Outputs["github_secrets_count"].Value.(float64)))

			secretNames := secrets.Value.([]any)
			require.Len(t, secretNames, 3)
			assert.Contains(t, secretNames, "TF_PROD_DB_CONNECTION_URL")
			assert.Contains(t, secretNames, "TF_PROD_DB_ID")
			assert.Contains(t, secretNames, "TF_PROD_DB_URN")
		}
	})

	t.Run("Digital Ocean - Postgres - Version 17", func(t *testing.T) {
		appDef := &appdef.Definition{
			Project: appdef.Project{
				Name: "version-test",
				Repo: appdef.GitHubRepo{
					Owner: "ainsley-dev",
					Name:  "version-test",
				},
			},
			Resources: []appdef.Resource{
				{
					Name:     "db-v17",
					Type:     appdef.ResourceTypePostgres,
					Provider: appdef.ResourceProviderDigitalOcean,
					Config: map[string]any{
						"pg_version": "17",
						"size":       "db-s-1vcpu-1gb",
						"region":     "lon1",
						"node_count": 1,
					},
				},
			},
		}

		tf, teardown := setup(t, appDef)
		defer teardown()

		err := tf.Init(t.Context())
		require.NoError(t, err)

		got, err := tf.Plan(t.Context(), env.Production)
		require.NoError(t, err)
		require.NotNil(t, got)

		t.Log("Database Cluster with Version 17")
		{
			var dbCluster map[string]any
			for _, rc := range got.Plan.ResourceChanges {
				if rc.Type == "digitalocean_database_cluster" && rc.Name == "this" {
					dbCluster = rc.Change.After.(map[string]any)
					break
				}
			}
			require.NotNil(t, dbCluster, "Database cluster resource should be planned")

			assert.Equal(t, "17", dbCluster["version"])
			assert.Equal(t, "lon1", dbCluster["region"])
			assert.Equal(t, float64(1), dbCluster["node_count"])
			assert.Equal(t, "db-s-1vcpu-1gb", dbCluster["size"])
		}
	})

	t.Run("Digital Ocean - Postgres - With Allowed IPs", func(t *testing.T) {
		appDef := &appdef.Definition{
			Project: appdef.Project{
				Name: "secure-db",
				Repo: appdef.GitHubRepo{
					Owner: "ainsley-dev",
					Name:  "secure-db",
				},
			},
			Resources: []appdef.Resource{
				{
					Name:     "secure-db",
					Type:     appdef.ResourceTypePostgres,
					Provider: appdef.ResourceProviderDigitalOcean,
					Config: map[string]any{
						"pg_version":       "18",
						"size":             "db-s-1vcpu-1gb",
						"region":           "nyc3",
						"node_count":       1,
						"allowed_ips_addr": []string{"185.16.161.205", "159.65.87.97"},
					},
				},
			},
		}

		tf, teardown := setup(t, appDef)
		defer teardown()

		err := tf.Init(t.Context())
		require.NoError(t, err)

		got, err := tf.Plan(t.Context(), env.Production)
		require.NoError(t, err)
		require.NotNil(t, got)

		t.Log("Firewall Rules for Allowed IPs")
		{
			var firewallRules []map[string]any
			for _, rc := range got.Plan.ResourceChanges {
				if rc.Type == "digitalocean_database_firewall" && rc.Name == "this" {
					firewallData := rc.Change.After.(map[string]any)
					rules := firewallData["rule"].([]any)
					for _, rule := range rules {
						firewallRules = append(firewallRules, rule.(map[string]any))
					}
					break
				}
			}

			assert.NotEmpty(t, firewallRules, "Firewall rules should be planned")

			// Verify IP addresses are in the firewall rules.
			ipAddresses := make([]string, 0)
			for _, rule := range firewallRules {
				if rule["type"] == "ip_addr" {
					ipAddresses = append(ipAddresses, rule["value"].(string))
				}
			}

			assert.Contains(t, ipAddresses, "185.16.161.205")
			assert.Contains(t, ipAddresses, "159.65.87.97")
		}
	})

	t.Run("Digital Ocean - Spaces - Basic", func(t *testing.T) {
		appDef := &appdef.Definition{
			Project: appdef.Project{
				Name: "project",
				Repo: appdef.GitHubRepo{
					Owner: "ainsley-dev",
					Name:  "project",
				},
			},
			Resources: []appdef.Resource{
				{
					Name:     "storage",
					Type:     appdef.ResourceTypeS3,
					Provider: appdef.ResourceProviderDigitalOcean,
					Config: map[string]any{
						"acl": "private",
					},
				},
			},
		}

		tf, teardown := setup(t, appDef)
		defer teardown()

		err := tf.Init(t.Context())
		require.NoError(t, err)

		got, err := tf.Plan(t.Context(), env.Production)
		require.NoError(t, err)
		require.NotNil(t, got)
		assert.True(t, got.HasChanges, "Plan should have changes")

		t.Log("Plan Summary")
		{
			require.Len(t, got.Plan.ResourceChanges, 10, "Should plan to create 10 resources")
		}

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

		t.Log("Resource Outputs")
		{
			resources := got.Plan.PlannedValues.Outputs["resources"]
			require.NotNil(t, resources)
			assert.True(t, resources.Sensitive, "Resources output should be sensitive")

			resourceNames := got.Plan.PlannedValues.Outputs["resource_names"]
			require.NotNil(t, resourceNames)
			names := resourceNames.Value.([]any)
			assert.Len(t, names, 1)
			assert.Equal(t, "storage", names[0])
		}

		t.Log("GitHub Secrets")
		{
			secrets := got.Plan.PlannedValues.Outputs["github_secrets_created"]
			require.NotNil(t, secrets)
			assert.Equal(t, 6, int(got.Plan.PlannedValues.Outputs["github_secrets_count"].Value.(float64)))

			secretNames := secrets.Value.([]any)
			assert.Len(t, secretNames, 6)
			assert.Contains(t, secretNames, "TF_PROD_STORAGE_BUCKET_NAME")
			assert.Contains(t, secretNames, "TF_PROD_STORAGE_BUCKET_URL")
			assert.Contains(t, secretNames, "TF_PROD_STORAGE_ENDPOINT")
			assert.Contains(t, secretNames, "TF_PROD_STORAGE_ID")
			assert.Contains(t, secretNames, "TF_PROD_STORAGE_REGION")
			assert.Contains(t, secretNames, "TF_PROD_STORAGE_URN")
		}

		t.Log("GitHub Secret Values")
		{
			// Verify that the bucket name secret has the correct plaintext value
			for _, rc := range got.Plan.ResourceChanges {
				if rc.Type == "github_actions_secret" && rc.Name == "resource_outputs" {
					after := rc.Change.After.(map[string]any)
					secretName := after["secret_name"].(string)

					if secretName == "TF_PROD_STORAGE_BUCKET_NAME" {
						assert.Equal(t, "project-storage", after["plaintext_value"])
					}
					if secretName == "TF_PROD_STORAGE_REGION" {
						assert.Equal(t, "ams3", after["plaintext_value"])
					}
				}
			}
		}
	})

	t.Run("Digital Ocean - Spaces - Public ACL", func(t *testing.T) {
		appDef := &appdef.Definition{
			Project: appdef.Project{
				Name: "public-storage",
				Repo: appdef.GitHubRepo{
					Owner: "ainsley-dev",
					Name:  "public-storage",
				},
			},
			Resources: []appdef.Resource{
				{
					Name:     "assets",
					Type:     appdef.ResourceTypeS3,
					Provider: appdef.ResourceProviderDigitalOcean,
					Config: map[string]any{
						"acl":    "public-read",
						"region": "nyc3",
					},
				},
			},
		}

		tf, teardown := setup(t, appDef)
		defer teardown()

		err := tf.Init(t.Context())
		require.NoError(t, err)

		got, err := tf.Plan(t.Context(), env.Production)
		require.NoError(t, err)
		require.NotNil(t, got)

		t.Log("Spaces Bucket with Public ACL")
		{
			var bucket map[string]any
			for _, rc := range got.Plan.ResourceChanges {
				if rc.Type == "digitalocean_spaces_bucket" && rc.Name == "this" {
					bucket = rc.Change.After.(map[string]any)
					break
				}
			}
			require.NotNil(t, bucket, "Spaces bucket resource should be planned")

			assert.Equal(t, "public-read", bucket["acl"])
			assert.Equal(t, "public-storage-assets", bucket["name"])
			assert.Equal(t, "nyc3", bucket["region"])
		}
	})

	t.Run("Multiple Resources - Postgres and S3", func(t *testing.T) {
		appDef := &appdef.Definition{
			Project: appdef.Project{
				Name: "multi-resource",
				Repo: appdef.GitHubRepo{
					Owner: "ainsley-dev",
					Name:  "multi-resource",
				},
			},
			Resources: []appdef.Resource{
				{
					Name:     "db",
					Type:     appdef.ResourceTypePostgres,
					Provider: appdef.ResourceProviderDigitalOcean,
					Config: map[string]any{
						"pg_version": "18",
						"size":       "db-s-1vcpu-1gb",
						"region":     "ams3",
						"node_count": 1,
					},
				},
				{
					Name:     "storage",
					Type:     appdef.ResourceTypeS3,
					Provider: appdef.ResourceProviderDigitalOcean,
					Config: map[string]any{
						"acl":    "private",
						"region": "ams3",
					},
				},
			},
		}

		tf, teardown := setup(t, appDef)
		defer teardown()

		err := tf.Init(t.Context())
		require.NoError(t, err)

		got, err := tf.Plan(t.Context(), env.Production)
		require.NoError(t, err)
		require.NotNil(t, got)

		t.Log("Multiple Resources Planned")
		{
			resourceNames := got.Plan.PlannedValues.Outputs["resource_names"]
			require.NotNil(t, resourceNames)
			names := resourceNames.Value.([]any)
			assert.Len(t, names, 2)
			assert.Contains(t, names, "db")
			assert.Contains(t, names, "storage")
		}

		t.Log("Both Resource Types Present")
		{
			hasPostgres := false
			hasSpaces := false

			for _, rc := range got.Plan.ResourceChanges {
				if rc.Type == "digitalocean_database_cluster" {
					hasPostgres = true
				}
				if rc.Type == "digitalocean_spaces_bucket" {
					hasSpaces = true
				}
			}

			assert.True(t, hasPostgres, "Postgres cluster should be planned")
			assert.True(t, hasSpaces, "Spaces bucket should be planned")
		}

		t.Log("GitHub Secrets for Both Resources")
		{
			secrets := got.Plan.PlannedValues.Outputs["github_secrets_created"]
			require.NotNil(t, secrets)

			secretNames := secrets.Value.([]any)
			// Should have secrets for both resources.
			assert.Contains(t, secretNames, "TF_PROD_DB_CONNECTION_URL")
			assert.Contains(t, secretNames, "TF_PROD_STORAGE_BUCKET_NAME")
		}
	})

	t.Run("BackBlaze - B2 Bucket Custom", func(t *testing.T) {
		appDef := &appdef.Definition{
			Project: appdef.Project{
				Name: "b2-project",
				Repo: appdef.GitHubRepo{
					Owner: "ainsley-dev",
					Name:  "b2-project",
				},
			},
			Resources: []appdef.Resource{
				{
					Name:     "backups",
					Type:     appdef.ResourceTypeS3,
					Provider: appdef.ResourceProviderBackBlaze,
					Config: map[string]any{
						"bucket_type": "allPrivate",
					},
				},
			},
		}

		tf, teardown := setup(t, appDef)
		defer teardown()

		err := tf.Init(t.Context())
		require.NoError(t, err)

		got, err := tf.Plan(t.Context(), env.Production)
		require.NoError(t, err)
		require.NotNil(t, got)

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

			assert.Equal(t, "b2-project-backups", b2Bucket["bucket_name"])
			assert.Equal(t, "allPrivate", b2Bucket["bucket_type"])
		}
	})

	t.Run("Mixed Providers - DO and BackBlaze", func(t *testing.T) {
		appDef := &appdef.Definition{
			Project: appdef.Project{
				Name: "mixed-providers",
				Repo: appdef.GitHubRepo{
					Owner: "ainsley-dev",
					Name:  "mixed-providers",
				},
			},
			Resources: []appdef.Resource{
				{
					Name:     "db",
					Type:     appdef.ResourceTypePostgres,
					Provider: appdef.ResourceProviderDigitalOcean,
					Config: map[string]any{
						"pg_version": "18",
						"size":       "db-s-1vcpu-1gb",
						"region":     "ams3",
					},
				},
				{
					Name:     "backups",
					Type:     appdef.ResourceTypeS3,
					Provider: appdef.ResourceProviderBackBlaze,
					Config: map[string]any{
						"bucket_type": "allPrivate",
					},
				},
			},
		}

		tf, teardown := setup(t, appDef)
		defer teardown()

		err := tf.Init(t.Context())
		require.NoError(t, err)

		got, err := tf.Plan(t.Context(), env.Production)
		require.NoError(t, err)
		require.NotNil(t, got)

		t.Log("Mixed Provider Resources")
		{
			hasDigitalOcean := false
			hasBackBlaze := false

			for _, rc := range got.Plan.ResourceChanges {
				if rc.Type == "digitalocean_database_cluster" {
					hasDigitalOcean = true
				}
				if rc.Type == "b2_bucket" {
					hasBackBlaze = true
				}
			}

			assert.True(t, hasDigitalOcean, "DigitalOcean resource should be planned")
			assert.True(t, hasBackBlaze, "BackBlaze resource should be planned")
		}
	})
}

func TestTerraform_DefaultB2Bucket(t *testing.T) {
	appDef := &appdef.Definition{
		Project: appdef.Project{
			Name: "project",
			Repo: appdef.GitHubRepo{
				Owner: "ainsley-dev",
				Name:  "project",
			},
		},
		Resources: []appdef.Resource{},
		Apps:      []appdef.App{},
	}

	tf, teardown := setup(t, appDef)
	defer teardown()

	err := tf.Init(t.Context())
	require.NoError(t, err)

	got, err := tf.Plan(t.Context(), env.Production)
	require.NoError(t, err)
	require.NotNil(t, got)

	t.Log("Default B2 Bucket Configuration")
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

		// Verify lifecycle rules for single version
		lifecycleRules := b2Bucket["lifecycle_rules"].([]any)
		require.Len(t, lifecycleRules, 1, "Should have exactly one lifecycle rule")

		rule := lifecycleRules[0].(map[string]any)
		assert.Equal(t, float64(1), rule["days_from_hiding_to_deleting"], "Should delete old versions after 1 day")
		assert.Equal(t, float64(0), rule["days_from_uploading_to_hiding"], "Should hide old versions immediately")
		assert.Equal(t, "", rule["file_name_prefix"], "Should apply to all files")
	}

	t.Log("Default B2 Bucket Output")
	{
		defaultBucket := got.Plan.PlannedValues.Outputs["default_b2_bucket"]
		require.NotNil(t, defaultBucket, "Default B2 bucket output should exist")
		assert.True(t, defaultBucket.Sensitive, "Default B2 bucket output should be sensitive")
	}
}

func TestTerraform_Apps(t *testing.T) {
	t.Run("Digital Ocean - SvelteKit App", func(t *testing.T) {
		appDef := &appdef.Definition{
			Project: appdef.Project{
				Name: "project",
				Repo: appdef.GitHubRepo{
					Owner: "ainsley-dev",
					Name:  "project",
				},
			},
			Apps: []appdef.App{
				{
					Name: "api",
					Type: appdef.AppTypeSvelteKit, // or whatever app type
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
					//EnvVars: []appdef.EnvVar{
					//	{
					//		Key:   "DATABASE_URL",
					//		Value: "{{db.connection_url}}",
					//	},
					//},
				},
			},
		}

		tf, teardown := setup(t, appDef)
		defer teardown()

		err := tf.Init(t.Context())
		require.NoError(t, err)

		got, err := tf.Plan(t.Context(), env.Production)
		require.NoError(t, err)
		require.NotNil(t, got)
		assert.True(t, got.HasChanges, "Plan should have changes")

		t.Log("Plan Summary")
		{
			// Note: Creates app + default B2 bucket.
			require.GreaterOrEqual(t, len(got.Plan.ResourceChanges), 1, "Should plan to create at least 1 resource")
		}

		t.Log("Database Resources")
		{
			// Verify database was created (reuse assertions from Postgres test)
			var dbCluster map[string]any
			for _, rc := range got.Plan.ResourceChanges {
				if rc.Type == "digitalocean_database_cluster" && rc.Name == "this" {
					dbCluster = rc.Change.After.(map[string]any)
					break
				}
			}
			require.NotNil(t, dbCluster, "Database cluster should be planned")
			assert.Equal(t, "project-db", dbCluster["name"])
		}

		t.Log("App Platform Configuration")
		{
			var app map[string]any
			for _, rc := range got.Plan.ResourceChanges {
				if rc.Type == "digitalocean_app" && rc.Name == "this" {
					app = rc.Change.After.(map[string]any)
					break
				}
			}
			require.NotNil(t, app, "App Platform app should be planned")

			// Check spec
			spec := app["spec"].([]any)[0].(map[string]any)
			assert.Equal(t, "project-api", spec["name"])
			assert.Equal(t, "lon", spec["region"])

			// Check alerts
			alerts := spec["alert"].([]any)
			assert.Len(t, alerts, 2)

			alertRules := make([]string, len(alerts))
			for i, alert := range alerts {
				alertRules[i] = alert.(map[string]any)["rule"].(string)
			}
			assert.Contains(t, alertRules, "DEPLOYMENT_FAILED")
			assert.Contains(t, alertRules, "DEPLOYMENT_LIVE")
		}

		t.Log("App Service Configuration")
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
			assert.Equal(t, "api", service["name"])
			assert.Equal(t, float64(3000), service["http_port"])
			assert.Equal(t, float64(1), service["instance_count"])
			assert.Equal(t, "apps-s-1vcpu-1gb", service["instance_size_slug"])

			// Check health check
			healthChecks := service["health_check"].([]any)
			require.Len(t, healthChecks, 1)
			healthCheck := healthChecks[0].(map[string]any)
			assert.Equal(t, "/health", healthCheck["http_path"])
			assert.Equal(t, float64(10), healthCheck["failure_threshold"])
			assert.Equal(t, float64(90), healthCheck["initial_delay_seconds"])
			assert.Equal(t, float64(5), healthCheck["period_seconds"])

			// Check service alerts
			serviceAlerts := service["alert"].([]any)
			assert.Len(t, serviceAlerts, 3)

			serviceAlertRules := make(map[string]map[string]any)
			for _, alert := range serviceAlerts {
				alertMap := alert.(map[string]any)
				rule := alertMap["rule"].(string)
				serviceAlertRules[rule] = alertMap
			}

			// CPU alert
			assert.Contains(t, serviceAlertRules, "CPU_UTILIZATION")
			assert.Equal(t, "GREATER_THAN", serviceAlertRules["CPU_UTILIZATION"]["operator"])
			assert.Equal(t, float64(80), serviceAlertRules["CPU_UTILIZATION"]["value"])
			assert.Equal(t, "FIVE_MINUTES", serviceAlertRules["CPU_UTILIZATION"]["window"])

			// Memory alert
			assert.Contains(t, serviceAlertRules, "MEM_UTILIZATION")
			assert.Equal(t, float64(80), serviceAlertRules["MEM_UTILIZATION"]["value"])

			// Restart count alert
			assert.Contains(t, serviceAlertRules, "RESTART_COUNT")
			assert.Equal(t, float64(3), serviceAlertRules["RESTART_COUNT"]["value"])
		}

		t.Log("Container Image Configuration")
		{
			var app map[string]any
			for _, rc := range got.Plan.ResourceChanges {
				if rc.Type == "digitalocean_app" && rc.Name == "this" {
					app = rc.Change.After.(map[string]any)
					break
				}
			}

			spec := app["spec"].([]any)[0].(map[string]any)
			service := spec["service"].([]any)[0].(map[string]any)
			images := service["image"].([]any)
			require.Len(t, images, 1)

			image := images[0].(map[string]any)
			assert.Equal(t, "ghcr.io", image["registry"])
			assert.Equal(t, "GHCR", image["registry_type"])
			assert.Equal(t, "latest", image["tag"])

			// Check registry credentials reference
			assert.NotNil(t, image["registry_credentials"])
		}

		t.Log("GitHub Secrets")
		{
			secrets := got.Plan.PlannedValues.Outputs["github_secrets_created"]
			require.NotNil(t, secrets)

			secretNames := secrets.Value.([]any)
			// Should have 3 DB secrets (connection_url, id, urn)
			assert.Len(t, secretNames, 3)
			assert.Contains(t, secretNames, "TF_PROD_DB_CONNECTION_URL")
			assert.Contains(t, secretNames, "TF_PROD_DB_ID")
			assert.Contains(t, secretNames, "TF_PROD_DB_URN")
		}

		t.Log("Resource Outputs")
		{
			resourceNames := got.Plan.PlannedValues.Outputs["resource_names"]
			require.NotNil(t, resourceNames)
			names := resourceNames.Value.([]any)
			assert.Len(t, names, 1)
			assert.Equal(t, "db", names[0])
		}
	})

	t.Run("Digital Ocean - GoLang App", func(t *testing.T) {
		appDef := &appdef.Definition{
			Project: appdef.Project{
				Name: "golang-project",
				Repo: appdef.GitHubRepo{
					Owner: "ainsley-dev",
					Name:  "golang-project",
				},
			},
			Apps: []appdef.App{
				{
					Name: "api",
					Type: appdef.AppTypeGoLang,
					Path: "apps/api",
					Infra: appdef.Infra{
						Provider: appdef.ResourceProviderDigitalOcean,
						Type:     "container",
						Config: map[string]any{
							"region":            "lon",
							"size":              "apps-s-1vcpu-1gb",
							"instance_count":    1,
							"port":              8080,
							"health_check_path": "/health",
						},
					},
				},
			},
		}

		tf, teardown := setup(t, appDef)
		defer teardown()

		err := tf.Init(t.Context())
		require.NoError(t, err)

		got, err := tf.Plan(t.Context(), env.Production)
		require.NoError(t, err)
		require.NotNil(t, got)

		t.Log("App Platform App Configuration")
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
			assert.Equal(t, "golang-project-api", spec["name"])

			services := spec["service"].([]any)
			require.Len(t, services, 1)

			service := services[0].(map[string]any)
			assert.Equal(t, "api", service["name"])
			assert.Equal(t, float64(8080), service["http_port"])
		}
	})

	t.Run("Digital Ocean - Payload CMS App", func(t *testing.T) {
		appDef := &appdef.Definition{
			Project: appdef.Project{
				Name: "cms-project",
				Repo: appdef.GitHubRepo{
					Owner: "ainsley-dev",
					Name:  "cms-project",
				},
			},
			Apps: []appdef.App{
				{
					Name: "cms",
					Type: appdef.AppTypePayload,
					Path: "apps/cms",
					Infra: appdef.Infra{
						Provider: appdef.ResourceProviderDigitalOcean,
						Type:     "container",
						Config: map[string]any{
							"region":         "nyc3",
							"size":           "apps-s-1vcpu-2gb",
							"instance_count": 2,
						},
					},
				},
			},
		}

		tf, teardown := setup(t, appDef)
		defer teardown()

		err := tf.Init(t.Context())
		require.NoError(t, err)

		got, err := tf.Plan(t.Context(), env.Production)
		require.NoError(t, err)
		require.NotNil(t, got)

		t.Log("Payload CMS Configuration")
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
			services := spec["service"].([]any)
			require.Len(t, services, 1)

			service := services[0].(map[string]any)
			assert.Equal(t, "cms", service["name"])
			assert.Equal(t, float64(3000), service["http_port"])
			assert.Equal(t, float64(2), service["instance_count"])
			assert.Equal(t, "apps-s-1vcpu-2gb", service["instance_size_slug"])
		}
	})

	t.Run("App With Environment Variables", func(t *testing.T) {
		appDef := &appdef.Definition{
			Project: appdef.Project{
				Name: "env-app",
				Repo: appdef.GitHubRepo{
					Owner: "ainsley-dev",
					Name:  "env-app",
				},
			},
			Apps: []appdef.App{
				{
					Name: "web",
					Type: appdef.AppTypeSvelteKit,
					Path: "apps/web",
					Infra: appdef.Infra{
						Provider: appdef.ResourceProviderDigitalOcean,
						Type:     "container",
						Config: map[string]any{
							"region": "lon",
						},
					},
					Env: appdef.Environment{
						Production: map[string]appdef.EnvValue{
							"API_URL":    {Value: "https://api.example.com", Source: appdef.EnvSourceValue},
							"SECRET_KEY": {Value: "supersecret", Source: appdef.EnvSourceSOPS},
						},
					},
				},
			},
		}

		tf, teardown := setup(t, appDef)
		defer teardown()

		err := tf.Init(t.Context())
		require.NoError(t, err)

		got, err := tf.Plan(t.Context(), env.Production)
		require.NoError(t, err)
		require.NotNil(t, got)

		t.Log("App With Environment Variables Planned")
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
			assert.Equal(t, "env-app-web", spec["name"])

			// Note: Environment variable propagation to Terraform plan may vary.
			// This test verifies the app is created; env var details are tested elsewhere.
		}
	})

	t.Run("Multiple Apps", func(t *testing.T) {
		appDef := &appdef.Definition{
			Project: appdef.Project{
				Name: "multi-app",
				Repo: appdef.GitHubRepo{
					Owner: "ainsley-dev",
					Name:  "multi-app",
				},
			},
			Apps: []appdef.App{
				{
					Name: "frontend",
					Type: appdef.AppTypeSvelteKit,
					Path: "apps/frontend",
					Infra: appdef.Infra{
						Provider: appdef.ResourceProviderDigitalOcean,
						Type:     "container",
						Config: map[string]any{
							"region": "lon",
						},
					},
				},
				{
					Name: "backend",
					Type: appdef.AppTypeGoLang,
					Path: "apps/backend",
					Infra: appdef.Infra{
						Provider: appdef.ResourceProviderDigitalOcean,
						Type:     "container",
						Config: map[string]any{
							"region": "lon",
						},
					},
				},
			},
		}

		tf, teardown := setup(t, appDef)
		defer teardown()

		err := tf.Init(t.Context())
		require.NoError(t, err)

		got, err := tf.Plan(t.Context(), env.Production)
		require.NoError(t, err)
		require.NotNil(t, got)

		t.Log("Multiple Apps Planned")
		{
			appCount := 0
			for _, rc := range got.Plan.ResourceChanges {
				if rc.Type == "digitalocean_app" {
					appCount++
				}
			}
			assert.Equal(t, 2, appCount, "Should plan 2 apps")
		}
	})
}

func TestTerraform_Integration(t *testing.T) {
	t.Run("App With Postgres Resource", func(t *testing.T) {
		appDef := &appdef.Definition{
			Project: appdef.Project{
				Name: "fullstack",
				Repo: appdef.GitHubRepo{
					Owner: "ainsley-dev",
					Name:  "fullstack",
				},
			},
			Resources: []appdef.Resource{
				{
					Name:     "db",
					Type:     appdef.ResourceTypePostgres,
					Provider: appdef.ResourceProviderDigitalOcean,
					Config: map[string]any{
						"pg_version": "18",
						"size":       "db-s-1vcpu-1gb",
						"region":     "ams3",
						"node_count": 1,
					},
				},
			},
			Apps: []appdef.App{
				{
					Name: "api",
					Type: appdef.AppTypeGoLang,
					Path: "apps/api",
					Infra: appdef.Infra{
						Provider: appdef.ResourceProviderDigitalOcean,
						Type:     "container",
						Config: map[string]any{
							"region": "ams3",
						},
					},
					Env: appdef.Environment{
						Production: map[string]appdef.EnvValue{
							"DB_HOST": {Value: "db.example.com", Source: appdef.EnvSourceValue},
						},
					},
				},
			},
		}

		tf, teardown := setup(t, appDef)
		defer teardown()

		err := tf.Init(t.Context())
		require.NoError(t, err)

		got, err := tf.Plan(t.Context(), env.Production)
		require.NoError(t, err)
		require.NotNil(t, got)

		t.Log("Both App and Resource Planned")
		{
			hasApp := false
			hasDatabase := false

			for _, rc := range got.Plan.ResourceChanges {
				if rc.Type == "digitalocean_app" {
					hasApp = true
				}
				if rc.Type == "digitalocean_database_cluster" {
					hasDatabase = true
				}
			}

			assert.True(t, hasApp, "App should be planned")
			assert.True(t, hasDatabase, "Database should be planned")
		}

		t.Log("GitHub Secrets Include Both")
		{
			secrets := got.Plan.PlannedValues.Outputs["github_secrets_created"]
			require.NotNil(t, secrets)

			secretNames := secrets.Value.([]any)
			// Should have DB secrets.
			assert.Contains(t, secretNames, "TF_PROD_DB_CONNECTION_URL")
		}
	})

	t.Run("Multiple Apps With Shared Resources", func(t *testing.T) {
		appDef := &appdef.Definition{
			Project: appdef.Project{
				Name: "complex",
				Repo: appdef.GitHubRepo{
					Owner: "ainsley-dev",
					Name:  "complex",
				},
			},
			Resources: []appdef.Resource{
				{
					Name:     "db",
					Type:     appdef.ResourceTypePostgres,
					Provider: appdef.ResourceProviderDigitalOcean,
					Config: map[string]any{
						"pg_version": "18",
						"size":       "db-s-1vcpu-1gb",
						"region":     "ams3",
					},
				},
				{
					Name:     "storage",
					Type:     appdef.ResourceTypeS3,
					Provider: appdef.ResourceProviderDigitalOcean,
					Config: map[string]any{
						"acl":    "private",
						"region": "ams3",
					},
				},
			},
			Apps: []appdef.App{
				{
					Name: "api",
					Type: appdef.AppTypeGoLang,
					Path: "apps/api",
					Infra: appdef.Infra{
						Provider: appdef.ResourceProviderDigitalOcean,
						Type:     "container",
						Config: map[string]any{
							"region": "ams3",
						},
					},
				},
				{
					Name: "web",
					Type: appdef.AppTypeSvelteKit,
					Path: "apps/web",
					Infra: appdef.Infra{
						Provider: appdef.ResourceProviderDigitalOcean,
						Type:     "container",
						Config: map[string]any{
							"region": "ams3",
						},
					},
				},
			},
			Shared: appdef.Shared{
				Env: appdef.Environment{
					Production: map[string]appdef.EnvValue{
						"SHARED_KEY": {Value: "shared", Source: appdef.EnvSourceValue},
					},
				},
			},
		}

		tf, teardown := setup(t, appDef)
		defer teardown()

		err := tf.Init(t.Context())
		require.NoError(t, err)

		got, err := tf.Plan(t.Context(), env.Production)
		require.NoError(t, err)
		require.NotNil(t, got)

		t.Log("Complex Infrastructure Planned")
		{
			hasApps := 0
			hasDatabase := false
			hasStorage := false

			for _, rc := range got.Plan.ResourceChanges {
				if rc.Type == "digitalocean_app" {
					hasApps++
				}
				if rc.Type == "digitalocean_database_cluster" {
					hasDatabase = true
				}
				if rc.Type == "digitalocean_spaces_bucket" {
					hasStorage = true
				}
			}

			assert.Equal(t, 2, hasApps, "Should plan 2 apps")
			assert.True(t, hasDatabase, "Database should be planned")
			assert.True(t, hasStorage, "Storage should be planned")
		}

		t.Log("Resource and App Names")
		{
			resourceNames := got.Plan.PlannedValues.Outputs["resource_names"]
			require.NotNil(t, resourceNames)
			names := resourceNames.Value.([]any)
			assert.Len(t, names, 2, "Should have 2 resources")
			assert.Contains(t, names, "db")
			assert.Contains(t, names, "storage")

			// Note: Shared environment variable propagation to Terraform plan may vary.
			// This test verifies the infrastructure is created correctly.
		}
	})
}

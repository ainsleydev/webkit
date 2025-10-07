package integration

import (
	"testing"

	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

func TestTerraform_Resources_Postgres(t *testing.T) {
	t.Parallel()

	options := setupTerraform(t, "postgres.tfvars")
	plan, err := terraform.InitAndPlanAndShowWithStructE(t, options)

	t.Log("Plan Succeeded")
	{
		assert.NoError(t, err, "Terraform plan should be successful")
		assert.Equal(t, 6, len(plan.ResourcePlannedValuesMap))
	}

	t.Log("Database")
	{
		resource, err := findResource("digitalocean_database_db", plan.ResourcePlannedValuesMap)
		assert.NoError(t, err)
		assert.Equal(t, "postgres", resource.AttributeValues["name"])
	}

	t.Log("Database Cluster")
	{
		resource, err := findResource("digitalocean_database_cluster", plan.ResourcePlannedValuesMap)
		assert.NoError(t, err)

		values := resource.AttributeValues
		assert.Equal(t, "postgres", values["name"])
		assert.Equal(t, "pg", values["engine"])
		assert.Equal(t, "16", values["version"])
		assert.Equal(t, "db-s-2vcpu-1gb", values["size"])
		assert.Equal(t, "lon2", values["region"])
		assert.EqualValues(t, 2, values["node_count"])
		assert.ElementsMatch(t, []string{"development", "terraform", "test"}, values["tags"])
	}

	t.Log("Database User")
	{
		resource, err := findResource("digitalocean_database_user", plan.ResourcePlannedValuesMap)
		assert.NoError(t, err)
		assert.Equal(t, "postgres_admin", resource.AttributeValues["name"])
	}

	t.Log("Database Connection Pool")
	{
		resource, err := findResource("digitalocean_database_connection_pool", plan.ResourcePlannedValuesMap)
		assert.NoError(t, err)

		values := resource.AttributeValues
		assert.Equal(t, "postgres_pool", values["name"])
		assert.Equal(t, "transaction", values["mode"])
		assert.EqualValues(t, 20, values["size"])
		assert.Equal(t, "postgres", values["db_name"])
		assert.Equal(t, "postgres_admin", values["user"])
	}

	t.Log("Database Firewall")
	{
		resource, err := findResource("digitalocean_database_firewall", plan.ResourcePlannedValuesMap)
		assert.NoError(t, err)

		values := resource.AttributeValues

		rules, ok := values["rule"].([]any)
		assert.True(t, ok, "Rule should be an array")
		assert.Len(t, rules, 1)

		rule, ok := rules[0].(map[string]any)
		assert.True(t, ok)
		assert.Equal(t, "ip_addr", rule["type"])
		assert.Equal(t, "192.168.1.1", rule["value"])
	}
}

func TestTerraform_Resources_Spaces(t *testing.T) {
	t.Parallel()

	options := setupTerraform(t, "spaces.tfvars")
	plan, err := terraform.InitAndPlanAndShowWithStructE(t, options)

	t.Log("Plan Succeeded")
	{
		assert.NoError(t, err, "Terraform plan should be successful")
		assert.Equal(t, 3, len(plan.ResourcePlannedValuesMap))
	}

	t.Log("Spaces Bucket")
	{
		resource, err := findResource("digitalocean_spaces_bucket", plan.ResourcePlannedValuesMap)
		assert.NoError(t, err)

		values := resource.AttributeValues
		assert.Equal(t, "test-bucket", values["name"])
		assert.Equal(t, "ams3", values["region"])
		assert.Equal(t, "public-read", values["acl"])
	}

	t.Log("CORS Configuration")
	{
		resource, err := findResource("digitalocean_spaces_bucket_cors_configuration", plan.ResourcePlannedValuesMap)
		assert.NoError(t, err)

		values := resource.AttributeValues
		assert.Equal(t, "ams3", values["region"])

		// CORS
		corsRules, ok := values["cors_rule"].([]any)
		assert.True(t, ok, "cors_rule should be an array")
		assert.Len(t, corsRules, 1)

		rule, ok := corsRules[0].(map[string]any)
		assert.True(t, ok)

		// Headers
		allowedHeaders, ok := rule["allowed_headers"].([]any)
		assert.True(t, ok)
		assert.ElementsMatch(t, []any{"*"}, allowedHeaders)

		// Methods
		allowedMethods, ok := rule["allowed_methods"].([]any)
		assert.True(t, ok)
		assert.ElementsMatch(t, []any{"GET"}, allowedMethods)

		// Allowed origins
		allowedOrigins, ok := rule["allowed_origins"].([]any)
		assert.True(t, ok)
		assert.ElementsMatch(t, []any{"*"}, allowedOrigins)

		// Max age
		assert.EqualValues(t, 31536000, rule["max_age_seconds"])
	}

	t.Log("CDN")
	{
		resource, err := findResource("digitalocean_cdn", plan.ResourcePlannedValuesMap)
		assert.NoError(t, err)
		assert.NotNil(t, resource)
	}
}

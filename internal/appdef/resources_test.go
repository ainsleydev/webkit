package appdef

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ainsleydev/webkit/pkg/env"
	"github.com/ainsleydev/webkit/pkg/util/ptr"
)

func TestResourceType_String(t *testing.T) {
	t.Parallel()

	got := ResourceTypePostgres.String()
	assert.Equal(t, "postgres", got)
	assert.IsType(t, "", got)
}

func TestResourceProvider_String(t *testing.T) {
	t.Parallel()

	got := ResourceProviderDigitalOcean.String()
	assert.Equal(t, "digitalocean", got)
	assert.IsType(t, "", got)
}

func TestRequiredOutputs(t *testing.T) {
	t.Parallel()

	tt := map[string]struct {
		input ResourceType
		want  []string
	}{
		"Postgres": {
			input: ResourceTypePostgres,
			want:  []string{"id", "connection_url", "host", "port", "database", "user", "password"},
		},
		"S3": {
			input: ResourceTypeS3,
			want:  []string{"id", "bucket_name", "bucket_url", "region"},
		},
		"UnknownType": {
			input: ResourceType("unknown"),
			want:  nil,
		},
	}

	for name, test := range tt {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			got := test.input.Outputs()
			assert.Equal(t, test.want, got)
		})
	}
}

func TestGitHubSecretName(t *testing.T) {
	t.Parallel()

	tt := []struct {
		resource Resource
		env      env.Environment
		output   string
		want     string
	}{
		{resource: Resource{Name: "db"}, env: env.Production, output: "connection_url", want: "TF_PROD_DB_CONNECTION_URL"},
		{resource: Resource{Name: "my-app-db"}, env: env.Staging, output: "password", want: "TF_STAGING_MY_APP_DB_PASSWORD"},
		{resource: Resource{Name: "cache"}, env: env.Development, output: "host", want: "TF_DEV_CACHE_HOST"},
	}

	for _, test := range tt {
		t.Run(test.want, func(t *testing.T) {
			t.Parallel()
			got := test.resource.GitHubSecretName(test.env, test.output)
			assert.Equal(t, test.want, got)
		})
	}
}

func TestResourceApplyDefaults(t *testing.T) {
	t.Parallel()

	tt := map[string]struct {
		input Resource
		want  Resource
	}{
		"Nil Config And Outputs": {
			input: Resource{},
			want: Resource{
				Config:     make(map[string]any),
				Backup:     ResourceBackupConfig{Enabled: ptr.BoolPtr(true)},
				Monitoring: ptr.BoolPtr(true),
			},
		},
		"Existing Config And Outputs": {
			input: Resource{
				Config: map[string]any{"size": "small"},
			},
			want: Resource{
				Config:     map[string]any{"size": "small"},
				Backup:     ResourceBackupConfig{Enabled: ptr.BoolPtr(true)},
				Monitoring: ptr.BoolPtr(true),
			},
		},
	}

	for name, test := range tt {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			test.input.applyDefaults()
			assert.Equal(t, test.want, test.input)
		})
	}
}

func TestResource_IsTerraformManaged(t *testing.T) {
	t.Parallel()

	tt := map[string]struct {
		terraformManaged *bool
		want             bool
	}{
		"Nil defaults to true": {terraformManaged: nil, want: true},
		"Explicit false":       {terraformManaged: ptr.BoolPtr(false), want: false},
		"Explicit true":        {terraformManaged: ptr.BoolPtr(true), want: true},
	}

	for name, test := range tt {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			resource := Resource{TerraformManaged: test.terraformManaged}
			got := resource.IsTerraformManaged()
			assert.Equal(t, test.want, got)
		})
	}
}

func TestResource_IsBackupEnabled(t *testing.T) {
	t.Parallel()

	tt := map[string]struct {
		backup ResourceBackupConfig
		want   bool
	}{
		"Nil defaults to true": {backup: ResourceBackupConfig{Enabled: nil}, want: true},
		"Explicit false":       {backup: ResourceBackupConfig{Enabled: ptr.BoolPtr(false)}, want: false},
		"Explicit true":        {backup: ResourceBackupConfig{Enabled: ptr.BoolPtr(true)}, want: true},
	}

	for name, test := range tt {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			resource := Resource{Backup: test.backup}
			got := resource.IsBackupEnabled()
			assert.Equal(t, test.want, got)
		})
	}
}

func TestResource_IsMonitoringEnabled(t *testing.T) {
	t.Parallel()

	tt := map[string]struct {
		monitoring *bool
		want       bool
	}{
		"Nil defaults to true": {monitoring: nil, want: true},
		"Explicit false":       {monitoring: ptr.BoolPtr(false), want: false},
		"Explicit true":        {monitoring: ptr.BoolPtr(true), want: true},
	}

	for name, test := range tt {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			resource := Resource{Monitoring: test.monitoring}
			got := resource.IsMonitoringEnabled()
			assert.Equal(t, test.want, got)
		})
	}
}

func TestResource_GenerateBackupMonitor(t *testing.T) {
	t.Parallel()

	tt := map[string]struct {
		resource     Resource
		projectTitle string
		want         *Monitor
	}{
		"Backup and monitoring enabled": {
			resource: Resource{
				Name:       "db",
				Title:      "Database",
				Backup:     ResourceBackupConfig{Enabled: ptr.BoolPtr(true)},
				Monitoring: ptr.BoolPtr(true),
			},
			projectTitle: "Test Project",
			want: &Monitor{
				Name: "Test Project - Database Backup",
				Type: MonitorTypePush,
			},
		},
		"Backup disabled": {
			resource: Resource{
				Name:       "db",
				Title:      "Database",
				Backup:     ResourceBackupConfig{Enabled: ptr.BoolPtr(false)},
				Monitoring: ptr.BoolPtr(true),
			},
			projectTitle: "Test Project",
			want:         nil,
		},
		"Monitoring disabled": {
			resource: Resource{
				Name:       "db",
				Title:      "Database",
				Backup:     ResourceBackupConfig{Enabled: ptr.BoolPtr(true)},
				Monitoring: ptr.BoolPtr(false),
			},
			projectTitle: "Test Project",
			want:         nil,
		},
		"Both disabled": {
			resource: Resource{
				Name:       "db",
				Title:      "Database",
				Backup:     ResourceBackupConfig{Enabled: ptr.BoolPtr(false)},
				Monitoring: ptr.BoolPtr(false),
			},
			projectTitle: "Test Project",
			want:         nil,
		},
	}

	for name, test := range tt {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			got := test.resource.GenerateBackupMonitor(test.projectTitle)
			assert.Equal(t, test.want, got)
		})
	}
}

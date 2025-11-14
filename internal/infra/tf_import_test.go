package infra

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ainsleydev/webkit/internal/appdef"
)

func TestBuildImportAddresses(t *testing.T) {
	t.Parallel()

	tt := map[string]struct {
		projectName string
		resource    *appdef.Resource
		baseID      string
		want        []importAddress
		wantErr     bool
	}{
		"DigitalOcean Postgres without firewall": {
			projectName: "test-project",
			resource: &appdef.Resource{
				Name:     "db",
				Type:     appdef.ResourceTypePostgres,
				Provider: appdef.ResourceProviderDigitalOcean,
				Config:   map[string]any{},
			},
			baseID: "cluster-123",
			want: []importAddress{
				{
					Address: "module.resources[\"db\"].module.do_postgres[0].digitalocean_database_cluster.this",
					ID:      "cluster-123",
				},
				{
					Address: "module.resources[\"db\"].module.do_postgres[0].digitalocean_database_user.this",
					ID:      "cluster-123,test_project_db_admin",
				},
				{
					Address: "module.resources[\"db\"].module.do_postgres[0].digitalocean_database_db.this",
					ID:      "cluster-123,test_project_db",
				},
				{
					Address: "module.resources[\"db\"].module.do_postgres[0].digitalocean_database_connection_pool.this",
					ID:      "cluster-123,test_project_db_pool",
				},
			},
			wantErr: false,
		},
		"DigitalOcean Postgres with firewall (search-spares example)": {
			projectName: "search-spares",
			resource: &appdef.Resource{
				Name:     "db",
				Type:     appdef.ResourceTypePostgres,
				Provider: appdef.ResourceProviderDigitalOcean,
				Config: map[string]any{
					"allowed_ips_addr": []any{"185.16.161.205", "159.65.87.97"},
				},
			},
			baseID: "cluster-456",
			want: []importAddress{
				{
					Address: "module.resources[\"db\"].module.do_postgres[0].digitalocean_database_cluster.this",
					ID:      "cluster-456",
				},
				{
					Address: "module.resources[\"db\"].module.do_postgres[0].digitalocean_database_user.this",
					ID:      "cluster-456,search_spares_db_admin",
				},
				{
					Address: "module.resources[\"db\"].module.do_postgres[0].digitalocean_database_db.this",
					ID:      "cluster-456,search_spares_db",
				},
				{
					Address: "module.resources[\"db\"].module.do_postgres[0].digitalocean_database_connection_pool.this",
					ID:      "cluster-456,search_spares_db_pool",
				},
				{
					Address: "module.resources[\"db\"].module.do_postgres[0].digitalocean_database_firewall.this[0]",
					ID:      "cluster-456",
				},
			},
			wantErr: false,
		},
		"DigitalOcean S3 bucket without CDN": {
			projectName: "test-project",
			resource: &appdef.Resource{
				Name:     "media-bucket",
				Type:     appdef.ResourceTypeS3,
				Provider: appdef.ResourceProviderDigitalOcean,
				Config:   map[string]any{},
			},
			baseID: "bucket-789",
			want: []importAddress{
				{
					Address: "module.resources[\"media-bucket\"].module.do_bucket[0].digitalocean_spaces_bucket.this",
					ID:      "ams3,bucket-789",
				},
				{
					Address: "module.resources[\"media-bucket\"].module.do_bucket[0].digitalocean_spaces_bucket_cors_configuration.this",
					ID:      "ams3,bucket-789",
				},
			},
			wantErr: false,
		},
		"DigitalOcean S3 bucket with CDN": {
			projectName: "test-project",
			resource: &appdef.Resource{
				Name:     "store",
				Type:     appdef.ResourceTypeS3,
				Provider: appdef.ResourceProviderDigitalOcean,
				Config: map[string]any{
					"region": "ams3",
					"cdn_id": "26226874-569d-41d9-9191-e18a09a2d134",
				},
			},
			baseID: "s-clark-store",
			want: []importAddress{
				{
					Address: "module.resources[\"store\"].module.do_bucket[0].digitalocean_spaces_bucket.this",
					ID:      "ams3,s-clark-store",
				},
				{
					Address: "module.resources[\"store\"].module.do_bucket[0].digitalocean_spaces_bucket_cors_configuration.this",
					ID:      "ams3,s-clark-store",
				},
				{
					Address: "module.resources[\"store\"].module.do_bucket[0].digitalocean_cdn.this",
					ID:      "26226874-569d-41d9-9191-e18a09a2d134",
				},
			},
			wantErr: false,
		},
		"Unsupported provider": {
			projectName: "test-project",
			resource: &appdef.Resource{
				Name:     "cache",
				Type:     appdef.ResourceTypePostgres,
				Provider: "unsupported",
				Config:   map[string]any{},
			},
			baseID:  "id-123",
			want:    nil,
			wantErr: true,
		},
		"Unsupported resource type": {
			projectName: "test-project",
			resource: &appdef.Resource{
				Name:     "unknown",
				Type:     "redis",
				Provider: appdef.ResourceProviderDigitalOcean,
				Config:   map[string]any{},
			},
			baseID:  "id-456",
			want:    nil,
			wantErr: true,
		},
		"Turso SQLite database": {
			projectName: "test-project",
			resource: &appdef.Resource{
				Name:     "db",
				Type:     appdef.ResourceTypeSQLite,
				Provider: appdef.ResourceProviderTurso,
				Config: map[string]any{
					"organisation": "my-org",
					"group":        "default",
				},
			},
			baseID: "my-org/my-database",
			want: []importAddress{
				{
					Address: "module.resources[\"db\"].module.turso_database[0].turso_database.this",
					ID:      "my-org/my-database",
				},
				{
					Address: "module.resources[\"db\"].module.turso_database[0].turso_database_token.this",
					ID:      "my-org/my-database",
				},
			},
			wantErr: false,
		},
		"Turso with unsupported resource type": {
			projectName: "test-project",
			resource: &appdef.Resource{
				Name:     "cache",
				Type:     appdef.ResourceTypePostgres,
				Provider: appdef.ResourceProviderTurso,
				Config:   map[string]any{},
			},
			baseID:  "my-org/my-db",
			want:    nil,
			wantErr: true,
		},
	}

	for name, test := range tt {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got, err := buildImportAddresses(test.projectName, test.resource, test.baseID)
			assert.Equal(t, test.wantErr, err != nil)

			if !test.wantErr {
				assert.Equal(t, test.want, got)
			}
		})
	}
}

func TestBuildPostgresImports(t *testing.T) {
	t.Parallel()

	tt := map[string]struct {
		projectName string
		resource    *appdef.Resource
		clusterID   string
		want        []importAddress
	}{
		"Simple database name": {
			projectName: "my-project",
			resource: &appdef.Resource{
				Name:   "db",
				Config: map[string]any{},
			},
			clusterID: "abc123",
			want: []importAddress{
				{
					Address: "module.resources[\"db\"].module.do_postgres[0].digitalocean_database_cluster.this",
					ID:      "abc123",
				},
				{
					Address: "module.resources[\"db\"].module.do_postgres[0].digitalocean_database_user.this",
					ID:      "abc123,my_project_db_admin",
				},
				{
					Address: "module.resources[\"db\"].module.do_postgres[0].digitalocean_database_db.this",
					ID:      "abc123,my_project_db",
				},
				{
					Address: "module.resources[\"db\"].module.do_postgres[0].digitalocean_database_connection_pool.this",
					ID:      "abc123,my_project_db_pool",
				},
			},
		},
		"Hyphenated database name": {
			projectName: "my-company",
			resource: &appdef.Resource{
				Name:   "prod-db",
				Config: map[string]any{},
			},
			clusterID: "xyz789",
			want: []importAddress{
				{
					Address: "module.resources[\"prod-db\"].module.do_postgres[0].digitalocean_database_cluster.this",
					ID:      "xyz789",
				},
				{
					Address: "module.resources[\"prod-db\"].module.do_postgres[0].digitalocean_database_user.this",
					ID:      "xyz789,my_company_prod_db_admin",
				},
				{
					Address: "module.resources[\"prod-db\"].module.do_postgres[0].digitalocean_database_db.this",
					ID:      "xyz789,my_company_prod_db",
				},
				{
					Address: "module.resources[\"prod-db\"].module.do_postgres[0].digitalocean_database_connection_pool.this",
					ID:      "xyz789,my_company_prod_db_pool",
				},
			},
		},
		"With firewall rules": {
			projectName: "secure-app",
			resource: &appdef.Resource{
				Name: "db",
				Config: map[string]any{
					"allowed_ips_addr": []any{"192.168.1.1"},
				},
			},
			clusterID: "secure123",
			want: []importAddress{
				{
					Address: "module.resources[\"db\"].module.do_postgres[0].digitalocean_database_cluster.this",
					ID:      "secure123",
				},
				{
					Address: "module.resources[\"db\"].module.do_postgres[0].digitalocean_database_user.this",
					ID:      "secure123,secure_app_db_admin",
				},
				{
					Address: "module.resources[\"db\"].module.do_postgres[0].digitalocean_database_db.this",
					ID:      "secure123,secure_app_db",
				},
				{
					Address: "module.resources[\"db\"].module.do_postgres[0].digitalocean_database_connection_pool.this",
					ID:      "secure123,secure_app_db_pool",
				},
				{
					Address: "module.resources[\"db\"].module.do_postgres[0].digitalocean_database_firewall.this[0]",
					ID:      "secure123",
				},
			},
		},
	}

	for name, test := range tt {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := buildPostgresImports(test.projectName, test.resource, test.clusterID)
			assert.Equal(t, test.want, got)
		})
	}
}

func TestBuildS3Imports(t *testing.T) {
	t.Parallel()

	tt := map[string]struct {
		resource *appdef.Resource
		bucketID string
		want     []importAddress
	}{
		"Bucket without CDN": {
			resource: &appdef.Resource{
				Name:   "media",
				Config: map[string]any{},
			},
			bucketID: "bucket-123",
			want: []importAddress{
				{
					Address: "module.resources[\"media\"].module.do_bucket[0].digitalocean_spaces_bucket.this",
					ID:      "ams3,bucket-123",
				},
				{
					Address: "module.resources[\"media\"].module.do_bucket[0].digitalocean_spaces_bucket_cors_configuration.this",
					ID:      "ams3,bucket-123",
				},
			},
		},
		"Bucket with CDN": {
			resource: &appdef.Resource{
				Name: "store",
				Config: map[string]any{
					"region": "ams3",
					"cdn_id": "26226874-569d-41d9-9191-e18a09a2d134",
				},
			},
			bucketID: "s-clark-store",
			want: []importAddress{
				{
					Address: "module.resources[\"store\"].module.do_bucket[0].digitalocean_spaces_bucket.this",
					ID:      "ams3,s-clark-store",
				},
				{
					Address: "module.resources[\"store\"].module.do_bucket[0].digitalocean_spaces_bucket_cors_configuration.this",
					ID:      "ams3,s-clark-store",
				},
				{
					Address: "module.resources[\"store\"].module.do_bucket[0].digitalocean_cdn.this",
					ID:      "26226874-569d-41d9-9191-e18a09a2d134",
				},
			},
		},
		"Bucket with custom region": {
			resource: &appdef.Resource{
				Name: "nyc-bucket",
				Config: map[string]any{
					"region": "nyc3",
				},
			},
			bucketID: "my-bucket",
			want: []importAddress{
				{
					Address: "module.resources[\"nyc-bucket\"].module.do_bucket[0].digitalocean_spaces_bucket.this",
					ID:      "nyc3,my-bucket",
				},
				{
					Address: "module.resources[\"nyc-bucket\"].module.do_bucket[0].digitalocean_spaces_bucket_cors_configuration.this",
					ID:      "nyc3,my-bucket",
				},
			},
		},
	}

	for name, test := range tt {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := buildS3Imports(test.resource, test.bucketID)
			assert.Equal(t, test.want, got)
		})
	}
}

func TestBuildAppImportAddresses(t *testing.T) {
	t.Parallel()

	tt := map[string]struct {
		projectName string
		app         *appdef.App
		appID       string
		want        []importAddress
		wantErr     bool
	}{
		"DigitalOcean App Platform (container)": {
			projectName: "search-spares",
			app: &appdef.App{
				Name: "web",
				Infra: appdef.Infra{
					Provider: appdef.ResourceProviderDigitalOcean,
					Type:     "container",
					Config:   map[string]any{},
				},
			},
			appID: "app-abc123",
			want: []importAddress{
				{
					Address: "module.apps[\"web\"].module.do_app[0].digitalocean_app.this",
					ID:      "app-abc123",
				},
			},
			wantErr: false,
		},
		"DigitalOcean Droplet (vm)": {
			projectName: "test-project",
			app: &appdef.App{
				Name: "api",
				Infra: appdef.Infra{
					Provider: appdef.ResourceProviderDigitalOcean,
					Type:     "vm",
					Config:   map[string]any{},
				},
			},
			appID: "droplet-123",
			want: []importAddress{
				{
					Address: "module.apps[\"api\"].module.do_droplet[0].digitalocean_droplet.this",
					ID:      "droplet-123",
				},
			},
			wantErr: false,
		},
		"Unsupported provider": {
			projectName: "test-project",
			app: &appdef.App{
				Name: "app",
				Infra: appdef.Infra{
					Provider: "unsupported",
					Type:     "container",
					Config:   map[string]any{},
				},
			},
			appID:   "app-123",
			want:    nil,
			wantErr: true,
		},
		"Unsupported platform type": {
			projectName: "test-project",
			app: &appdef.App{
				Name: "app",
				Infra: appdef.Infra{
					Provider: appdef.ResourceProviderDigitalOcean,
					Type:     "serverless",
					Config:   map[string]any{},
				},
			},
			appID:   "app-456",
			want:    nil,
			wantErr: true,
		},
	}

	for name, test := range tt {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got, err := buildAppImportAddresses(test.projectName, test.app, test.appID)
			assert.Equal(t, test.wantErr, err != nil)

			if !test.wantErr {
				assert.Equal(t, test.want, got)
			}
		})
	}
}

func TestBuildDigitalOceanAppImports(t *testing.T) {
	t.Parallel()

	tt := map[string]struct {
		projectName string
		app         *appdef.App
		appID       string
		want        []importAddress
		wantErr     bool
	}{
		"Container app (App Platform)": {
			projectName: "my-project",
			app: &appdef.App{
				Name: "web",
				Infra: appdef.Infra{
					Type:   "container",
					Config: map[string]any{},
				},
			},
			appID: "app-123",
			want: []importAddress{
				{
					Address: "module.apps[\"web\"].module.do_app[0].digitalocean_app.this",
					ID:      "app-123",
				},
			},
			wantErr: false,
		},
		"VM app (Droplet)": {
			projectName: "my-project",
			app: &appdef.App{
				Name: "api",
				Infra: appdef.Infra{
					Type:   "vm",
					Config: map[string]any{},
				},
			},
			appID: "droplet-456",
			want: []importAddress{
				{
					Address: "module.apps[\"api\"].module.do_droplet[0].digitalocean_droplet.this",
					ID:      "droplet-456",
				},
			},
			wantErr: false,
		},
	}

	for name, test := range tt {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got, err := buildDigitalOceanAppImports(test.projectName, test.app, test.appID)
			assert.Equal(t, test.wantErr, err != nil)

			if !test.wantErr {
				assert.Equal(t, test.want, got)
			}
		})
	}
}

func TestBuildTursoImports(t *testing.T) {
	t.Parallel()

	tt := map[string]struct {
		resource   *appdef.Resource
		databaseID string
		want       []importAddress
		wantErr    bool
	}{
		"SQLite database": {
			resource: &appdef.Resource{
				Name: "db",
				Type: appdef.ResourceTypeSQLite,
				Config: map[string]any{
					"organisation": "my-org",
					"group":        "default",
				},
			},
			databaseID: "my-org/my-database",
			want: []importAddress{
				{
					Address: "module.resources[\"db\"].module.turso_database[0].turso_database.this",
					ID:      "my-org/my-database",
				},
				{
					Address: "module.resources[\"db\"].module.turso_database[0].turso_database_token.this",
					ID:      "my-org/my-database",
				},
			},
			wantErr: false,
		},
		"Unsupported resource type": {
			resource: &appdef.Resource{
				Name: "cache",
				Type: appdef.ResourceTypePostgres,
			},
			databaseID: "my-org/my-db",
			want:       nil,
			wantErr:    true,
		},
	}

	for name, test := range tt {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got, err := buildTursoImports(test.resource, test.databaseID)
			assert.Equal(t, test.wantErr, err != nil)

			if !test.wantErr {
				assert.Equal(t, test.want, got)
			}
		})
	}
}

func TestBuildTursoSQLiteImports(t *testing.T) {
	t.Parallel()

	tt := map[string]struct {
		resource   *appdef.Resource
		databaseID string
		want       []importAddress
	}{
		"Simple database name": {
			resource: &appdef.Resource{
				Name: "db",
			},
			databaseID: "my-org/my-database",
			want: []importAddress{
				{
					Address: "module.resources[\"db\"].module.turso_database[0].turso_database.this",
					ID:      "my-org/my-database",
				},
				{
					Address: "module.resources[\"db\"].module.turso_database[0].turso_database_token.this",
					ID:      "my-org/my-database",
				},
			},
		},
		"Hyphenated resource name": {
			resource: &appdef.Resource{
				Name: "prod-db",
			},
			databaseID: "acme-corp/production-db",
			want: []importAddress{
				{
					Address: "module.resources[\"prod-db\"].module.turso_database[0].turso_database.this",
					ID:      "acme-corp/production-db",
				},
				{
					Address: "module.resources[\"prod-db\"].module.turso_database[0].turso_database_token.this",
					ID:      "acme-corp/production-db",
				},
			},
		},
	}

	for name, test := range tt {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := buildTursoSQLiteImports(test.resource, test.databaseID)
			assert.Equal(t, test.want, got)
		})
	}
}

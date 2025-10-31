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
		"DigitalOcean S3 bucket": {
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
					ID:      "bucket-789",
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
		"Simple bucket": {
			resource: &appdef.Resource{
				Name: "media",
			},
			bucketID: "bucket-123",
			want: []importAddress{
				{
					Address: "module.resources[\"media\"].module.do_bucket[0].digitalocean_spaces_bucket.this",
					ID:      "bucket-123",
				},
			},
		},
		"Hyphenated bucket name": {
			resource: &appdef.Resource{
				Name: "user-uploads",
			},
			bucketID: "bucket-456",
			want: []importAddress{
				{
					Address: "module.resources[\"user-uploads\"].module.do_bucket[0].digitalocean_spaces_bucket.this",
					ID:      "bucket-456",
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

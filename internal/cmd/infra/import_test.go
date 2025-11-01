package infra

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/ainsleydev/webkit/internal/appdef"
	"github.com/ainsleydev/webkit/internal/infra"
	mockinfra "github.com/ainsleydev/webkit/internal/infra/mocks"
	"github.com/ainsleydev/webkit/pkg/env"
)

func TestImport(t *testing.T) {
	t.SkipNow()

	t.Run("Init Error", func(t *testing.T) {
		mock := mockinfra.NewMockManager(gomock.NewController(t))
		mock.EXPECT().
			Cleanup().
			Times(1)

		input, teardown := setup(t, &appdef.Definition{}, mock, true)
		defer teardown()

		err := Import(t.Context(), input)
		assert.Error(t, err)
		assert.ErrorContains(t, err, "init error")
	})

	t.Run("Import Error", func(t *testing.T) {
		mock := mockinfra.NewMockManager(gomock.NewController(t))
		mock.EXPECT().
			Import(gomock.Any(), gomock.Any()).
			Return(infra.ImportOutput{
				ImportedResources: []string{},
				Output:            "Error: resource not found",
			}, errors.New("import failed"))
		mock.EXPECT().
			Cleanup().
			Times(1)

		def := &appdef.Definition{
			Project: appdef.Project{Name: "test-project"},
			Resources: []appdef.Resource{
				{Name: "db", Type: appdef.ResourceTypePostgres, Provider: appdef.ResourceProviderDigitalOcean},
			},
		}

		input, buf, teardown := setupWithPrinter(t, def, mock, false)
		defer teardown()

		err := Import(t.Context(), input)
		assert.Error(t, err)
		assert.ErrorContains(t, err, "executing terraform import")
		assert.Contains(t, buf.String(), "Error: resource not found")
	})

	t.Run("Success Single Resource", func(t *testing.T) {
		mock := mockinfra.NewMockManager(gomock.NewController(t))
		mock.EXPECT().
			Import(gomock.Any(), infra.ImportInput{
				ResourceName: "db",
				ResourceID:   "cluster-123",
				Environment:  env.Production,
			}).
			Return(infra.ImportOutput{
				ImportedResources: []string{
					"module.resources[\"db\"].module.do_postgres[0].digitalocean_database_cluster.this",
					"module.resources[\"db\"].module.do_postgres[0].digitalocean_database_user.this",
				},
				Output: "Import successful",
			}, nil)
		mock.EXPECT().
			Cleanup().
			Times(1)

		def := &appdef.Definition{
			Project: appdef.Project{Name: "test-project"},
			Resources: []appdef.Resource{
				{Name: "db", Type: appdef.ResourceTypePostgres, Provider: appdef.ResourceProviderDigitalOcean},
			},
		}

		input, buf, teardown := setupWithPrinter(t, def, mock, false)
		defer teardown()

		err := Import(t.Context(), input)
		assert.NoError(t, err)
		assert.Contains(t, buf.String(), "Successfully imported 2 resource(s)")
		assert.Contains(t, buf.String(), "digitalocean_database_cluster.this")
		assert.Contains(t, buf.String(), "Next steps")
	})

	t.Run("Success Multiple Resources", func(t *testing.T) {
		mock := mockinfra.NewMockManager(gomock.NewController(t))
		mock.EXPECT().
			Import(gomock.Any(), gomock.Any()).
			Return(infra.ImportOutput{
				ImportedResources: []string{
					"module.resources[\"db\"].module.do_postgres[0].digitalocean_database_cluster.this",
					"module.resources[\"db\"].module.do_postgres[0].digitalocean_database_user.this",
					"module.resources[\"db\"].module.do_postgres[0].digitalocean_database_db.this",
					"module.resources[\"db\"].module.do_postgres[0].digitalocean_database_connection_pool.this",
					"module.resources[\"db\"].module.do_postgres[0].digitalocean_database_firewall.this[0]",
				},
				Output: "Import complete",
			}, nil)
		mock.EXPECT().
			Cleanup().
			Times(1)

		def := &appdef.Definition{
			Project: appdef.Project{Name: "test-project"},
			Resources: []appdef.Resource{
				{
					Name:     "db",
					Type:     appdef.ResourceTypePostgres,
					Provider: appdef.ResourceProviderDigitalOcean,
					Config: map[string]any{
						"allowed_ips_addr": []any{"192.168.1.1"},
					},
				},
			},
		}

		input, buf, teardown := setupWithPrinter(t, def, mock, false)
		defer teardown()

		err := Import(t.Context(), input)
		assert.NoError(t, err)
		assert.Contains(t, buf.String(), "Successfully imported 5 resource(s)")
		assert.Contains(t, buf.String(), "digitalocean_database_firewall.this[0]")
	})

	t.Run("Invalid Environment", func(t *testing.T) {
		mock := mockinfra.NewMockManager(gomock.NewController(t))
		mock.EXPECT().
			Cleanup().
			Times(1)

		def := &appdef.Definition{
			Project: appdef.Project{Name: "test-project"},
			Resources: []appdef.Resource{
				{Name: "db", Type: appdef.ResourceTypePostgres, Provider: appdef.ResourceProviderDigitalOcean},
			},
		}

		input, teardown := setup(t, def, mock, false)
		defer teardown()

		err := Import(t.Context(), input)
		assert.Error(t, err)
		assert.ErrorContains(t, err, "parsing environment")
	})

	t.Run("Success Import App", func(t *testing.T) {
		mock := mockinfra.NewMockManager(gomock.NewController(t))
		mock.EXPECT().
			Import(gomock.Any(), infra.ImportInput{
				ResourceName: "web",
				ResourceID:   "app-123",
				Environment:  env.Production,
				IsApp:        true,
			}).
			Return(infra.ImportOutput{
				ImportedResources: []string{
					"module.apps[\"web\"].module.do_app[0].digitalocean_app.this",
				},
				Output: "Import successful",
			}, nil)
		mock.EXPECT().
			Cleanup().
			Times(1)

		def := &appdef.Definition{
			Project: appdef.Project{Name: "test-project"},
			Apps: []appdef.App{
				{
					Name: "web",
					Infra: appdef.Infra{
						Provider: appdef.ResourceProviderDigitalOcean,
						Type:     "container",
					},
				},
			},
		}

		input, buf, teardown := setupWithPrinter(t, def, mock, false)
		defer teardown()

		err := Import(t.Context(), input)
		assert.NoError(t, err)
		assert.Contains(t, buf.String(), "Successfully imported 1 Terraform resource(s)")
		assert.Contains(t, buf.String(), "digitalocean_app.this")
		assert.Contains(t, buf.String(), "Next steps")
	})

	t.Run("Error Both App And Resource Provided", func(t *testing.T) {
		mock := mockinfra.NewMockManager(gomock.NewController(t))
		mock.EXPECT().
			Cleanup().
			Times(1)

		def := &appdef.Definition{
			Project: appdef.Project{Name: "test-project"},
		}

		input, teardown := setup(t, def, mock, false)
		defer teardown()

		err := Import(t.Context(), input)
		assert.Error(t, err)
		assert.ErrorContains(t, err, "mutually exclusive")
	})

	t.Run("Error Neither App Nor Resource Provided", func(t *testing.T) {
		mock := mockinfra.NewMockManager(gomock.NewController(t))
		mock.EXPECT().
			Cleanup().
			Times(1)

		def := &appdef.Definition{
			Project: appdef.Project{Name: "test-project"},
		}

		input, teardown := setup(t, def, mock, false)
		defer teardown()

		err := Import(t.Context(), input)
		assert.Error(t, err)
		assert.ErrorContains(t, err, "either --resource or --app must be specified")
	})
}

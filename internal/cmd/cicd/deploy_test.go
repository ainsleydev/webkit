package cicd

import (
	"path/filepath"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ainsleydev/webkit/internal/appdef"
)

func TestDeployContainerWorkflow(t *testing.T) {
	t.Parallel()

	tt := map[string]struct {
		apps       []appdef.App
		wantExists bool
	}{
		"Generates workflow for container app": {
			apps: []appdef.App{
				{
					Name: "web",
					Type: appdef.AppTypeSvelteKit,
					Build: appdef.Build{
						Dockerfile: "Dockerfile",
						Port:       3000,
					},
					Infra: appdef.Infra{
						Provider: "digitalocean",
						Type:     "container",
					},
				},
			},
			wantExists: true,
		},
		"Skips workflow when no container apps": {
			apps: []appdef.App{
				{
					Name: "api",
					Type: appdef.AppTypeGoLang,
					Build: appdef.Build{
						Dockerfile: "Dockerfile",
						Port:       8080,
					},
					Infra: appdef.Infra{
						Provider: "digitalocean",
						Type:     "vm",
					},
				},
			},
			wantExists: false,
		},
		"Skips workflow when app has no dockerfile": {
			apps: []appdef.App{
				{
					Name: "web",
					Type: appdef.AppTypeSvelteKit,
					Build: appdef.Build{
						Port: 3000,
					},
					Infra: appdef.Infra{
						Provider: "digitalocean",
						Type:     "container",
					},
				},
			},
			wantExists: false,
		},
	}

	for name, test := range tt {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			appDef := &appdef.Definition{Apps: test.apps}
			input := setup(t, afero.NewMemMapFs(), appDef)

			err := DeployContainerWorkflow(t.Context(), input)
			require.NoError(t, err)

			path := filepath.Join(workflowsPath, "deploy-container.yaml")
			exists, err := afero.Exists(input.FS, path)
			require.NoError(t, err)
			assert.Equal(t, test.wantExists, exists)

			if test.wantExists {
				content, err := afero.ReadFile(input.FS, path)
				require.NoError(t, err)
				assert.Contains(t, string(content), "Deploy Container App")
				assert.Contains(t, string(content), "workflow_dispatch")
				assert.Contains(t, string(content), "image_tag")
			}
		})
	}
}

func TestDeployVMWorkflow(t *testing.T) {
	t.Parallel()

	tt := map[string]struct {
		apps       []appdef.App
		wantExists bool
	}{
		"Generates workflow for VM app": {
			apps: []appdef.App{
				{
					Name: "cms",
					Type: appdef.AppTypePayload,
					Build: appdef.Build{
						Dockerfile: "Dockerfile",
						Port:       3000,
					},
					Infra: appdef.Infra{
						Provider: "digitalocean",
						Type:     "vm",
					},
				},
			},
			wantExists: true,
		},
		"Skips workflow when no VM apps": {
			apps: []appdef.App{
				{
					Name: "web",
					Type: appdef.AppTypeSvelteKit,
					Build: appdef.Build{
						Dockerfile: "Dockerfile",
						Port:       3000,
					},
					Infra: appdef.Infra{
						Provider: "digitalocean",
						Type:     "container",
					},
				},
			},
			wantExists: false,
		},
		"Skips workflow when app has no dockerfile": {
			apps: []appdef.App{
				{
					Name: "cms",
					Type: appdef.AppTypePayload,
					Build: appdef.Build{
						Port: 3000,
					},
					Infra: appdef.Infra{
						Provider: "digitalocean",
						Type:     "vm",
					},
				},
			},
			wantExists: false,
		},
	}

	for name, test := range tt {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			appDef := &appdef.Definition{Apps: test.apps}
			input := setup(t, afero.NewMemMapFs(), appDef)

			err := DeployVMWorkflow(t.Context(), input)
			require.NoError(t, err)

			path := filepath.Join(workflowsPath, "deploy-vm.yaml")
			exists, err := afero.Exists(input.FS, path)
			require.NoError(t, err)
			assert.Equal(t, test.wantExists, exists)

			if test.wantExists {
				content, err := afero.ReadFile(input.FS, path)
				require.NoError(t, err)
				assert.Contains(t, string(content), "Deploy VM App")
				assert.Contains(t, string(content), "workflow_dispatch")
				assert.Contains(t, string(content), "image_tag")
			}
		})
	}
}

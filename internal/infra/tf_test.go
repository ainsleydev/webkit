//go:build !race

package infra

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/hashicorp/terraform-exec/tfexec"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/ainsleydev/webkit/internal/appdef"
	"github.com/ainsleydev/webkit/internal/infra/internal/tfmocks"
	"github.com/ainsleydev/webkit/internal/manifest"
	"github.com/ainsleydev/webkit/internal/mocks"
	"github.com/ainsleydev/webkit/internal/util/executil"
	"github.com/ainsleydev/webkit/pkg/env"
)

func setup(t *testing.T, appDef *appdef.Definition) (*Terraform, func()) {
	t.Helper()

	if !executil.Exists("terraform") {
		t.Skip("terraform not found in PATH")
	}

	setupEnv(t)

	tf, err := NewTerraform(t.Context(), appDef, manifest.NewTracker())
	require.NoError(t, err)

	tf.useLocalBackend = true

	return tf, func() {
		tf.Cleanup()
		teardownEnv(t)
	}
}

func setupEnv(t *testing.T) {
	t.Helper()

	t.Setenv("DO_API_KEY", "key")
	t.Setenv("DO_SPACES_ACCESS_KEY", "access")
	t.Setenv("DO_SPACES_SECRET_KEY", "secret")
	t.Setenv("BACK_BLAZE_BUCKET", "bucket")
	t.Setenv("BACK_BLAZE_KEY_ID", "id")
	t.Setenv("BACK_BLAZE_APPLICATION_KEY", "appkey")
	t.Setenv("GITHUB_TOKEN", "token")
}

func teardownEnv(t *testing.T) {
	t.Helper()

	envVars := []string{
		"DO_API_KEY",
		"DO_SPACES_ACCESS_KEY",
		"DO_SPACES_SECRET_KEY",
		"BACK_BLAZE_BUCKET",
		"BACK_BLAZE_KEY_ID",
		"BACK_BLAZE_APPLICATION_KEY",
		"GITHUB_TOKEN",
	}

	for _, key := range envVars {
		err := os.Unsetenv(key)
		assert.NoError(t, err, fmt.Sprintf("failed to unset env var %s", key))
	}
}

func TestNewTerraform(t *testing.T) {
	if !executil.Exists("terraform") {
		t.Skip("terraform not found in PATH")
	}

	t.Run("TerraformNotInPath", func(t *testing.T) {
		t.Setenv("PATH", "/nonexistent")
		defer teardownEnv(t)

		_, err := NewTerraform(t.Context(), &appdef.Definition{}, manifest.NewTracker())
		assert.Error(t, err)
	})

	t.Run("Invalid Environment", func(t *testing.T) {
		got, err := NewTerraform(t.Context(), &appdef.Definition{}, manifest.NewTracker())
		defer teardownEnv(t)

		assert.Nil(t, got)
		assert.Error(t, err)
		assert.ErrorContains(t, err, "parsing terraform environment")
	})

	t.Run("Success", func(t *testing.T) {
		setupEnv(t)
		defer teardownEnv(t)

		got, err := NewTerraform(t.Context(), &appdef.Definition{}, manifest.NewTracker())
		require.NoError(t, err)
		assert.NotNil(t, got)
		assert.NotEmpty(t, got.env)
		assert.NotEmpty(t, got.path)
		assert.Contains(t, got.path, "terraform")
	})
}

func TestTerraform_Init(t *testing.T) {
	if !executil.Exists("terraform") {
		t.Skip("terraform not found in PATH")
	}

	t.Run("Temp Dir Error", func(t *testing.T) {
		tf, teardown := setup(t, &appdef.Definition{})
		tf.fs = afero.NewReadOnlyFs(tf.fs)
		defer teardown()

		err := tf.Init(t.Context())
		assert.Error(t, err)
		assert.ErrorContains(t, err, "creating tf tmp dir")
	})

	t.Run("Success", func(t *testing.T) {
		tf, teardown := setup(t, &appdef.Definition{})
		defer teardown()

		err := tf.Init(t.Context())
		require.NoError(t, err)

		t.Log("Verify Temp Dir")
		{
			assert.NotEmpty(t, tf.tmpDir)
			assert.DirExists(t, tf.tmpDir)
		}

		t.Log("Verify Base")
		{
			baseDir := filepath.Join(tf.tmpDir, "base")
			assert.DirExists(t, baseDir)
		}

		t.Log("Verify Terraform")
		{
			terraformDir := filepath.Join(tf.tmpDir, "base", ".terraform")
			assert.DirExists(t, terraformDir)
		}
	})

	t.Run("Init Twice", func(t *testing.T) {
		tf, teardown := setup(t, &appdef.Definition{})
		defer teardown()

		err := tf.Init(t.Context())
		require.NoError(t, err)

		firstTmpDir := tf.tmpDir

		// Init again should create new temp dir
		err = tf.Init(t.Context())
		require.NoError(t, err)
		assert.NotEqual(t, firstTmpDir, tf.tmpDir)
	})
}

func TestTerraform_Plan(t *testing.T) {
	if !executil.Exists("terraform") {
		t.Skip("terraform not found in PATH")
	}

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
			},
		},
	}

	t.Run("Plan Without Init", func(t *testing.T) {
		tf, teardown := setup(t, appDef)
		defer teardown()

		_, err := tf.Plan(t.Context(), env.Production)
		assert.Error(t, err)
		assert.ErrorContains(t, err, "terraform not initialized")
	})

	t.Run("Vars FS Error", func(t *testing.T) {
		tf, teardown := setup(t, appDef)
		defer teardown()

		err := tf.Init(t.Context())
		require.NoError(t, err)

		tf.fs = afero.NewReadOnlyFs(tf.fs)

		_, err = tf.Plan(t.Context(), env.Production)
		assert.Error(t, err)
		assert.ErrorContains(t, err, "failed to write tf vars file")
	})

	t.Run("Plan Error", func(t *testing.T) {
		tf, teardown := setup(t, appDef)
		defer teardown()

		err := tf.Init(t.Context())
		require.NoError(t, err)

		ctrl := gomock.NewController(t)
		mock := tfmocks.NewMockterraformExecutor(ctrl)
		tf.tf = mock

		mock.EXPECT().
			Plan(gomock.Any(), gomock.Any()).
			Return(false, errors.New("plan error"))

		_, err = tf.Plan(t.Context(), env.Production)
		assert.Error(t, err)
		assert.ErrorContains(t, err, "plan error")
	})

	t.Run("ShowPlanFileRaw Error", func(t *testing.T) {
		tf, teardown := setup(t, appDef)
		defer teardown()

		err := tf.Init(t.Context())
		require.NoError(t, err)

		ctrl := gomock.NewController(t)
		mock := tfmocks.NewMockterraformExecutor(ctrl)
		tf.tf = mock

		mock.EXPECT().
			Plan(gomock.Any(), gomock.Any()).
			Return(false, nil)
		mock.EXPECT().
			ShowPlanFileRaw(gomock.Any(), gomock.Any()).
			Return("", errors.New("show plan file raw error"))

		_, err = tf.Plan(t.Context(), env.Production)
		assert.Error(t, err)
		assert.ErrorContains(t, err, "show plan file raw error")
	})

	t.Run("ShowPlanFile Error", func(t *testing.T) {
		tf, teardown := setup(t, appDef)
		defer teardown()

		err := tf.Init(t.Context())
		require.NoError(t, err)

		ctrl := gomock.NewController(t)
		mock := tfmocks.NewMockterraformExecutor(ctrl)
		tf.tf = mock

		mock.EXPECT().
			Plan(gomock.Any(), gomock.Any()).
			Return(false, nil)
		mock.EXPECT().
			ShowPlanFileRaw(gomock.Any(), gomock.Any()).
			Return("", nil)
		mock.EXPECT().
			ShowPlanFile(gomock.Any(), gomock.Any()).
			Return(nil, errors.New("show plan file error"))

		_, err = tf.Plan(t.Context(), env.Production)
		assert.Error(t, err)
		assert.ErrorContains(t, err, "show plan file error")
	})

	t.Run("Success", func(t *testing.T) {
		tf, teardown := setup(t, appDef)
		defer teardown()

		err := tf.Init(t.Context())
		require.NoError(t, err)

		got, err := tf.Plan(t.Context(), env.Production)
		assert.NoError(t, err)
		assert.NotNil(t, got)
		assert.Contains(t, got.Output, "module.resources[\"db\"].module.do_postgres[0].digitalocean_database_cluster.this will be created")
		assert.NotNil(t, got.Plan)
	})
}

func TestTerraform_Apply(t *testing.T) {
	if !executil.Exists("terraform") {
		t.Skip("terraform not found in PATH")
	}

	t.Run("Apply Without Init", func(t *testing.T) {
		tf, teardown := setup(t, &appdef.Definition{})
		defer teardown()

		_, err := tf.Apply(t.Context(), env.Production)
		assert.Error(t, err)
		assert.ErrorContains(t, err, "terraform not initialized")
	})

	t.Run("Vars FS Error", func(t *testing.T) {
		tf, teardown := setup(t, &appdef.Definition{
			Resources: []appdef.Resource{
				{
					Name:     "db",
					Type:     appdef.ResourceTypePostgres,
					Provider: appdef.ResourceProviderDigitalOcean,
				},
			},
		})
		defer teardown()

		err := tf.Init(t.Context())
		require.NoError(t, err)

		tf.fs = afero.NewReadOnlyFs(tf.fs)

		_, err = tf.Apply(t.Context(), env.Production)
		assert.Error(t, err)
		assert.ErrorContains(t, err, "writing tfvars file")
	})

	// Note: If we don't stub the following functions, it will actually
	// provision resources in DigitalOcean which will probably cost
	// a shitload. Let's mock it instead.
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
			},
		},
	}

	t.Run("Success", func(t *testing.T) {
		tf, teardown := setup(t, appDef)
		defer teardown()

		err := tf.Init(t.Context())
		require.NoError(t, err)

		ctrl := gomock.NewController(t)
		mock := tfmocks.NewMockterraformExecutor(ctrl)
		tf.tf = mock

		mock.EXPECT().
			Apply(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Return(nil).Times(1)
		mock.EXPECT().SetStdout(gomock.Any()).Times(1)
		mock.EXPECT().SetStderr(gomock.Any()).Times(1)

		_, err = tf.Apply(t.Context(), env.Production)
		assert.NoError(t, err)
	})

	t.Run("Apply Failure", func(t *testing.T) {
		tf, teardown := setup(t, appDef)
		defer teardown()

		err := tf.Init(t.Context())
		require.NoError(t, err)

		ctrl := gomock.NewController(t)
		mock := tfmocks.NewMockterraformExecutor(ctrl)
		mock.EXPECT().
			Apply(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Return(errors.New("authentication failed")).Times(1)
		mock.EXPECT().SetStdout(gomock.Any()).Times(1)
		mock.EXPECT().SetStderr(gomock.Any()).Times(1)

		tf.tf = mock

		_, err = tf.Apply(t.Context(), env.Production)
		assert.Error(t, err)
		assert.ErrorContains(t, err, "terraform apply failed")
	})
}

func TestTerraform_Destroy(t *testing.T) {
	if !executil.Exists("terraform") {
		t.Skip("terraform not found in PATH")
	}

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
			},
		},
	}

	t.Run("Destroy Without Init", func(t *testing.T) {
		tf, teardown := setup(t, appDef)
		defer teardown()

		_, err := tf.Destroy(t.Context(), env.Production)
		assert.Error(t, err)
		assert.ErrorContains(t, err, "terraform not initialized")
	})

	t.Run("Vars FS Error", func(t *testing.T) {
		tf, teardown := setup(t, appDef)
		defer teardown()

		err := tf.Init(t.Context())
		require.NoError(t, err)

		tf.fs = afero.NewReadOnlyFs(tf.fs)

		_, err = tf.Destroy(t.Context(), env.Production)
		assert.Error(t, err)
		assert.ErrorContains(t, err, "writing tfvars file")
	})

	t.Run("Success", func(t *testing.T) {
		tf, teardown := setup(t, appDef)
		defer teardown()

		err := tf.Init(t.Context())
		require.NoError(t, err)

		ctrl := gomock.NewController(t)
		mock := tfmocks.NewMockterraformExecutor(ctrl)
		tf.tf = mock

		mock.EXPECT().Destroy(gomock.Any(), gomock.Any()).Return(nil).Times(1)
		mock.EXPECT().SetStdout(gomock.Any()).Times(1)
		mock.EXPECT().SetStderr(gomock.Any()).Times(1)

		_, err = tf.Destroy(t.Context(), env.Production)
		assert.NoError(t, err)
	})

	t.Run("Destroy Failure", func(t *testing.T) {
		tf, teardown := setup(t, appDef)
		defer teardown()

		err := tf.Init(t.Context())
		require.NoError(t, err)

		ctrl := gomock.NewController(t)
		mock := tfmocks.NewMockterraformExecutor(ctrl)
		mock.EXPECT().Destroy(gomock.Any(), gomock.Any()).Return(errors.New("destroy failed")).Times(1)
		mock.EXPECT().SetStdout(gomock.Any()).Times(1)
		mock.EXPECT().SetStderr(gomock.Any()).Times(1)

		tf.tf = mock

		_, err = tf.Destroy(t.Context(), env.Production)
		assert.Error(t, err)
		assert.ErrorContains(t, err, "terraform destroy failed")
	})
}

func TestTerraform_Output(t *testing.T) {
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
			},
		},
		Apps: []appdef.App{
			{
				Name: "cms",
				Type: appdef.AppTypePayload,
			},
		},
	}

	t.Run("Output Without Init", func(t *testing.T) {
		tf := &Terraform{}
		_, err := tf.Output(t.Context(), env.Production)
		assert.Error(t, err)
		assert.ErrorContains(t, err, "terraform not initialized")
	})

	t.Run("Output Error", func(t *testing.T) {
		tf, teardown := setup(t, appDef)
		defer teardown()
		require.NoError(t, tf.Init(t.Context()))

		ctrl := gomock.NewController(t)
		mock := tfmocks.NewMockterraformExecutor(ctrl)
		tf.tf = mock

		mock.EXPECT().
			Output(gomock.Any()).
			Return(nil, errors.New("output failed"))

		_, err := tf.Output(t.Context(), env.Production)
		assert.Error(t, err)
		assert.ErrorContains(t, err, "output failed")
	})

	t.Run("Malformed Resources JSON", func(t *testing.T) {
		tf, teardown := setup(t, appDef)
		defer teardown()
		require.NoError(t, tf.Init(t.Context()))

		ctrl := gomock.NewController(t)
		mock := tfmocks.NewMockterraformExecutor(ctrl)
		tf.tf = mock

		mockOutputs := map[string]tfexec.OutputMeta{
			"resources": {
				Value: []byte(`{malformed-json}`),
			},
		}

		mock.EXPECT().Output(gomock.Any()).Return(mockOutputs, nil)

		_, err := tf.Output(t.Context(), env.Production)
		assert.Error(t, err)
		assert.ErrorContains(t, err, "unmarshalling resources")
	})

	t.Run("Malformed Apps JSON", func(t *testing.T) {
		tf, teardown := setup(t, appDef)
		defer teardown()
		require.NoError(t, tf.Init(t.Context()))

		ctrl := gomock.NewController(t)
		mock := tfmocks.NewMockterraformExecutor(ctrl)
		tf.tf = mock

		mockOutputs := map[string]tfexec.OutputMeta{
			"apps": {
				Value: []byte(`{malformed-json}`),
			},
		}

		mock.EXPECT().Output(gomock.Any()).Return(mockOutputs, nil)

		_, err := tf.Output(t.Context(), env.Production)
		assert.Error(t, err)
		assert.ErrorContains(t, err, "unmarshalling apps")
	})

	t.Run("Extra Value Unmarshal Fails", func(t *testing.T) {
		tf, teardown := setup(t, appDef)
		defer teardown()
		require.NoError(t, tf.Init(t.Context()))

		ctrl := gomock.NewController(t)
		mock := tfmocks.NewMockterraformExecutor(ctrl)
		tf.tf = mock

		// Invalid JSON string to trigger fallback to string
		mockOutputs := map[string]tfexec.OutputMeta{
			"extra_key": {Value: []byte("{invalid-json}")},
		}

		mock.EXPECT().Output(gomock.Any()).Return(mockOutputs, nil)

		result, err := tf.Output(t.Context(), env.Production)
		require.NoError(t, err)
		assert.Equal(t, "{invalid-json}", result.Extra["extra_key"])
	})

	t.Run("Success", func(t *testing.T) {
		tf, teardown := setup(t, appDef)
		defer teardown()
		require.NoError(t, tf.Init(t.Context()))

		ctrl := gomock.NewController(t)
		mock := tfmocks.NewMockterraformExecutor(ctrl)
		tf.tf = mock

		mockOutputs := map[string]tfexec.OutputMeta{
			"resources": {
				Value: []byte(`{"store":{"bucket_name":"my-store"}}`),
			},
			"apps": {
				Value: []byte(`{"web-app":{"app_url":"https://web.example.com"}}`),
			},
			"project_name": {
				Value: []byte(`"my-project"`),
			},
		}

		mock.EXPECT().Output(gomock.Any()).Return(mockOutputs, nil)

		result, err := tf.Output(t.Context(), env.Production)
		require.NoError(t, err)

		t.Log("Resources")
		{
			assert.NotNil(t, result.Resources["store"])
			assert.Equal(t, "my-store", result.Resources["store"]["bucket_name"])
		}

		t.Log("Apps")
		{
			assert.NotNil(t, result.Apps["web-app"])
			assert.Equal(t, "https://web.example.com", result.Apps["web-app"]["app_url"])
		}

		t.Log("Extra")
		{
			assert.Equal(t, "my-project", result.Extra["project_name"])
		}
	})
}

func TestTerraform_Cleanup(t *testing.T) {
	if !executil.Exists("terraform") {
		t.Skip("terraform not found in PATH")
	}

	t.Run("Removes Dir", func(t *testing.T) {
		tf, teardown := setup(t, &appdef.Definition{})
		defer teardown()

		err := tf.Init(t.Context())
		require.NoError(t, err)

		tmpDir := tf.tmpDir
		assert.DirExists(t, tmpDir)

		tf.Cleanup()

		assert.NoDirExists(t, tmpDir)
	})

	t.Run("NoOpIfNotInitialized", func(t *testing.T) {
		tf := &Terraform{
			path: "/usr/bin/terraform",
		}

		assert.NotPanics(t, func() {
			tf.Cleanup()
		})
	})
}

func TestTerraform_DetermineImageTag(t *testing.T) {
	t.Run("Uses GITHUB_SHA when set", func(t *testing.T) {
		appDef := &appdef.Definition{
			Project: appdef.Project{
				Repo: appdef.GitHubRepo{
					Owner: "test-owner",
					Name:  "test-repo",
				},
			},
		}

		ctrl := gomock.NewController(t)
		mockClient := mocks.NewGHClient(ctrl)

		tf := &Terraform{
			appDef:   appDef,
			ghClient: mockClient,
		}

		t.Setenv("GITHUB_SHA", "abc123def456")

		tag := tf.determineImageTag(context.Background(), "web")
		assert.Equal(t, "sha-abc123def456", tag)
	})

	t.Run("Queries GHCR when GITHUB_SHA not set", func(t *testing.T) {
		t.Setenv("GITHUB_SHA", "")

		appDef := &appdef.Definition{
			Project: appdef.Project{
				Repo: appdef.GitHubRepo{
					Owner: "test-owner",
					Name:  "test-repo",
				},
			},
		}

		ctrl := gomock.NewController(t)
		mockClient := mocks.NewGHClient(ctrl)
		mockClient.EXPECT().
			GetLatestSHATag(gomock.Any(), "test-owner", "test-repo", "web").
			Return("sha-xyz789")

		tf := &Terraform{
			appDef:   appDef,
			ghClient: mockClient,
		}

		tag := tf.determineImageTag(context.Background(), "web")
		assert.Equal(t, "sha-xyz789", tag)
	})

	t.Run("Falls back to latest when GHCR returns empty", func(t *testing.T) {
		t.Setenv("GITHUB_SHA", "")

		appDef := &appdef.Definition{
			Project: appdef.Project{
				Repo: appdef.GitHubRepo{
					Owner: "test-owner",
					Name:  "test-repo",
				},
			},
		}

		ctrl := gomock.NewController(t)
		mockClient := mocks.NewGHClient(ctrl)
		mockClient.EXPECT().
			GetLatestSHATag(gomock.Any(), "test-owner", "test-repo", "web").
			Return("")

		tf := &Terraform{
			appDef:   appDef,
			ghClient: mockClient,
		}

		tag := tf.determineImageTag(context.Background(), "web")
		assert.Equal(t, "latest", tag)
	})
}

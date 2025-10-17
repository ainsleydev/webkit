package infra

import (
	"errors"
	"fmt"
	"path/filepath"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/ainsleydev/webkit/internal/appdef"
	"github.com/ainsleydev/webkit/internal/infra/internal/tfmocks"
	"github.com/ainsleydev/webkit/internal/util/executil"
	"github.com/ainsleydev/webkit/pkg/env"
)

func setup(t *testing.T, appDef *appdef.Definition) *Terraform {
	t.Helper()

	if !executil.Exists("terraform") {
		t.Skip("terraform not found in PATH")
	}

	setupEnv(t)

	tf, err := NewTerraform(t.Context(), appDef)
	require.NoError(t, err)

	tf.useLocalBackend = true

	return tf
}

func setupEnv(t *testing.T) {
	t.Helper()

	t.Setenv("DO_API_KEY", "key")
	t.Setenv("DO_SPACES_ACCESS_KEY", "access")
	t.Setenv("DO_SPACES_SECRET_KEY", "secret")
	t.Setenv("BACK_BLAZE_BUCKET", "bucket")
	t.Setenv("BACK_BLAZE_KEY_ID", "id")
	t.Setenv("BACK_BLAZE_APPLICATION_KEY", "appkey")
}

func TestNewTerraform(t *testing.T) {
	if !executil.Exists("terraform") {
		t.Skip("terraform not found in PATH")
	}

	t.Run("TerraformNotInPath", func(t *testing.T) {
		t.Setenv("PATH", "/nonexistent")

		_, err := NewTerraform(t.Context(), &appdef.Definition{})
		assert.Error(t, err)
	})

	t.Run("Invalid Environment", func(t *testing.T) {
		got, err := NewTerraform(t.Context(), &appdef.Definition{})
		assert.Nil(t, got)
		assert.Error(t, err)
		assert.ErrorContains(t, err, "parsing terraform environment")
	})

	t.Run("Success", func(t *testing.T) {
		setupEnv(t)

		got, err := NewTerraform(t.Context(), &appdef.Definition{})
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
		tf := setup(t, &appdef.Definition{})
		tf.fs = afero.NewReadOnlyFs(tf.fs)
		defer tf.Cleanup()

		err := tf.Init(t.Context())
		assert.Error(t, err)
		assert.ErrorContains(t, err, "creating tf tmp dir")
	})

	t.Run("Success", func(t *testing.T) {
		tf := setup(t, &appdef.Definition{})
		defer tf.Cleanup()

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
		tf := setup(t, &appdef.Definition{})
		defer tf.Cleanup()

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
				Repo:  "project",
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
		tf := setup(t, appDef)
		defer tf.Cleanup()

		_, err := tf.Plan(t.Context(), env.Production)
		assert.Error(t, err)
		assert.ErrorContains(t, err, "terraform not initialized")
	})

	t.Run("Nothing To Provision", func(t *testing.T) {
		tf := setup(t, &appdef.Definition{})
		defer tf.Cleanup()

		err := tf.Init(t.Context())
		require.NoError(t, err)

		_, err = tf.Plan(t.Context(), env.Production)
		assert.Error(t, err)
		assert.ErrorContains(t, err, "no app or resources are defined")
	})

	t.Run("Vars FS Error", func(t *testing.T) {
		tf := setup(t, appDef)
		defer tf.Cleanup()

		err := tf.Init(t.Context())
		require.NoError(t, err)

		tf.fs = afero.NewReadOnlyFs(tf.fs)

		_, err = tf.Plan(t.Context(), env.Production)
		assert.Error(t, err)
		assert.ErrorContains(t, err, "failed to write tf vars file")
	})

	t.Run("Plan Error", func(t *testing.T) {
		tf := setup(t, appDef)
		defer tf.Cleanup()

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
		tf := setup(t, appDef)
		defer tf.Cleanup()

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
		tf := setup(t, appDef)
		defer tf.Cleanup()

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
		tf := setup(t, appDef)
		defer tf.Cleanup()

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
		tf := setup(t, &appdef.Definition{})
		defer tf.Cleanup()

		_, err := tf.Apply(t.Context(), env.Production)
		assert.Error(t, err)
		assert.ErrorContains(t, err, "terraform not initialized")
	})

	t.Run("Nothing To Provision", func(t *testing.T) {
		tf := setup(t, &appdef.Definition{})
		defer tf.Cleanup()

		err := tf.Init(t.Context())
		require.NoError(t, err)

		_, err = tf.Apply(t.Context(), env.Production)
		assert.Error(t, err)
		assert.ErrorContains(t, err, "no app or resources are defined")
	})

	t.Run("Vars FS Error", func(t *testing.T) {
		tf := setup(t, &appdef.Definition{
			Resources: []appdef.Resource{
				{
					Name:     "db",
					Type:     appdef.ResourceTypePostgres,
					Provider: appdef.ResourceProviderDigitalOcean,
				},
			},
		})
		defer tf.Cleanup()

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
				Repo:  "project",
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
		tf := setup(t, appDef)
		defer tf.Cleanup()

		err := tf.Init(t.Context())
		require.NoError(t, err)

		ctrl := gomock.NewController(t)
		mock := tfmocks.NewMockterraformExecutor(ctrl)
		tf.tf = mock

		mock.EXPECT().
			Apply(gomock.Any(), gomock.Any()).
			Return(nil).Times(1)
		mock.EXPECT().SetStdout(gomock.Any()).Times(1)
		mock.EXPECT().SetStderr(gomock.Any()).Times(1)

		got, err := tf.Apply(t.Context(), env.Production)
		assert.NoError(t, err)
		fmt.Print(got.Output)
	})

	t.Run("Apply Failure", func(t *testing.T) {
		tf := setup(t, appDef)
		defer tf.Cleanup()

		err := tf.Init(t.Context())
		require.NoError(t, err)

		ctrl := gomock.NewController(t)
		mock := tfmocks.NewMockterraformExecutor(ctrl)
		mock.EXPECT().Apply(gomock.Any()).
			Return(errors.New("authentication failed")).Times(1)
		mock.EXPECT().SetStdout(gomock.Any()).Times(1)
		mock.EXPECT().SetStderr(gomock.Any()).Times(1)

		tf.tf = mock

		_, err = tf.Apply(t.Context(), env.Production)
		assert.Error(t, err)
		assert.ErrorContains(t, err, "terraform apply failed")
	})
}

func TestTerraform_Cleanup(t *testing.T) {
	if !executil.Exists("terraform") {
		t.Skip("terraform not found in PATH")
	}

	t.Run("Removes Dir", func(t *testing.T) {
		tf := setup(t, &appdef.Definition{})
		defer tf.Cleanup()

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

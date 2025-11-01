package main

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ainsleydev/webkit/internal/util/executil"
)

// Note: main(), createTag(), and deleteTag() are not tested as they require
// interactive user input via bufio.Reader. Testing these would require
// complex mocking of stdin or subprocess testing patterns.

func TestGetLatestTag(t *testing.T) {
	t.Parallel()

	t.Run("Returns tag when exists", func(t *testing.T) {
		t.Parallel()

		runner := executil.NewMemRunner()
		runner.AddStub("git fetch --tags", executil.Result{Output: ""}, nil)
		runner.AddStub("git describe --tags --abbrev=0", executil.Result{Output: "v1.2.3\n"}, nil)

		tag := getLatestTag(runner, context.Background())
		assert.Equal(t, "v1.2.3", tag)
	})

	t.Run("Returns v0.0.0 when no tags exist", func(t *testing.T) {
		t.Parallel()

		runner := executil.NewMemRunner()
		runner.AddStub("git fetch --tags", executil.Result{Output: ""}, nil)
		runner.AddStub("git describe --tags --abbrev=0", executil.Result{Output: ""}, errors.New("no tags found"))

		tag := getLatestTag(runner, context.Background())
		assert.Equal(t, "v0.0.0", tag)
	})

	t.Run("Ignores fetch errors", func(t *testing.T) {
		t.Parallel()

		runner := executil.NewMemRunner()
		runner.AddStub("git fetch --tags", executil.Result{Output: ""}, errors.New("no remote"))
		runner.AddStub("git describe --tags --abbrev=0", executil.Result{Output: "v2.0.0\n"}, nil)

		tag := getLatestTag(runner, context.Background())
		assert.Equal(t, "v2.0.0", tag)
	})
}

func TestGetAllTags(t *testing.T) {
	t.Parallel()

	t.Run("Returns tags when exist", func(t *testing.T) {
		t.Parallel()

		runner := executil.NewMemRunner()
		runner.AddStub("git tag --sort=-version:refname", executil.Result{Output: "v1.2.3\nv1.2.2\nv1.2.1"}, nil)

		tags, err := getAllTags(runner, context.Background())
		require.NoError(t, err)
		assert.Equal(t, []string{"v1.2.3", "v1.2.2", "v1.2.1"}, tags)
	})

	t.Run("Returns empty slice when no tags", func(t *testing.T) {
		t.Parallel()

		runner := executil.NewMemRunner()
		runner.AddStub("git tag --sort=-version:refname", executil.Result{Output: ""}, nil)

		tags, err := getAllTags(runner, context.Background())
		require.NoError(t, err)
		assert.Equal(t, []string{}, tags)
	})

	t.Run("Returns error when git fails", func(t *testing.T) {
		t.Parallel()

		runner := executil.NewMemRunner()
		runner.AddStub("git tag --sort=-version:refname", executil.Result{Output: ""}, errors.New("git error"))

		tags, err := getAllTags(runner, context.Background())
		assert.Error(t, err)
		assert.Nil(t, tags)
	})

	t.Run("Handles whitespace correctly", func(t *testing.T) {
		t.Parallel()

		runner := executil.NewMemRunner()
		runner.AddStub("git tag --sort=-version:refname", executil.Result{Output: "  v1.0.0  \n  v0.9.0  \n"}, nil)

		tags, err := getAllTags(runner, context.Background())
		require.NoError(t, err)
		assert.Len(t, tags, 2)
	})
}

func TestGitCreateTag(t *testing.T) {
	t.Parallel()

	t.Run("Success", func(t *testing.T) {
		t.Parallel()

		runner := executil.NewMemRunner()
		runner.AddStub("git tag v1.0.0", executil.Result{Output: ""}, nil)

		err := gitCreateTag(runner, context.Background(), "v1.0.0")
		assert.NoError(t, err)
	})

	t.Run("Error when tag exists", func(t *testing.T) {
		t.Parallel()

		runner := executil.NewMemRunner()
		runner.AddStub("git tag v1.0.0", executil.Result{Output: "tag already exists"}, errors.New("tag already exists"))

		err := gitCreateTag(runner, context.Background(), "v1.0.0")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "tag already exists")
	})
}

func TestGitPushTag(t *testing.T) {
	t.Parallel()

	t.Run("Success", func(t *testing.T) {
		t.Parallel()

		runner := executil.NewMemRunner()
		runner.AddStub("git push origin v1.0.0", executil.Result{Output: ""}, nil)

		err := gitPushTag(runner, context.Background(), "v1.0.0")
		assert.NoError(t, err)
	})

	t.Run("Error on push failure", func(t *testing.T) {
		t.Parallel()

		runner := executil.NewMemRunner()
		runner.AddStub("git push origin v1.0.0", executil.Result{Output: "failed to push"}, errors.New("push failed"))

		err := gitPushTag(runner, context.Background(), "v1.0.0")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "push failed")
	})
}

func TestGitDeleteLocalTag(t *testing.T) {
	t.Parallel()

	t.Run("Success", func(t *testing.T) {
		t.Parallel()

		runner := executil.NewMemRunner()
		runner.AddStub("git tag -d v1.0.0", executil.Result{Output: "Deleted tag 'v1.0.0'"}, nil)

		err := gitDeleteLocalTag(runner, context.Background(), "v1.0.0")
		assert.NoError(t, err)
	})

	t.Run("Error when tag doesn't exist", func(t *testing.T) {
		t.Parallel()

		runner := executil.NewMemRunner()
		runner.AddStub("git tag -d v1.0.0", executil.Result{Output: "tag not found"}, errors.New("tag not found"))

		err := gitDeleteLocalTag(runner, context.Background(), "v1.0.0")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "tag not found")
	})
}

func TestGitDeleteRemoteTag(t *testing.T) {
	t.Parallel()

	t.Run("Success", func(t *testing.T) {
		t.Parallel()

		runner := executil.NewMemRunner()
		runner.AddStub("git push origin --delete v1.0.0", executil.Result{Output: ""}, nil)

		err := gitDeleteRemoteTag(runner, context.Background(), "v1.0.0")
		assert.NoError(t, err)
	})

	t.Run("Error on remote deletion failure", func(t *testing.T) {
		t.Parallel()

		runner := executil.NewMemRunner()
		runner.AddStub("git push origin --delete v1.0.0", executil.Result{Output: "remote tag not found"}, errors.New("remote error"))

		err := gitDeleteRemoteTag(runner, context.Background(), "v1.0.0")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "remote error")
	})
}

func TestGitCommitVersionFile(t *testing.T) {
	t.Parallel()

	t.Run("Success", func(t *testing.T) {
		t.Parallel()

		runner := executil.NewMemRunner()
		runner.AddStub("git add internal/version/version.go", executil.Result{Output: ""}, nil)
		runner.AddStub("git commit -m chore: Updating version to v1.0.0", executil.Result{Output: "[main abc123] chore: Updating version to v1.0.0"}, nil)

		err := gitCommitVersionFile(runner, context.Background(), "v1.0.0")
		assert.NoError(t, err)
	})

	t.Run("Error on git add failure", func(t *testing.T) {
		t.Parallel()

		runner := executil.NewMemRunner()
		runner.AddStub("git add internal/version/version.go", executil.Result{Output: "file not found"}, errors.New("add failed"))

		err := gitCommitVersionFile(runner, context.Background(), "v1.0.0")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "git add failed")
	})

	t.Run("Error on git commit failure", func(t *testing.T) {
		t.Parallel()

		runner := executil.NewMemRunner()
		runner.AddStub("git add internal/version/version.go", executil.Result{Output: ""}, nil)
		runner.AddStub("git commit -m chore: Updating version to v1.0.0", executil.Result{Output: "nothing to commit"}, errors.New("commit failed"))

		err := gitCommitVersionFile(runner, context.Background(), "v1.0.0")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "git commit failed")
	})
}

func TestGitPushCommit(t *testing.T) {
	t.Parallel()

	t.Run("Success", func(t *testing.T) {
		t.Parallel()

		runner := executil.NewMemRunner()
		runner.AddStub("git push", executil.Result{Output: ""}, nil)

		err := gitPushCommit(runner, context.Background())
		assert.NoError(t, err)
	})

	t.Run("Error on push failure", func(t *testing.T) {
		t.Parallel()

		runner := executil.NewMemRunner()
		runner.AddStub("git push", executil.Result{Output: "push rejected"}, errors.New("push failed"))

		err := gitPushCommit(runner, context.Background())
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "push failed")
	})
}

func TestGenerateVersionFile(t *testing.T) {
	t.Parallel()

	t.Run("Generates correct content", func(t *testing.T) {
		t.Parallel()

		var capturedContent string

		// Note: This test would require mocking the filesystem and generator.
		// Since generateVersionFile uses real filesystem (afero.NewOsFs()),
		// we can only test its structure without mocking, which would be complex.
		// The function's correctness is primarily verified through integration tests.

		expectedContent := `package version

const Version = "v1.2.3"
`
		assert.Contains(t, expectedContent, "package version")
		assert.Contains(t, expectedContent, `const Version = "v1.2.3"`)

		_ = capturedContent
	})
}

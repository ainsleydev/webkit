package infra

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ainsleydev/webkit/internal/appdef"
	"github.com/ainsleydev/webkit/internal/printer"
)

func TestOutputText(t *testing.T) {
	t.Parallel()

	t.Run("Skip with reason", func(t *testing.T) {
		t.Parallel()

		analysis := appdef.ChangeAnalysis{
			Skip:   true,
			Reason: "app.json unchanged",
		}
		buf := &bytes.Buffer{}
		p := printer.New(buf)

		err := outputText(analysis, p, false)

		output := buf.String()
		assert.Contains(t, output, "Decision: app.json unchanged")
		assert.Contains(t, output, "Terraform apply can be skipped")
		require.NoError(t, err)
	})

	t.Run("Skip silent mode", func(t *testing.T) {
		t.Parallel()

		analysis := appdef.ChangeAnalysis{
			Skip:   true,
			Reason: "app.json unchanged",
		}
		buf := &bytes.Buffer{}
		p := printer.New(buf)

		err := outputText(analysis, p, true)

		output := buf.String()
		assert.Empty(t, output)
		require.NoError(t, err)
	})

	t.Run("Terraform needed", func(t *testing.T) {
		t.Parallel()

		analysis := appdef.ChangeAnalysis{
			Skip:   false,
			Reason: "Infrastructure config changed",
		}
		buf := &bytes.Buffer{}
		p := printer.New(buf)

		err := outputText(analysis, p, false)

		output := buf.String()
		assert.Contains(t, output, "Decision: Infrastructure config changed")
		assert.Contains(t, output, "Terraform apply is needed")
		require.Error(t, err)
	})

	t.Run("With changed apps", func(t *testing.T) {
		t.Parallel()

		analysis := appdef.ChangeAnalysis{
			Skip:   false,
			Reason: "DigitalOcean container app env values changed",
			ChangedApps: []appdef.AppChange{
				{Name: "web", EnvChanged: true},
				{Name: "api", InfraChanged: true},
			},
		}
		buf := &bytes.Buffer{}
		p := printer.New(buf)

		err := outputText(analysis, p, false)

		output := buf.String()
		assert.Contains(t, output, "Changed apps:")
		assert.Contains(t, output, "web: env changed")
		assert.Contains(t, output, "api: infrastructure changed")
		require.Error(t, err)
	})
}

func TestOutputJSON(t *testing.T) {
	t.Parallel()

	t.Run("Skip", func(t *testing.T) {
		t.Parallel()

		analysis := appdef.ChangeAnalysis{
			Skip:   true,
			Reason: "app.json unchanged",
		}
		buf := &bytes.Buffer{}

		err := outputJSON(analysis, buf)

		var result appdef.ChangeAnalysis
		jsonErr := json.Unmarshal(buf.Bytes(), &result)
		require.NoError(t, jsonErr)
		assert.Equal(t, analysis.Skip, result.Skip)
		assert.Equal(t, analysis.Reason, result.Reason)
		require.NoError(t, err)
	})

	t.Run("Terraform needed", func(t *testing.T) {
		t.Parallel()

		analysis := appdef.ChangeAnalysis{
			Skip:   false,
			Reason: "Infrastructure config changed",
			ChangedApps: []appdef.AppChange{
				{Name: "web", EnvChanged: true},
			},
		}
		buf := &bytes.Buffer{}

		err := outputJSON(analysis, buf)

		var result appdef.ChangeAnalysis
		jsonErr := json.Unmarshal(buf.Bytes(), &result)
		require.NoError(t, jsonErr)
		assert.Equal(t, analysis.Skip, result.Skip)
		assert.Equal(t, analysis.Reason, result.Reason)
		require.Error(t, err)
	})
}

func TestOutputGitHub(t *testing.T) {
	t.Parallel()

	t.Run("Skip", func(t *testing.T) {
		t.Parallel()

		analysis := appdef.ChangeAnalysis{
			Skip:   true,
			Reason: "app.json unchanged",
		}
		buf := &bytes.Buffer{}

		err := outputGitHub(analysis, buf)

		output := buf.String()
		assert.Contains(t, output, "skip_terraform=true")
		assert.Contains(t, output, "reason=app.json unchanged")
		assert.Contains(t, output, "::notice::app.json unchanged")
		require.NoError(t, err)
	})

	t.Run("Terraform needed", func(t *testing.T) {
		t.Parallel()

		analysis := appdef.ChangeAnalysis{
			Skip:   false,
			Reason: "Infrastructure config changed",
		}
		buf := &bytes.Buffer{}

		err := outputGitHub(analysis, buf)

		output := buf.String()
		assert.Contains(t, output, "skip_terraform=false")
		assert.Contains(t, output, "reason=Infrastructure config changed")
		assert.Contains(t, output, "::notice::Infrastructure config changed")
		require.Error(t, err)
	})
}

package infra

import (
	"bytes"
	"encoding/json"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ainsleydev/webkit/internal/appdef"
	"github.com/ainsleydev/webkit/internal/printer"
)

func TestOutputText(t *testing.T) {
	t.Parallel()

	tt := map[string]struct {
		analysis   appdef.ChangeAnalysis
		silent     bool
		wantOutput []string
		wantExit   bool
	}{
		"Skip with reason": {
			analysis: appdef.ChangeAnalysis{
				Skip:   true,
				Reason: "app.json unchanged",
			},
			silent:     false,
			wantOutput: []string{"Decision: app.json unchanged", "Terraform apply can be skipped"},
			wantExit:   false,
		},
		"Skip silent mode": {
			analysis: appdef.ChangeAnalysis{
				Skip:   true,
				Reason: "app.json unchanged",
			},
			silent:     true,
			wantOutput: []string{},
			wantExit:   false,
		},
		"Terraform needed": {
			analysis: appdef.ChangeAnalysis{
				Skip:   false,
				Reason: "Infrastructure config changed",
			},
			silent:     false,
			wantOutput: []string{"Decision: Infrastructure config changed", "Terraform apply is needed"},
			wantExit:   true,
		},
		"With changed apps": {
			analysis: appdef.ChangeAnalysis{
				Skip:   false,
				Reason: "DigitalOcean container app env values changed",
				ChangedApps: []appdef.AppChange{
					{Name: "web", EnvChanged: true},
					{Name: "api", InfraChanged: true},
				},
			},
			silent:     false,
			wantOutput: []string{"Changed apps:", "web: env changed", "api: infrastructure changed"},
			wantExit:   true,
		},
	}

	for name, test := range tt {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			buf := &bytes.Buffer{}
			p := printer.New(buf)

			err := outputText(test.analysis, p, test.silent)

			output := buf.String()
			for _, want := range test.wantOutput {
				assert.Contains(t, output, want)
			}

			if test.wantExit {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestOutputJSON(t *testing.T) {
	tt := map[string]struct {
		analysis appdef.ChangeAnalysis
		wantExit bool
	}{
		"Skip": {
			analysis: appdef.ChangeAnalysis{
				Skip:   true,
				Reason: "app.json unchanged",
			},
			wantExit: false,
		},
		"Terraform needed": {
			analysis: appdef.ChangeAnalysis{
				Skip:   false,
				Reason: "Infrastructure config changed",
				ChangedApps: []appdef.AppChange{
					{Name: "web", EnvChanged: true},
				},
			},
			wantExit: true,
		},
	}

	for name, test := range tt {
		t.Run(name, func(t *testing.T) {
			// Capture stdout.
			oldStdout := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w

			err := outputJSON(test.analysis)

			// Restore stdout.
			w.Close()
			os.Stdout = oldStdout

			// Read captured output.
			buf := make([]byte, 4096)
			n, _ := r.Read(buf)
			output := string(buf[:n])

			// Verify JSON is valid and contains expected data.
			var result appdef.ChangeAnalysis
			jsonErr := json.Unmarshal([]byte(output), &result)
			require.NoError(t, jsonErr)
			assert.Equal(t, test.analysis.Skip, result.Skip)
			assert.Equal(t, test.analysis.Reason, result.Reason)

			if test.wantExit {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestOutputGitHub(t *testing.T) {
	tt := map[string]struct {
		analysis   appdef.ChangeAnalysis
		wantOutput []string
		wantExit   bool
	}{
		"Skip": {
			analysis: appdef.ChangeAnalysis{
				Skip:   true,
				Reason: "app.json unchanged",
			},
			wantOutput: []string{
				"skip_terraform=true",
				"reason=app.json unchanged",
				"::notice::app.json unchanged",
			},
			wantExit: false,
		},
		"Terraform needed": {
			analysis: appdef.ChangeAnalysis{
				Skip:   false,
				Reason: "Infrastructure config changed",
			},
			wantOutput: []string{
				"skip_terraform=false",
				"reason=Infrastructure config changed",
				"::notice::Infrastructure config changed",
			},
			wantExit: true,
		},
	}

	for name, test := range tt {
		t.Run(name, func(t *testing.T) {
			// Capture stdout.
			oldStdout := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w

			err := outputGitHub(test.analysis)

			// Restore stdout.
			w.Close()
			os.Stdout = oldStdout

			// Read captured output.
			buf := make([]byte, 4096)
			n, _ := r.Read(buf)
			output := string(buf[:n])

			for _, want := range test.wantOutput {
				assert.Contains(t, output, want)
			}

			if test.wantExit {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

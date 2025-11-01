package templates

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGithubVariable(t *testing.T) {
	t.Parallel()

	tt := map[string]struct {
		input string
		want  string
	}{
		"Simple variable":  {input: "github.sha", want: "${{ github.sha }}"},
		"Runner variable":  {input: "runner.os", want: "${{ runner.os }}"},
		"Empty string":     {input: "", want: "${{  }}"},
		"Complex variable": {input: "github.event.pull_request.number", want: "${{ github.event.pull_request.number }}"},
	}

	for name, test := range tt {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := githubVariable(test.input)
			assert.Equal(t, test.want, got)
		})
	}
}

func TestGithubSecret(t *testing.T) {
	t.Parallel()

	tt := map[string]struct {
		input string
		want  string
	}{
		"Simple secret":   {input: "API_KEY", want: "${{ secrets.API_KEY }}"},
		"Empty secret":    {input: "", want: "${{ secrets. }}"},
		"Complex secret":  {input: "DEPLOYMENT_TOKEN", want: "${{ secrets.DEPLOYMENT_TOKEN }}"},
		"Lowercase name":  {input: "my_secret", want: "${{ secrets.my_secret }}"},
		"Mixed case name": {input: "My_Secret", want: "${{ secrets.My_Secret }}"},
	}

	for name, test := range tt {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := githubSecret(test.input)
			assert.Equal(t, test.want, got)
		})
	}
}

func TestGithubInput(t *testing.T) {
	t.Parallel()

	tt := map[string]struct {
		input string
		want  string
	}{
		"Simple input":    {input: "environment", want: "${{ inputs.environment }}"},
		"Empty input":     {input: "", want: "${{ inputs. }}"},
		"Complex input":   {input: "deploy-branch", want: "${{ inputs.deploy-branch }}"},
		"Uppercase input": {input: "BRANCH_NAME", want: "${{ inputs.BRANCH_NAME }}"},
	}

	for name, test := range tt {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := githubInput(test.input)
			assert.Equal(t, test.want, got)
		})
	}
}

func TestGithubEnv(t *testing.T) {
	t.Parallel()

	tt := map[string]struct {
		input string
		want  string
	}{
		"Simple env":       {input: "NODE_ENV", want: "${{ env.NODE_ENV }}"},
		"Empty env":        {input: "", want: "${{ env. }}"},
		"Complex env":      {input: "DEPLOYMENT_URL", want: "${{ env.DEPLOYMENT_URL }}"},
		"Lowercase env":    {input: "path", want: "${{ env.path }}"},
		"Mixed case env":   {input: "Api_Key", want: "${{ env.Api_Key }}"},
		"With underscores": {input: "MY_CUSTOM_VAR", want: "${{ env.MY_CUSTOM_VAR }}"},
	}

	for name, test := range tt {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := githubEnv(test.input)
			assert.Equal(t, test.want, got)
		})
	}
}

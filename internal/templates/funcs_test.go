package templates

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGithubExpression(t *testing.T) {
	t.Parallel()

	tt := map[string]struct {
		input string
		want  string
	}{
		"Simple expression":  {input: "github.sha", want: "${{ github.sha }}"},
		"Runner expression":  {input: "runner.os", want: "${{ runner.os }}"},
		"Empty string":       {input: "", want: "${{  }}"},
		"Complex expression": {input: "github.event.pull_request.number", want: "${{ github.event.pull_request.number }}"},
		"Step output":        {input: "steps.version.outputs.version", want: "${{ steps.version.outputs.version }}"},
	}

	for name, test := range tt {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := githubExpression(test.input)
			assert.Equal(t, test.want, got)
		})
	}
}

func TestGithubVariable(t *testing.T) {
	t.Parallel()

	tt := map[string]struct {
		input string
		want  string
	}{
		"Simple variable":    {input: "MY_VAR", want: "${{ vars.MY_VAR }}"},
		"Empty variable":     {input: "", want: "${{ vars. }}"},
		"Complex variable":   {input: "PROD_DATABASE_URL", want: "${{ vars.PROD_DATABASE_URL }}"},
		"Lowercase variable": {input: "my_var", want: "${{ vars.my_var }}"},
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

func TestPrettyConfigKey(t *testing.T) {
	t.Parallel()

	tt := map[string]struct {
		input string
		want  string
	}{
		"Single word":            {input: "location", want: "Location"},
		"Snake case":             {input: "server_type", want: "Server Type"},
		"Multiple underscores":   {input: "foo_bar_baz", want: "Foo Bar Baz"},
		"Empty string":           {input: "", want: ""},
		"Already capitalised":    {input: "Server_Type", want: "Server Type"},
		"All caps":               {input: "SSH_KEYS", want: "SSH KEYS"},
		"Mixed case":             {input: "Api_Key", want: "Api Key"},
		"Common config keys":     {input: "mount_point", want: "Mount Point"},
		"IP range":               {input: "ip_range", want: "Ip Range"},
		"Database size":          {input: "size", want: "Size"},
		"Long key":               {input: "very_long_config_key_name", want: "Very Long Config Key Name"},
		"Numeric":                {input: "port_8080", want: "Port 8080"},
		"With numbers":           {input: "server_1_type", want: "Server 1 Type"},
		"Single character":       {input: "a", want: "A"},
		"Two words":              {input: "user_name", want: "User Name"},
		"Hetzner location":       {input: "nbg1", want: "Nbg1"},
		"DigitalOcean droplet":   {input: "droplet_size", want: "Droplet Size"},
	}

	for name, test := range tt {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := prettyConfigKey(test.input)
			assert.Equal(t, test.want, got)
		})
	}
}

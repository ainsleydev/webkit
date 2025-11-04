//go:build integration

package ansible

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ainsleydev/webkit/internal/util/executil"
)

func TestAnsibleVMDeployment(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	if !executil.Exists("docker") {
		t.Skip("docker not found in PATH; skipping integration test")
	}

	if !executil.Exists("ansible-playbook") {
		t.Skip("ansible-playbook not found in PATH; skipping integration test")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	containerName := fmt.Sprintf("webkit-vm-test-%d", time.Now().Unix())
	t.Logf("Creating test container: %s", containerName)

	cmd := exec.CommandContext(ctx, "docker", "run", "-d",
		"--name", containerName,
		"--privileged",
		"ubuntu:22.04",
		"sleep", "3600")

	err := cmd.Run()
	require.NoError(t, err, "Failed to create test container")

	t.Cleanup(func() {
		cleanupCmd := exec.Command("docker", "rm", "-f", containerName)
		_ = cleanupCmd.Run()
	})

	t.Log("Installing Python in test container")
	installPython := exec.CommandContext(ctx, "docker", "exec", containerName,
		"bash", "-c", "apt-get update -qq && apt-get install -y -qq python3 python3-pip curl sudo")
	err = installPython.Run()
	require.NoError(t, err, "Failed to install Python in container")

	tmpDir := t.TempDir()
	fs := afero.NewOsFs()

	t.Log("Creating test fixtures")
	fixturesDir := filepath.Join(tmpDir, "fixtures")
	err = fs.MkdirAll(filepath.Join(fixturesDir, "resources", "secrets"), 0755)
	require.NoError(t, err)

	appJSON := `{
  "webkit_version": "v0.0.15",
  "apps": [
    {
      "name": "test-app",
      "path": "./test-app",
      "build": {
        "dockerfile": "Dockerfile",
        "port": 3000
      },
      "infra": {
        "provider": "digitalocean",
        "type": "vm",
        "config": {
          "domain": "test.example.com"
        }
      },
      "env": {
        "production": {
          "FOO": {
            "source": "value",
            "value": "bar"
          }
        }
      }
    }
  ]
}`

	err = afero.WriteFile(fs, filepath.Join(fixturesDir, "app.json"), []byte(appJSON), 0644)
	require.NoError(t, err)

	secretsYAML := "TEST_SECRET: test_value\n"
	err = afero.WriteFile(fs, filepath.Join(fixturesDir, "resources", "secrets", "production.yaml"), []byte(secretsYAML), 0644)
	require.NoError(t, err)

	inventoryPath := filepath.Join(tmpDir, "inventory.ini")
	inventory := fmt.Sprintf("[all]\n%s ansible_connection=docker ansible_user=root\n", containerName)
	err = afero.WriteFile(fs, inventoryPath, []byte(inventory), 0644)
	require.NoError(t, err)

	ansibleDir, err := filepath.Abs("../../..")
	require.NoError(t, err)

	playbookPath := filepath.Join(tmpDir, "test_playbook.yaml")
	playbook := fmt.Sprintf(`---
- name: Test webkit VM deployment
  hosts: all
  become: true
  vars:
    webkit_version: latest
    age_secret_key: "AGE-SECRET-KEY-1TEST123456789"
    app_definition_path: "%s"
    secrets_path: "%s"
    app_name: test-app
    environment: production
    git_sha: abc123
    github_user: testuser
    github_token: testtoken
    domain: test.example.com
    docker_image: test-image
    docker_image_tag: sha-abc123
    docker_port: 3000

  roles:
    - webkit

  post_tasks:
    - name: Verify webkit is installed
      command: /usr/local/bin/webkit version
      register: webkit_version_output
      changed_when: false

    - name: Display webkit version
      debug:
        msg: "{{ webkit_version_output.stdout }}"

    - name: Verify app.json was copied
      stat:
        path: /etc/webkit/app.json
      register: app_json_stat
      failed_when: not app_json_stat.stat.exists

    - name: Verify AGE key was written
      stat:
        path: /root/.config/webkit/age.key
      register: age_key_stat
      failed_when: not age_key_stat.stat.exists or age_key_stat.stat.mode != '0600'

    - name: Create test env directory
      file:
        path: /opt/test-app
        state: directory
        mode: '0755'

    - name: Generate env file using webkit
      command: >
        /usr/local/bin/webkit env generate
        --app test-app
        --environment production
        --output /opt/test-app/.env
      args:
        chdir: /etc/webkit
      environment:
        SOPS_AGE_KEY: "{{ age_secret_key }}"
      register: webkit_generate_output

    - name: Display webkit generate output
      debug:
        msg: "{{ webkit_generate_output.stdout }}"

    - name: Verify env file was created
      stat:
        path: /opt/test-app/.env
      register: env_file_stat
      failed_when: not env_file_stat.stat.exists

    - name: Read env file contents
      slurp:
        src: /opt/test-app/.env
      register: env_file_contents

    - name: Verify env file contains expected vars
      assert:
        that:
          - "'FOO=' in env_file_contents.content | b64decode"
        fail_msg: "Env file does not contain expected variables"
`, filepath.Join(fixturesDir, "app.json"), filepath.Join(fixturesDir, "resources", "secrets"))

	err = afero.WriteFile(fs, playbookPath, []byte(playbook), 0644)
	require.NoError(t, err)

	t.Log("Running Ansible playbook")
	ansibleCmd := exec.CommandContext(ctx, "ansible-playbook",
		"-i", inventoryPath,
		playbookPath,
		"-v")
	ansibleCmd.Env = append(os.Environ(), fmt.Sprintf("ANSIBLE_ROLES_PATH=%s/roles", ansibleDir))
	ansibleCmd.Dir = tmpDir

	output, err := ansibleCmd.CombinedOutput()
	t.Logf("Ansible output:\n%s", string(output))

	require.NoError(t, err, "Ansible playbook failed")

	assert.Contains(t, string(output), "PLAY RECAP")
	assert.NotContains(t, strings.ToLower(string(output)), "failed=")

	t.Log("Verifying webkit installation in container")
	{
		verifyCmd := exec.CommandContext(ctx, "docker", "exec", containerName, "/usr/local/bin/webkit", "version")
		verifyOutput, err := verifyCmd.CombinedOutput()
		require.NoError(t, err, "webkit version command failed")
		t.Logf("WebKit version: %s", string(verifyOutput))
	}

	t.Log("Verifying env file in container")
	{
		catCmd := exec.CommandContext(ctx, "docker", "exec", containerName, "cat", "/opt/test-app/.env")
		envOutput, err := catCmd.CombinedOutput()
		require.NoError(t, err, "Failed to read env file from container")
		assert.Contains(t, string(envOutput), "FOO=", "Env file should contain FOO variable")
		t.Logf("Env file contents:\n%s", string(envOutput))
	}
}

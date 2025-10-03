package appdef

import (
	"io"

	"github.com/goccy/go-json"
	"github.com/spf13/afero"
)

const (
	// JsonFileName defines the file name of the app manifest,
	// that should appear in the root of each webkit dir.
	JsonFileName = "app.json"
)

type (
	Definition struct {
		WebkitVersion string     `json:"webkit_version"`
		Project       Project    `json:"project"`
		Shared        Shared     `json:"shared"`
		Resources     []Resource `json:"resources"`
		Apps          []App      `json:"apps"`
	}
	Project struct {
		Name        string `json:"name"`
		Title       string `json:"title"`
		Description string `json:"description"`
		Repo        string `json:"repo"`
	}
	Resource struct {
		Name     string         `json:"name"`
		Type     string         `json:"type"`
		Provider string         `json:"provider"`
		Config   map[string]any `json:"config"` // Conforms to Terraform
		Outputs  []string       `json:"outputs"`
	}
	App struct {
		Name        string   `json:"name"`
		Type        AppType  `json:"type"`
		Description string   `json:"description,omitempty"`
		Path        string   `json:"path"`
		Build       Build    `json:"build"`
		Infra       Infra    `json:"infra"`
		Env         Env      `json:"env"`
		DependsOn   []string `json:"depends_on,omitempty"`
	}
	Build struct {
		Dockerfile string `json:"dockerfile"`
	}
	Infra struct {
		Provider string `json:"provider"`
		Type     string `json:"type"`
		Config   struct {
			Size          string   `json:"size,omitempty"`
			Region        string   `json:"region"`
			Domain        string   `json:"domain"`
			SshKeys       []string `json:"ssh_keys,omitempty"`
			InstanceCount int      `json:"instance_count,omitempty"`
			EnvFromShared bool     `json:"env_from_shared,omitempty"`
		} `json:"config"`
	}
	Shared struct {
		Env Env `json:"env"`
	}
	Env struct {
		Dev        []EnvValue `json:"dev"`
		Staging    []EnvValue `json:"staging"`
		Production []EnvValue `json:"production"`
	}
	EnvValue struct {
		Key   string `json:"key"`
		Type  string `json:"type"`
		From  string `json:"from,omitempty"`
		Value string `json:"value,omitempty"`
	}
)

type (
	AppType string
)

const (
	AppTypeSvelteKit AppType = "svelte-kit"
	AppTypeGoLang    AppType = "golang"
	AppTypePayload   AppType = "payload-cms"
)

// String implements fmt.Stringer on the AppType.
func (a AppType) String() string {
	return string(a)
}

func Read(root afero.Fs) (*Definition, error) {
	file, err := root.Open(JsonFileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	// TODO: Apply defaults and return validation errors if the user has fucked it.
	var def Definition
	if err := json.Unmarshal(data, &def); err != nil {
		return nil, err
	}

	return &def, nil
}

func (d Definition) GithubLabels() []string {
	labels := []string{"webkit"}

	for _, v := range d.Apps {
		labels = append(labels, v.Type.String())
	}

	return labels
}

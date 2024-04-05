package webkit

type Schema struct {
	Name        string `yaml:"name"`
	Title       string `yaml:"title"`
	Description string `yaml:"description"`
	Version     string `yaml:"version"`
	GitHubURL   string `yaml:"github_url"`
	CMS         struct {
		Enabled bool `yaml:"enabled"`
	} `yaml:"payload"`
}

func ParseSchema() (*Schema, error) {
	return nil, nil
}

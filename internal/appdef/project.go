package appdef

type (
	// Project defines root metadata about the project such as business
	// names and descriptions.
	Project struct {
		Name        string     `json:"name"`
		Title       string     `json:"title"`
		Description string     `json:"description"`
		Repo        GitHubRepo `json:"repo"`
	}
	// GitHubRepo defines the metadata for GitHub repositories within
	// an app declaration.
	GitHubRepo struct {
		Owner string `json:"owner"`
		Name  string `json:"name"`
	}
)

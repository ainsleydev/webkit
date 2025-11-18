package appdef

type (
	// Project defines root metadata about the project such as business
	// names and descriptions. This information is used throughout webkit
	// for identification, documentation generation, and CI/CD configuration.
	Project struct {
		Name        string     `json:"name" validate:"required,lowercase,alphanumdash" description:"Unique identifier for the project (lowercase, hyphenated)"`
		Title       string     `json:"title" validate:"required" description:"Human-readable project name displayed in documentation and UIs"`
		Description string     `json:"description" validate:"required,max=200" description:"Brief description of the project's purpose and functionality"`
		Repo        GitHubRepo `json:"repo" validate:"required" description:"GitHub repository information for the project"`
	}
	// GitHubRepo defines the metadata for GitHub repositories.
	// This information is used for CI/CD integration, secrets management,
	// and linking documentation to the source repository.
	GitHubRepo struct {
		Owner string `json:"owner" validate:"required" description:"GitHub username or organisation that owns the repository"`
		Name  string `json:"name" validate:"required" description:"Repository name on GitHub"`
	}
)

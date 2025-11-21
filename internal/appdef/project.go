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
		Brand       Brand      `json:"brand,omitempty" description:"Branding configuration for the project including colours and logo"`
	}
	// GitHubRepo defines the metadata for GitHub repositories.
	// This information is used for CI/CD integration, secrets management,
	// and linking documentation to the source repository.
	GitHubRepo struct {
		Owner string `json:"owner" validate:"required" description:"GitHub username or organisation that owns the repository"`
		Name  string `json:"name" validate:"required" description:"Repository name on GitHub"`
	}
	// Brand defines branding information for the project.
	// This information is used for monitor tags, status pages, and other
	// visual representations of the project.
	Brand struct {
		PrimaryColour   string `json:"primaryColour,omitempty" validate:"omitempty,hexcolor" description:"Primary brand colour in hex format (e.g., #3B82F6)"`
		SecondaryColour string `json:"secondaryColour,omitempty" validate:"omitempty,hexcolor" description:"Secondary brand colour in hex format (e.g., #10B981)"`
		LogoURL         string `json:"logoUrl,omitempty" validate:"omitempty,url" description:"URL to the project's logo image"`
	}
)

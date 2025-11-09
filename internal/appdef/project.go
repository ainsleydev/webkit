package appdef

type (
	// Project defines root metadata about the project such as business
	// names and descriptions. This information is used throughout webkit
	// for identification, documentation generation, and CI/CD configuration.
	Project struct {
		Name        string     `json:"name" required:"true" pattern:"^[a-z][a-z0-9-]*$" description:"Unique identifier for the project (lowercase, hyphenated)"`
		Title       string     `json:"title" required:"true" description:"Human-readable project name displayed in documentation and UIs"`
		Description string     `json:"description" required:"true" maxLength:"200" description:"Brief description of the project's purpose and functionality"`
		Repo        GitHubRepo `json:"repo" required:"true" description:"GitHub repository information for the project"`
	}
	// GitHubRepo defines the metadata for GitHub repositories.
	// This information is used for CI/CD integration, secrets management,
	// and linking documentation to the source repository.
	GitHubRepo struct {
		Owner string `json:"owner" required:"true" description:"GitHub username or organisation that owns the repository"`
		Name  string `json:"name" required:"true" description:"Repository name on GitHub"`
	}
	// Notifications defines alert and notification settings for the project.
	// Provider-agnostic configuration supporting Slack, Discord, and other webhook-based services.
	Notifications struct {
		WebhookURL string `json:"webhook_url,omitzero" format:"uri" description:"Webhook URL for sending notifications (e.g., Slack webhook)"`
	}
)

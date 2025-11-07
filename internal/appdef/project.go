package appdef

type (
	// Project defines root metadata about the project such as business
	// names and descriptions.
	Project struct {
		Name          string        `json:"name"`
		Title         string        `json:"title"`
		Description   string        `json:"description"`
		Repo          GitHubRepo    `json:"repo"`
		Notifications Notifications `json:"notifications"`
	}
	// GitHubRepo defines the metadata for GitHub repositories within
	// an app declaration.
	GitHubRepo struct {
		Owner string `json:"owner"`
		Name  string `json:"name"`
	}
	// Notifications defines alert and notification settings for the project.
	// Provider-agnostic configuration supporting Slack, Discord, and other webhook-based services.
	Notifications struct {
		WebhookURL string `json:"webhook_url"`
	}
)

package github

// RepoSettings represents the root configuration structure for GitHub
// repository settings.
//
// See: https://github.com/apps/settings
// And: https://github.com/repository-settings/app/blob/master/docs/configuration.md
type RepoSettings struct {
	Repository    Repository     `yaml:"repository"`
	Labels        []Label        `yaml:"labels,omitempty"`
	Milestones    []Milestone    `yaml:"milestones,omitempty"`
	Collaborators []Collaborator `yaml:"collaborators,omitempty"`
	Teams         []Team         `yaml:"teams,omitempty"`
	Branches      []Branch       `yaml:"branches,omitempty"`
}

// Repository contains all repository-level settings
type Repository struct {
	Name                         string   `yaml:"name,omitempty"`
	Description                  string   `yaml:"description,omitempty"`
	Homepage                     string   `yaml:"homepage,omitempty"`
	Topics                       []string `yaml:"topics,omitempty"`
	Private                      bool     `yaml:"private,omitempty"`
	HasIssues                    bool     `yaml:"has_issues,omitempty"`
	HasProjects                  bool     `yaml:"has_projects,omitempty"`
	HasWiki                      bool     `yaml:"has_wiki,omitempty"`
	HasDownloads                 bool     `yaml:"has_downloads,omitempty"`
	DefaultBranch                string   `yaml:"default_branch,omitempty"`
	AllowSquashMerge             bool     `yaml:"allow_squash_merge,omitempty"`
	AllowMergeCommit             bool     `yaml:"allow_merge_commit,omitempty"`
	AllowRebaseMerge             bool     `yaml:"allow_rebase_merge,omitempty"`
	DeleteBranchOnMerge          bool     `yaml:"delete_branch_on_merge,omitempty"`
	EnableAutomatedSecurityFixes bool     `yaml:"enable_automated_security_fixes,omitempty"`
	EnableVulnerabilityAlerts    bool     `yaml:"enable_vulnerability_alerts,omitempty"`
}

// Label represents a GitHub label configuration
type Label struct {
	Name        string `yaml:"name"`
	Color       string `yaml:"color"`
	Description string `yaml:"description,omitempty"`
	NewName     string `yaml:"new_name,omitempty"`
}

// Milestone represents a GitHub milestone configuration
type Milestone struct {
	Title       string `yaml:"title"`
	Description string `yaml:"description,omitempty"`
	State       string `yaml:"state"`
}

// Collaborator represents a repository collaborator
type Collaborator struct {
	Username   string `yaml:"username"`
	Permission string `yaml:"permission"`
}

// Team represents a team with repository access
type Team struct {
	Name       string `yaml:"name"`
	Permission string `yaml:"permission"`
}

// Branch represents branch protection settings
type Branch struct {
	Name       string            `yaml:"name"`
	Protection *BranchProtection `yaml:"protection,omitempty"`
}

// BranchProtection contains branch protection rules
type BranchProtection struct {
	RequiredPullRequestReviews *RequiredPullRequestReviews `yaml:"required_pull_request_reviews,omitempty"`
	RequiredStatusChecks       *RequiredStatusChecks       `yaml:"required_status_checks,omitempty"`
	EnforceAdmins              bool                        `yaml:"enforce_admins,omitempty"`
	RequiredLinearHistory      bool                        `yaml:"required_linear_history,omitempty"`
	Restrictions               *Restrictions               `yaml:"restrictions,omitempty"`
}

// RequiredPullRequestReviews contains PR review requirements
type RequiredPullRequestReviews struct {
	RequiredApprovingReviewCount int                    `yaml:"required_approving_review_count"`
	DismissStaleReviews          bool                   `yaml:"dismiss_stale_reviews,omitempty"`
	RequireCodeOwnerReviews      bool                   `yaml:"require_code_owner_reviews,omitempty"`
	DismissalRestrictions        *DismissalRestrictions `yaml:"dismissal_restrictions,omitempty"`
}

// DismissalRestrictions specifies who can dismiss reviews
type DismissalRestrictions struct {
	Users []string `yaml:"users,omitempty"`
	Teams []string `yaml:"teams,omitempty"`
}

// RequiredStatusChecks contains status check requirements
type RequiredStatusChecks struct {
	Strict   bool     `yaml:"strict"`
	Contexts []string `yaml:"contexts,omitempty"`
}

// Restrictions contains push restrictions
type Restrictions struct {
	Apps  []string `yaml:"apps,omitempty"`
	Users []string `yaml:"users,omitempty"`
	Teams []string `yaml:"teams,omitempty"`
}

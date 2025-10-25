package files

import (
	"context"
	"strings"

	"github.com/ainsleydev/webkit/internal/cmd/internal/schemas/github"
	"github.com/ainsleydev/webkit/internal/cmdtools"
	"github.com/ainsleydev/webkit/internal/manifest"
	"github.com/ainsleydev/webkit/internal/scaffold"
	"github.com/ainsleydev/webkit/internal/templates"
)

var gitSettingsTemplates = map[string]string{
	".gitignore":              ".gitignore",
	".github/dependabot.yaml": ".github/dependabot.yaml.tmpl",
}

// GitSettings scaffolds the repo settings and ignore files.
//
// TODO: Stale, Pull Request Template.
func GitSettings(_ context.Context, input cmdtools.CommandInput) error {
	app := input.AppDef()

	for file, template := range gitSettingsTemplates {
		err := input.Generator().Template(file,
			templates.MustLoadTemplate(template),
			app,
			scaffold.WithTracking(manifest.SourceProject()),
		)
		if err != nil {
			return err
		}
	}

	return input.Generator().YAML(".github/settings.yml",
		repoSettings(input),
		scaffold.WithTracking(manifest.SourceProject()),
	)
}

func repoSettings(input cmdtools.CommandInput) github.RepoSettings {
	return github.RepoSettings{
		Repository: github.Repository{
			AllowMergeCommit:    false,
			DeleteBranchOnMerge: true,
			// settings.yml files doesn't like list objects for some reason.
			Topics:       strings.Join(input.AppDef().GithubLabels(), ", "),
			Private:      true,
			HasWiki:      false,
			HasDownloads: false,
		},
		Teams: []github.Team{
			{Name: "core", Permission: "admin"},
		},
		Branches: []github.Branch{
			{
				Name: "main",
				Protection: &github.BranchProtection{
					// Add this when growing.
					//RequiredPullRequestReviews: &github.RequiredPullRequestReviews{
					//	DismissStaleReviews:          true,
					//	RequireCodeOwnerReviews:      false,
					//	RequiredApprovingReviewCount: 1,
					//},
					Restrictions: &github.Restrictions{
						// Only this team can merge
						Teams: []string{"core"},
						Users: make([]string, 0),
						Apps:  make([]string, 0),
					},
					EnforceAdmins:  false,
					AllowForcePush: false,
					AllowDeletions: false,
				},
			},
		},
	}
}

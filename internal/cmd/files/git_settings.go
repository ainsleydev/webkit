package files

import (
	"context"

	"github.com/ainsleydev/webkit/internal/cmd/internal/cmdtools"
	"github.com/ainsleydev/webkit/internal/cmd/internal/schemas/github"
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
	gen := scaffold.New(input.FS, input.Manifest)
	app := input.AppDef()

	for file, template := range gitSettingsTemplates {
		err := gen.Template(file,
			templates.MustLoadTemplate(template),
			app,
			scaffold.WithTracking("files.GitSettings", "project:root", true),
		)
		if err != nil {
			return err
		}
	}

	return gen.YAML(".github/settings.yaml", repoSettings(input))
}

func repoSettings(input cmdtools.CommandInput) github.RepoSettings {
	return github.RepoSettings{
		Repository: github.Repository{
			AllowMergeCommit:    false,
			DeleteBranchOnMerge: true,
			Topics:              input.AppDef().GithubLabels(),
			Private:             true,
			HasWiki:             false,
			HasDownloads:        false,
		},
		Branches: []github.Branch{
			{
				Name: "main",
				Protection: &github.BranchProtection{
					RequiredPullRequestReviews: &github.RequiredPullRequestReviews{
						DismissStaleReviews:          true,
						RequireCodeOwnerReviews:      true,
						RequiredApprovingReviewCount: 1,
					},
					EnforceAdmins: true,
				},
			},
		},
	}
}

package operations

import (
	"context"

	"github.com/ainsleydev/webkit/internal/cmd/internal/cmdtools"
	"github.com/ainsleydev/webkit/internal/cmd/internal/schemas/github"
	"github.com/ainsleydev/webkit/internal/scaffold"
	"github.com/ainsleydev/webkit/internal/templates"
)

// CreateGitSettings scaffolds the repo settings and ignore files.
func CreateGitSettings(_ context.Context, input cmdtools.CommandInput) error {
	gen := scaffold.New(input.FS)

	err := gen.Template(".gitignore", templates.MustLoadTemplate(".gitignore"), nil)
	if err != nil {
		return err
	}

	// TODO:
	// Dependabot

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

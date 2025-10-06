package cmd

import (
	"context"

	"github.com/urfave/cli/v3"

	"github.com/ainsleydev/webkit/internal/cmd/internal"
	"github.com/ainsleydev/webkit/internal/github"
	"github.com/ainsleydev/webkit/internal/scaffold"
	"github.com/ainsleydev/webkit/internal/templates"
)

var createGithubSettingsCmd = &cli.Command{
	Name:   "git",
	Action: cmdtools.WrapCommand(createGitSettings),
}

func createGitSettings(_ context.Context, input cmdtools.CommandInput) error {
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

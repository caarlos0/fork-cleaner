package forkcleaner

import (
	"context"
	"fmt"
	"net/http"

	"github.com/google/go-github/v33/github"
)

type RepositoryWithDetails struct {
	Repo              *github.Repository
	Issues            []*github.Issue
	Commits           *github.CommitsComparison
	ParentRepoMissing bool
	ParentRepoDMCA    bool
}

// FindAllForks lists all the forks for the current user.
func FindAllForks(
	ctx context.Context,
	client *github.Client,
) ([]*RepositoryWithDetails, error) {
	var forks []*RepositoryWithDetails
	repos, err := getAllRepos(ctx, client)
	if err != nil {
		return forks, nil
	}
	for _, repo := range repos {
		if !repo.GetFork() {
			continue
		}

		var login = repo.GetOwner().GetLogin()
		var name = repo.GetName()

		// Get repository as List omits parent information.
		frepo, _, err := client.Repositories.Get(ctx, login, name)
		if err != nil {
			return forks, fmt.Errorf("failed to get repository: %s: %w", repo.GetFullName(), err)
		}

		var parent = frepo.GetParent()

		// get fork's Issues
		issues, _, err := client.Issues.ListByRepo(
			ctx,
			parent.GetOwner().GetLogin(),
			parent.GetName(),
			&github.IssueListByRepoOptions{
				ListOptions: github.ListOptions{
					PerPage: pageSize,
				},
				Creator: login,
			},
		)
		if err != nil {
			return forks, fmt.Errorf("failed to get repository's issues: %s: %w", repo.GetFullName(), err)
		}

		// compare Commits with upstream
		commits, resp, err := client.Repositories.CompareCommits(
			ctx,
			parent.GetOwner().GetLogin(),
			parent.GetName(),
			parent.GetDefaultBranch(),
			fmt.Sprintf("%s:%s", login, repo.GetDefaultBranch()),
		)
		if err != nil {
			return forks, fmt.Errorf("failed to compare repository with upstream: %s: %w", repo.GetFullName(), err)
		}

		forks = append(forks, &RepositoryWithDetails{
			Repo:              frepo,
			Issues:            issues,
			Commits:           commits,
			ParentRepoMissing: resp.StatusCode == http.StatusNotFound,
			ParentRepoDMCA:    resp.StatusCode == http.StatusUnavailableForLegalReasons,
		})
	}
	return forks, nil
}

func getAllRepos(
	ctx context.Context,
	client *github.Client,
) ([]*github.Repository, error) {
	var allRepos []*github.Repository
	var opts = &github.RepositoryListOptions{
		ListOptions: github.ListOptions{PerPage: pageSize},
		Affiliation: "owner",
	}
	for {
		repos, resp, err := client.Repositories.List(ctx, "", opts)
		if err != nil {
			return allRepos, err
		}
		allRepos = append(allRepos, repos...)
		if resp.NextPage == 0 {
			break
		}
		opts.ListOptions.Page = resp.NextPage
	}
	return allRepos, nil
}

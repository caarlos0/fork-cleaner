// Package forkcleaner provides functions to find and remove unused forks.
package forkcleaner

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/google/go-github/v33/github"
)

const pageSize = 100

type RepositoryWithDetails struct {
	Name               string
	RepoURL            string
	Private            bool
	ParentDeleted      bool
	ParentDMCATakeDown bool
	Forks              int
	Stars              int
	OpenPRs            int
	CommitsAhead       int
	LastUpdate         time.Time
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
	for _, r := range repos {
		if !r.GetFork() {
			continue
		}

		var login = r.GetOwner().GetLogin()
		var name = r.GetName()

		// Get repository as List omits parent information.
		repo, _, err := client.Repositories.Get(ctx, login, name)
		if err != nil {
			return forks, fmt.Errorf("failed to get repository: %s: %w", repo.GetFullName(), err)
		}

		var parent = repo.GetParent()

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

		forks = append(forks, buildDetails(repo, issues, commits, resp.StatusCode))
	}
	return forks, nil
}

func buildDetails(repo *github.Repository, issues []*github.Issue, commits *github.CommitsComparison, code int) *RepositoryWithDetails {
	var openPrs int
	for _, issue := range issues {
		if issue.IsPullRequest() {
			openPrs++
		}
	}
	return &RepositoryWithDetails{
		Name:               repo.GetFullName(),
		RepoURL:            repo.GetURL(),
		Private:            repo.GetPrivate(),
		ParentDeleted:      code == http.StatusNotFound,
		ParentDMCATakeDown: code == http.StatusUnavailableForLegalReasons,
		Forks:              repo.GetForksCount(),
		Stars:              repo.GetStargazersCount(),
		OpenPRs:            openPrs,
		CommitsAhead:       commits.GetAheadBy(),
		LastUpdate:         repo.GetUpdatedAt().Time,
	}
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

// Delete delete the given list of forks.
func Delete(
	ctx context.Context,
	client *github.Client,
	deletions []*github.Repository,
) error {
	for _, repo := range deletions {
		_, err := client.Repositories.Delete(ctx, *repo.Owner.Login, *repo.Name)
		if err != nil {
			return err
		}
	}
	return nil
}

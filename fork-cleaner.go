// Package forkcleaner provides functions to find and remove unused forks.
package forkcleaner

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/google/go-github/v50/github"
)

const pageSize = 100

type RepositoryWithDetails struct {
	Name               string
	ParentName         string
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
	login string,
) ([]*RepositoryWithDetails, error) {
	var forks []*RepositoryWithDetails
	repos, err := getAllRepos(ctx, client, login)
	if err != nil {
		return forks, nil
	}
	for _, r := range repos {
		if !r.GetFork() {
			continue
		}

		login := r.GetOwner().GetLogin()
		name := r.GetName()

		// Get repository as List omits parent information.
		repo, resp, err := client.Repositories.Get(ctx, login, name)
		switch resp.StatusCode {
		case http.StatusForbidden:
			// no access, ignore
			continue
		case http.StatusNotFound, http.StatusUnavailableForLegalReasons:
			forks = append(forks, buildDetails(r, nil, nil, resp.StatusCode))
			continue
		}

		if err != nil {
			return forks, fmt.Errorf("failed to get repository: %s: %w", repo.GetFullName(), err)
		}

		parent := repo.GetParent()

		// get parent's Issues
		issues, err := getIssues(ctx, client, login, parent)
		if err != nil {
			return forks, fmt.Errorf("failed to get repository's issues: %s: %w", parent.GetFullName(), err)
		}

		// compare Commits with parent
		commits, resp, err := client.Repositories.CompareCommits(
			ctx,
			parent.GetOwner().GetLogin(),
			parent.GetName(),
			parent.GetDefaultBranch(),
			fmt.Sprintf("%s:%s", login, repo.GetDefaultBranch()),
			&github.ListOptions{},
		)
		if err != nil && resp.StatusCode != 404 {
			return forks, fmt.Errorf("failed to compare repository with parent: %s: %w", repo.GetFullName(), err)
		}

		forks = append(forks, buildDetails(repo, issues, commits, resp.StatusCode))
	}
	return forks, nil
}

func buildDetails(repo *github.Repository, issues []*github.Issue, commits *github.CommitsComparison, code int) *RepositoryWithDetails {
	var openPrs, aheadBy int
	for _, issue := range issues {
		if issue.IsPullRequest() {
			openPrs++
		}
	}
	if commits != nil {
		aheadBy = commits.GetAheadBy()
	}
	return &RepositoryWithDetails{
		Name:               repo.GetFullName(),
		ParentName:         repo.GetParent().GetFullName(),
		RepoURL:            repo.GetURL(),
		Private:            repo.GetPrivate(),
		ParentDeleted:      code == http.StatusNotFound,
		ParentDMCATakeDown: code == http.StatusUnavailableForLegalReasons,
		Forks:              repo.GetForksCount(),
		Stars:              repo.GetStargazersCount(),
		OpenPRs:            openPrs,
		CommitsAhead:       aheadBy,
		LastUpdate:         repo.GetUpdatedAt().Time,
	}
}

func getAllRepos(
	ctx context.Context,
	client *github.Client,
	login string,
) ([]*github.Repository, error) {
	var allRepos []*github.Repository
	opts := &github.RepositoryListOptions{
		ListOptions: github.ListOptions{PerPage: pageSize},
		Affiliation: "owner",
	}
	for {
		repos, resp, err := client.Repositories.List(ctx, login, opts)
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

func getIssues(
	ctx context.Context,
	client *github.Client,
	login string,
	repo *github.Repository,
) ([]*github.Issue, error) {
	var allIssues []*github.Issue
	opts := &github.IssueListByRepoOptions{
		ListOptions: github.ListOptions{
			PerPage: pageSize,
		},
		Creator: login,
	}
	for {
		issues, resp, err := client.Issues.ListByRepo(
			ctx,
			repo.GetOwner().GetLogin(),
			repo.GetName(),
			opts,
		)
		if err != nil {
			return allIssues, err
		}
		allIssues = append(allIssues, issues...)
		if resp.NextPage == 0 {
			break
		}
		opts.ListOptions.Page = resp.NextPage
	}
	return allIssues, nil
}

// Delete delete the given list of forks.
func Delete(
	ctx context.Context,
	client *github.Client,
	deletions []*RepositoryWithDetails,
) error {
	for _, repo := range deletions {
		parts := strings.Split(repo.Name, "/")
		log.Println("deleting repository:", repo.Name)
		_, err := client.Repositories.Delete(ctx, parts[0], parts[1])
		if err != nil {
			return fmt.Errorf("couldn't delete repository: %s: %w", repo.Name, err)
		}
	}
	return nil
}

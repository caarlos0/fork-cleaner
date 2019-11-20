// Package forkcleaner provides functions to find and remove unused forks.
package forkcleaner

import (
	"context"
	"fmt"
	"time"

	"github.com/google/go-github/github"
)

// Filter applied to the repositories list
type Filter struct {
	Blacklist           []string
	IncludePrivate      bool
	Since               time.Duration
	ExcludeCommitsAhead bool
}

// Delete delete the given list of forks
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

// Find list the forks from a given owner that could be deleted
func Find(
	ctx context.Context,
	client *github.Client,
	filter Filter,
) ([]*github.Repository, error) {
	lopt := github.ListOptions{PerPage: 100}
	ropt := &github.RepositoryListOptions{
		ListOptions: lopt,
		Affiliation: "owner",
	}
	iopt := &github.IssueListByRepoOptions{
		ListOptions: lopt,
	}
	var deletions []*github.Repository
	var login string
	for {
		repos, resp, err := client.Repositories.List(ctx, "", ropt)
		if err != nil {
			return deletions, err
		}
		for _, repo := range repos {
			if login == "" {
				login = repo.GetOwner().GetLogin()
				iopt.Creator = login
			}
			if !repo.GetFork() {
				continue
			}
			rn := repo.GetName()
			// Get repository as List omits parent information.
			repo, _, err = client.Repositories.Get(ctx, login, rn)
			if err != nil {
				return deletions, err
			}
			parent := repo.GetParent()
			po := parent.GetOwner().GetLogin()
			pn := parent.GetName()
			issues, _, err := client.Issues.ListByRepo(ctx, po, pn, iopt)
			if err != nil {
				return deletions, err
			}
			commits, _, compareErr := client.Repositories.CompareCommits(ctx, po, rn, *parent.DefaultBranch, fmt.Sprintf("%s:%s", login, *repo.DefaultBranch))
			if compareErr != nil {
				return deletions, compareErr
			}

			if shouldDelete(repo, filter, issues, commits) {
				deletions = append(deletions, repo)
			}
		}
		if resp.NextPage == 0 {
			break
		}
		ropt.ListOptions.Page = resp.NextPage
	}
	return deletions, nil
}

func shouldDelete(
	repo *github.Repository,
	filter Filter,
	issues []*github.Issue,
	commitComparison *github.CommitsComparison,
) bool {
	for _, r := range filter.Blacklist {
		if r == repo.GetName() {
			return false
		}
	}
	if !filter.IncludePrivate && repo.GetPrivate() {
		return false
	}
	if repo.GetForksCount() > 0 ||
		repo.GetStargazersCount() > 0 ||
		!time.Now().Add(-filter.Since).After((repo.GetUpdatedAt()).Time) {
		return false
	}
	for _, issue := range issues {
		if issue.IsPullRequest() {
			return false
		}
	}

	// check if the fork has commits ahead of the parent repo
	if filter.ExcludeCommitsAhead {
		if *commitComparison.AheadBy > 0 {
			return false
		}
	}

	return true
}

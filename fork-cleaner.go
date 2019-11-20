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
	ShowExcludeReason   bool
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
) ([]*github.Repository, []string, error) {
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
	exclusionReasons := make([]string, 0)
	for {
		repos, resp, err := client.Repositories.List(ctx, "", ropt)
		if err != nil {
			return deletions, exclusionReasons, err
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
				return deletions, exclusionReasons, err
			}
			parent := repo.GetParent()
			po := parent.GetOwner().GetLogin()
			pn := parent.GetName()
			issues, _, err := client.Issues.ListByRepo(ctx, po, pn, iopt)
			if err != nil {
				return deletions, exclusionReasons, err
			}
			commits, _, compareErr := client.Repositories.CompareCommits(ctx, po, rn, *parent.DefaultBranch, fmt.Sprintf("%s:%s", login, *repo.DefaultBranch))
			if compareErr != nil {
				return deletions, exclusionReasons, compareErr
			}

			ok, reason := shouldDelete(repo, filter, issues, commits)
			if ok {
				deletions = append(deletions, repo)
			} else {
				exclusionReasons = append(exclusionReasons, reason)
			}
		}
		if resp.NextPage == 0 {
			break
		}
		ropt.ListOptions.Page = resp.NextPage
	}
	return deletions, exclusionReasons, nil
}

func shouldDelete(
	repo *github.Repository,
	filter Filter,
	issues []*github.Issue,
	commitComparison *github.CommitsComparison,
) (bool, string) {
	var reason string
	for _, r := range filter.Blacklist {
		if r == repo.GetName() {
			if filter.ShowExcludeReason {
				reason = fmt.Sprintf("%s excluded because: repo is blacklisted\n", *repo.HTMLURL)
			}
			return false, reason
		}
	}
	if !filter.IncludePrivate && repo.GetPrivate() {
		if filter.ShowExcludeReason {
			reason = fmt.Sprintf("%s excluded because: repo is private\n", *repo.HTMLURL)
		}
		return false, reason
	}
	if repo.GetForksCount() > 0 {
		if filter.ShowExcludeReason {
			reason = fmt.Sprintf("%s excluded because: repo has %d forks\n", *repo.HTMLURL, *repo.ForksCount)
		}
		return false, reason
	}
	if repo.GetStargazersCount() > 0 {
		if filter.ShowExcludeReason {
			reason = fmt.Sprintf("%s excluded because: repo has %d stars\n", *repo.HTMLURL, *repo.StargazersCount)
		}
		return false, reason
	}
	if !time.Now().Add(-filter.Since).After((repo.GetUpdatedAt()).Time) {
		if filter.ShowExcludeReason {
			reason = fmt.Sprintf("%s excluded because: repo has recent activity (last update on %s)\n", *repo.HTMLURL, repo.GetUpdatedAt().Format("1/2/2006"))
		}
		return false, reason
	}
	for _, issue := range issues {
		if issue.IsPullRequest() {
			if filter.ShowExcludeReason {
				reason = fmt.Sprintf("%s excluded because: repo has a pull request\n", *repo.HTMLURL)
			}
			return false, reason
		}
	}

	// check if the fork has commits ahead of the parent repo
	if filter.ExcludeCommitsAhead {
		if *commitComparison.AheadBy > 0 {
			if filter.ShowExcludeReason {
				reason = fmt.Sprintf("%s excluded because: repo is %d commits ahead of parent\n", *repo.HTMLURL, *commitComparison.AheadBy)
			}
			return false, reason
		}
	}

	return true, reason
}

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

// Filter applied to the repositories list.
type Filter struct {
	Blacklist           []string
	Since               time.Duration
	IncludePrivate      bool
	IncludeStarred      bool
	IncludeForked       bool
	ExcludeCommitsAhead bool
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

// Find list the forks from a given owner that could be deleted.
func Find(
	ctx context.Context,
	client *github.Client,
	filter Filter,
) ([]*github.Repository, []string, error) {
	lopt := github.ListOptions{PerPage: pageSize}
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
			commits, resp, compareErr := client.Repositories.CompareCommits(ctx, po, pn, *parent.DefaultBranch, fmt.Sprintf("%s:%s", login, *repo.DefaultBranch))
			if resp.StatusCode == http.StatusNotFound {
				exclusionReasons = append(exclusionReasons, fmt.Sprintf("%s excluded because: parent repo doesn't exist anymore\n", *repo.HTMLURL))
				continue
			}
			if resp.StatusCode == http.StatusUnavailableForLegalReasons {
				exclusionReasons = append(exclusionReasons, fmt.Sprintf("%s excluded because: DMCA take down\n", *repo.HTMLURL))
				continue
			}
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
	for _, r := range filter.Blacklist {
		if r == repo.GetName() {
			return false, fmt.Sprintf("%s excluded because: repo is blacklisted\n", *repo.HTMLURL)
		}
	}
	if !filter.IncludePrivate && repo.GetPrivate() {
		return false, fmt.Sprintf("%s excluded because: repo is private\n", *repo.HTMLURL)
	}
	if !filter.IncludeForked && repo.GetForksCount() > 0 {
		return false, fmt.Sprintf("%s excluded because: repo has %d forks\n", *repo.HTMLURL, *repo.ForksCount)
	}
	if !filter.IncludeStarred && repo.GetStargazersCount() > 0 {
		return false, fmt.Sprintf("%s excluded because: repo has %d stars\n", *repo.HTMLURL, *repo.StargazersCount)
	}
	if !time.Now().Add(-filter.Since).After((repo.GetUpdatedAt()).Time) {
		return false, fmt.Sprintf("%s excluded because: repo has recent activity (last update on %s)\n", *repo.HTMLURL, repo.GetUpdatedAt().Format("1/2/2006"))
	}
	for _, issue := range issues {
		if issue.IsPullRequest() {
			return false, fmt.Sprintf("%s excluded because: repo has a pull request\n", *repo.HTMLURL)
		}
	}

	// check if the fork has commits ahead of the parent repo
	if filter.ExcludeCommitsAhead && *commitComparison.AheadBy > 0 {
		return false, fmt.Sprintf("%s excluded because: repo is %d commits ahead of parent\n", *repo.HTMLURL, *commitComparison.AheadBy)
	}

	return true, ""
}

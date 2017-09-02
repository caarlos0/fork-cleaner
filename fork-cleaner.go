// Package forkcleaner provides functions to find and remove unused forks.
package forkcleaner

import (
	"context"
	"time"

	"github.com/google/go-github/github"
)

// Filter applied to the repositories list
type Filter struct {
	Blacklist      []string
	IncludePrivate bool
	Since          time.Duration
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
	opt := &github.RepositoryListOptions{
		ListOptions: github.ListOptions{PerPage: 100},
		Affiliation: "owner",
	}
	var deletions []*github.Repository
	for {
		repos, resp, err := client.Repositories.List(ctx, "", opt)
		if err != nil {
			return deletions, err
		}
		for _, repo := range repos {
			if shouldDelete(repo, filter) {
				deletions = append(deletions, repo)
			}
		}
		if resp.NextPage == 0 {
			break
		}
		opt.ListOptions.Page = resp.NextPage
	}
	return deletions, nil
}

func shouldDelete(repo *github.Repository, filter Filter) bool {
	for _, r := range filter.Blacklist {
		if r == repo.GetName() {
			return false
		}
	}
	if !filter.IncludePrivate && repo.GetPrivate() {
		return false
	}
	return repo.GetFork() &&
		repo.GetForksCount() == 0 &&
		repo.GetStargazersCount() == 0 &&
		time.Now().Add(-filter.Since).After((repo.GetUpdatedAt()).Time)
}

package cleaner

import (
	"context"
	"time"

	"github.com/google/go-github/github"
)

// DeleteForks delete the given list of forks
func DeleteForks(
	ctx context.Context,
	deletions []*github.Repository,
	client *github.Client,
) error {
	for _, repo := range deletions {
		_, err := client.Repositories.Delete(ctx, *repo.Owner.Login, *repo.Name)
		if err != nil {
			return err
		}
	}
	return nil
}

// Repos list the forks from a given owner that could be deleted
func Repos(
	ctx context.Context,
	owner string,
	client *github.Client,
) ([]*github.Repository, error) {
	opt := &github.RepositoryListOptions{
		ListOptions: github.ListOptions{PerPage: 50},
	}
	var deletions []*github.Repository
	for {
		repos, resp, err := client.Repositories.List(ctx, owner, opt)
		if err != nil {
			return deletions, err
		}
		for _, repo := range repos {
			if shouldDelete(repo) {
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

func shouldDelete(repo *github.Repository) bool {
	return *repo.Fork &&
		*repo.ForksCount == 0 &&
		*repo.StargazersCount == 0 &&
		!*repo.Private &&
		time.Now().AddDate(0, -1, 0).After((*repo.UpdatedAt).Time)
}

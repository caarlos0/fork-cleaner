package forkcleaner

import (
	"context"

	"github.com/google/go-github/v33/github"
)

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

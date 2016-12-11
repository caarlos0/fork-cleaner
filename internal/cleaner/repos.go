package cleaner

import (
	"fmt"
	"time"

	"github.com/google/go-github/github"
)

func DeleteForks(deletions []*github.Repository, client *github.Client) error {
	for _, repo := range deletions {
		fmt.Println("Deleting fork", *repo.FullName+"...")
		_, err := client.Repositories.Delete(*repo.Owner.Login, *repo.Name)
		if err != nil {
			return err
		}
	}
	return nil
}

func Repos(owner string, client *github.Client) ([]*github.Repository, error) {
	repos, err := allRepos(owner, client)
	if err != nil {
		panic(err)
	}
	fmt.Println("Repos that could be deleted:")
	var deletions []*github.Repository
	for _, repo := range repos {
		if shouldDelete(repo) {
			deletions = append(deletions, repo)
			fmt.Println(*repo.HTMLURL)
		}
	}
	return deletions, err
}

func allRepos(owner string, client *github.Client) ([]*github.Repository, error) {
	opt := &github.RepositoryListOptions{
		ListOptions: github.ListOptions{PerPage: 50},
	}

	var allRepos []*github.Repository
	for {
		repos, resp, err := client.Repositories.List(owner, opt)
		if err != nil {
			return allRepos, err
		}
		allRepos = append(allRepos, repos...)
		if resp.NextPage == 0 {
			break
		}
		opt.ListOptions.Page = resp.NextPage
	}
	return allRepos, nil
}

func shouldDelete(repo *github.Repository) bool {
	return *repo.Fork &&
		*repo.ForksCount == 0 &&
		*repo.StargazersCount == 0 &&
		!*repo.Private &&
		time.Now().AddDate(0, -1, 0).After((*repo.UpdatedAt).Time)
}

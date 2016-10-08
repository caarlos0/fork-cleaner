package main

import (
	"bufio"
	"fmt"
	"os"

	"time"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

func repos(org string, client *github.Client) ([]*github.Repository, error) {
	opt := &github.RepositoryListOptions{
		ListOptions: github.ListOptions{PerPage: 50},
	}

	var allRepos []*github.Repository
	for {
		repos, resp, err := client.Repositories.List(org, opt)
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

func main() {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv("GITHUB_TOKEN")},
	)
	tc := oauth2.NewClient(oauth2.NoContext, ts)
	client := github.NewClient(tc)
	user, _, err := client.Users.Get("")
	if err != nil {
		panic(err)
	}

	allRepos, err := repos(*user.Login, client)
	if err != nil {
		panic(err)
	}
	fmt.Println("Repos that could be deleted:")
	var deletions []*github.Repository
	for _, repo := range allRepos {
		if shouldDelete(repo) {
			deletions = append(deletions, repo)
			fmt.Println(*repo.HTMLURL)
		}
	}
	fmt.Print("\nDelete these ", len(deletions), " forks? [y/n] ")
	reply, err := bufio.NewReader(os.Stdin).ReadString('\n')
	if err != nil {
		panic(err)
	}
	if reply == "y\n" || reply == "Y\n" {
		for _, repo := range deletions {
			fmt.Println("Deleting fork", *repo.FullName+"...")
			_, err := client.Repositories.Delete(*repo.Owner.Login, *repo.Name)
			if err != nil {
				panic(err)
			}
		}
	} else {
		fmt.Println("Replied", reply, ". Exiting.")
	}
}

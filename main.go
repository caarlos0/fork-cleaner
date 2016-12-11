package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/caarlos0/fork-cleaner/internal/cleaner"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

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

	deletions, err := cleaner.Repos(*user.Login, client)
	if err != nil {
		panic(err)
	}

	fmt.Print("\nDelete all ", len(deletions), " listed forks? [y/n] ")
	reply, err := bufio.NewReader(os.Stdin).ReadString('\n')
	if err != nil {
		panic(err)
	}
	if reply == "y\n" || reply == "Y\n" {
		if err := cleaner.DeleteForks(deletions, client); err != nil {
			panic(err)
		}
	} else {
		fmt.Println("OK, exiting.")
	}
}

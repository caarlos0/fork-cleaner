package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/caarlos0/fork-cleaner/internal/cleaner"
	"github.com/caarlos0/spin"
	"github.com/google/go-github/github"
	"github.com/urfave/cli"
	"golang.org/x/oauth2"
)

var version = "master"

func main() {
	app := cli.NewApp()
	app.Name = "fork-cleaner"
	app.Version = version
	app.Author = "Carlos Alexandro Becker (caarlos0@gmail.com)"
	app.Usage = "Delete old, unused forks"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			EnvVar: "GITHUB_TOKEN",
			Name:   "token",
			Usage:  "Your GitHub token",
		},
		cli.StringFlag{
			Name:  "owner",
			Usage: "GitHub user or organization to clean up.",
		},
	}
	app.Action = func(c *cli.Context) error {
		ts := oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: c.String("token")},
		)
		tc := oauth2.NewClient(oauth2.NoContext, ts)
		client := github.NewClient(tc)
		owner := c.String("owner")
		if owner == "" {
			user, _, err := client.Users.Get("")
			if err != nil {
				return err
			}
			owner = *user.Login
		}

		sg := spin.New("\033[36m %s Gathering data for '" + owner + "'...\033[m")
		sg.Start()
		deletions, err := cleaner.Repos(owner, client)
		sg.Stop()

		if err != nil {
			return err
		}
		for _, repo := range deletions {
			fmt.Println(*repo.HTMLURL)
		}

		fmt.Print("\nDelete all ", len(deletions), " listed forks? [y/n] ")
		reply, err := bufio.NewReader(os.Stdin).ReadString('\n')
		if err != nil {
			return err
		}
		if reply == "y\n" || reply == "Y\n" {
			sd := spin.New(fmt.Sprintf(
				"\033[36m %s Deleting %d forks...\033[m", "%s", len(deletions),
			))
			sd.Start()
			err = cleaner.DeleteForks(deletions, client)
			sd.Stop()
			if err != nil {
				return err
			}
		} else {
			fmt.Println("OK, exiting.")
		}
		return nil
	}

	app.Run(os.Args)
}

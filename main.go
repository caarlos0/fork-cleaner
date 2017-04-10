package main

import (
	"context"
	"fmt"
	"log"
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
		cli.BoolFlag{
			Name:  "dry-run, d",
			Usage: "Only list forks to be deleted, but don't delete them",
		},
	}
	app.Action = func(c *cli.Context) error {
		log.SetFlags(0)
		ctx := context.Background()
		ts := oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: c.String("token")},
		)
		tc := oauth2.NewClient(ctx, ts)
		client := github.NewClient(tc)
		owner := c.String("owner")
		if owner == "" {
			user, _, err := client.Users.Get(ctx, "")
			if err != nil {
				return cli.NewExitError(err.Error(), 1)
			}
			owner = *user.Login
		}

		sg := spin.New("\033[36m %s Gathering data for '" + owner + "'...\033[m")
		sg.Start()
		deletions, err := cleaner.Repos(ctx, owner, client)
		sg.Stop()
		if err != nil {
			return cli.NewExitError(err.Error(), 1)
		}
		if len(deletions) == 0 {
			log.Println("0 forks to delete!")
			return nil
		}
		log.Println(len(deletions), "forks to delete:")
		log.SetPrefix(" --> ")
		for _, repo := range deletions {
			log.Println(*repo.HTMLURL)
		}
		log.SetPrefix("")

		if c.Bool("dry-run") {
			log.Println("\nDry-Run (-d) is set! No action taken.")
			return nil
		}
		fmt.Printf("\n\n")
		sd := spin.New(fmt.Sprintf(
			"\033[36m %s Deleting %d forks...\033[m", "%s", len(deletions),
		))
		sd.Start()
		err = cleaner.DeleteForks(ctx, deletions, client)
		sd.Stop()
		if err != nil {
			return cli.NewExitError(err.Error(), 1)
		}

		return nil
	}

	log.Fatalln(app.Run(os.Args))
}

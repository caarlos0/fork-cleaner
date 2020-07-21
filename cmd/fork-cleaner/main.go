package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/Songmu/prompter"
	forkcleaner "github.com/caarlos0/fork-cleaner"
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
			Name:   "token, t",
			Usage:  "Your GitHub token",
		},
		cli.BoolFlag{
			Name:  "force, f",
			Usage: "Don't ask to remove the forks",
		},
		cli.BoolFlag{
			Name:  "include-private, p",
			Usage: "Include private repositories",
		},
		cli.BoolFlag{
			Name:  "exclude-commits-ahead, a",
			Usage: "Exclude repositories with commits ahead of parent",
		},
		cli.BoolFlag{
			Name:  "show-exclusion-reason",
			Usage: "Show the reason a fork is excluded",
		},
		cli.StringSliceFlag{
			Name:  "blacklist, exclude, b",
			Usage: "Blacklist of repos that shouldn't be removed (names only)",
		},
		cli.DurationFlag{
			Name:  "no-activity-since, since",
			Usage: "Time to check for activity",
			Value: 30 * 24 * time.Hour,
		},
	}
	app.Action = func(c *cli.Context) error {
		log.SetFlags(0)
		token := c.String("token")
		blacklist := c.StringSlice("blacklist")
		ctx := context.Background()
		ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
		tc := oauth2.NewClient(ctx, ts)
		client := github.NewClient(tc)
		if token == "" {
			return cli.NewExitError("missing github token", 1)
		}

		sg := spin.New("\033[36m %s Gathering data...\033[m")
		sg.Start()
		filter := forkcleaner.Filter{
			Blacklist:           blacklist,
			IncludePrivate:      c.Bool("include-private"),
			Since:               c.Duration("since"),
			ExcludeCommitsAhead: c.Bool("exclude-commits-ahead"),
		}
		forks, excludedForks, err := forkcleaner.Find(ctx, client, filter)
		sg.Stop()
		if err != nil {
			return cli.NewExitError(err.Error(), 1)
		}
		if c.Bool("show-exclusion-reason") && len(excludedForks) > 0 {
			log.Println(len(excludedForks), "forks excluded from deletion:")
			for _, f := range excludedForks {
				log.Print(f)
			}
			log.Println()
		}
		if len(forks) == 0 {
			log.Println("0 forks to delete!")
			return nil
		}
		log.Println(len(forks), "forks to delete:")
		log.SetPrefix(" --> ")
		for _, repo := range forks {
			log.Println(*repo.HTMLURL)
		}
		log.SetPrefix("")

		remove := true
		if !c.Bool("force") {
			remove = prompter.YN("Remove the above listed forks?", false)
		}
		if !remove {
			log.Println("OK, exiting")
			return nil
		}
		fmt.Printf("\n\n")
		sd := spin.New(fmt.Sprintf(
			"\033[36m %s Deleting %d forks...\033[m", "%s", len(forks),
		))
		sd.Start()
		err = forkcleaner.Delete(ctx, client, forks)
		sd.Stop()
		if err == nil {
			log.Println("Forks removed!")
			return nil
		}
		return cli.NewExitError(err.Error(), 1)
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatalln(err)
	}
}

package main

import (
	"context"
	"fmt"
	"log"
	"os"

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
		cli.StringFlag{
			Name:  "owner, o",
			Usage: "GitHub user or organization to clean up",
		},
		cli.BoolFlag{
			Name:  "force, f",
			Usage: "Don't ask to remove the forks",
		},
		cli.StringSliceFlag{
			Name:  "blacklist, exclude, b",
			Usage: "Blacklist of repos that shouldn't be removed",
		},
	}
	app.Action = func(c *cli.Context) error {
		log.SetFlags(0)
		var token = c.String("token")
		var owner = c.String("owner")
		var blacklist = c.StringSlice("blacklist")
		var ctx = context.Background()
		var ts = oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
		var tc = oauth2.NewClient(ctx, ts)
		var client = github.NewClient(tc)
		if token == "" {
			return cli.NewExitError("missing github token", 1)
		}
		if owner == "" {
			user, _, err := client.Users.Get(ctx, "")
			if err != nil {
				return cli.NewExitError(err.Error(), 1)
			}
			owner = *user.Login
		}

		var sg = spin.New("\033[36m %s Gathering data for '" + owner + "'...\033[m")
		sg.Start()
		forks, err := forkcleaner.Find(ctx, client, owner, blacklist)
		sg.Stop()
		if err != nil {
			return cli.NewExitError(err.Error(), 1)
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

		var remove = true
		if !c.Bool("force") {
			remove = prompter.YN("Remove the above listed forks?", false)
		}
		if !remove {
			log.Println("OK, exiting")
			return nil
		}
		fmt.Printf("\n\n")
		var sd = spin.New(fmt.Sprintf(
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

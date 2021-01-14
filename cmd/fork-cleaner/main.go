package main

import (
	"context"
	"log"
	"os"

	"github.com/caarlos0/fork-cleaner/v2/internal/ui"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/google/go-github/v33/github"
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
			EnvVar: "GITHUB_URL",
			Name:   "github-url, g",
			Usage:  "Base GitHub URL",
			Value:  "https://api.github.com/",
		},
	}

	app.Action = func(c *cli.Context) error {
		log.SetFlags(0)
		f, err := tea.LogToFile("fork-cleaner.log", "")
		if err != nil {
			return cli.NewExitError(err.Error(), 1)
		}
		defer func() { _ = f.Close() }()

		token := c.String("token")
		ghurl := c.String("github-url")

		ctx := context.Background()
		ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
		tc := oauth2.NewClient(ctx, ts)
		client, err := github.NewEnterpriseClient(ghurl, ghurl, tc)
		if err != nil {
			return cli.NewExitError(err.Error(), 1)
		}

		if token == "" {
			return cli.NewExitError("missing github token", 1)
		}

		var p = tea.NewProgram(ui.NewInitialModel(client))
		p.EnterAltScreen()
		defer p.ExitAltScreen()
		if err = p.Start(); err != nil {
			return cli.NewExitError(err.Error(), 1)
		}
		return nil
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatalln(err)
	}
}

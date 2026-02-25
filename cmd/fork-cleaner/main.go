package main

import (
	"context"
	"log"
	"os"
	"path/filepath"

	tea "charm.land/bubbletea/v2"
	"github.com/caarlos0/fork-cleaner/v2/internal/ui"
	"github.com/google/go-github/v83/github"
	"github.com/urfave/cli/v2"
	"golang.org/x/oauth2"
)

var version = "main"

func main() {
	app := cli.NewApp()
	app.Name = "fork-cleaner"
	app.Version = version
	app.Authors = []*cli.Author{{
		Name:  "Carlos Alexandro Becker",
		Email: "carlos@becker.software",
	}}
	app.Usage = "Archive or delete old, unused forks"
	app.Flags = []cli.Flag{
		&cli.StringFlag{
			EnvVars: []string{"GITHUB_TOKEN"},
			Name:    "token",
			Usage:   "Your GitHub token",
			Aliases: []string{"t"},
		},
		&cli.StringFlag{
			EnvVars: []string{"GITHUB_URL"},
			Name:    "github-url",
			Usage:   "Base GitHub URL",
			Value:   "https://api.github.com/",
			Aliases: []string{"g"},
		},
		&cli.StringFlag{
			Name:    "user",
			Usage:   "GitHub username or organization name. Defaults to current user.",
			Aliases: []string{"u"},
		},
		&cli.BoolFlag{
			Name:    "skip-upstream-check",
			Usage:   "Skip checking and pulling details from the parent/upstream repository",
			Aliases: []string{"skip-upstream"},
			Value:   false,
		},
	}

	app.Action = func(c *cli.Context) error {
		log.SetFlags(0)
		f, err := tea.LogToFile(filepath.Join("/tmp", "fork-cleaner.log"), "")
		if err != nil {
			return cli.Exit(err.Error(), 1)
		}
		defer func() { _ = f.Close() }()

		token := c.String("token")
		ghurl := c.String("github-url")
		login := c.String("user")
		skipUpstream := c.Bool("skip-upstream-check")

		ctx := context.Background()
		ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
		tc := oauth2.NewClient(ctx, ts)
		client, err := github.NewEnterpriseClient(ghurl, ghurl, tc)
		if err != nil {
			return cli.Exit(err.Error(), 1)
		}

		if token == "" {
			return cli.Exit("missing github token", 1)
		}

		p := tea.NewProgram(ui.NewAppModel(client, login, skipUpstream))
		if _, err = p.Run(); err != nil {
			return cli.Exit(err.Error(), 1)
		}
		return nil
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatalln(err)
	}
}

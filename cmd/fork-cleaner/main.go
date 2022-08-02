package main

import (
	"context"
	"log"
	"os"

	"github.com/caarlos0/fork-cleaner/v2/internal/ui"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/google/go-github/v45/github"
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
	app.Usage = "Delete old, unused forks"
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
	}

	app.Action = func(c *cli.Context) error {
		log.SetFlags(0)
		f, err := tea.LogToFile("fork-cleaner.log", "")
		if err != nil {
			return cli.Exit(err.Error(), 1)
		}
		defer func() { _ = f.Close() }()

		token := c.String("token")
		ghurl := c.String("github-url")
		login := c.String("user")

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

		p := tea.NewProgram(ui.NewInitialModel(client, login), tea.WithAltScreen())
		if err = p.Start(); err != nil {
			return cli.Exit(err.Error(), 1)
		}
		return nil
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatalln(err)
	}
}

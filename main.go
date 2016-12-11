package main

import (
	"bufio"
	"fmt"
	"os"
	"time"

	"github.com/caarlos0/fork-cleaner/internal/cleaner"
	"github.com/google/go-github/github"
	spin "github.com/tj/go-spin"
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

		active := true
		s := spin.New()
		s.Set(`⦾⦿`)
		go func() {
			for active {
				fmt.Printf("\r  \033[36m%s Gathering data for '%s'...\033[m", s.Next(), owner)
				time.Sleep(100 * time.Millisecond)
			}
		}()
		deletions, err := cleaner.Repos(owner, client)
		active = false
		fmt.Printf("\r")
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
			active = true
			go func() {
				for active {
					fmt.Printf("\r  \033[36m%s Removing forks...\033[m", s.Next())
					time.Sleep(100 * time.Millisecond)
				}
			}()
			err = cleaner.DeleteForks(deletions, client)
			active = false
			fmt.Printf("\r")
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

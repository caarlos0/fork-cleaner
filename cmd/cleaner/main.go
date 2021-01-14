package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/caarlos0/fork-cleaner/internal/ui"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/google/go-github/v33/github"
	"golang.org/x/oauth2"
)

var version = "master"

func main() {
	log.SetFlags(0)
	f, err := os.OpenFile("fork-cleaner.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0677)
	if err != nil {
		log.Fatalln(err)
	}
	defer func() { _ = f.Close() }()
	log.SetOutput(f)

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: os.Getenv("GITHUB_TOKEN")})
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	var p = tea.NewProgram(ui.NewInitialModel(client))
	p.EnterAltScreen()
	err = p.Start()
	p.ExitAltScreen()
	if err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}

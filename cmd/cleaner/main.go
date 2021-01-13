package main

import (
	"context"
	"fmt"
	"os"

	forkcleaner "github.com/caarlos0/fork-cleaner"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/google/go-github/v33/github"
	"github.com/muesli/termenv"
	"golang.org/x/oauth2"
)

var version = "master"

func main() {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: os.Getenv("GITHUB_TOKEN")})
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	var s = spinner.NewModel()
	s.Spinner = spinner.MiniDot
	var p = tea.NewProgram(model{
		client:  client,
		spinner: s,
		loading: true,
	})
	p.EnterAltScreen()
	err := p.Start()
	p.ExitAltScreen()
	if err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}

type model struct {
	repos   []*forkcleaner.RepositoryWithDetails
	spinner spinner.Model
	loading bool
	err     error
	client  *github.Client
}

func (m model) Init() tea.Cmd {
	return tea.Batch(getRepos(m.client), spinner.Tick)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case errMsg:
		m.loading = false
		m.err = msg.error
		return m, nil
	case gotRepoListMsg:
		m.repos = msg.repos
		m.loading = false
		return m, nil
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q", "esc":
			return m, tea.Quit
		}
	default:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}
	return m, nil
}

func (m model) View() string {
	if m.loading {
		return termenv.String(m.spinner.View()).Foreground(termenv.ColorProfile().Color("205")).String() + " Loading..."
	}
	if m.err != nil {
		return fmt.Sprintf("Error: %s", m.err.Error())
	}
	return fmt.Sprintf("got %v repos", len(m.repos))
}

type gotRepoListMsg struct {
	repos []*forkcleaner.RepositoryWithDetails
}

type errMsg struct{ error }

func (e errMsg) Error() string { return e.Error() }

func getRepos(client *github.Client) tea.Cmd {
	return func() tea.Msg {
		repos, err := forkcleaner.FindAllForks(context.Background(), client)
		if err != nil {
			return errMsg{err}
		}
		return gotRepoListMsg{repos}
	}
}

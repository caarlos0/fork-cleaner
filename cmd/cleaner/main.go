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
var mainColor = termenv.ColorProfile().Color("205")

func main() {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: os.Getenv("GITHUB_TOKEN")})
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	var s = spinner.NewModel()
	s.Spinner = spinner.MiniDot
	var p = tea.NewProgram(initialModel{
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

type listModel struct {
	err      error
	client   *github.Client
	repos    []*forkcleaner.RepositoryWithDetails
	cursor   int
	selected map[int]struct{}
}

func (m listModel) Init() tea.Cmd {
	return nil
}

func (m listModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q", "esc":
			return m, tea.Quit
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.repos)-1 {
				m.cursor++
			}
		case "a":
			for i := range m.repos {
				m.selected[i] = struct{}{}
			}
		case "n":
			for i := range m.selected {
				delete(m.selected, i)
			}
		case " ":
			_, ok := m.selected[m.cursor]
			if ok {
				delete(m.selected, m.cursor)
			} else {
				m.selected[m.cursor] = struct{}{}
			}
		}
	}
	return m, nil
}

func (m listModel) View() string {
	var s = termenv.String("Which forks you want to delete?\n\n").Bold().String()

	for i, choice := range m.repos {
		var line = choice.Repo.GetFullName() + "\n"

		if _, ok := m.selected[i]; ok {
			line = termenv.String(line).CrossOut().Faint().String()
		}

		if m.cursor == i {
			line = termenv.String(line).Foreground(mainColor).Bold().String()
		}

		s += line
	}

	return s + "\nPress " + bold("q") + " to quit, " +
		bold("space") + " to select the current item, " +
		bold("a") + " to select all items, " +
		bold("n") + " to deselect all items, " +
		bold("enter") + " to delete selected.\n"
}

func bold(s string) string {
	return termenv.String(s).Foreground(mainColor).Bold().String()
}

type initialModel struct {
	err     error
	client  *github.Client
	spinner spinner.Model
	loading bool
}

func (m initialModel) Init() tea.Cmd {
	return tea.Batch(getRepos(m.client), spinner.Tick)
}

func (m initialModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case errMsg:
		m.loading = false
		m.err = msg.error
		return m, nil
	case gotRepoListMsg:
		nm := listModel{
			repos:    msg.repos,
			client:   m.client,
			selected: map[int]struct{}{},
		}
		return nm, nm.Init()
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

func (m initialModel) View() string {
	if m.loading {
		return termenv.String(m.spinner.View()).Foreground(mainColor).String() + " Loading list of forks..."
	}
	if m.err != nil {
		return fmt.Sprintf("Error: %s", m.err.Error())
	}
	return "oops..."
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

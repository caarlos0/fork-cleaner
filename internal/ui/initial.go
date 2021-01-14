package ui

import (
	"context"
	"fmt"

	forkcleaner "github.com/caarlos0/fork-cleaner"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/google/go-github/v33/github"
	"github.com/muesli/termenv"
)

func NewInitialModel(client *github.Client) InitialModel {
	var s = spinner.NewModel()
	s.Spinner = spinner.MiniDot

	return InitialModel{
		client:  client,
		spinner: s,
		loading: true,
	}
}

type InitialModel struct {
	err     error
	client  *github.Client
	spinner spinner.Model
	loading bool
}

func (m InitialModel) Init() tea.Cmd {
	return tea.Batch(getRepos(m.client), spinner.Tick)
}

func (m InitialModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case errMsg:
		m.loading = false
		m.err = msg.error
		return m, nil
	case gotRepoListMsg:
		var list = NewListModel(m.client, msg.repos)
		return list, list.Init()
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

func (m InitialModel) View() string {
	if m.loading {
		return termenv.String(m.spinner.View()).Foreground(primaryColor).String() + " Loading list of forks..."
	}
	if m.err != nil {
		return fmt.Sprintf("Error: %s", m.err.Error())
	}
	return "oops..."
}

type gotRepoListMsg struct {
	repos []*forkcleaner.RepositoryWithDetails
}

func getRepos(client *github.Client) tea.Cmd {
	return func() tea.Msg {
		repos, err := forkcleaner.FindAllForks(context.Background(), client)
		if err != nil {
			return errMsg{err}
		}
		return gotRepoListMsg{repos}
	}
}

package ui

import (
	"context"

	forkcleaner "github.com/caarlos0/fork-cleaner/v2"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/google/go-github/v47/github"
)

// NewInitialModel creates a new InitialModel with required fields.
func NewInitialModel(client *github.Client, login string) InitialModel {
	s := spinner.NewModel()
	s.Spinner = spinner.MiniDot

	return InitialModel{
		client:  client,
		login:   login,
		spinner: s,
		loading: true,
	}
}

// InitialModel is the UI when the CLI starts, basically loading the repos.
type InitialModel struct {
	err     error
	login   string
	client  *github.Client
	spinner spinner.Model
	loading bool
}

func (m InitialModel) Init() tea.Cmd {
	return tea.Batch(getRepos(m.client, m.login), spinner.Tick)
}

func (m InitialModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case errMsg:
		m.loading = false
		m.err = msg.error
		return m, nil
	case gotRepoListMsg:
		list := NewListModel(m.client, msg.repos)
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
		return boldPrimaryForeground(m.spinner.View()) + " Gathering repositories..." + singleOptionHelp("q", "quit")
	}
	if m.err != nil {
		return errorView("Error gathering the repository list", m.err)
	}
	return ""
}

type gotRepoListMsg struct {
	repos []*forkcleaner.RepositoryWithDetails
}

func getRepos(client *github.Client, login string) tea.Cmd {
	return func() tea.Msg {
		repos, err := forkcleaner.FindAllForks(context.Background(), client, login)
		if err != nil {
			return errMsg{err}
		}
		return gotRepoListMsg{repos}
	}
}

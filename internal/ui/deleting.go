package ui

import (
	"context"

	forkcleaner "github.com/caarlos0/fork-cleaner/v2"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/google/go-github/v45/github"
)

// NewDeletingModel creates a DeletingModel with required fields.
func NewDeletingModel(client *github.Client, repos []*forkcleaner.RepositoryWithDetails, previous ListModel) DeletingModel {
	s := spinner.NewModel()
	s.Spinner = spinner.MiniDot

	return DeletingModel{
		client:   client,
		repos:    repos,
		spinner:  s,
		previous: previous,
	}
}

// DeletingModel is the UI in which the user can review the repos they
// selected to be deleted and either finally delete them or cancel.
type DeletingModel struct {
	client   *github.Client
	repos    []*forkcleaner.RepositoryWithDetails
	cursor   int
	spinner  spinner.Model
	loading  bool
	previous ListModel
}

func (m DeletingModel) Init() tea.Cmd {
	return nil
}

func (m DeletingModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case reposDeletedMsg:
		return NewDeleteEndModelSucceed(msg.total), nil
	case errMsg:
		return NewDeleteEndModelFailed(msg.error), nil
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "q", "esc", "n":
			return m.previous, m.previous.Init()
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.repos)-1 {
				m.cursor++
			}
		case "y":
			m.loading = true
			return m, tea.Batch(deleteRepos(m.client, m.repos), spinner.Tick)
		}
	default:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}
	return m, nil
}

func (m DeletingModel) View() string {
	if m.loading {
		return redFaintForeground(m.spinner.View()) + redForeground(" Deleting repositories...")
	}

	s := redForeground("Are you sure you want to delete the selected repositories? (y/N)\n\n")
	for i, repo := range m.repos {
		line := faint(iconSelected+" "+repo.Name) + "\n"
		if m.cursor == i {
			line = "\n" + boldRedForeground(line) + viewRepositoryDetails(repo)
		}
		s += line
	}
	return s + helpView([]helpOption{
		{"q/esc/n", "abort", true},
		{"up/down", "navigate", false},
		{"y", "delete items", false},
	})
}

type reposDeletedMsg struct {
	total int
}

func deleteRepos(client *github.Client, repos []*forkcleaner.RepositoryWithDetails) tea.Cmd {
	return func() tea.Msg {
		if err := forkcleaner.Delete(context.Background(), client, repos); err != nil {
			return errMsg{err}
		}
		return reposDeletedMsg{len(repos)}
	}
}

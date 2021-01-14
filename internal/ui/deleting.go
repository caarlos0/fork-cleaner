package ui

import (
	"context"
	"fmt"
	"strconv"

	forkcleaner "github.com/caarlos0/fork-cleaner"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/google/go-github/v33/github"
	"github.com/muesli/termenv"
)

func NewDeletingModel(client *github.Client, repos []*forkcleaner.RepositoryWithDetails) DeletingModel {
	var s = spinner.NewModel()
	s.Spinner = spinner.MiniDot

	return DeletingModel{
		client: client,
		repos:  repos,
		spinner: s,
	}
}

type DeletingModel struct {
	err    error
	client *github.Client
	repos  []*forkcleaner.RepositoryWithDetails
	cursor int
	deleted int
	spinner spinner.Model
	loading bool
}

func (m DeletingModel) Init() tea.Cmd {
	return nil
}

func (m DeletingModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case reposDeletedMsg:
		m.loading = false
		m.deleted = msg.total
	case errMsg:
		m.loading = false
		m.err = msg.error
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q", "esc", "n":
			return m, tea.Quit
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
		return redFaintForeground(m.spinner.View()) + redForeground("Deleting repositories...")
	}
	if m.deleted > 0 {
		return redFaintForeground("Successfully deleted ") + redForeground(strconv.Itoa(m.deleted)) + redFaintForeground(" forks.")
	}
	if m.err != nil {
		return redForeground(fmt.Sprintf("Error getting the repository list: %s", m.err.Error()))
	}

	var s = redForeground("Are you sure you want to delete the selected repositories? (y/N)\n\n")
	for i, repo := range m.repos {
		var line = termenv.String(iconSelected+" "+repo.Name).Faint().String() + "\n"
		if m.cursor == i {
			line = boldPrimaryForeground(line) + viewRepositoryDetails(repo)
		}
		s += line
	}
	return s + helpView()
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

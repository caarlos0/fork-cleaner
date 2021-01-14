package ui

import (
	forkcleaner "github.com/caarlos0/fork-cleaner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/google/go-github/v33/github"
	"github.com/muesli/termenv"
)

func NewListModel(client *github.Client, repos []*forkcleaner.RepositoryWithDetails) ListModel {
	return ListModel{
		client:   client,
		repos:    repos,
		selected: map[int]struct{}{},
	}
}

type ListModel struct {
	err      error
	client   *github.Client
	repos    []*forkcleaner.RepositoryWithDetails
	cursor   int
	selected map[int]struct{}
}

func (m ListModel) Init() tea.Cmd {
	return nil
}

func (m ListModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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

const (
	iconSelected    = "●"
	iconNotSelected = "○"
)

func (m ListModel) View() string {
	var s = termenv.String("Which forks you want to delete?\n\n").Foreground(secondaryColor).Bold().String()

	for i, choice := range m.repos {
		var line = choice.Repo.GetFullName() + "\n"

		if _, ok := m.selected[i]; ok {
			line = termenv.String(iconSelected + " " + line).Faint().String()
		} else {
			line = iconNotSelected + " " + line
		}

		if m.cursor == i {
			line = termenv.String(line).Foreground(primaryColor).Bold().String()
		}

		s += line
	}

	return s + "\nPress " + bold("q") + " to quit, " +
		bold("space") + " to select the current item, " +
		bold("a") + " to select all items, " +
		bold("n") + " to deselect all items, " +
		bold("enter") + " to delete selected.\n"
}

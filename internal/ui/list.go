package ui

import (
	"fmt"
	"time"

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

	for i, repo := range m.repos {
		var line = repo.Name + "\n"

		if _, ok := m.selected[i]; ok {
			line = termenv.String(iconSelected + " " + line).Faint().String()
		} else {
			line = iconNotSelected + " " + line
		}

		if m.cursor == i {
			line = termenv.String(line).Foreground(primaryColor).Bold().String() + viewRepositoryDetails(repo)
		}

		s += line
	}

	return s + "\nPress " + bold("q") + " to quit, " +
		bold("space") + " to select the current item, " +
		bold("a") + " to select all items, " +
		bold("n") + " to deselect all items, " +
		bold("enter") + " to delete selected.\n"
}

func viewRepositoryDetails(repo *forkcleaner.RepositoryWithDetails) string {
	var details []string
	if repo.ParentDeleted {
		details = append(details, "Parent repository was deleted")
	}
	if repo.ParentDMCATakeDown {
		details = append(details, "Parent repository was taken down by DMCA")
	}
	if repo.Private {
		details = append(details, "Is private")
	}
	if repo.CommitsAhead > 0 {
		details = append(details, fmt.Sprintf("Has %d commits ahead of upstream", repo.CommitsAhead))
	}
	if repo.Forks > 0 {
		details = append(details, fmt.Sprintf("Has %d forks", repo.Forks))
	}
	if repo.Stars > 0 {
		details = append(details, fmt.Sprintf("Has %d stars", repo.Stars))
	}
	if repo.OpenPRs > 0 {
		details = append(details, fmt.Sprintf("Has %d open PRs to upstream", repo.OpenPRs))
	}
	if time.Now().Add(-30 * 24 * time.Hour).Before(repo.LastUpdate) {
		details = append(details, fmt.Sprintf("Was updated recently (%s)", repo.LastUpdate))
	}

	if len(details) == 0 {
		return ""
	}

	var s string
	for _, d := range details {
		s += "    * " + d + "\n"
	}
	s += "\n"
	return termenv.String(s).Faint().Italic().String()
}

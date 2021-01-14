package ui

import (
	"fmt"
	"time"

	forkcleaner "github.com/caarlos0/fork-cleaner/v2"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/google/go-github/v33/github"
	"github.com/muesli/termenv"
)

// NewListModel creates a new ListModel with the required fields.
func NewListModel(client *github.Client, repos []*forkcleaner.RepositoryWithDetails) ListModel {
	return ListModel{
		client:   client,
		repos:    repos,
		selected: map[int]struct{}{},
	}
}

// ListModel is the UI in which the user can select which forks should be
// deleted if any, and see details on each of them.
type ListModel struct {
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
		case "d":
			var deleteable []*forkcleaner.RepositoryWithDetails
			for k := range m.selected {
				deleteable = append(deleteable, m.repos[k])
			}
			var dm = NewDeletingModel(m.client, deleteable, m)
			return dm, dm.Init()
		}
	}
	return m, nil
}

func (m ListModel) View() string {
	var s = boldSecondaryForeground("Which of these forks do you want to delete?\n\n")

	for i, repo := range m.repos {
		var line = repo.Name
		if _, ok := m.selected[i]; ok {
			line = iconSelected + " " + line
		} else {
			line = faint(iconNotSelected + " " + line)
		}
		line += "\n"

		if m.cursor == i {
			line = "\n" + boldPrimaryForeground(line) + viewRepositoryDetails(repo)
		}

		s += line
	}

	return s + helpView([]helpOption{
		{"up/down", "navigate", false},
		{"space", "toggle selection", false},
		{"d", "delete selected", true},
		{"a", "select all", false},
		{"n", "deselect all", false},
		{"q/esc", "quit", false},
	})
}

func viewRepositoryDetails(repo *forkcleaner.RepositoryWithDetails) string {
	var details []string
	if repo.ParentName != "" {
		details = append(details, fmt.Sprintf("Forked from %s", repo.ParentName))
	}
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
		details = append(details, fmt.Sprintf("Has %d commit%s ahead of upstream", repo.CommitsAhead, maybePlural(repo.CommitsAhead)))
	}
	if repo.Forks > 0 {
		details = append(details, fmt.Sprintf("Has %d fork%s", repo.Forks, maybePlural(repo.Forks)))
	}
	if repo.Stars > 0 {
		details = append(details, fmt.Sprintf("Has %d star%s", repo.Stars, maybePlural(repo.Stars)))
	}
	if repo.OpenPRs > 0 {
		details = append(details, fmt.Sprintf("Has %d open PR%s to upstream", repo.OpenPRs, maybePlural(repo.OpenPRs)))
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

func maybePlural(n int) string {
	if n == 1 {
		return ""
	}
	return "s"
}

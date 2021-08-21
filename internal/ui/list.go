package ui

import (
	"fmt"
	"strings"
	"time"

	forkcleaner "github.com/caarlos0/fork-cleaner/v2"
	timeago "github.com/caarlos0/timea.go"
)

// import (
// 	"fmt"
// 	"time"

// 	forkcleaner "github.com/caarlos0/fork-cleaner/v2"
// 	tea "github.com/charmbracelet/bubbletea"
// 	"github.com/google/go-github/v33/github"
// 	"github.com/muesli/termenv"
// )

// // NewListModel creates a new ListModel with the required fields.
// func NewListModel(client *github.Client, repos []*forkcleaner.RepositoryWithDetails) ListModel {
// 	return ListModel{
// 		client:   client,
// 		repos:    repos,
// 		selected: map[int]struct{}{},
// 	}
// }

// // ListModel is the UI in which the user can select which forks should be
// // deleted if any, and see details on each of them.
// type ListModel struct {
// 	client   *github.Client
// 	repos    []*forkcleaner.RepositoryWithDetails
// 	cursor   int
// 	selected map[int]struct{}
// }

// func (m ListModel) Init() tea.Cmd {
// 	return nil
// }

// func (m ListModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
// 	switch msg := msg.(type) {
// 	case tea.KeyMsg:
// 		switch msg.String() {
// 		case "ctrl+c", "q", "esc":
// 			return m, tea.Quit
// 		case "up", "k":
// 			if m.cursor > 0 {
// 				m.cursor--
// 			}
// 		case "down", "j":
// 			if m.cursor < len(m.repos)-1 {
// 				m.cursor++
// 			}
// 		case "a":
// 			for i := range m.repos {
// 				m.selected[i] = struct{}{}
// 			}
// 		case "n":
// 			for i := range m.selected {
// 				delete(m.selected, i)
// 			}
// 		case " ":
// 			_, ok := m.selected[m.cursor]
// 			if ok {
// 				delete(m.selected, m.cursor)
// 			} else {
// 				m.selected[m.cursor] = struct{}{}
// 			}
// 		case "d":
// 			var deleteable []*forkcleaner.RepositoryWithDetails
// 			for k := range m.selected {
// 				deleteable = append(deleteable, m.repos[k])
// 			}
// 			dm := NewDeletingModel(m.client, deleteable, m)
// 			return dm, dm.Init()
// 		}
// 	}
// 	return m, nil
// }

// func (m ListModel) View() string {
// 	s := boldSecondaryForeground("Which of these forks do you want to delete?\n\n")

// 	for i, repo := range m.repos {
// 		line := repo.Name
// 		if _, ok := m.selected[i]; ok {
// 			line = iconSelected + " " + line
// 		} else {
// 			line = faint(iconNotSelected + " " + line)
// 		}
// 		line += "\n"

// 		if m.cursor == i {
// 			nl := ""
// 			if i > 0 {
// 				nl = "\n"
// 			}
// 			line = nl + boldPrimaryForeground(line) + viewRepositoryDetails(repo)
// 		}

// 		s += line
// 	}

// 	return s + helpView([]helpOption{
// 		{"up/down", "navigate", false},
// 		{"space", "toggle selection", false},
// 		{"d", "delete selected", true},
// 		{"a", "select all", false},
// 		{"n", "deselect all", false},
// 		{"q/esc", "quit", false},
// 	})
// }

func viewRepositoryDetails(repo *forkcleaner.RepositoryWithDetails) string {
	var details []string
	if repo.ParentDeleted {
		details = append(details, "parent was deleted")
	}
	if repo.ParentDMCATakeDown {
		details = append(details, "parent was taken down by DMCA")
	}
	if repo.Private {
		details = append(details, "is private")
	}
	if repo.CommitsAhead > 0 {
		details = append(details, fmt.Sprintf("%d commit%s ahead", repo.CommitsAhead, maybePlural(repo.CommitsAhead)))
	}
	if repo.Forks > 0 {
		details = append(details, fmt.Sprintf("has %d fork%s", repo.Forks, maybePlural(repo.Forks)))
	}
	if repo.Stars > 0 {
		details = append(details, fmt.Sprintf("has %d star%s", repo.Stars, maybePlural(repo.Stars)))
	}
	if repo.OpenPRs > 0 {
		details = append(details, fmt.Sprintf("has %d open PR%s to upstream", repo.OpenPRs, maybePlural(repo.OpenPRs)))
	}
	if time.Now().Add(-30 * 24 * time.Hour).Before(repo.LastUpdate) {
		details = append(details, fmt.Sprintf("recently updated (%s)", timeago.Of(repo.LastUpdate)))
	}

	return strings.Join(details, separator)
}

func maybePlural(n int) string {
	if n == 1 {
		return ""
	}
	return "s"
}

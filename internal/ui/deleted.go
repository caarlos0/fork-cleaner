package ui

// import (
// 	"strconv"

// 	tea "github.com/charmbracelet/bubbletea"
// )

// // NewDeleteEndModelSucceed creates a DeleteEndModel with a success result.
// func NewDeleteEndModelSucceed(deleted int) DeleteEndModel {
// 	return DeleteEndModel{
// 		deleted: deleted,
// 	}
// }

// // NewDeleteEndModelFailed creates a DeleteEndModel with a failed result.
// func NewDeleteEndModelFailed(err error) DeleteEndModel {
// 	return DeleteEndModel{
// 		err: err,
// 	}
// }

// // DeleteEndModel is the UI for when the forks were either deleted or failed
// // to do so.
// type DeleteEndModel struct {
// 	err     error
// 	deleted int
// }

// func (m DeleteEndModel) Init() tea.Cmd {
// 	return nil
// }

// func (m DeleteEndModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
// 	switch msg := msg.(type) {
// 	case tea.KeyMsg:
// 		switch msg.String() {
// 		case "ctrl+c", "q", "esc":
// 			return m, tea.Quit
// 		}
// 	}
// 	return m, nil
// }

// func (m DeleteEndModel) View() string {
// 	if m.deleted > 0 {
// 		return redFaintForeground("Successfully deleted ") + redForeground(strconv.Itoa(m.deleted)) + redFaintForeground(" forks.") +
// 			singleOptionHelp("q", "quit")
// 	}
// 	if m.err != nil {
// 		return errorView("Error deleting repositories", m.err)
// 	}
// 	return ""
// }

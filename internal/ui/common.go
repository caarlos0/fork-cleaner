package ui

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/lipgloss"
)

var (
	defaultStyles = list.NewDefaultItemStyles()

	primaryColor   = defaultStyles.SelectedTitle.GetForeground()
	secondaryColor = defaultStyles.NormalTitle.GetForeground()
	errorColor     = lipgloss.AdaptiveColor{
		Light: "#e94560",
		Dark:  "#f05945",
	}
	listStyle    = lipgloss.NewStyle().Margin(2)
	detailsStyle = lipgloss.NewStyle().PaddingLeft(2)

	boldPrimaryForeground   = lipgloss.NewStyle().Foreground(primaryColor).Bold(true)
	boldSecondaryForeground = lipgloss.NewStyle().Foreground(secondaryColor).Bold(true)
	errorStyle              = lipgloss.NewStyle().Foreground(errorColor)
)

const (
	iconSelected    = "●"
	iconNotSelected = "○"
	separator       = " • "
)

type errMsg struct{ error }

func (e errMsg) Error() string { return e.error.Error() }

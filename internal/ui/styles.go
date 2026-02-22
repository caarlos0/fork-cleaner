package ui

import "charm.land/lipgloss/v2"

var (
	listStyle    = lipgloss.NewStyle().Margin(2)
	detailsStyle = lipgloss.NewStyle().PaddingLeft(2)
)

func errorStyle(fn lipgloss.LightDarkFunc) lipgloss.Style {
	errorColor := fn(lipgloss.Color("#e94560"), lipgloss.Color("#f05945"))
	return lipgloss.NewStyle().Foreground(errorColor)
}

const (
	iconSelected    = "●"
	iconNotSelected = "○"
	separator       = " • "
)

package ui

import "github.com/charmbracelet/bubbles/key"

var (
	keySelectAll       = key.NewBinding(key.WithKeys("a"), key.WithHelp("a", "select all"))
	keySelectNone      = key.NewBinding(key.WithKeys("n"), key.WithHelp("n", "select none"))
	keySelectToggle    = key.NewBinding(key.WithKeys(" "), key.WithHelp("space", "toggle selected item"))
	keyDeletedSelected = key.NewBinding(key.WithKeys("d"), key.WithHelp("d", "delete selected forks"))
)

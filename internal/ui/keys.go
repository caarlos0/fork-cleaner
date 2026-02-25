package ui

import "charm.land/bubbles/v2/key"

var (
	keySelectAll       = key.NewBinding(key.WithKeys("a"), key.WithHelp("a", "select all"))
	keySelectNone      = key.NewBinding(key.WithKeys("n"), key.WithHelp("n", "select none"))
	keySelectToggle    = key.NewBinding(key.WithKeys(" "), key.WithHelp("space", "toggle selected item"))
	keyDeletedSelected = key.NewBinding(key.WithKeys("d"), key.WithHelp("d", "delete selected forks"))
	keyArchiveSelected = key.NewBinding(key.WithKeys("r"), key.WithHelp("r", "archive selected forks"))
)

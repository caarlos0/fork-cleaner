package ui

import "github.com/muesli/termenv"

var (
	primaryColor   = termenv.ColorProfile().Color("205")
	secondaryColor = termenv.ColorProfile().Color("#89F0CB")
)

func bold(s string) string {
	return termenv.String(s).Foreground(primaryColor).Bold().String()
}


type errMsg struct{ error }

func (e errMsg) Error() string { return e.Error() }

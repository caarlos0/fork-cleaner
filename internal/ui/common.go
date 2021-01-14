package ui

import "github.com/muesli/termenv"

var (
	primaryColor   = termenv.ColorProfile().Color("205")
	secondaryColor = termenv.ColorProfile().Color("#89F0CB")
	gray           = termenv.ColorProfile().Color("#626262")
	midGray        = termenv.ColorProfile().Color("#4A4A4A")
)

func boldPrimaryForeground(s string) string {
	return termenv.String(s).Foreground(primaryColor).Bold().String()
}

func boldSecondaryForeground(s string) string {
	return termenv.String(s).Foreground(secondaryColor).Bold().String()
}

func grayForeground(s string) string {
	return termenv.String(s).Foreground(gray).String()
}

func midGrayForeground(s string) string {
	return termenv.String(s).Foreground(midGray).String()
}

type errMsg struct{ error }

func (e errMsg) Error() string { return e.Error() }

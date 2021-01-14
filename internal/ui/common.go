package ui

import (
	"fmt"

	"github.com/muesli/termenv"
)

var (
	primary   = termenv.ColorProfile().Color("205")
	secondary = termenv.ColorProfile().Color("#89F0CB")
	gray      = termenv.ColorProfile().Color("#626262")
	midGray   = termenv.ColorProfile().Color("#4A4A4A")
	red       = termenv.ColorProfile().Color("#ED567A")
)

const (
	iconSelected    = "●"
	iconNotSelected = "○"
)

func boldPrimaryForeground(s string) string {
	return termenv.String(s).Foreground(primary).Bold().String()
}

func boldSecondaryForeground(s string) string {
	return termenv.String(s).Foreground(secondary).Bold().String()
}

func boldRedForeground(s string) string {
	return termenv.String(s).Foreground(red).Bold().String()
}

func redForeground(s string) string {
	return termenv.String(s).Foreground(red).String()
}

func redFaintForeground(s string) string {
	return termenv.String(s).Foreground(red).Faint().String()
}

func grayForeground(s string) string {
	return termenv.String(s).Foreground(gray).String()
}

func midGrayForeground(s string) string {
	return termenv.String(s).Foreground(midGray).String()
}

func faint(s string) string {
	return termenv.String(s).Faint().String()
}

type errMsg struct{ error }

func (e errMsg) Error() string { return e.Error() }

func errorView(action string, err error) string {
	return redForeground(fmt.Sprintf(action+": %s.\nCheck the log file for more details.", err.Error())) + singleOptionHelp("q", "quit")
}

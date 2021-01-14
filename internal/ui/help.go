package ui

import (
	"strings"

	"github.com/muesli/termenv"
)

func helpView() string {
	var s []string

	for _, help := range []struct {
		k, h string
	}{
		{"q/esc", "quit"},
		{"up/down", "navigate"},
		{"space", "toggle selection on current item"},
		{"a", "select all items"},
		{"n", "deselect all items"},
	} {
		if help.k == "d" {
		}
		s = append(s, grayForeground(help.k)+" "+midGrayForeground(help.h))
	}
	s = append(s, grayForeground("d")+" "+ termenv.String("delete selected items").Foreground(secondary).Faint().String())

	var separator = midGrayForeground(" • ")
	return "\n\n" + strings.Join(s, separator)
}

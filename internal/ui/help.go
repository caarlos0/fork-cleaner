package ui

import (
	"strings"

	"github.com/muesli/termenv"
)

func singleOptionHelp(k, v string) string {
	return helpView([]helpOption{
		{k, v, true},
	})
}

func helpView(options []helpOption) string {
	var s []string

	for _, help := range options {
		if help.primary {
			s = append(s, grayForeground(help.key)+" "+termenv.String(help.help).Foreground(secondary).Faint().String())
			continue
		}
		s = append(s, grayForeground(help.key)+" "+midGrayForeground(help.help))
	}

	var separator = midGrayForeground(" • ")
	return "\n\n" + strings.Join(s, separator)
}

type helpOption struct {
	key, help string
	primary   bool
}

package ui

// import (
// 	"strings"

// 	"github.com/muesli/termenv"
// )

// func singleOptionHelp(k, v string) string {
// 	return helpView([]helpOption{
// 		{k, v, true},
// 	})
// }

// var separator = midGrayForeground(" • ")

// func helpView(options []helpOption) string {
// 	var lines []string

// 	var line []string
// 	for i, help := range options {
// 		if help.primary {
// 			line = append(line, grayForeground(help.key)+" "+termenv.String(help.help).Foreground(secondary).Faint().String())
// 		} else {
// 			line = append(line, grayForeground(help.key)+" "+midGrayForeground(help.help))
// 		}
// 		// splits in rows of 3 options max
// 		if (i+1)%3 == 0 {
// 			lines = append(lines, strings.Join(line, separator))
// 			line = []string{}
// 		}
// 	}

// 	// append remainder
// 	lines = append(lines, strings.Join(line, separator))

// 	return "\n\n" + strings.Join(lines, "\n")
// }

// type helpOption struct {
// 	key, help string
// 	primary   bool
// }

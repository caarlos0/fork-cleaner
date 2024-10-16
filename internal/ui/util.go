package ui

func maybePlural(n int) string {
	if n == 1 {
		return ""
	}
	return "s"
}

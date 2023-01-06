package strings

import "strings"

func CenterStringWithSpaces(str string, desiredLength int) string {
	padAmount := desiredLength - len(str)
	if padAmount < 1 {
		return str
	}
	padLeft := padAmount / 2
	padRight := padAmount / 2
	if padAmount%2 == 1 {
		padLeft++
	}
	return strings.Repeat(" ", padLeft) + str + strings.Repeat(" ", padRight)
}

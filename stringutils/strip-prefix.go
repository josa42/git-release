package stringutils

import (
	"strings"
)

func min(x, y int) int {
	if x < y {
		return x
	}
	return y
}

// StripPrefix :
func StripPrefix(str string) string {

	lines := []string{}

	var prefix int

	for idx, line := range strings.Split(str, "\n") {
		if idx == 0 {
			continue
		}
		if idx == 1 {
			prefix = countLeadingSpace(line)
		}

		remove := min(prefix, len(line))

		lines = append(lines, line[remove:])
	}

	return strings.Join(lines, "\n")
}

func countLeadingSpace(line string) int {
	count := 0
	for _, runeValue := range line {
		if runeValue != ' ' {
			break
		}

		count++
	}
	return count
}

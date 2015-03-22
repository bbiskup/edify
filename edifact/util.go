package edifact

import (
//"fmt"
)

// Parser utilities

// Split a string, respecting the specified EDIFACT escape character
func SplitEDIFACT(str string, sep rune, escapeChar rune) []string {
	result := []string{}
	isEscape := false
	var current []rune

	for _, c := range str {
		// fmt.Printf("Pos: %d, rune: %c\n", i, c)

		if c == escapeChar {
			isEscape = true
		}

		if c == sep {
			if !isEscape {
				result = append(result, string(current))
				current = []rune{}
				continue
			} else {
				isEscape = false
			}
		}
		current = append(current, c)
	}
	if current != nil {
		result = append(result, string(current))
	}

	return result
}

// Get indentation
func GetIndent(str string) int {
	var i int
	for _, c := range str {
		if c != ' ' {
			break
		}
		i++
	}
	return i
}

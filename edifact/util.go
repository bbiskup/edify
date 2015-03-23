package edifact

import (
	//"fmt"
	"log"
	"strings"
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

// Splits array of lines by hanging indent
func SplitByHangingIndent(lines []string) [][]string {
	result := [][]string{}
	numLines := len(lines)
	oldIndent := 9999999999
	var currentSection []string
	for i := 0; i < numLines; i++ {
		line := lines[i]
		if len(strings.TrimSpace(line)) == 0 {
			continue
		}
		indent := GetIndent(line)
		log.Printf("indent: %d", indent)

		if indent < oldIndent {
			if currentSection != nil {
				result = append(result, currentSection)
			}
			currentSection = []string{line}
		} else {
			currentSection = append(currentSection, line)
		}
		oldIndent = indent
	}

	if currentSection != nil && len(currentSection) > 0 {
		result = append(result, currentSection)
	}

	return result
}

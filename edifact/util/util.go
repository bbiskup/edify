// Package util provides miscellaneous utilities for handling EDIFACT data
package util

import (
	"errors"
	"strings"
)

const (
	EllipsisStr    = "..."
	lenEllipsisStr = len(EllipsisStr)
)

// SplitEDIFACT splits a string, respecting the specified EDIFACT escape character.
func SplitEDIFACT(str string, sep rune, escapeChar rune) []string {
	result := []string{}
	isEscape := false
	var current []rune

	for _, c := range str {
		// fmt.Printf("Pos: %d, rune: %c\n", i, c)

		if c == '\n' {
			c = ' '
		}

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

// GetIndent determines the indentation of str (i.e. the number of leading
// spaces).
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

// SplitByHangingIndent splits array of lines by hanging indent.
func SplitByHangingIndent(lines []string, splitIndent int) [][]string {
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
		// log.Printf("line: '%s'; indent: %d", line, indent)

		if indent < oldIndent || indent == splitIndent {
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

// RemoveLeadingAndTrailingEmptyLines removes leading and trailing
// empty lines. A line is considered empty if it has zero length.
func RemoveLeadingAndTrailingEmptyLines(lines []string) []string {
	if len(lines) == 0 {
		return lines
	}

	if len(lines[0]) == 0 {
		lines = lines[1:]
	}

	lenLines := len(lines)
	if lenLines > 0 {
		if len(lines[lenLines-1]) == 0 {
			lines = lines[0 : lenLines-1]
		}
	}
	return lines
}

// SplitMultipleLinesByEmptyLines splits non-zero-length lines
// at zero-length lines
func SplitMultipleLinesByEmptyLines(lines []string) [][]string {
	result := [][]string{}

	if len(lines) == 0 {
		return result
	}

	var current []string = []string{}
	for _, line := range lines {
		if len(line) > 0 {
			current = append(current, line)
		} else {
			result = append(result, current)
			current = []string{}
		}
	}
	result = append(result, current)
	return result
}

// CustBoolStr returns custom strings for boolean value.
func CustBoolStr(value bool, trueStr string, falseStr string) string {
	if value {
		return trueStr
	} else {
		return falseStr
	}
}

// Ellipsis converts a string into substring of specified max. length if too long.
func Ellipsis(str string, maxLen int) string {
	lenStr := len(str)
	if lenStr <= maxLen {
		return str
	} else {
		if lenStr <= lenEllipsisStr {
			return EllipsisStr
		} else {
			return str[:maxLen-lenEllipsisStr] + "..."
		}

	}
}

// JoinByHangingIndent joins lines indented beyond the specified base
// indent with the previous line
func JoinByHangingIndent(lines []string, baseIndent int, collapseSpaces bool) []string {
	result := []string{}
	current := []string{}

	concat := func(tokens []string) string {
		if collapseSpaces {
			trimmed := []string{}
			for _, token := range tokens {
				trimmed = append(trimmed, strings.TrimSpace(token))
			}
			return strings.Join(trimmed, " ")
		} else {
			return strings.Join(tokens, "")
		}
	}

	for _, line := range lines {
		indent := GetIndent(line)
		if indent <= baseIndent {
			if len(current) > 0 {
				result = append(result, concat(current))
			}
			current = []string{line}
		} else {
			current = append(current, line)
		}
	}
	if len(current) > 0 {
		result = append(result, concat(current))
	}
	return result
}

// TrimWhiteSpaceAndJoin trims each string passed to it and then
// joins the trimmed strings using joinStr
func TrimWhiteSpaceAndJoin(lines []string, joinStr string) string {
	trimmed := []string{}
	for _, line := range lines {
		trimmed = append(trimmed, strings.TrimSpace(line))
	}
	return strings.Join(trimmed, joinStr)
}

// CheckNotNil returns nil if none of the passed values is nil,
// otherwise an error
func CheckNotNil(values ...interface{}) error {
	for _, val := range values {
		if val == nil {
			return errors.New("Not-nil check failed")
		}
	}
	return nil
}

// Unused can be used to avoid compiler errors for unused variables
func Unused(values ...interface{}) {

}

// GetIndentStr creates an indentation strings (consisting of 'indent' spaces)
func GetIndentStr(indent int) string {
	return strings.Repeat("  ", indent)
}

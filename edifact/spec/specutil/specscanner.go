package specutil

import (
	"bufio"
	"fmt"
	"os"
	"runtime"
	"strings"
)

const (
	// Separator between specifications (partial)
	specSep = "--------------------"
)

// scanner for spec files with entries separated by "-----...." lines
type SpecScanner struct {
	file        *os.File
	scanner     *bufio.Scanner
	HasMore     bool
	HeaderLines []string
}

// Creates a new spec scanner, given a file name
func NewSpecScanner(fileName string) (*SpecScanner, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}

	result, err := NewSpecScannerFromReader(bufio.NewReader(file))
	runtime.SetFinalizer(result, finalizer)
	return result, err
}

func (s *SpecScanner) Err() error {
	if s != nil {
		return s.scanner.Err()
	} else {
		return nil
	}
}

func (s *SpecScanner) String() string {
	var name string
	if s.file != nil {
		name = s.file.Name()
	} else {
		name = fmt.Sprintf("%#v", s.scanner)
	}
	return "SpecScanner [" + name + "]"
}

// fetch all lines up to next spec separator
func (s *SpecScanner) GetNextSpecLines(skipEmptyLines bool) (lines []string, err error) {
	for {
		scanResult := s.scanner.Scan()
		if !scanResult {
			if s.scanner.Err() == nil {
				// EOF
				if s.file != nil {
					s.file.Close()
					s.file = nil
				}
				s.HasMore = false
				return lines, nil
			}
		}
		err := s.scanner.Err()
		if err != nil {
			s.HasMore = false
			return nil, err
		}

		line := s.scanner.Text()
		strippedLine := strings.TrimSpace(line)
		if skipEmptyLines && len(strippedLine) == 0 {
			continue
		}

		if strings.HasPrefix(line, specSep) {
			s.HasMore = true
			return lines, nil
		}

		lines = append(lines, line)
	}
}

// return all spec lines at once
func (s *SpecScanner) GetAllSpecLines(skipEmptyLines bool) (linesGroups [][]string, err error) {
	result := [][]string{}
	for s.HasMore {
		specLines, err := s.GetNextSpecLines(skipEmptyLines)
		if err != nil {
			return nil, err
		}
		result = append(result, specLines)
	}
	return result, nil
}

// Create a scanner from a provided reader (e.g. for testing)
func NewSpecScannerFromReader(reader *bufio.Reader) (scanner *SpecScanner, err error) {
	result := &SpecScanner{
		file:    nil,
		scanner: bufio.NewScanner(reader),
		HasMore: true,
	}

	result.HeaderLines, err = result.GetNextSpecLines(true)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func finalizer(s *SpecScanner) {
	// log.Printf("Running finalizer for %#v", s)
	if s.file != nil {
		s.file.Close()
		s.file = nil
	}
}

// Signature for callback to parse
type ParseSection func(lines []string) error

func ParseSpecFile(fileName string, parseSection ParseSection) error {
	scanner, err := NewSpecScanner(fileName)
	if err != nil {
		return err
	}

	for {
		// read specification parts
		specLines, err := scanner.GetNextSpecLines(false)

		if err != nil {
			return err
		}

		if !scanner.HasMore && len(specLines) == 0 {
			// log.Println("No more lines")
			break
		}

		// log.Printf("specLines: \n%s\n", specLines)
		err = parseSection(specLines)
		if err != nil {
			return err
		}
	}
	return nil
}

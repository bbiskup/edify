package specparser

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"runtime"
	"strings"
)

// scanner for spec files with entries separated by "-----...." lines
type SpecScanner struct {
	file    *os.File
	scanner *bufio.Scanner
	hasMore bool
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
func (s *SpecScanner) GetNextSpecLines() (lines []string, err error) {
	for {
		if !s.hasMore {
			return nil, errors.New("No more data")
		}
		scanResult := s.scanner.Scan()
		if !scanResult {
			if s.scanner.Err() == nil {
				// EOF
				if s.file != nil {
					s.file.Close()
					s.file = nil
				}
				s.hasMore = false
				return lines, nil
			}
		}
		err := s.scanner.Err()
		if err != nil {
			s.hasMore = false
			return nil, err
		}

		line := s.scanner.Text()
		strippedLine := strings.TrimSpace(line)
		if len(strippedLine) == 0 {
			continue
		}

		if strings.HasPrefix(line, specSep) {
			s.hasMore = true
			return lines, nil
		}

		lines = append(lines, line)
	}
	s.hasMore = true
	return lines, nil
}

// return all spec lines at once
func (s *SpecScanner) GetAllSpecLines() (linesGroups [][]string, err error) {
	result := [][]string{}
	for s.hasMore {
		specLines, err := s.GetNextSpecLines()
		if err != nil {
			return nil, err
		}
		result = append(result, specLines)
	}
	return result, nil
}

// Creates a new spec scanner, given a file name
func NewSpecScanner(fileName string) (*SpecScanner, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}

	result, err := NewSpecScannerFromReader(bufio.NewReader(file)), err
	runtime.SetFinalizer(result, finalizer)
	return result, err
}

// Create a scanner from a provided reader (e.g. for testing)
func NewSpecScannerFromReader(reader *bufio.Reader) *SpecScanner {
	return &SpecScanner{
		file:    nil,
		scanner: bufio.NewScanner(reader),
		hasMore: true,
	}
}

func finalizer(s *SpecScanner) {
	log.Printf("Running finalizer for %#v", s)
	if s.file != nil {
		s.file.Close()
		s.file = nil
	}
}

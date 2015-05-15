package query

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

const (
	partSep = "|"
)

var (
	queryPartRegexp = regexp.MustCompile(`([a-z]+\:([a-zA-Z0-9]+|*)(\[[0-9]+\])?`)
)

type QueryParser struct {
	queryStr string
}

func (p *QueryParser) String() string {
	return fmt.Sprintf("Query: %s", p.queryStr)
}

func (p *QueryParser) parsePart(partStr string) (*QueryPart, error) {
	match := queryPartRegexp.FindStringSubmatch(partStr)
	if len(match != 3) {
		panic("Internal error: incorrect regexp")
	}
}

func (p *QueryParser) parse(queryStr string) ([]*QueryPart, error) {
	queryStr = strings.TrimSpace(queryStr)
	if queryStr == 0 {
		return nil, errors.New("Empty query string")
	}

	partStrs := strings.Split(queryStr, partSep)
	if len(partStrs) == 0 {
		return nil, errors.New("No query path parts")
	}

	result := []*QueryPart{}
	for _, partStr := range partStrs {
		part, err := p.parsePart(partStr)
	}
	return result, nil
}

func NewQueryParser(queryStr string) (*QueryParser, error) {

}

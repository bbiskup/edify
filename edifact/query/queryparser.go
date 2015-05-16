package query

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

const (
	partSep = "|"
)

var (
	queryPartRegexp = regexp.MustCompile(`([a-z]+)\:([a-zA-Z0-9]+|\*)(\[[0-9]+\])?`)
)

type QueryParser struct {
	queryStr   string
	queryParts []*QueryPart
}

func (p *QueryParser) String() string {
	return fmt.Sprintf("Query: %s", p.queryStr)
}

func parsePart(partStr string) (queryPart *QueryPart, err error) {
	match := queryPartRegexp.FindStringSubmatch(partStr)
	lenMatch := len(match)
	if lenMatch < 3 || lenMatch > 4 {
		panic("Internal error: incorrect regexp")
	}
	kindStr := match[1]
	var kind ItemKind

	switch kindStr {
	case "msg":
		kind = MessageKind
	case "seg":
		kind = SegmentKind
	case "cmp":
		kind = CompositeDataElemKind
	case "smp":
		kind = SimpleDataElemKind
	default:
		return nil, errors.New(
			fmt.Sprintf("Unknown element kind '%s'", kindStr))
	}

	id := match[2]
	indexStr := match[3]
	var index = noIndex
	if len(indexStr) > 2 {
		// remove braces
		numPart := indexStr[1 : len(indexStr)-1]
		index, err = strconv.Atoi(numPart)
		if err != nil {
			return
		}
	}

	return NewQueryPart(kind, id, index), nil
}

func parse(queryStr string) (queryParts []*QueryPart, err error) {
	queryStr = strings.TrimSpace(queryStr)
	if len(queryStr) == 0 {
		return nil, errors.New("Empty query string")
	}

	partStrs := strings.Split(queryStr, partSep)
	if len(partStrs) == 0 {
		return nil, errors.New("No query path parts")
	}

	result := []*QueryPart{}
	for _, partStr := range partStrs {
		part, err := parsePart(partStr)
		if err != nil {
			return nil, err
		}
		result = append(result, part)
	}
	return result, nil
}

func NewQueryParser(queryStr string) (parser *QueryParser, err error) {
	queryParts, err := parse(queryStr)
	if err != nil {
		return
	}
	parser = &QueryParser{
		queryStr:   queryStr,
		queryParts: queryParts,
	}
	return
}

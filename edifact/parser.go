package edifact

import (
	"errors"
	"strings"
)

const (
	UNH               = "UNH"
	UNT               = "UT"
	SegTerm           = "'"
	SegTagDataElemSep = "+"
)

type Parser struct {
	err error
}

func NewParser() *Parser {
	return &Parser{err: nil}
}

func (p *Parser) ParseElements(elementStrs []string) (element []*Element) {
	if p.err != nil {
		return nil
	}
	return nil
}

func (p *Parser) ParseSegment(segmentStr string) (segment *Segment) {
	if p.err != nil {
		return nil
	}
	parts := strings.Split(segmentStr, SegTagDataElemSep)
	if len(parts) < 2 {
		p.err = errors.New("Segment too short")
		return nil
	}
	segmentName := parts[0]
	segment = NewSegment(segmentName)

	elements := p.ParseElements(parts[1:])
	segment.AddElements(elements)
	return segment
}

func (p *Parser) ParseSegments(segmentStrs []string) []*Segment {
	if p.err != nil {
		return nil
	}
	result := []*Segment{}

	for _, segmentStr := range segmentStrs {
		result = append(result, p.ParseSegment(segmentStr))
		if p.err != nil {
			return nil
		}
	}

	return result
}

func (p *Parser) ParseMessage(edifactMessage string) (message *Message, err error) {
	// reset error
	p.err = nil

	segmentStrs := strings.Split(edifactMessage, SegTerm)
	segments := p.ParseSegments(segmentStrs)

	if len(segments) < 2 {
		p.err = errors.New("Message header and/or tail missing")
		return nil, p.err
	}

	if segments[0].Name != UNH {
		p.err = errors.New("No header message header (UNH)")
		return nil, p.err
	}

	if segments[len(segments)-1].Name != UNT {
		p.err = errors.New("No header message tail (UNT)")
		return nil, p.err
	}

	message = NewMessage("dummyname")

	return message, p.err
}

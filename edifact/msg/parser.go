package msg

import (
	"errors"
	"fmt"
	"github.com/bbiskup/edify/edifact/util"
	"log"
	"strings"
)

type Parser struct {
	err error
}

func NewParser() *Parser {
	return &Parser{err: nil}
}

func (p *Parser) ParseElement(elementStr string) (element *DataElement) {
	if strings.Index(elementStr, RepetitionSep) != -1 {
		p.err = errors.New("Data element repetition currently not supported")
		return nil
	}
	parts := strings.Split(elementStr, CompDataElemSepStr)
	return NewDataElement(parts)
}

func (p *Parser) ParseElements(elementStrs []string) (elements []*DataElement) {
	fmt.Printf("ParseElements %s", elementStrs)
	if p.err != nil {
		return nil
	}

	elements = []*DataElement{}
	for _, elementStr := range elementStrs {
		element := p.ParseElement(elementStr)
		if element == nil {
			return nil
		}
		if p.err != nil {
			return nil
		}
		elements = append(elements, element)
	}
	return
}

func (p *Parser) ParseSegment(segmentStr string) (segment *Segment) {
	if p.err != nil {
		return nil
	}

	segmentStr = strings.Trim(segmentStr, " \t\n")
	if len(segmentStr) == 0 {
		// empty segment not treated as error
		return nil
	}

	parts := util.SplitEDIFACT(segmentStr, SegTagDataElemSep, ReleaseChar)
	if len(parts) < 2 {
		p.err = errors.New(fmt.Sprintf("Segment too short (%#v)", parts))
		return nil
	}
	segmentId := parts[0]
	segment = NewSegment(segmentId)

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
		segment := p.ParseSegment(segmentStr)
		if segment == nil {
			continue
		}
		result = append(result, segment)
		if p.err != nil {
			return nil
		}
	}

	return result
}

func (p *Parser) ParseMessage(edifactMessage string) (message *Message, err error) {
	log.Printf("Parsing message")
	// reset error
	p.err = nil

	segmentStrs := util.SplitEDIFACT(edifactMessage, SegTerm, ReleaseChar)
	log.Print("segmentStrs: ", segmentStrs)
	segments := p.ParseSegments(segmentStrs)

	if p.err != nil {
		return nil, p.err
	}

	// log.Printf("Segments: %s", segments)

	if len(segments) < 2 {
		p.err = errors.New("Message header and/or tail missing")
		return nil, p.err
	}

	if segments[0].Id != UNH {
		p.err = errors.New("No message header (UNH)")
		return nil, p.err
	}

	tailName := segments[len(segments)-1].Id
	if tailName != UNT {
		log.Print("tail segment: ", tailName)

		p.err = errors.New("No message tail (UNT)")
		return nil, p.err
	}

	message = NewMessage("dummyname", segments)

	return message, p.err
}

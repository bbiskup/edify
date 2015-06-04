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

func (p *Parser) ParseElem(elementStr string) (element *DataElem) {
	if strings.Index(elementStr, RepetitionSep) != -1 {
		p.err = errors.New("Data element repetition currently not supported")
		return nil
	}
	parts := strings.Split(elementStr, CompDataElemSepStr)
	return NewDataElem(parts)
}

func (p *Parser) ParseElems(elementStrs []string) (elements []*DataElem) {
	fmt.Printf("ParseElems %s", elementStrs)
	if p.err != nil {
		return nil
	}

	elements = []*DataElem{}
	for _, elementStr := range elementStrs {
		element := p.ParseElem(elementStr)
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

func (p *Parser) ParseSeg(segmentStr string) (segment *Seg) {
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
		p.err = errors.New(fmt.Sprintf("Seg too short (%#v)", parts))
		return nil
	}
	segmentId := parts[0]
	segment = NewSeg(segmentId)

	elements := p.ParseElems(parts[1:])
	segment.AddElems(elements)
	return segment
}

func (p *Parser) ParseSegs(segmentStrs []string) []*Seg {
	if p.err != nil {
		return nil
	}
	result := []*Seg{}

	for _, segmentStr := range segmentStrs {
		segment := p.ParseSeg(segmentStr)
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

func (p *Parser) ParseRawMsg(edifactMessage string) (rawMessage *RawMsg, err error) {
	log.Printf("Parsing raw message")
	// reset error
	p.err = nil

	segmentStrs := util.SplitEDIFACT(edifactMessage, SegTerm, ReleaseChar)
	log.Print("segmentStrs: ", segmentStrs)
	segments := p.ParseSegs(segmentStrs)

	if p.err != nil {
		return nil, p.err
	}

	// log.Printf("Segs: %s", segments)

	if len(segments) < 2 {
		p.err = errors.New("Raw message header and/or tail missing")
		return nil, p.err
	}

	if segments[0].Id() != UNH {
		p.err = errors.New("No message header (UNH)")
		return nil, p.err
	}

	tailName := segments[len(segments)-1].Id()
	if tailName != UNT {
		log.Print("tail segment: ", tailName)

		p.err = errors.New("No message tail (UNT)")
		return nil, p.err
	}

	rawMessage = NewRawMsg("dummyname", segments)

	return rawMessage, p.err
}

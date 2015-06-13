package rawmsg

import (
	"errors"
	"fmt"
	"github.com/bbiskup/edify/edifact/util"
	"io/ioutil"
	"log"
	"strings"
)

type Parser struct {
	err error
}

func NewParser() *Parser {
	return &Parser{err: nil}
}

func (p *Parser) ParseElem(elementStr string) (element *RawDataElem) {
	if strings.Index(elementStr, RepetitionSep) != -1 {
		p.err = errors.New("Data element repetition currently not supported")
		return nil
	}
	parts := strings.Split(elementStr, CompDataElemSepStr)
	return NewRawDataElem(parts)
}

func (p *Parser) ParseElems(elementStrs []string) (elements []*RawDataElem) {
	fmt.Printf("ParseElems %s", elementStrs)
	if p.err != nil {
		return nil
	}

	elements = []*RawDataElem{}
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

func (p *Parser) ParseRawSeg(rawSegStr string) (rawSeg *RawSeg) {
	if p.err != nil {
		return nil
	}

	rawSegStr = strings.Trim(rawSegStr, " \t\n")
	if len(rawSegStr) == 0 {
		// empty segment not treated as error
		return nil
	}

	parts := util.SplitEDIFACT(rawSegStr, SegTagDataElemSep, ReleaseChar)
	if len(parts) < 2 {
		p.err = fmt.Errorf("RawSeg too short (%#v)", parts)
		return nil
	}
	rawSegId := parts[0]
	rawSeg = NewRawSeg(rawSegId)

	elements := p.ParseElems(parts[1:])
	rawSeg.AddElems(elements)
	return rawSeg
}

func (p *Parser) ParseRawSegs(rawSegStrs []string) []*RawSeg {
	if p.err != nil {
		return nil
	}
	result := []*RawSeg{}

	for _, rawSegStr := range rawSegStrs {
		rawSeg := p.ParseRawSeg(rawSegStr)
		if rawSeg == nil {
			continue
		}
		result = append(result, rawSeg)
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

	rawSegStrs := util.SplitEDIFACT(edifactMessage, SegTerm, ReleaseChar)
	log.Print("rawSegStrs: ", rawSegStrs)
	segments := p.ParseRawSegs(rawSegStrs)

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

func (p *Parser) ParseRawMsgFile(fileName string) (rawMessage *RawMsg, err error) {
	msgFileContents, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, err
	}
	return p.ParseRawMsg(string(msgFileContents))
}

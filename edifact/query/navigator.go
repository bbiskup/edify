package query

import (
	"errors"
	"fmt"
	"github.com/bbiskup/edify/edifact/msg"
	"log"
)

// Provides navigation within EDIFACT message
type Navigator struct {
}

// return segment, segment group or data element
func (n *Navigator) navigate(queryStr string, nestedMsg *msg.NestedMsg) (msgPart msg.NestedMsgPart, err error) {

	log.Printf("navigate msg: %s, query %s", nestedMsg.Name, queryStr)

	//var currentSeg *msg.Seg
	currentGrp := nestedMsg.GetTopLevelGrp()
	var currentMsgPart msg.NestedMsgPart = currentGrp

	queryParser, err := NewQueryParser(queryStr)
	if err != nil {
		return nil, errors.New(fmt.Sprintf(
			"Query failed: %s", queryStr))
	}
	for _, queryPart := range queryParser.queryParts {
		switch queryPart.ItemKind {
		case SegGrpKind:
			repGrp := currentGrp.GetPartByKey(queryPart.Id)
			if repGrp == nil {
				return nil, errors.New(fmt.Sprintf(
					"Segment group '%s' not found in '%s'", queryPart.Id, currentGrp.Id()))
			}

			repGrp, ok := repGrp.(*msg.RepSegGrp)
			if !ok {
				return nil, errors.New(fmt.Sprintf(
					"Part %s is not a segment group, but '%T'", repGrp.Id(), repGrp))
			}
			currentMsgPart = repGrp

		case SegKind:
			seg, err := currentGrp.FindNthOccurrenceOfSeg(queryPart.Id, 0)
			if err != nil {
				return nil, errors.New(fmt.Sprintf(
					"Segment %s (Occurrence #%d) could not be found in group %s",
					queryPart.Id, queryPart.Index, currentGrp.Id()))
			}
			numSegs := seg.Count()
			if queryPart.Index >= numSegs {
				return nil, errors.New(fmt.Sprintf(
					"Repeat index %d out of range for segment %s (max: %d)",
					queryPart.Index, queryPart.Id, numSegs))
			}
			nthSeg := seg.GetSeg(queryPart.Index)
			//currentSeg = nthSeg
			currentMsgPart = nthSeg

		case CompositeDataElemKind:
			panic("Not implemented")

		case SimpleDataElemKind:
			panic("Not implemented")
		}
	}
	return currentMsgPart, nil
}

func (n *Navigator) GetSegGrp(
	queryStr string, message *msg.NestedMsg) (group *msg.SegGrp, err error) {

	msgPart, err := n.navigate(queryStr, message)
	if err != nil {
		return nil, err
	}
	group, ok := msgPart.(*msg.SegGrp)
	if !ok {
		return nil, errors.New(fmt.Sprintf("Unexpected type %T", group))
	}
	return group, nil
}

func (n *Navigator) GetSeg(queryStr string, message *msg.NestedMsg) (*msg.Seg, error) {
	msgPart, err := n.navigate(queryStr, message)
	if err != nil {
		return nil, err
	}
	segment, ok := msgPart.(*msg.Seg)
	if !ok {
		return nil, errors.New(fmt.Sprintf("Unexpected type %T", segment))
	}
	return segment, nil
}

func (n *Navigator) GetSegDataElem(queryStr string, message *msg.NestedMsg) (msg.DataElem, error) {
	msgPart, err := n.navigate(queryStr, message)
	if err != nil {
		return nil, err
	}
	dataElem, ok := msgPart.(msg.DataElem)
	if !ok {
		return nil, errors.New(fmt.Sprintf("Unexpected type %T", dataElem))
	}
	return dataElem, nil
}

func NewNavigator() *Navigator {
	return &Navigator{}
}

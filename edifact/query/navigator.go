package query

import (
	"fmt"
	"github.com/bbiskup/edify/edifact/msg"
	"log"
)

type SimpleDataElemGetter interface {
	Id() string
	GetSimpleDataElemById(dataElemId string) (*msg.SimpleDataElem, error)
}

// Provides navigation within EDIFACT message
type Navigator struct {
}

func NewNavigator() *Navigator {
	return &Navigator{}
}

// return segment, segment group or data element
func (n *Navigator) Navigate(queryStr string, nestedMsg *msg.NestedMsg) (msgPart msg.NestedMsgPart, err error) {

	log.Printf("Navigate msg: %s, query %s", nestedMsg.Name, queryStr)

	var currentSeg *msg.Seg
	currentGrp := nestedMsg.GetTopLevelGrp()
	var currentMsgPart msg.NestedMsgPart = currentGrp

	queryParser, err := NewQueryParser(queryStr)
	if err != nil {
		return nil, fmt.Errorf("Query '%s' failed: %s", queryStr, err)
	}
	for _, queryPart := range queryParser.queryParts {
		switch queryPart.ItemKind {
		case SegGrpKind:
			repGrpPart := currentGrp.GetPartByKey(queryPart.Id)
			if repGrpPart == nil {
				return nil, fmt.Errorf(
					"Segment group '%s' not found in '%s'", queryPart.Id, currentGrp.Id())
			}

			repGrp, ok := repGrpPart.(*msg.RepSegGrp)
			if !ok {
				return nil, fmt.Errorf("Part %s is not a segment group, but '%T'", repGrp.Id(), repGrp)
			}
			numSegGrps := repGrp.Count()
			if queryPart.Index >= numSegGrps {
				return nil, fmt.Errorf(
					"Repeat index %d out of range for segment group %s (max: %d)",
					queryPart.Index, queryPart.Id, numSegGrps)
			}
			nthGrp := repGrp.GetSegGrp(queryPart.Index)
			currentMsgPart = nthGrp
			currentGrp = nthGrp // Descend
			currentSeg = nil

		case SegKind:
			seg, err := currentGrp.FindNthOccurrenceOfSeg(queryPart.Id, 0)
			if err != nil {
				return nil, fmt.Errorf(
					"Segment %s (Occurrence #%d) could not be found in group %s",
					queryPart.Id, queryPart.Index, currentGrp.Id())
			}
			numSegs := seg.Count()
			if queryPart.Index >= numSegs {
				return nil, fmt.Errorf(
					"Repeat index %d out of range for segment %s (max: %d)",
					queryPart.Index, queryPart.Id, numSegs)
			}
			nthSeg := seg.GetSeg(queryPart.Index)
			currentSeg = nthSeg
			currentMsgPart = nthSeg

		case CompositeDataElemKind:
			_, ok := currentMsgPart.(*msg.Seg)
			if !ok {
				return nil, fmt.Errorf(
					"Parent %s is not a segment", currentMsgPart.Id())
			}

			numDataElems := len(currentSeg.DataElems)
			if queryPart.Index >= numDataElems {
				return nil, fmt.Errorf(
					"data element index %d out of range for segment %s (max: %d)",
					queryPart.Index, queryPart.Id, numDataElems)
			}

			compositeDataElem, err := currentSeg.GetCompositeDataElemById(queryPart.Id)
			if err != nil {
				return nil, err
			}
			currentMsgPart = compositeDataElem

		case SimpleDataElemKind:
			// Could either be part of a segment or of a composite data element

			currentMsgPartSub, ok := currentMsgPart.(SimpleDataElemGetter)
			if !ok {
				return nil, fmt.Errorf(
					"Parent element %s of simple data element %s neither segment nor composite data element",
					currentMsgPart.Id(), queryPart.Id)
			}
			currentMsgPart, err = currentMsgPartSub.GetSimpleDataElemById(queryPart.Id)
			if err != nil {
				return nil, err
			}

		}
	}
	return currentMsgPart, nil
}

func (n *Navigator) GetSegGrp(
	queryStr string, message *msg.NestedMsg) (group *msg.SegGrp, err error) {

	msgPart, err := n.Navigate(queryStr, message)
	if err != nil {
		return nil, err
	}
	group, ok := msgPart.(*msg.SegGrp)
	if !ok {
		return nil, fmt.Errorf("Unexpected type %T", group)
	}
	return group, nil
}

func (n *Navigator) GetSeg(queryStr string, message *msg.NestedMsg) (*msg.Seg, error) {
	msgPart, err := n.Navigate(queryStr, message)
	if err != nil {
		return nil, err
	}
	segment, ok := msgPart.(*msg.Seg)
	if !ok {
		return nil, fmt.Errorf("Unexpected type %T", segment)
	}
	return segment, nil
}

func (n *Navigator) GetSegDataElem(queryStr string, message *msg.NestedMsg) (msg.DataElem, error) {
	msgPart, err := n.Navigate(queryStr, message)
	if err != nil {
		return nil, err
	}
	dataElem, ok := msgPart.(msg.DataElem)
	if !ok {
		return nil, fmt.Errorf("Unexpected type %T", dataElem)
	}
	return dataElem, nil
}

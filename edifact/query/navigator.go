package query

import (
	"errors"
	"fmt"
	"github.com/bbiskup/edify/edifact/msg"
)

// Provides navigation within EDIFACT message
type Navigator struct {
}

func (n *Navigator) navigate(queryStr string, message *msg.NestedMsg) (msgPart msg.NestedMsgPart, err error) {
	// return segment, segment group or data element
	panic("Not implemented")
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

func (n *Navigator) GetSegDataElem(queryStr string, message *msg.NestedMsg) (*msg.DataElem, error) {
	msgPart, err := n.navigate(queryStr, message)
	if err != nil {
		return nil, err
	}
	dataElem, ok := msgPart.(*msg.DataElem)
	if !ok {
		return nil, errors.New(fmt.Sprintf("Unexpected type %T", dataElem))
	}
	return dataElem, nil
}

func NewNavigator() *Navigator {
	return &Navigator{}
}
